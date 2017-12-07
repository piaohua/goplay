package main

import (
	"goplay/data"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerUser(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CPing:
		arg := msg.(*pb.CPing)
		rsp := handler.Ping(arg)
		ws.Send(rsp)
	case *pb.CConfig:
		arg := msg.(*pb.CConfig)
		rsp := handler.Config(arg)
		ws.Send(rsp)
	case *pb.CVipList:
		arg := msg.(*pb.CVipList)
		rsp := handler.VipList(arg)
		ws.Send(rsp)
	case *pb.CClassicList:
		arg := msg.(*pb.CClassicList)
		rsp := handler.ClassicList(arg)
		ws.Send(rsp)
	case *pb.CPrizeBox:
		arg := msg.(*pb.CPrizeBox)
		rsp, rtype, amount := handler.PrizeBox(arg, ws.User)
		ws.addPrize(rtype, data.LogType22, amount)
		ws.Send(rsp)
	case *pb.CPrizeDraw:
		arg := msg.(*pb.CPrizeDraw)
		rsp, rtype, amount := handler.PrizeDraw(arg, ws.User)
		ws.addPrize(rtype, data.LogType21, amount)
		ws.Send(rsp)
	case *pb.CPrizeList:
		arg := msg.(*pb.CPrizeList)
		rsp := handler.PrizeList(arg)
		ws.Send(rsp)
	case *pb.CBankrupts:
		arg := msg.(*pb.CBankrupts)
		rsp, coin := handler.Bankrupt(arg, ws.User)
		ws.addCurrency(0, coin, data.LogType11)
		ws.Send(rsp)
	case *pb.CBuildAgent:
		arg := msg.(*pb.CBuildAgent)
		msg1 := &pb.BuildAgent{
			Userid: arg.GetUserid(),
			Agent:  ws.User.GetAgent(),
			Uid:    ws.User.GetUserid(),
		}
		ws.rolePid.Request(msg1, ctx.Self())
	case *pb.BuiltAgent:
		arg := msg.(*pb.BuiltAgent)
		if arg.Result == 0 {
			//绑定
			ws.User.SetAgent(arg.Agent)
			//日志
			msg2 := &pb.LogBuildAgency{
				Userid: ws.User.GetUserid(),
				Agent:  arg.Agent,
			}
			ws.dbmsPid.Tell(msg2)
			//赠送
			var diamond int32 = config.GetEnv(data.ENV3)
			ws.addCurrency(diamond, 0, data.LogType19)
		}
		rsp := new(pb.SBuildAgent)
		rsp.Result = arg.Result
		ws.Send(rsp)
	case *pb.CBank:
		arg := msg.(*pb.CBank)
		rtype := arg.GetRtype()
		switch rtype {
		case 1:
			rsp, coin, _ := handler.Bank(arg, ws.User)
			ws.addCurrency(0, coin, data.LogType12)
			ws.Send(rsp)
		case 2:
			rsp, coin, tax := handler.Bank(arg, ws.User)
			ws.addCurrency(0, coin, data.LogType13)
			ws.addCurrency(0, tax, data.LogType14)
			ws.Send(rsp)
		case 3:
			res := handler.BankGive(arg, ws.User)
			if res.Error != pb.OK {
				ws.Send(res)
				return
			}
			msg2 := new(pb.BankGive)
			msg2.Userid = arg.Userid
			msg2.Amount = arg.Amount
			ws.rolePid.Request(msg2, ctx.Self())
		case 4:
		}
	case *pb.BankGave:
		arg := msg.(*pb.BankGave)
		rsp, coin, tax := handler.BankGave(arg, ws.User)
		ws.bank(arg, coin, tax)
		ws.Send(rsp)
	case *pb.CGetCurrency:
		arg := msg.(*pb.CGetCurrency)
		//响应
		rsp := handler.GetCurrency(arg, ws.User)
		ws.Send(rsp)
	case *pb.CBuy:
		arg := msg.(*pb.CBuy)
		//优化
		rsp, diamond, coin := handler.Buy(arg, ws.User)
		//同步兑换
		ws.addCurrency(diamond, coin, data.LogType18)
		//响应
		ws.Send(rsp)
	case *pb.CShop:
		arg := msg.(*pb.CShop)
		//响应
		rsp := handler.Shop(arg, ws.User)
		ws.Send(rsp)
	case *pb.CUserData:
		arg := msg.(*pb.CUserData)
		glog.Debugf("CUserData %#v", arg)
		userid := arg.GetUserid()
		if userid != ws.User.GetUserid() && userid != "" {
			msg1 := new(pb.GetUserData)
			msg1.Userid = userid
			ws.rolePid.Request(msg1, ctx.Self())
		} else {
			//TODO 添加房间数据返回
			rsp := handler.GetUserData3(arg, ws.User)
			ws.Send(rsp)
		}
	case *pb.GotUserData:
		arg := msg.(*pb.GotUserData)
		glog.Debugf("GotUserData %#v", arg)
		rsp := handler.GetUserData2(arg)
		ws.Send(rsp)
	default:
		//glog.Errorf("unknown message %v", msg)
		ws.HandlerPay(msg, ctx)
	}
}

func (ws *WSConn) bank(arg *pb.BankGave, coin, tax int32) {
	//受赠者货币变更, TODO 在线更新消息
	msg4 := &pb.ChangeCurrency{
		Userid: arg.Userid,
		Coin:   coin,
		Type:   int32(data.LogType15),
		Upsert: true,
	}
	ws.rolePid.Tell(msg4)
	//受赠者日志
	msg3 := &pb.LogCoin{
		Userid: arg.Userid,
		Type:   int32(data.LogType15),
		Num:    coin,
		Rest:   arg.Coin + uint32(coin),
	}
	ws.dbmsPid.Tell(msg3)
	//赠送者日志
	msg1 := &pb.LogCoin{
		Userid: ws.User.GetUserid(),
		Type:   int32(data.LogType15),
		Num:    (-1 * coin),
		Rest:   ws.User.GetCoin(),
	}
	ws.dbmsPid.Tell(msg1)
	//扣税日志
	msg2 := &pb.LogCoin{
		Userid: ws.User.GetUserid(),
		Type:   int32(data.LogType16),
		Num:    tax,
		Rest:   ws.User.GetBank(),
	}
	ws.dbmsPid.Tell(msg2)
}
