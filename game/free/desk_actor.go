package main

import (
	"goplay/glog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

func (a *DeskFree) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		glog.Notice("Starting, initialize actor here")
		//ctx.Self().Tell(new(pb.ServeStart))
	case *actor.Stopping:
		glog.Notice("Stopping, actor is about to shut down")
		//ctx.Self().Tell(new(pb.ServeStop))
	case *actor.Stopped:
		glog.Notice("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		glog.Notice("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
	case proto.Message:
		a.Handler(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//初始化
func (a *DeskFree) newDesk() *actor.PID {
	props := actor.FromInstance(a) //实例
	return actor.Spawn(props)      //启动一个进程
}
