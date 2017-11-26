package main

import (
	"goplay/glog"
	"goplay/pb"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (a *RoomActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HallConnect:
		//初始化建立连接
		glog.Infof("room init: %v", ctx.Self().String())
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
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *RoomActor) HandlerStop(ctx actor.Context) {
	glog.Debugf("HandlerStop: %s", a.Name)
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
