package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskFree) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	default:
		//glog.Errorf("unknown message %v", msg)
		a.HandlerLogic(msg, ctx)
	}
}

//启动服务
func (a *DeskFree) start(ctx actor.Context) {
	glog.Infof("desk free start: %v", ctx.Self().String())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *DeskFree) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("desk ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *DeskFree) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//逻辑处理
	a.tickerHandler()
}

//关闭时钟
func (a *DeskFree) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
	//离开消息
	for k, p := range a.pids {
		msg2 := &pb.LeaveDesk{
			Roomid: a.id,
			Userid: k,
		}
		nodePid.Tell(msg2)
		p.Tell(msg2)
	}
	//逻辑处理
	a.close()
}

func (a *DeskFree) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//断开处理
}
