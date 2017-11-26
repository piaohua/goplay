package main

import (
	"fmt"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

//大厅服务
type HallActor struct {
	Name string
	//登录服务
	loginPid *actor.PID
	//数据中心服务
	dbmsPid *actor.PID
	//房间服务
	roomPid *actor.PID
	//角色服务
	rolePid *actor.PID
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

func newHallActor() actor.Actor {
	a := new(HallActor)
	//name
	a.Name = cfg.Section("hall").Name()
	a.gates = make(map[string]*actor.PID)
	a.roles = make(map[string]string)
	a.count = make(map[string]uint32)
	return a
}

func NewRemote(bind, name string) {
	remote.Start(bind)
	remote.Register(name, actor.FromProducer(newHallActor))
}
