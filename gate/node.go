package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	nodePid *actor.PID
)

//网关服务
type GateActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//数据中心服务
	dbmsPid *actor.PID
	//房间服务
	roomPid *actor.PID
	//角色服务
	rolePid *actor.PID
	//投注服务
	betsPid *actor.PID
	//邮箱服务
	mailPid *actor.PID
	//节点角色进程
	roles map[string]*actor.PID
}

func (a *GateActor) Receive(ctx actor.Context) {
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

func newGateActor() actor.Producer {
	return func() actor.Actor {
		a := new(GateActor)
		a.Name = cfg.Section("gate.node1").Name()
		//roles key=userid
		a.roles = make(map[string]*actor.PID)
		return a
	}
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	props := actor.
		FromProducer(newGateActor()).
		WithMailbox(mailbox.Bounded(20000))
	remote.Register(name, props)
	nodePid, err = actor.SpawnNamed(props, name)
	if err != nil {
		glog.Fatalf("nodePid err %v", err)
	}
	glog.Infof("nodePid %s", nodePid.String())
	nodePid.Tell(new(pb.HallConnect))
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
