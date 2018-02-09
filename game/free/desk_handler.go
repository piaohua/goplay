package main

import (
	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskFree) HandlerLogic(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//TODO
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//TODO
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.PrintDesk:
		//打印牌局状态信息,test
		a.print()
	case *pb.EnterDesk:
		arg := msg.(*pb.EnterDesk)
		glog.Debugf("EnterDesk %#v", arg)
		a.enterDesk(arg, ctx)
	case *pb.CEnterFreeRoom:
		arg := msg.(*pb.CEnterFreeRoom)
		glog.Debugf("CEnterFreeRoom %#v", arg)
		//检测重复进入
		msg1 := a.res_reEnter()
		ctx.Respond(msg1)
	case *pb.CChatText:
		arg := msg.(*pb.CChatText)
		glog.Debugf("CChatText %#v", arg)
		userid := a.getRouter(ctx)
		seat := a.seats[userid]
		//房间消息广播,聊天
		a.broadcast(handler.ChatMsg(seat, userid, arg.Content))
	case *pb.CChatVoice:
		arg := msg.(*pb.CChatVoice)
		glog.Debugf("CChatVoice %#v", arg)
		userid := a.getRouter(ctx)
		seat := a.seats[userid]
		//房间消息广播,聊天
		a.broadcast(handler.ChatMsg2(seat, userid, arg.Content))
	case *pb.CFreeDealer:
		arg := msg.(*pb.CFreeDealer)
		glog.Debugf("CFreeDealer %#v", arg)
		//userid := a.router[ctx.Sender().String()]
		userid := a.getRouter(ctx)
		var state uint32 = arg.GetState()
		var num uint32 = arg.GetCoin()
		errcode := a.BeDealer(userid, state, num)
		if errcode == pb.OK {
			return
		}
		//响应
		rsp := new(pb.SFreeDealer)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.CDealerList:
		arg := msg.(*pb.CDealerList)
		glog.Debugf("CDealerList %#v", arg)
		//userid := a.getRouter(ctx)
		//上庄列表
		rsp := a.res_bedealerlist()
		ctx.Respond(rsp)
	case *pb.CFreeSit:
		arg := msg.(*pb.CFreeSit)
		glog.Debugf("CFreeSit %#v", arg)
		userid := a.getRouter(ctx)
		var seat uint32 = arg.GetSeat()
		var state bool = arg.GetState()
		errcode := a.SitDown(userid, seat, state)
		if errcode == pb.OK {
			return
		}
		//响应
		rsp := new(pb.SFreeSit)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.CFreeBet:
		arg := msg.(*pb.CFreeBet)
		glog.Debugf("CFreeBet %#v", arg)
		userid := a.getRouter(ctx)
		value := arg.GetValue()
		seatBet := arg.GetSeat()
		errcode := a.ChoiceBet(userid, seatBet, value)
		if errcode == pb.OK {
			return
		}
		//响应
		rsp := new(pb.SFreeBet)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.CLeave:
		arg := msg.(*pb.CLeave)
		glog.Debugf("CLeave %#v", arg)
		userid := a.getRouter(ctx)
		errcode := a.Leave(userid)
		if errcode == pb.OK {
			return
		}
		//离开消息
		msg2 := &pb.LeaveDesk{
			Roomid: a.id,
			Userid: userid,
		}
		nodePid.Tell(msg2)
		ctx.Respond(msg2)
		//响应
		rsp := new(pb.SLeave)
		rsp.Error = errcode
		ctx.Respond(rsp)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		//充值或购买同步
		a.changeCurrency(arg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//获取路由
func (a *DeskFree) getRouter(ctx actor.Context) string {
	return a.router[ctx.Sender().String()]
}

//进入
func (a *DeskFree) enterDesk(arg *pb.EnterDesk, ctx actor.Context) {
	rsp := new(pb.EnteredDesk)
	user := new(data.User)
	err2 := json.Unmarshal([]byte(arg.Data), user)
	if err2 != nil {
		glog.Errorf("user Unmarshal err %v", err2)
		rsp.Error = pb.RoomNotExist
		ctx.Respond(rsp)
		return
	}
	errcode := a.Enter(user)
	if errcode != pb.OK {
		rsp.Error = errcode
		ctx.Respond(rsp)
		return
	}
	ctx.Respond(rsp)
	//加入游戏
	a.pids[user.Userid] = ctx.Sender()
	//设置路由
	a.router[ctx.Sender().String()] = user.Userid
	//进入消息
	msg3 := new(pb.JoinDesk)
	msg3.Roomid = a.id
	msg3.Rtype = a.data.Rtype
	msg3.Userid = user.Userid
	msg3.Sender = ctx.Sender()
	nodePid.Request(msg3, ctx.Self())
}

//收益
func (a *DeskFree) sendFreeCoin(p *data.User, num int32, rtype int) {
	if num == 0 {
		return
	}
	if p == nil {
		return
	}
	p.AddCoin(num)
	a.syncCoin(p.GetUserid(), num, rtype)
}

//同步收益
func (a *DeskFree) syncCoin(userid string, coin int32, ltype int) {
	//货币变更及时同步
	msg2 := &pb.ChangeCurrency{
		Userid: userid,
		Coin:   coin,
		Type:   int32(ltype),
	}
	if v, ok := a.pids[userid]; ok {
		v.Tell(msg2)
	} else {
		nodePid.Tell(msg2)
	}
}

//更新货币
func (a *DeskFree) changeCurrency(arg *pb.ChangeCurrency) {
	p := a.getPlayer(arg.Userid)
	if p == nil {
		return
	}
	p.AddCoin(arg.Coin)
	p.AddDiamond(arg.Diamond)
}
