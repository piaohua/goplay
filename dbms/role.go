package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

var (
	rolePid *actor.PID
)

//角色列表服务
type RoleActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//角色数据
	roles map[string]*data.User
	//离线数据
	offline map[string]*data.User
	//角色PID-userid
	router map[string]string
	//唯一id生成
	uniqueid *data.IDGen
}

func (a *RoleActor) Receive(ctx actor.Context) {
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

func newRoleActor() actor.Actor {
	a := new(RoleActor)
	a.Name = cfg.Section("role").Name()
	a.roles = make(map[string]*data.User)
	a.offline = make(map[string]*data.User)
	a.router = make(map[string]string)
	//唯一id初始化
	a.uniqueid = data.InitIDGen(data.USERID_KEY)
	return a
}