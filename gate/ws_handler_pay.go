package main

import (
	"math"
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerPay(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CApplePay:
		arg := msg.(*pb.CApplePay)
		glog.Debugf("CApplePay %#v", arg)
		ws.applePay(arg, ctx)
	case *pb.CLogin:
		//登录消息
		arg := msg.(*pb.CLogin)
		glog.Debugf("CLogin %#v", arg)
		ws.login(arg, ctx)
	case *pb.CWxLogin:
		//登录消息
		arg := msg.(*pb.CWxLogin)
		glog.Debugf("CWxLogin %#v", arg)
		ws.wxlogin(arg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (ws *WSConn) applePay(arg *pb.CApplePay, ctx actor.Context) {
	rsp, record, trade := handler.AppleOrder(arg, ws.User)
	if rsp.Error != pb.OK {
		ws.Send(rsp)
		return
	}
	//验证
	msg1 := new(pb.ApplePay)
	msg1.Trade = trade
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("ApplePay err: %v", err1)
		rsp.Error = pb.AppleOrderFail
		ws.Send(rsp)
		return
	}
	response1 := res1.(*pb.ApplePaid)
	if response1 == nil {
		glog.Error("ApplePay fail")
		rsp.Error = pb.AppleOrderFail
		ws.Send(rsp)
		return
	}
	if !response1.Result {
		glog.Error("ApplePay fail")
		rsp.Error = pb.AppleOrderFail
		ws.Send(rsp)
		return
	}
	ws.sendGoods(record.Diamond, record.Money, record.First)
	ws.Send(rsp)
}

//发货
func (ws *WSConn) sendGoods(diamond, money uint32, first int) {
	ws.User.AddMoney(money)
	//消息
	ws.addCurrency(int32(diamond), 0, data.LogType4)
	//消息
	stoc := &pb.SGetCurrency{
		Coin:    ws.User.GetCoin(),
		Diamond: ws.User.GetDiamond(),
	}
	ws.Send(stoc)
	//奖励
	ws.firstPay(diamond, first)
	ws.vipGive(diamond, money)
}

//vip赠送
func (ws *WSConn) vipGive(num, money uint32) {
	level := ws.User.GetVipLevel()
	if level == 0 {
		lev := config.GetVipLevel(money)
		ws.User.SetVip(lev, money)
		return
	}
	diamond := ws.getVipGive(level, int32(num))
	if diamond > 0 {
		//消息
		ws.addCurrency(int32(diamond), 0, data.LogType28)
	}
	lev2 := config.GetVipLevel(ws.User.GetVip() + money)
	ws.User.SetVip(lev2, money)
	//vip改变通知
	stoc := &pb.SPushVip{
		Level:  uint32(lev2),
		Number: ws.User.GetVip(),
	}
	ws.Send(stoc)
}

//vip赠送
func (ws *WSConn) getVipGive(level int, num int32) int32 {
	if level <= 0 {
		return 0
	}
	pay := config.GetVipPay(level)
	return int32(math.Ceil(float64(num) * (float64(pay) / 100)))
}

//first == 1 首充
func (ws *WSConn) firstPay(num uint32, first int) {
	if first != 1 {
		return
	}
	var diamond int32 = config.GetEnv(data.ENV4)
	var coin int32 = config.GetEnv(data.ENV5)
	if diamond == 0 && coin == 0 {
		return
	}
	ws.addCurrency(diamond, coin, data.LogType20)
}
