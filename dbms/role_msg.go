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
		//ip := arg.GetIpaddr()
		rsp, user := login.Regist(arg, a.uniqueid)
		ctx.Respond(rsp)
		a.handlerLogin(user, true, ctx)
	case *pb.CLogin:
		arg := msg.(*pb.CLogin)
		glog.Debugf("CLogin %#v", arg)
		//ip := arg.GetIpaddr()
		rsp, user := login.Login(arg)
		ctx.Respond(rsp)
		a.handlerLogin(user, false, ctx)
	case *pb.CWxLogin:
		arg := msg.(*pb.CWxLogin)
		glog.Debugf("CWxLogin %#v", arg)
		//ip := arg.GetIpaddr()
		rsp, user := login.WxLogin(arg, a.uniqueid)
		isRegist := rsp.GetIsreg()
		ctx.Respond(rsp)
		//顺序不能错
		a.handlerLogin(user, isRegist, ctx)
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
		hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
		glog.Infof("a.hallPid: %s", a.hallPid.String())
		if err != nil {
			glog.Fatalf("remote hall err %v", err)
		}
		a.hallPid = hallPid.Pid
		connect := &pb.HallConnect{
			Sender: ctx.Self(),
			Name:   a.Name,
		}
		a.hallPid.Tell(connect)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.SyncUser:
		arg := msg.(*pb.SyncUser)
		a.respUser(arg, ctx)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		a.respCurrency(arg, ctx)
	case *pb.CBuy:
		arg := msg.(*pb.CBuy)
		user := a.getUser(ctx)
		//优化
		rsp, diamond, coin := handler.Buy(arg, user)
		//同步兑换
		a.sendCurrency(user, diamond, coin, data.LogType18, ctx)
		//响应
		ctx.Respond(rsp)
	case *pb.CShop:
		arg := msg.(*pb.CShop)
		user := a.getUser(ctx)
		//响应
		rsp := handler.Shop(arg, user)
		ctx.Respond(rsp)
	case *pb.GetUserid:
		arg := msg.(*pb.GetUserid)
		//响应登录
		rsp := new(pb.GotUserid)
		rsp.Userid = a.router[arg.Sender.String()]
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
func (a *RoleActor) handlerLogin(user *data.User,
	isRegist bool, ctx actor.Context) {
	if user == nil {
		//登录失败
		return
	}
	a.logined(user)
	a.sendUser(user, ctx)
	glog.Debugf("user %#v", user)
	if isRegist {
		//注册奖励发放
		var diamond int32 = config.GetEnv(data.ENV1)
		var coin int32 = config.GetEnv(data.ENV2)
		a.sendCurrency(user, diamond, coin, data.LogType1, ctx)
	}
}

func (a *RoleActor) logined(user *data.User) {
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
	glog.Debugf("login userid: %s", user.Userid)
	glog.Debugf("roles len: %d", len(a.roles))
	glog.Debugf("offline len: %d", len(a.offline))
	//router 可直接在这里处理
}

func (a *RoleActor) sendUser(user *data.User, ctx actor.Context) {
	//同步数据,只有登录时才向节点同步数据
	//其它时候为节点向role同步,避免数据覆盖
	msg3 := new(pb.SyncUser)
	msg3.Userid = user.Userid
	result, err := json.Marshal(user)
	if err != nil {
		glog.Errorf("user Marshal err %v", err)
		return
	}
	msg3.Data = string(result)
	ctx.Sender().Tell(msg3)
}

func (a *RoleActor) respUser(arg *pb.SyncUser, ctx actor.Context) {
	glog.Debugf("SyncUser %#v", arg.Userid)
	user := a.roles[arg.Userid]
	if user == nil {
		glog.Errorf("respUser user empty %s", arg.Userid)
		return
	}
	err := json.Unmarshal([]byte(arg.Data), user)
	if err != nil {
		glog.Errorf("user Unmarshal err %v", err)
	}
	glog.Debugf("user %#v", user)
	//定时回存数据
	user.Save()
}

func (a *RoleActor) addCurrency(user *data.User,
	rtype, ltype int, amount int32, ctx actor.Context) {
	switch uint32(rtype) {
	case data.DIAMOND:
		a.sendCurrency(user, amount, 0, ltype, ctx)
	case data.COIN:
		a.sendCurrency(user, 0, amount, ltype, ctx)
	}
}

func (a *RoleActor) sendCurrency(user *data.User,
	diamond, coin int32, ltype int, ctx actor.Context) {
	if user == nil {
		glog.Errorf("sendCurrency user empty: %d", ltype)
	}
	msg3 := new(pb.ChangeCurrency)
	msg3.Userid = user.Userid
	msg3.Type = int32(ltype)
	msg3.Coin = coin
	msg3.Diamond = diamond
	ctx.Sender().Tell(msg3)
}

func (a *RoleActor) respCurrency(arg *pb.ChangeCurrency,
	ctx actor.Context) {
	userid := arg.Userid
	diamond := arg.Diamond
	coin := arg.Coin
	ltype := int(arg.Type)
	user := a.roles[userid]
	if user == nil {
		glog.Errorf("respCurrency user empty %s, type %d", userid, ltype)
		return
	}
	user.AddDiamond(diamond)
	user.AddCoin(coin)
}

func (a *RoleActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
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
