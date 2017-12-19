package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *GateActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		//连接成功
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		//成功断开
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
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
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		//初始化建立连接
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
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

//启动服务
func (a *GateActor) start(ctx actor.Context) {
	glog.Infof("gate start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	role := cfg.Section("cookie").Key("role").Value()
	bets := cfg.Section("cookie").Key("bets").Value()
	mail := cfg.Section("cookie").Key("mail").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	a.rolePid = actor.NewPID(bind, role)
	a.betsPid = actor.NewPID(bind, bets)
	a.mailPid = actor.NewPID(bind, mail)
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.dbmsPid.Request(connect, ctx.Self())
	a.hallPid.Request(connect, ctx.Self())
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	glog.Infof("a.rolePid: %s", a.rolePid.String())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *GateActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("gate ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("gate ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *GateActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *GateActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *GateActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//msg := new(pb.ServeClose)
	//for k, v := range a.roles {
	//	glog.Debugf("Stop role: %s", k)
	//	v.Tell(msg)
	//}
	//TODO 关闭消息
	//for k, v := range a.roles {
	//}
	//断开处理
	msg := &pb.Disconnect{
		Name: a.Name,
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Request(msg, ctx.Self())
	}
	if a.roomPid != nil {
		a.roomPid.Request(msg, ctx.Self())
	}
	if a.rolePid != nil {
		a.rolePid.Request(msg, ctx.Self())
	}
	if a.hallPid != nil {
		a.hallPid.Request(msg, ctx.Self())
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
