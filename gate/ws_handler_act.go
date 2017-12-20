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
			ws.Send(rsp)
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
		ws.rolePid.Request(arg, ctx.Self())
	case *pb.CBetting:
		arg := msg.(*pb.CBetting)
		number := arg.GetNumber()
		glog.Debugf("CBetting %#v", arg)
		if ws.User.GetDiamond() < number {
			msg1 := new(pb.SBetting)
			msg1.Error = pb.NotEnoughDiamond
			ws.Send(msg1)
			return
		}
		msg2 := new(pb.BetsOn)
		msg2.Userid = ws.User.GetUserid()
		msg2.Seat = arg.GetSeat()
		msg2.Number = arg.GetNumber()
		ws.rolePid.Request(msg2, ctx.Self())
	case *pb.SBetting:
		arg := msg.(*pb.SBetting)
		glog.Debugf("SBetting %#v", arg)
		if arg.Error == pb.OK {
			ws.expend(arg.Number, data.LogType36)
		}
		ws.Send(arg)
	case *pb.CBettingRecord:
		arg := msg.(*pb.CBettingRecord)
		glog.Debugf("CBettingRecord %#v", arg)
		msg2 := new(pb.BetsRecord)
		msg2.Userid = ws.User.GetUserid()
		msg2.Page = arg.GetPage()
		ws.rolePid.Request(msg2, ctx.Self())
	default:
		//glog.Errorf("unknown message %v", msg)
		ws.HandlerDesk(msg, ctx)
	}
}
