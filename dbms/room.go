package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

var (
	roomPid *actor.PID
)

//房间列表服务
type RoomActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//房间列表
	rooms map[string]*data.DeskData
	//唯一id生成
	uniqueid *data.IDGen
}

func (a *RoomActor) Receive(ctx actor.Context) {
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

func newRoomActor() actor.Actor {
	a := new(RoomActor)
	a.Name = cfg.Section("room").Name()
	a.rooms = make(map[string]*data.DeskData)
	//唯一id初始化
	a.uniqueid = data.InitIDGen(data.ROOMID_KEY)
	return a
}
