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
		glog.Infof("SetLogined %v", arg.Message)
	case *pb.CRegist:
		//注册消息
		arg := msg.(*pb.CRegist)
		arg.Ipaddr = ws.GetIPAddr()
		ws.rolePid.Request(arg, ctx.Self())
		glog.Debugf("CRegist %#v", arg)
	case *pb.CLogin:
		//登录消息
		arg := msg.(*pb.CLogin)
		arg.Ipaddr = ws.GetIPAddr()
		ws.rolePid.Request(arg, ctx.Self())
		glog.Debugf("CLogin %#v", arg)
	case *pb.CWxLogin:
		//登录消息
		arg := msg.(*pb.CWxLogin)
		arg.Ipaddr = ws.GetIPAddr()
		ws.rolePid.Request(arg, ctx.Self())
		glog.Debugf("CWxLogin %#v", arg)
	case *pb.SRegist:
		arg := msg.(*pb.SRegist)
		glog.Debugf("SRegist %#v", arg)
		if arg.GetUserid() != "" {
			//登录成功
			ctx.SetReceiveTimeout(0) //login Successfully, timeout off
		}
		ws.Send(msg)
	case *pb.SLogin:
		arg := msg.(*pb.SLogin)
		glog.Debugf("SLogin %#v", arg)
		if arg.GetUserid() != "" {
			//登录成功
			ctx.SetReceiveTimeout(0) //login Successfully, timeout off
		}
		ws.Send(msg)
	case *pb.SWxLogin:
		arg := msg.(*pb.SWxLogin)
		glog.Debugf("SWxLogin %#v", arg)
		if arg.GetUserid() != "" {
			//登录成功
			ctx.SetReceiveTimeout(0) //login Successfully, timeout off
			ws.login(arg.GetUserid(), ctx)
			ws.logined(arg.GetIsreg(), ctx)
		}
		ws.Send(msg)
	case *pb.SyncUser:
		arg := msg.(*pb.SyncUser)
		glog.Debugf("SyncUser %#v", arg.Userid)
		err := json.Unmarshal([]byte(arg.Data), ws.User)
		if err != nil {
			glog.Errorf("user Unmarshal err %v", err)
		}
		glog.Debugf("User %#v", ws.User)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		diamond := arg.Diamond
		coin := arg.Coin
		ltype := int(arg.Type)
		ws.addCurrency(diamond, coin, ltype, ctx)
	case *pb.LoginElse:
		ws.loginElse() //别处登录
	case *pb.ServeStop:
		arg := new(pb.SLoginOut)
		arg.Rtype = 2 //停服
		ws.Send(arg)
		//断开连接
		ws.Close()
	case *pb.CBuy, *pb.CShop:
		ws.rolePid.Request(msg, ctx.Self())
	case proto.Message:
		//响应消息
		ws.Send(msg)
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

//登录处理
func (ws *WSConn) login(userid string, ctx actor.Context) {
	ws.User = new(data.User)
	ws.User.Userid = userid
	//重连处理直接断开旧连接,登录成功再连接节点和大厅
	//TODO 不在同一节点就关闭旧的,如果在就接替或直接关闭
	msg1 := new(pb.LoginHall)
	msg1.Userid = userid
	msg1.NodeName = nodePid.String()
	timeout := 3 * time.Second
	res1, err1 := ws.hallPid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("LoginHall err: %v", err1)
	}
	response1 := res1.(*pb.LoginedHall)
	glog.Debugf("response1: %#v", response1)
	//成功后登录网关
	msg2 := new(pb.LoginGate)
	msg2.Sender = ctx.Self()
	msg2.Userid = userid
	res2, err2 := nodePid.RequestFuture(msg2, timeout).Result()
	if err2 != nil {
		glog.Errorf("LoginGate err: %v", err2)
	}
	response2 := res2.(*pb.LoginedGate)
	glog.Debugf("response2: %#v", response2)
	//TODO 查看房间数据
	//msg3 := &pb.LoginGate{}
	//res3, err3 := ws.roomPid.RequestFuture(msg3, timeout).Result()
	//response3 := res3.(*pb.LoginedGate)
	//登录成功响应
	msg4 := new(pb.Login)
	msg4.Userid = userid
	res4, err4 := ws.rolePid.RequestFuture(msg4, timeout).Result()
	if err4 != nil {
		//TODO 断开
		glog.Errorf("Login err: %v", err4)
	}
	response4 := res4.(*pb.Logined)
	glog.Debugf("response4: %#v", response4)
}

//登录处理
func (ws *WSConn) logined(isRegist bool, ctx actor.Context) {
	if ws.User == nil {
		//登录失败
		return
	}
	if isRegist {
		//注册日志
		msg1 := &pb.LogRegist{
			Userid:   ws.User.Userid,
			Nickname: ws.User.Nickname,
			Atype:    ws.User.Atype,
			Ip:       ws.GetIPAddr(),
		}
		ws.dbmsPid.Tell(msg1)
	}
	//登录日志
	msg2 := &pb.LogLogin{
		Userid: ws.User.Userid,
		Atype:  ws.User.Atype,
		Ip:     ws.GetIPAddr(),
	}
	ws.dbmsPid.Tell(msg2)
}

//奖励发放
func (ws *WSConn) addCurrency(diamond, coin int32,
	ltype int, ctx actor.Context) {
	if ws.User == nil {
		glog.Errorf("add currency user empty: %d", ltype)
		return
	}
	ws.User.AddDiamond(diamond)
	ws.User.AddCoin(coin)
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
