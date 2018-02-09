package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
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
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerMsg(msg, ctx)
	}
}

//启动服务
func (a *DeskActor) start(ctx actor.Context) {
	glog.Infof("desk start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//主动同步配置
	msg2 := &pb.GetConfig{
		Type: pb.CONFIG_ENV,
	}
	a.dbmsPid.Request(msg2, ctx.Self())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *DeskActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *DeskActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *DeskActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *DeskActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//关闭消息
	msg1 := new(pb.ServeStop)
	for k, v := range a.desks {
		//关闭房间消息
		msg2 := new(pb.CloseDesk)
		msg2.Roomid = k
		a.roomPid.Request(msg2, ctx.Self())
		a.hallPid.Request(msg2, ctx.Self())
		//关闭房间服务
		glog.Debugf("Stop desk: %s", k)
		v.Request(msg1, ctx.Self())
		//停掉服务
		v.Stop()
	}
	//断开处理
	msg := &pb.Disconnect{
		Name: a.Name,
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Tell(msg)
	}
	if a.hallPid != nil {
		a.hallPid.Tell(msg)
	}
}
