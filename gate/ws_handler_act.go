package main

import (
	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家活动请求处理
func (ws *WSConn) HandlerAct(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CLotteryInfo:
		arg := msg.(*pb.CLotteryInfo)
		glog.Debugf("CLotteryInfo %#v", arg)
		rsp := handler.LotteryInfo(arg)
		ws.Send(rsp)
	case *pb.CLottery:
		arg := msg.(*pb.CLottery)
		glog.Debugf("CLottery %#v", arg)
		rsp, number, prize, ok := handler.Lottery(arg, ws.User)
		if rsp.Error != pb.OK {
			ws.Send(res)
			return
		}
		ws.expend(number, data.LogType37)
		if ok {
			ws.addCurrency(int32(prize), 0, data.LogType37)
		}
		ws.Send(rsp)
	case *pb.CBettingInfo:
		arg := msg.(*pb.CBettingInfo)
		glog.Debugf("CBettingInfo %#v", arg)
		//rsp := betting.GetInfo()
		//ws.Send(rsp)
	case *pb.CBetting:
		arg := msg.(*pb.CBetting)
		glog.Debugf("CBetting %#v", arg)
	case *pb.CBettingRecord:
		arg := msg.(*pb.CBettingRecord)
		glog.Debugf("CBettingRecord %#v", arg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
