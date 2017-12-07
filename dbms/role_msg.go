package main

import (
	"encoding/json"

	"goplay/data"
	"goplay/game/handler"
	"goplay/game/login"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *RoleActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HallConnect:
		//初始化建立连接
		glog.Infof("role init: %v", ctx.Self().String())
		//连接
		bind := cfg.Section("hall").Key("bind").Value()
		name := cfg.Section("cookie").Key("name").Value()
		//timeout := 3 * time.Second
		//hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
		//if err != nil {
		//	glog.Fatalf("remote hall err %v", err)
		//}
		//a.hallPid = hallPid.Pid
		a.hallPid = actor.NewPID(bind, name)
		glog.Infof("a.hallPid: %s", a.hallPid.String())
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
		a.syncUser(arg, ctx)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		a.changeCurrency(arg, ctx)
	case *pb.Login:
		//登录成功
		arg := msg.(*pb.Login)
		a.logined(arg, ctx)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		a.logouted(arg, ctx)
	case *pb.GetUserid:
		arg := msg.(*pb.GetUserid)
		//响应登录
		rsp := new(pb.GotUserid)
		rsp.Userid = a.router[arg.Sender.String()]
		ctx.Respond(rsp)
	case *pb.RoleRegist:
		arg := msg.(*pb.RoleRegist)
		glog.Debugf("RoleRegist %#v", arg)
		var phone string = arg.GetPhone()
		//在线表中查找
		if _, ok := a.players[phone]; ok {
			rsp := new(pb.RoleRegisted)
			rsp.Error = pb.PhoneRegisted
			ctx.Respond(rsp)
			return
		}
		//数据库中查找
		rsp := login.Regist(arg, a.uniqueid)
		ctx.Respond(rsp)
	case *pb.RoleLogin:
		arg := msg.(*pb.RoleLogin)
		glog.Debugf("RoleLogin %#v", arg)
		var phone string = arg.GetPhone()
		//在线表中查找
		user := a.getUser(phone)
		//数据库中查找
		rsp := login.Login(arg, user)
		ctx.Respond(rsp)
	case *pb.WxLogin:
		arg := msg.(*pb.WxLogin)
		glog.Debugf("WxLogin %#v", arg)
		var wxuid string = arg.GetWxuid()
		//在线表中查找
		user := a.getUser(wxuid)
		//数据库中查找
		rsp := login.WxLogin(arg, user, a.uniqueid)
		ctx.Respond(rsp)
	case *pb.BuildAgent:
		arg := msg.(*pb.BuildAgent)
		rsp := handler.BuildAgent(arg)
		ctx.Respond(rsp)
	case *pb.BankGive:
		arg := msg.(*pb.BankGive)
		user := a.getUserById(arg.Userid)
		rsp := new(pb.BankGave)
		if user == nil || user.Userid == "" {
			rsp.Error = pb.GiveUseridError
		} else if user != nil {
			rsp.Userid = user.GetUserid()
			rsp.Coin = user.GetCoin()
		}
		ctx.Respond(rsp)
	case *pb.GetUserData:
		arg := msg.(*pb.GetUserData)
		user := a.getUserById(arg.Userid)
		rsp := handler.GetUserData1(user)
		ctx.Respond(rsp)
	case *pb.ApplePay:
		arg := msg.(*pb.ApplePay)
		rsp := handler.AppleVerify(arg)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
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

//在线表中查找,不存在时离线表中获取
func (a *RoleActor) getUserById(userid string) *data.User {
	if user, ok := a.roles[userid]; ok {
		return user
	}
	if user, ok := a.offline[userid]; ok {
		return user
	}
	user := new(data.User)
	user.GetById(userid) //数据库中取
	return user
}

//在线表中查找
func (a *RoleActor) getUser(account string) *data.User {
	user := new(data.User)
	if v, ok := a.players[account]; ok {
		if v2, ok := a.roles[v]; ok {
			user = v2
		} else if v2, ok := a.offline[v]; ok {
			user = v2
		}
	}
	return user
}

//登录处理
func (a *RoleActor) logined(arg *pb.Login, ctx actor.Context) {
	user := new(data.User)
	//进程id映射
	a.router[ctx.Sender().String()] = arg.Userid
	//已经在线,用在线数据
	if v, ok := a.roles[arg.Userid]; ok {
		user = v
		return
	}
	//已经离线,用离线数据
	if v, ok := a.offline[arg.Userid]; ok {
		user = v
		delete(a.offline, arg.Userid)
	}
	//解析
	err := json.Unmarshal([]byte(arg.Data), user)
	if err != nil {
		glog.Errorf("user %s Unmarshal err %v", arg.Userid, err)
		return
	}
	//映射
	if user.Wxuid != "" {
		a.players[user.Wxuid] = arg.Userid
	} else if user.Phone != "" {
		a.players[user.Phone] = arg.Userid
	}
	//登录成功
	a.roles[arg.Userid] = user
	glog.Debugf("login userid: %s", arg.Userid)
	glog.Debugf("roles len: %d", len(a.roles))
	glog.Debugf("offline len: %d", len(a.offline))
	//响应登录
	rsp := new(pb.Logined)
	ctx.Respond(rsp)
}

//登出处理
func (a *RoleActor) logouted(arg *pb.Logout, ctx actor.Context) {
	glog.Debugf("Logout userid: %s", arg.Userid)
	if v, ok := a.roles[arg.Userid]; ok {
		//离线
		a.offline[arg.Userid] = v
		//移除
		delete(a.roles, arg.Userid)
		//TODO 定期离线数据清理,移除,存储
		//if v.Wxuid != "" {
		//	delete(a.players, v.Wxuid)
		//} else if user.Phone != "" {
		//	delete(a.players, v.Phone)
		//}
	}
	delete(a.router, arg.Sender.String())
}

//在线同步数据
func (a *RoleActor) syncUser(arg *pb.SyncUser, ctx actor.Context) {
	glog.Debugf("SyncUser %#v", arg.Userid)
	user := a.roles[arg.Userid]
	if user == nil {
		glog.Errorf("syncUser user empty %s", arg.Userid)
		return
	}
	err := json.Unmarshal([]byte(arg.Data), user)
	if err != nil {
		glog.Errorf("user %s Unmarshal err %v", arg.Userid, err)
		return
	}
	glog.Debugf("user %#v", user)
	//TODO 定时回存数据
	user.Save()
}

//更新货币
func (a *RoleActor) changeCurrency(arg *pb.ChangeCurrency,
	ctx actor.Context) {
	userid := arg.Userid
	diamond := arg.Diamond
	coin := arg.Coin
	bank := arg.Bank
	upsert := arg.Upsert
	ltype := int(arg.Type)
	user := a.getUser(userid)
	if user != nil {
		user.AddDiamond(diamond)
		user.AddCoin(coin)
		user.AddBank(bank)
		return
	}
	if !upsert {
		glog.Errorf("changeCurrency user empty %s, type %d", userid, ltype)
		return
	}
	//离线更新
	user = a.getUserById(userid)
	if user == nil || user.Userid == "" {
		glog.Errorf("changeCurrency user empty %s, type %d", userid, ltype)
		return
	}
	user.UpdateDiamond(diamond)
	user.UpdateCoin(coin)
	user.UpdateBank(bank)
}
