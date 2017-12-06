package main

import (
	"encoding/json"
	"niu/data"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
)

func (ws *WSConn) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.SetLogined:
		//设置连接,还未登录
		arg := msg.(*pb.SetLogined)
		ws.dbmsPid = arg.DbmsPid
		ws.roomPid = arg.RoomPid
		ws.rolePid = arg.RolePid
		ws.hallPid = arg.HallPid
		glog.Infof("SetLogined %v", arg.Message)
	case *pb.ServeStop:
		arg := new(pb.SLoginOut)
		arg.Rtype = 2 //停服
		ws.Send(arg)
		//断开连接
		ws.Close()
	case *pb.LoginElse:
		ws.loginElse() //别处登录
	case *pb.SyncUser:
		//同步数据
		arg := msg.(*pb.SyncUser)
		glog.Debugf("SyncUser %#v", arg.Userid)
		err := json.Unmarshal([]byte(arg.Data), ws.User)
		if err != nil {
			glog.Errorf("user Unmarshal err %v", err)
		}
		glog.Debugf("User %#v", ws.User)
	case *pb.ChangeCurrency:
		//货币变更
		arg := msg.(*pb.ChangeCurrency)
		diamond := arg.Diamond
		coin := arg.Coin
		ltype := int(arg.Type)
		ws.addCurrency(diamond, coin, ltype)
	case proto.Message:
		//响应消息
		//ws.Send(msg)
		ws.HandlerLogin(msg, ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (ws *WSConn) loginElse() {
	arg := new(pb.SLoginOut)
	glog.Debugf("SLoginOut %s", ws.User.Userid)
	arg.Rtype = 1 //别处登录
	ws.Send(arg)
	//登出日志
	msg3 := &pb.LogLogout{
		Userid: ws.User.Userid,
		Event:  4, //别处登录
	}
	ws.dbmsPid.Tell(msg3)
	//表示已经断开
	ws.hallPid = nil
	//断开连接
	ws.Close()
}

func (ws *WSConn) addPrize(rtype, ltype int, amount int32) {
	switch uint32(rtype) {
	case data.DIAMOND:
		ws.addCurrency(amount, 0, ltype)
	case data.COIN:
		ws.addCurrency(0, amount, ltype)
	}
}

//奖励发放
func (ws *WSConn) addCurrency(diamond, coin int32, ltype int) {
	if ws.User == nil {
		glog.Errorf("add currency user empty: %d", ltype)
		return
	}
	ws.User.AddDiamond(diamond)
	ws.User.AddCoin(coin)
	//货币变更及时同步
	msg2 := &pb.ChangeCurrency{
		Userid:  ws.User.GetUserid(),
		Diamond: diamond,
		Coin:    coin,
		Type:    int32(ltype),
	}
	ws.rolePid.Tell(msg2)
	//消息
	msg := &pb.SPushCurrency{
		Rtype:   uint32(ltype),
		Diamond: diamond,
		Coin:    coin,
	}
	ws.Send(msg)
	//日志
	if diamond != 0 {
		msg1 := &pb.LogDiamond{
			Userid: ws.User.GetUserid(),
			Type:   int32(ltype),
			Num:    diamond,
			Rest:   ws.User.GetDiamond(),
		}
		ws.dbmsPid.Tell(msg1)
	}
	if coin != 0 {
		msg1 := &pb.LogCoin{
			Userid: ws.User.GetUserid(),
			Type:   int32(ltype),
			Num:    coin,
			Rest:   ws.User.GetCoin(),
		}
		ws.dbmsPid.Tell(msg1)
	}
}

//同步数据
func (ws *WSConn) syncUser() {
	if ws.User == nil {
		return
	}
	if ws.rolePid == nil {
		return
	}
	msg := new(pb.SyncUser)
	msg.Userid = ws.User.GetUserid()
	result, err := json.Marshal(ws.User)
	if err != nil {
		glog.Errorf("user %s Marshal err %v", ws.User.GetUserid(), err)
		return
	}
	msg.Data = string(result)
	ws.rolePid.Tell(msg)
}
