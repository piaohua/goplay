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
		ws.rolePid.Tell(msg1)
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
		case 4:
		}
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
		//userid := arg.GetUserid()
		//TODO userid != ws.GetUserid()
		//TODO room data
		rsp := handler.GetUserData(arg, ws.User)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
