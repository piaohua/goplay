package main

import (
	"fmt"
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

//数据库操作服务
type DBMSActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//日志服务
	logger *actor.PID
	//网关节点
	gates map[string]*actor.PID
}

func (a *DBMSActor) Receive(ctx actor.Context) {
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

func newDBMSActor() actor.Actor {
	a := new(DBMSActor)
	a.Name = cfg.Section("dbms").Name()
	a.gates = make(map[string]*actor.PID)
	a.logger = NewLogger()
	return a
}

func NewRemote(bind, name, room, role string) {
	remote.Start(bind)
	//
	//remote.Register(name, actor.FromProducer(newDBMSActor))
	dbmsProps := actor.
		FromInstance(newDBMSActor())
	remote.Register(name, dbmsProps)
	nodePid, err = actor.SpawnNamed(dbmsProps, name)
	if err != nil {
		glog.Fatalf("nodePid err %v", err)
	}
	glog.Infof("nodePid %s", nodePid.String())
	nodePid.Tell(new(pb.HallConnect))
	//
	//remote.Register(room, actor.FromProducer(newRoomActor))
	roomProps := actor.
		FromInstance(newRoomActor())
	remote.Register(room, roomProps)
	roomPid, err = actor.SpawnNamed(roomProps, room)
	if err != nil {
		glog.Fatalf("roomPid err %v", err)
	}
	glog.Infof("roomPid %s", roomPid.String())
	roomPid.Tell(new(pb.HallConnect))
	//
	//remote.Register(role, actor.FromProducer(newRoleActor))
	roleProps := actor.
		FromInstance(newRoleActor())
	remote.Register(role, roleProps)
	rolePid, err = actor.SpawnNamed(roleProps, role)
	if err != nil {
		glog.Fatalf("rolePid err %v", err)
	}
	glog.Infof("rolePid %s", rolePid.String())
	rolePid.Tell(new(pb.HallConnect))
}

//关闭
func Stop() {
	timeout := 3 * time.Second
	msg := new(pb.ServeStop)
	if rolePid != nil {
		res1, err1 := rolePid.RequestFuture(msg, timeout).Result()
		if err1 != nil {
			glog.Errorf("rolePid Stop err: %v", err1)
		}
		response1 := res1.(*pb.ServeStoped)
		glog.Debugf("response1: %#v", response1)
		rolePid.Stop()
	}
	if roomPid != nil {
		res1, err1 := roomPid.RequestFuture(msg, timeout).Result()
		if err1 != nil {
			glog.Errorf("roomPid Stop err: %v", err1)
		}
		response1 := res1.(*pb.ServeStoped)
		glog.Debugf("response1: %#v", response1)
		roomPid.Stop()
	}
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
