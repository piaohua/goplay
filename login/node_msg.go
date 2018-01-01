package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *LoginActor) Handler(msg interface{}, ctx actor.Context) {
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
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		glog.Debugf("WxpayCallback: %v", arg)
		a.hallPid.Tell(arg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *LoginActor) start(ctx actor.Context) {
	glog.Infof("desk start: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	a.dbmsPid = actor.NewPID(bind, name)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *LoginActor) ticker(ctx actor.Context) {
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
func (a *LoginActor) ding(ctx actor.Context) {
	//glog.Debugf("ding: %v", ctx.Self().String())
	//TODO 同步节点
	//a.hallPid.Request(msg, ctx.Self())
	//TODO 路由方式
	//err := cfg.Section("gate.node"+version).NewKey("host", "x:x")
}

//关闭时钟
func (a *LoginActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *LoginActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//TODO 关闭消息
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
