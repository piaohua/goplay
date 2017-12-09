package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

var (
	betsPid *actor.PID
)

//疯狂投注列表服务
type BetsActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//唯一id生成
	uniqueid *data.IDGen
}

func (a *BetsActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pb.Request:
		ctx.Respond(&pb.Response{})
	case *actor.Started:
		glog.Notice("Starting, initialize actor here")
	case *actor.Stopping:
		glog.Notice("Stopping, actor is about to shut down")
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

func (a *BetsActor) init(ctx actor.Context) {
	glog.Infof("ws init: %v", ctx.Self().String())
	ctx.SetReceiveTimeout(loop) //timeout set
}

func (a *BetsActor) timeout(ctx actor.Context) {
	glog.Debugf("timeout: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(0) //timeout off
	//TODO
}

func newBetsActor() actor.Actor {
	a := new(BetsActor)
	a.Name = cfg.Section("bets").Name()
	//唯一id初始化
	a.uniqueid = data.InitIDGen(data.BETTING_KEY)
	return a
}
