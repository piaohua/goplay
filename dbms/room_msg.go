package main

import (
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *RoomActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HallConnect:
		//初始化建立连接
		glog.Infof("room init: %v", ctx.Self().String())
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
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		rsp := new(pb.ServeStarted)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *RoomActor) start(ctx actor.Context) {
	glog.Infof("room start: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(loop) //timeout set
}

func (a *RoomActor) timeout(ctx actor.Context) {
	glog.Debugf("timeout: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(0) //timeout off
	//TODO
}

func (a *RoomActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, _ := range a.rooms {
		glog.Debugf("Stop room: %s", k)
		//TODO
		//v.Save()
	}
}
