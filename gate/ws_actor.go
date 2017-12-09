package main

import (
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

//test
type testWs struct{ Who string }

func (ws *WSConn) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		glog.Notice("Starting, initialize actor here")
		ws.init(ctx)
	case *actor.Stopping:
		glog.Notice("Stopping, actor is about to shut down")
		ws.disc(ctx)
	case *actor.Stopped:
		glog.Notice("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		glog.Notice("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		glog.Infof("ReceiveTimeout: %v", ctx.Self().String())
		ws.timeout(ctx)
	case *testWs:
		glog.Infof("self %s\n", ctx.Self().String())
		glog.Infof("msg.Who %v\n", msg.Who)
	case proto.Message:
		ws.Handler(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (ws *WSConn) init(ctx actor.Context) {
	glog.Infof("ws init: %v", ctx.Self().String())
	ctx.SetReceiveTimeout(waitForLogin) //login timeout set
	set := &pb.SetLogin{
		Sender: ctx.Self(),
	}
	nodePid.Tell(set)
}

func (ws *WSConn) disc(ctx actor.Context) {
	glog.Infof("ws disc: %v", ctx.Self().String())
	//已经断开,在别处登录
	if ws.hallPid == nil {
		return
	}
	if ws.User == nil {
		return
	}
	//回存数据
	ws.syncUser()
	//登出日志
	if ws.dbmsPid != nil {
		msg2 := &pb.LogLogout{
			Userid: ws.User.Userid,
			Event:  1,
		}
		ws.dbmsPid.Tell(msg2)
	}
	//断开处理
	msg := &pb.Logout{
		Sender: ctx.Self(),
		Userid: ws.User.Userid,
	}
	nodePid.Tell(msg)
}

func (ws *WSConn) timeout(ctx actor.Context) {
	if !ws.online {
		//断开连接
		ws.Close()
		return
	}
	if ws.status {
		//同步数据
		ws.syncUser()
		ws.status = false
	}
}

//初始化
func (ws *WSConn) initWs() *actor.PID {
	props := actor.FromInstance(ws) //实例
	return actor.Spawn(props)       //启动一个进程
}
