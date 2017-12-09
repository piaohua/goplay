package main

import (
	"goplay/glog"
	"goplay/pb"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	nodePid *actor.PID
	loop    = 30 * time.Second
)

//大厅服务
type HallActor struct {
	Name string
	//服务注册
	serve map[string]*actor.PID
	//所有网关节点
	gates map[string]*actor.PID
	//玩家所在节点
	roles map[string]string
	//节点人数
	count map[string]uint32
}

func (a *HallActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *pb.Request:
		ctx.Respond(&pb.Response{})
	case *actor.Started:
		glog.Notice("Starting, initialize actor here")
		a.init(ctx)
	case *actor.Stopping:
		glog.Notice("Stopping, actor is about to shut down")
	case *actor.Stopped:
		glog.Notice("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		glog.Notice("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
		a.timeout(ctx)
	case proto.Message:
		a.Handler(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func newHallActor() actor.Actor {
	a := new(HallActor)
	//name
	a.Name = cfg.Section("hall").Name()
	a.serve = make(map[string]*actor.PID)
	a.gates = make(map[string]*actor.PID)
	a.roles = make(map[string]string)
	a.count = make(map[string]uint32)
	return a
}

func (a *HallActor) init(ctx actor.Context) {
	glog.Infof("ws init: %v", ctx.Self().String())
	ctx.SetReceiveTimeout(loop) //timeout set
}

func (a *HallActor) timeout(ctx actor.Context) {
	glog.Debugf("timeout: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(0) //timeout off
	//TODO
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	remote.Register(name, actor.FromProducer(newHallActor))
	hallProps := actor.
		FromInstance(newHallActor())
	nodePid, err = actor.SpawnNamed(hallProps, name)
	if err != nil {
		glog.Fatalf("nodePid err %v", err)
	}
}

//关闭
func Stop() {
	timeout := 3 * time.Second
	msg := new(pb.ServeStop)
	if nodePid != nil {
		res1, err1 := nodePid.RequestFuture(msg, timeout).Result()
		if err1 != nil {
			glog.Errorf("nodePid Stop err: %v", err1)
		}
		response1 := res1.(*pb.ServeStoped)
		glog.Debugf("response1: %#v", response1)
		nodePid.Stop()
	}
}
