package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerUser(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CRegist:
		//注册消息
		arg := msg.(*pb.CRegist)
		glog.Debugf("CRegist %#v", arg)
		ws.regist(arg, ctx)
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

func (ws *WSConn) login(arg *pb.CLogin, ctx actor.Context) {
	msg1 := new(pb.RoleLogin)
	msg1.Phone = arg.GetPhone()
	msg1.Password = arg.GetPassword()
	msg1.Type = arg.GetType()
	//arg.Ipaddr = ws.GetIPAddr()
	//ws.rolePid.Request(arg, ctx.Self())
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("CLogin err: %v", err1)
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	response1 := res1.(*pb.RoleLogined)
	//TODO
	ctx.SetReceiveTimeout(0) //login Successfully, timeout off
}

func (ws *WSConn) regist(arg *pb.CRegist, ctx actor.Context) {
	stoc := RegistCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.RoleRegist)
	msg1.Phone = arg.GetPhone()
	msg1.Nickname = arg.GetNickname()
	msg1.Password = arg.GetPassword()
	msg1.Type = arg.GetType()
	//注册
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg1, timeout).Result()
	msg2 := new(pb.SRegist)
	if err1 != nil {
		glog.Errorf("CRegist err: %v", err1)
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	response1 := res1.(*pb.RoleLogined)
	if response1 == nil {
		glog.Error("CRegist fail")
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	err := json.Unmarshal([]byte(response1.Data), ws.User)
	if err != nil {
		glog.Errorf("user Unmarshal err %v", err)
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	if ws.User == nil || ws.User.GetUserid() == "" {
		glog.Error("CRegist fail")
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	//成功
	ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	ws.Send(stoc)
	//TODO
	//ws.login2(arg.GetUserid(), ctx)
	//ws.logined(arg.GetIsreg(), ctx)
	//ws.syncUser()
}

func (ws *WSConn) wxlogin(arg *pb.CWxLogin, ctx actor.Context) {
	//arg.Ipaddr = ws.GetIPAddr()
	//ws.rolePid.Request(arg, ctx.Self())
	if arg.GetUserid() != "" {
		//登录成功
		ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	}
	ws.Send(msg)
	if arg.GetUserid() != "" {
		//登录成功
		ctx.SetReceiveTimeout(0) //login Successfully, timeout off
		ws.login2(arg.GetUserid(), ctx)
		ws.logined(arg.GetIsreg(), ctx)
	}
	ws.Send(msg)
}

//登录处理
func (ws *WSConn) login2(userid string, ctx actor.Context) {
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
