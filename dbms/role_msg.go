package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/game/login"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (a *RoleActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CRegist:
		arg := msg.(*pb.CRegist)
		glog.Debugf("CRegist %#v", arg)
		ip := arg.GetIpaddr()
		rsp, user := login.Regist(arg, a.uniqueid)
		ctx.Respond(rsp)
		a.HandlerLogin(user, true, ip, ctx)
	case *pb.CLogin:
		arg := msg.(*pb.CLogin)
		glog.Debugf("CLogin %#v", arg)
		ip := arg.GetIpaddr()
		rsp, user := login.Login(arg)
		ctx.Respond(rsp)
		a.HandlerLogin(user, false, ip, ctx)
	case *pb.CWxLogin:
		arg := msg.(*pb.CWxLogin)
		glog.Debugf("CWxLogin %#v", arg)
		ip := arg.GetIpaddr()
		rsp, user := login.WxLogin(arg, a.uniqueid)
		isRegist := rsp.GetIsreg()
		ctx.Respond(rsp)
		a.HandlerLogin(user, isRegist, ip, ctx)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout userid: %s", arg.Userid)
		userid := arg.GetUserid()
		if v, ok := a.roles[userid]; ok {
			a.offline[userid] = v
			//移除
			delete(a.roles, userid)
		}
		delete(a.router, arg.Sender.String())
	case *pb.Login:
		//登录成功
		arg := msg.(*pb.Login)
		a.router[ctx.Sender().String()] = arg.Userid
		//响应登录
		rsp := new(pb.Logined)
		ctx.Respond(rsp)
	case *pb.HallConnect:
		//初始化建立连接
		glog.Infof("role init: %v", ctx.Self().String())
		//连接
		bind := cfg.Section("hall").Key("bind").Value()
		name := cfg.Section("cookie").Key("name").Value()
		timeout := 3 * time.Second
		a.hallPid, err = remote.SpawnNamed(bind, a.Name, name, timeout)
		glog.Infof("a.hallPid: %s", a.hallPid.String())
		if err != nil {
			glog.Fatalf("remote hall err %v", err)
		}
		connect := &pb.HallConnect{
			Sender: ctx.Self(),
			Name:   a.Name,
		}
		a.hallPid.Tell(connect)
	case *pb.ServeStop:
		//关闭服务
		a.HandlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.CBuy:
		arg := msg.(*pb.CBuy)
		user := a.getUser(ctx)
		//响应
		rsp := handler.Buy(arg, user)
		ctx.Respond(rsp)
	case *pb.CShop:
		arg := msg.(*pb.CShop)
		user := a.getUser(ctx)
		//响应
		rsp := handler.Shop(arg, user)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *RoleActor) getUser(ctx actor.Context) *data.User {
	userid := a.router[ctx.Sender().String()]
	user := a.roles[userid]
	return user
}

//登录处理
func (a *RoleActor) HandlerLogin(user *data.User,
	isRegist bool, ip string, ctx actor.Context) {
	if user == nil {
		//登录失败
		return
	}
	a.HandlerLogined(user)
	if isRegist {
		//注册奖励发放
		var diamond int32 = config.GetEnv(data.ENV1)
		var coin int32 = config.GetEnv(data.ENV2)
		a.HandlerPrize(user, diamond, coin, data.LogType1, ctx)
		//注册日志
		msg1 := &pb.LogRegist{
			Userid:   user.Userid,
			Nickname: user.Nickname,
			Atype:    user.Atype,
			Ip:       ip,
		}
		nodePid.Tell(msg1)
	}
	//TODO 登录日志有可能在登出消息之前,所以暂时这里处理
	msg3 := &pb.LogLogout{
		Userid: user.Userid,
		Event:  4,
	}
	nodePid.Tell(msg3)
	//登录日志
	msg2 := &pb.LogLogin{
		Userid: user.Userid,
		Atype:  user.Atype,
		Ip:     ip,
	}
	nodePid.Tell(msg2)
	glog.Debugf("user %#v", user)
	a.HandlerSync(user, ctx)
}

func (a *RoleActor) HandlerSync(user *data.User, ctx actor.Context) {
	//同步数据,只有登录时才向节点同步数据
	//其它时候为节点向role同步,避免数据覆盖
	//TODO 定时回存数据
	msg3 := new(pb.SyncUser)
	msg3.Userid = user.Userid
	result, err := json.Marshal(user)
	if err != nil {
		glog.Errorf("user Marshal err %v", err)
	}
	msg3.Data = result
	ctx.Sender().Tell(msg3)
}

func (a *RoleActor) HandlerLogined(user *data.User) {
	//已经在线,用在线数据
	if v, ok := a.roles[user.Userid]; ok && v != nil {
		user = v
		return
	}
	//已经离线,用离线数据
	if v, ok := a.offline[user.Userid]; ok && v != nil {
		user = v
		delete(a.offline, user.Userid)
	}
	//登录成功
	a.roles[user.Userid] = user
	glog.Debugf("Logoin userid: %s", user.Userid)
	glog.Debugf("roles len: %d", len(a.roles))
	glog.Debugf("offline len: %d", len(a.offline))
	//TODO router 可直接在这里处理
}

//奖励发放
func (a *RoleActor) HandlerPrize(user *data.User,
	diamond, coin int32, ltype int, ctx actor.Context) {
	user.AddDiamond(diamond)
	user.AddCoin(coin)
	//消息
	msg := &pb.SPushCurrency{
		Rtype:   uint32(ltype),
		Diamond: diamond,
		Coin:    coin,
	}
	ctx.Respond(msg)
	if diamond != 0 {
		//日志
		msg1 := &pb.LogDiamond{
			Userid: user.GetUserid(),
			Type:   int32(ltype),
			Num:    diamond,
			Rest:   user.GetDiamond(),
		}
		nodePid.Tell(msg1)
	}
	if coin != 0 {
		//日志
		msg1 := &pb.LogCoin{
			Userid: user.GetUserid(),
			Type:   int32(ltype),
			Num:    coin,
			Rest:   user.GetCoin(),
		}
		nodePid.Tell(msg1)
	}
}

func (a *RoleActor) HandlerStop(ctx actor.Context) {
	glog.Debugf("HandlerStop: %s", a.Name)
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, v := range a.offline {
		glog.Debugf("Stop offline: %s", k)
		v.Save()
	}
	for k, v := range a.roles {
		glog.Debugf("Stop online: %s", k)
		v.Save()
	}
}
