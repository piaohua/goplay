package main

import (
	"fmt"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

//中心服务
type CoreActor struct {
	Name string
}

func (a *CoreActor) Receive(ctx actor.Context) {
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
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func newCoreActor() actor.Actor {
	return new(CoreActor)
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	remote.Register(name, actor.FromProducer(newCoreActor))
}
