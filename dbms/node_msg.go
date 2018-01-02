package main

import (
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DBMSActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.Connect:
		arg := msg.(*pb.Connect)
		//网关注册
		a.gates[arg.Name] = ctx.Sender()
		//响应
		connected := &pb.Connected{
			Name: a.Name,
		}
		ctx.Respond(connected)
		glog.Infof("Connect %s", arg.Name)
		//同步配置
		a.syncConfig(arg.Name)
	case *pb.Disconnect:
		arg := msg.(*pb.Disconnect)
		//网关注销
		delete(a.gates, arg.Name)
		//响应
		disconnected := &pb.Disconnected{
			Name: a.Name,
		}
		ctx.Respond(disconnected)
		glog.Infof("Disconnect %s", arg.Name)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.SPushNewBetting,
		*pb.SPushJackpot,
		*pb.BetsResult:
		for _, v := range a.gates {
			v.Tell(msg)
		}
	case *pb.GetConfig:
		arg := msg.(*pb.GetConfig)
		glog.Infof("GetConfig %#v", arg)
	default:
		if a.logger == nil {
			glog.Errorf("unknown message %v", msg)
		} else {
			a.logger.Tell(msg)
		}
	}
}

//启动服务
func (a *DBMSActor) start(ctx actor.Context) {
	glog.Infof("dbms start: %v", ctx.Self().String())
	//初始化建立连接
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//timeout := 3 * time.Second
	//msg1 := new(pb.Connect)
	//res1, err1 := nodePid.RequestFuture(msg1, timeout).Result()
	//if err1 != nil || res1 == nil {
	//	glog.Fatalf("Hall Connect err: %v", err1)
	//	return
	//}
	//response1 := res1.(*pb.Connected)
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *DBMSActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("dbms ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("dbms ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *DBMSActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *DBMSActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *DBMSActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	for k, _ := range a.gates {
		glog.Debugf("Stop gate: %s", k)
	}
	if a.logger != nil {
		a.logger.Stop()
	}
}

//打包配置
func (a *DBMSActor) syncConfigMsg(ctype pb.ConfigType,
	d interface{}) interface{} {
	msg := new(pb.SyncConfig)
	msg.Type = ctype
	result, err := json.Marshal(d)
	if err != nil {
		glog.Errorf("syncConfig Marshal err %v, data %#v", err, d)
	}
	msg.Data = string(result)
	return msg
}

//同步配置
func (a *DBMSActor) syncConfig(key string) {
	if _, ok := a.gates[key]; !ok {
		glog.Errorf("gate not exists: %s", key)
		return
	}
	pid := a.gates[key]
	msg1 := a.syncConfigMsg(pb.CONFIG_BOX, config.GetBoxs())
	pid.Tell(msg1)
	msg2 := a.syncConfigMsg(pb.CONFIG_ENV, config.GetEnvs())
	pid.Tell(msg2)
	msg3 := a.syncConfigMsg(pb.CONFIG_LOTTERY, config.GetLotterys())
	pid.Tell(msg3)
	msg4 := a.syncConfigMsg(pb.CONFIG_NOTICE, config.GetNotices(data.NOTICE_TYPE1))
	pid.Tell(msg4)
	msg5 := a.syncConfigMsg(pb.CONFIG_PRIZE, config.GetPrizes())
	pid.Tell(msg5)
	msg6 := a.syncConfigMsg(pb.CONFIG_SHOP, config.GetShops2())
	pid.Tell(msg6)
	msg7 := a.syncConfigMsg(pb.CONFIG_VIP, config.GetVips())
	pid.Tell(msg7)
	msg8 := a.syncConfigMsg(pb.CONFIG_CLASSIC, config.GetClassics())
	pid.Tell(msg8)
}

//同步配置
func (a *DBMSActor) getConfig(arg *pb.GetConfig, ctx actor.Context) {
	switch arg.Type {
	case pb.CONFIG_ENV:
		msg2 := a.syncConfigMsg(pb.CONFIG_ENV, config.GetEnvs())
		ctx.Respond(msg2)
	}
}
