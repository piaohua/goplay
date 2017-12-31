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
	//角色userid-roomid
	router map[string]string
	//邀请码code-roomid
	codes map[string]string
	//房间人数roomid-numbers
	count map[string]uint32
	//TODO 匹配规则 rules
	//rules map[string]string
	//唯一id生成
	uniqueid *data.IDGen
	//关闭通道
	stopCh chan struct{}
	//更新状态
	status bool
	//计时
	timer int
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
	a.router = make(map[string]string)
	a.codes = make(map[string]string)
	a.count = make(map[string]uint32)
	a.stopCh = make(chan struct{})
	return a
}
