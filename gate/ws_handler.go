package main

import (
	"encoding/json"
	"time"

	"goplay/data"
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
		ws.betsPid = arg.BetsPid
		ws.mailPid = arg.MailPid
		glog.Infof("SetLogined %v", arg.Message)
	case *pb.ServeClose:
		arg := new(pb.SLoginOut)
		arg.Rtype = 2 //停服
		ws.Send(arg)
		//断开连接
		ws.Close()
	case *pb.ServeStop:
		ws.stop(ctx)
		//响应
		rsp := new(pb.ServeStarted)
		ctx.Respond(rsp)
	case *pb.ServeStoped:
	case *pb.ServeStart:
		ws.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.ServeStarted:
	case *pb.Tick:
		ws.ding(ctx)
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

func (ws *WSConn) start(ctx actor.Context) {
	glog.Infof("ws start: %v", ctx.Self().String())
	ctx.SetReceiveTimeout(waitForLogin) //login timeout set
	set := &pb.SetLogin{
		Sender: ctx.Self(),
	}
	nodePid.Tell(set)
}

//时钟
func (ws *WSConn) ticker(ctx actor.Context) {
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-ws.stopCh:
			glog.Info("ws ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-ws.stopCh:
			glog.Info("ws ticker closed")
			return
		case <-tick:
			if ws.pid != nil {
				ws.pid.Tell(msg)
			}
		}
	}
}

//30秒同步一次
func (ws *WSConn) ding(ctx actor.Context) {
	ws.timer += 1
	if ws.timer != 30 {
		return
	}
	ws.timer = 0
	if !ws.online {
		//断开连接
		ws.Close()
		return
	}
	if ws.status {
		//同步数据
		ws.syncUser()
		ws.status = false
	}
}

func (ws *WSConn) stop(ctx actor.Context) {
	glog.Infof("ws stop: %v", ctx.Self().String())
	//已经断开,在别处登录
	if ws.hallPid == nil {
		return
	}
	if ws.User == nil {
		return
	}
	if ws.gamePid != nil {
		//TODO 下线
	}
	//回存数据
	ws.syncUser()
	//登出日志
	if ws.dbmsPid != nil {
		msg2 := &pb.LogLogout{
			Userid: ws.User.Userid,
			Event:  1,
		}
		ws.dbmsPid.Tell(msg2)
	}
	//断开处理
	msg := &pb.Logout{
		Sender: ctx.Self(),
		Userid: ws.User.Userid,
	}
	nodePid.Tell(msg)
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
	//机器人不写日志
	if ws.User.GetPhone() != "" {
		return
	}
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

//消耗钻石
func (ws *WSConn) expend(cost uint32, ltype int) {
	diamond := -1 * int32(cost)
	ws.addCurrency(diamond, 0, ltype)
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
