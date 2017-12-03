package main

import (
	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerUser(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
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
		userid := ctos.GetUserid()
		//TODO userid != ws.GetUserid()
		//TODO room data
		rsp := handler.GetUserData(arg, ws.User)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
