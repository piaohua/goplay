package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

var (
	betsPid *actor.PID
)

//疯狂投注列表服务
type BetsActor struct {
	Name string
	//大厅服务
	hallPid *actor.PID
	//唯一id生成
	uniqueid *data.IDGen //期号,凌晨重置
	//
	closeCh chan bool        //关闭通道
	stopCh  chan struct{}    //关闭通道
	msgCh   chan interface{} //消息通道
	//
	betTime uint32   //时间
	betRest uint32   //时间
	timer   uint32   //时间
	state   uint32   //状态0投注,1等待
	today   string   //期号
	cards   []uint32 //牌
	niu     uint32   //牛
	//
	odds    map[uint32]float32 //赔率
	jackpot map[uint32]uint32  //奖池
	//
	ante   map[string]map[uint32]uint32 //玩家下注额
	result map[string]map[uint32]int32  //玩家输赢额
	//
	prize  map[string]int32 //个人中奖金额
	winner map[uint32]bool  //中奖位置
	lose   map[string]int32 //个人输赢总金额
}

func (a *BetsActor) Receive(ctx actor.Context) {
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

func (a *BetsActor) init(ctx actor.Context) {
	glog.Infof("ws init: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(loop) //timeout set
	a.initBetting(ctx)
}

func (a *BetsActor) timeout(ctx actor.Context) {
	glog.Debugf("timeout: %v", ctx.Self().String())
	//ctx.SetReceiveTimeout(0) //timeout off
	//TODO
}

func newBetsActor() actor.Actor {
	a := new(BetsActor)
	a.Name = cfg.Section("bets").Name()
	//唯一id初始化
	a.uniqueid = data.InitIDGen(data.BETTING_KEY)
	return a
}
