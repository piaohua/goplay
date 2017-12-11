package main

import (
	"encoding/json"

	"goplay/data"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *GateActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.GateConnected:
		//连接成功
		arg := msg.(*pb.GateConnected)
		glog.Infof("GateConnected %v", arg.Message)
	case *pb.GateDisconnected:
		//成功断开
		arg := msg.(*pb.GateDisconnected)
		glog.Infof("GateDisconnected %s", arg.Message)
	case *pb.LoginElse:
		//别处登录
		arg := msg.(*pb.LoginElse)
		userid := arg.GetUserid()
		if p, ok := a.roles[userid]; ok {
			p.Tell(arg)
		}
		glog.Debugf("LoginElse userid: %s", userid)
		//移除
		delete(a.roles, userid)
		//响应登录
		rsp := new(pb.LoginedElse)
		rsp.Userid = userid
		ctx.Respond(rsp)
	case *pb.SetLogin:
		arg := msg.(*pb.SetLogin)
		set := &pb.SetLogined{
			Message: ctx.Self().String(),
			DbmsPid: a.dbmsPid,
			RoomPid: a.roomPid,
			RolePid: a.rolePid,
			HallPid: a.hallPid,
			BetsPid: a.betsPid,
			MailPid: a.mailPid,
		}
		arg.Sender.Tell(set)
		glog.Infof("SetLogin %s", arg.Sender.String())
	case *pb.LoginGate:
		//登录成功
		arg := msg.(*pb.LoginGate)
		userid := arg.GetUserid()
		//添加
		a.roles[userid] = arg.Sender
		glog.Debugf("LoginGate userid: %s", userid)
		//响应登录
		rsp := new(pb.LoginedGate)
		rsp.Message = ctx.Self().String()
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		userid := arg.GetUserid()
		glog.Debugf("Logout userid: %s", userid)
		//移除
		delete(a.roles, userid)
		a.hallPid.Tell(arg)
		a.rolePid.Tell(arg)
		a.roomPid.Tell(arg)
	case *pb.HallConnect:
		//初始化建立连接
		a.init(ctx)
	case *pb.SyncConfig:
		//同步配置
		arg := msg.(*pb.SyncConfig)
		a.syncConfig(arg)
	case *pb.SPushNewBetting,
		*pb.SPushJackpot:
		for _, v := range a.roles {
			v.Tell(msg)
		}
	case *pb.BetsResult:
		arg := msg.(*pb.BetsResult)
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			msg1 := new(pb.SPushBetting)
			err1 := msg1.Unmarshal([]byte(arg.Message))
			if err1 != nil {
				glog.Errorf("BetsResult k %s err1 %v", userid, err1)
			}
			v.Tell(msg1)
		}
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			v.Tell(msg)
		} else {
			//离线
			a.rolePid.Tell(msg)
		}
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		if !handler.WxpayVerify(arg) {
			return
		}
		a.rolePid.Tell(arg)
	case *pb.WxpayGoods:
		arg := msg.(*pb.WxpayGoods)
		glog.Debugf("WxpayGoods: %v", arg)
		userid := arg.Userid
		if v, ok := a.roles[userid]; ok {
			v.Tell(arg)
		} else {
			glog.Errorf("WxpayGoods: %v", arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *GateActor) init(ctx actor.Context) {
	glog.Infof("gate init: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	role := cfg.Section("cookie").Key("role").Value()
	bets := cfg.Section("cookie").Key("bets").Value()
	mail := cfg.Section("cookie").Key("mail").Value()
	//timeout := 3 * time.Second
	//a.dbmsPid, err = remote.SpawnNamed(bind, a.Name, name, timeout)
	//if err != nil {
	//	glog.Fatalf("remote dbms err %v", err)
	//}
	//a.roomPid, err = remote.SpawnNamed(bind, a.Name, room, timeout)
	//if err != nil {
	//	glog.Fatalf("remote room err %v", err)
	//}
	//a.rolePid, err = remote.SpawnNamed(bind, a.Name, role, timeout)
	//if err != nil {
	//	glog.Fatalf("remote role err %v", err)
	//}
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	a.rolePid = actor.NewPID(bind, role)
	a.betsPid = actor.NewPID(bind, bets)
	a.mailPid = actor.NewPID(bind, mail)
	//a.dbmsPid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//a.roomPid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//a.rolePid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	////name
	//hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
	//if err != nil {
	//	glog.Fatalf("remote hall err %v", err)
	//}
	//a.hallPid = hallPid.Pid
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.GateConnect{
		Sender: ctx.Self(),
	}
	a.dbmsPid.Tell(connect)
	a.hallPid.Tell(connect)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	glog.Infof("a.rolePid: %s", a.rolePid.String())
}

func (a *GateActor) disc(ctx actor.Context) {
	//TODO 关闭消息
	//for k, v := range a.roles {
	//}
	//断开处理
	msg := &pb.GateDisconnect{
		Sender: ctx.Self(),
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Tell(msg)
	}
	if a.roomPid != nil {
		a.roomPid.Tell(msg)
	}
	if a.rolePid != nil {
		a.rolePid.Tell(msg)
	}
	if a.hallPid != nil {
		a.hallPid.Tell(msg)
	}
}

//同步配置
func (a *GateActor) syncConfig(arg *pb.SyncConfig) {
	switch arg.Type {
	case pb.CONFIG_BOX: //宝箱
		b := make([]data.Box, 0)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddBox(v)
		}
	case pb.CONFIG_ENV: //变量
		b := make(map[string]int32)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for k, v := range b {
			config.SetEnv2(k, v)
		}
	case pb.CONFIG_LOTTERY: //全民刮奖
		b := make(map[uint32]uint32)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for k, v := range b {
			config.SetLottery(k, v)
		}
	case pb.CONFIG_NOTICE: //公告
		b := make([]data.Notice, 0)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddNotice(v)
		}
	case pb.CONFIG_PRIZE: //抽奖
		b := make([]data.Prize, 0)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddPrize(v)
		}
	case pb.CONFIG_SHOP: //商城
		b := make(map[string]data.Shop)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddShop(v)
		}
	case pb.CONFIG_VIP: //VIP
		b := make(map[int]data.Vip)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddVip(v)
		}
	case pb.CONFIG_CLASSIC: //经典
		b := make(map[string]data.Classic)
		err = json.Unmarshal([]byte(arg.Data), &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddClassic(v)
		}
	default:
		glog.Errorf("syncConfig unknown type %d", arg.Type)
	}
}
