package main

import (
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家百人场请求处理
func (ws *WSConn) HandlerFree(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CEnterFreeRoom:
		arg := msg.(*pb.CEnterFreeRoom)
		glog.Debugf("CEnterFreeRoom %#v", arg)
		ws.freeEnter(arg, ctx)
	case *pb.CFreeDealer:
		arg := msg.(*pb.CFreeDealer)
		glog.Debugf("CFreeDealer %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SFreeDealer)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		var num uint32 = arg.GetCoin()
		if ws.User.GetCoin() < num {
			rsp := new(pb.SFreeDealer)
			rsp.Error = pb.NotEnoughCoin
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CDealerList:
		arg := msg.(*pb.CDealerList)
		glog.Debugf("CDealerList %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SDealerList)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CFreeSit:
		arg := msg.(*pb.CFreeSit)
		glog.Debugf("CFreeSit %#v", arg)
		ws.freeSit(arg, ctx)
	case *pb.CFreeBet:
		arg := msg.(*pb.CFreeBet)
		glog.Debugf("CFreeBet %#v", arg)
		ws.freeBet(arg, ctx)
	case *pb.CFreeTrend:
		arg := msg.(*pb.CFreeTrend)
		glog.Debugf("CFreeTrend %#v", arg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//坐下
func (ws *WSConn) freeSit(arg *pb.CFreeSit, ctx actor.Context) {
	if ws.gamePid == nil {
		rsp := new(pb.SFreeSit)
		rsp.Error = pb.NotInRoom
		ws.Send(rsp)
		return
	}
	var seat uint32 = arg.GetSeat()
	var state bool = arg.GetState()
	if !(seat >= 1 && seat <= 8) {
		rsp := new(pb.SFreeSit)
		rsp.Error = pb.OperateError
		ws.Send(rsp)
		return
	}
	//坐下限制
	limit_sit := config.GetEnv(data.ENV21)
	if state && int32(ws.User.GetCoin()) < limit_sit {
		rsp := new(pb.SFreeSit)
		rsp.Error = pb.NotEnoughCoin
		ws.Send(rsp)
		return
	}
	ws.gamePid.Request(arg, ctx.Self())
}

//下注
func (ws *WSConn) freeBet(arg *pb.CFreeBet, ctx actor.Context) {
	if ws.gamePid == nil {
		rsp := new(pb.SFreeBet)
		rsp.Error = pb.NotInRoom
		ws.Send(rsp)
		return
	}
	value := arg.GetValue()
	seatBet := arg.GetSeat()
	if !(seatBet >= 2 && seatBet <= 5) {
		rsp := new(pb.SFreeBet)
		rsp.Error = pb.OperateError
		ws.Send(rsp)
		return
	}
	if ws.User.GetCoin() < value {
		rsp := new(pb.SFreeBet)
		rsp.Error = pb.NotEnoughCoin
		ws.Send(rsp)
		return
	}
	ws.gamePid.Request(arg, ctx.Self())
}

//进入自由场
func (ws *WSConn) freeEnter(arg *pb.CEnterFreeRoom, ctx actor.Context) {
	if ws.gamePid != nil {
		//进入房间
		ws.entryRoom(ctx)
		return
	}
	stoc := new(pb.SEnterFreeRoom)
	//匹配可以进入的房间
	response1 := ws.matchRoom(data.ROOM_FREE)
	if response1 == nil {
		stoc.Error = pb.RoomNotExist
		ws.Send(stoc)
		return
	}
	ws.gamePid = response1.Desk
	//进入房间
	ws.entryRoom(ctx)
}
