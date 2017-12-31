package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	nodePid *actor.PID
)

//大厅服务
type HallActor struct {
	Name string
	//服务注册
	serve map[string]*actor.PID
	//所有网关节点
	//gates map[string]*actor.PID
	//玩家所在节点userid-gate
	roles map[string]string
	//节点人数gate-numbers
	count map[string]uint32
	//所有桌子roomid-deskPid
	desks map[string]*actor.PID
	//桌子人数roomid-numbers
	rnums map[string]uint32
	//桌子类型roomid-type
	rtype map[string]uint32
	//玩家所在桌子userid-roomid
	router map[string]string
	//关闭通道
	stopCh chan struct{}
	//更新状态
	status bool
	//计时
	timer int
}

func (a *HallActor) Receive(ctx actor.Context) {
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

func newHallActor() actor.Actor {
	a := new(HallActor)
	//name
	a.Name = cfg.Section("hall").Name()
	a.serve = make(map[string]*actor.PID)
	//a.gates = make(map[string]*actor.PID)
	a.roles = make(map[string]string)
	a.count = make(map[string]uint32)
	a.stopCh = make(chan struct{})
	return a
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
	nodePid.Tell(new(pb.ServeStart))
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
