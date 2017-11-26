package main

import (
	"fmt"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

//日志记录服务
type LoggerActor struct {
	Name string
}

func (a *LoggerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pb.Request:
		ctx.Respond(&pb.Response{})
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
	case proto.Message:
		a.Handler(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *LoggerActor) init(ctx actor.Context) {
	glog.Infof("logger init: %v", ctx.Self().String())
	//name
	a.Name = cfg.Section("logger").Name()
	glog.Infof("a.logger: %s", ctx.Self().String())
}

func NewLogger() *actor.PID {
	props := actor.FromInstance(new(LoggerActor))
	pid := actor.Spawn(props)
	return pid
}
