package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

var (
	mailPid *actor.PID
)

//邮件列表服务
type MailActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//邮件列表
	mails map[string]*data.Mail
	//唯一id生成
	uniqueid *data.IDGen
	//关闭通道
	stopCh chan struct{}
	//更新状态
	status bool
	//计时
	timer int
}

func (a *MailActor) Receive(ctx actor.Context) {
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

func newMailActor() actor.Actor {
	a := new(MailActor)
	a.Name = cfg.Section("mail").Name()
	a.mails = make(map[string]*data.Mail)
	//唯一id初始化
	a.uniqueid = data.InitIDGen(data.MAILID_KEY)
	a.stopCh = make(chan struct{})
	return a
}
