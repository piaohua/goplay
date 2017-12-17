package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/game/login"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家数据请求处理
func (ws *WSConn) HandlerLogin(msg interface{}, ctx actor.Context) {
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
		//glog.Errorf("unknown message %v", msg)
		if ws.User == nil {
			glog.Errorf("user empty message %v", msg)
			return
		}
		ws.HandlerUser(msg, ctx)
	}
}

func (ws *WSConn) regist(arg *pb.CRegist, ctx actor.Context) {
	stoc := login.RegistCheck(arg)
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
	if err1 != nil {
		glog.Errorf("CRegist err: %v", err1)
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	response1 := res1.(*pb.RoleRegisted)
	if response1 == nil {
		glog.Error("CRegist fail")
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	if response1.Error != pb.OK {
		glog.Error("CRegist fail")
		stoc.Error = response1.Error
		ws.Send(stoc)
		return
	}
	ws.User = new(data.User)
	err := json.Unmarshal([]byte(response1.Data), ws.User)
	if err != nil {
		glog.Errorf("user Unmarshal err %v", err)
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	if ws.User.GetUserid() == "" {
		glog.Error("CRegist fail")
		stoc.Error = pb.RegistError
		ws.Send(stoc)
		return
	}
	//成功
	//ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	stoc.Userid = ws.User.GetUserid()
	ws.Send(stoc)
	//成功后处理
	ws.logined(true, ctx)
}

func (ws *WSConn) login(arg *pb.CLogin, ctx actor.Context) {
	stoc := login.LoginCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.RoleLogin)
	msg1.Phone = arg.GetPhone()
	msg1.Password = arg.GetPassword()
	msg1.Type = arg.GetType()
	//登录
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("CLogin err: %v", err1)
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	response1 := res1.(*pb.RoleLogined)
	if response1 == nil {
		glog.Error("CLogin fail")
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	if response1.Error != pb.OK {
		glog.Error("CRegist fail")
		stoc.Error = response1.Error
		ws.Send(stoc)
		return
	}
	ws.User = new(data.User)
	err := json.Unmarshal([]byte(response1.Data), ws.User)
	if err != nil {
		glog.Errorf("user Unmarshal err %v", err)
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	if ws.User.GetUserid() == "" {
		glog.Error("CLogin fail")
		stoc.Error = pb.LoginError
		ws.Send(stoc)
		return
	}
	//成功
	//ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	stoc.Userid = ws.User.GetUserid()
	ws.Send(stoc)
	//成功后处理
	ws.logined(false, ctx)
}

func (ws *WSConn) wxlogin(arg *pb.CWxLogin, ctx actor.Context) {
	stoc, wxdata := login.WxLoginCheck(arg)
	if stoc.Error != pb.OK {
		ws.Send(stoc)
		return
	}
	msg1 := new(pb.WxLogin)
	msg1.Wxuid = wxdata.OpenId
	msg1.Nickname = wxdata.Nickname
	msg1.Photo = wxdata.HeadImagUrl
	msg1.Sex = uint32(wxdata.Sex)
	msg1.Type = arg.GetType()
	//登录
	timeout := 3 * time.Second
	res1, err1 := ws.rolePid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("CWxLogin err: %v", err1)
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	response1 := res1.(*pb.WxLogined)
	if response1 == nil {
		glog.Error("CWxLogin fail")
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	if response1.Error != pb.OK {
		glog.Error("CWxLogin fail")
		stoc.Error = response1.Error
		ws.Send(stoc)
		return
	}
	ws.User = new(data.User)
	err := json.Unmarshal([]byte(response1.Data), ws.User)
	if err != nil {
		glog.Errorf("user Unmarshal err %v", err)
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	if ws.User.GetUserid() == "" {
		glog.Error("CWxLogin fail")
		stoc.Error = pb.GetWechatUserInfoFail
		ws.Send(stoc)
		return
	}
	//成功
	//ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	stoc.Userid = ws.User.GetUserid()
	ws.Send(stoc)
	//成功后处理
	ws.logined(response1.IsRegist, ctx)
}

//登录成功处理
func (ws *WSConn) logined(isRegist bool, ctx actor.Context) {
	//重连处理直接断开旧连接,登录成功再连接节点和大厅
	msg1 := new(pb.LoginHall)
	msg1.Userid = ws.User.GetUserid()
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
	msg2.Userid = ws.User.GetUserid()
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
	result4, err4 := json.Marshal(ws.User)
	if err4 != nil {
		glog.Errorf("user Marshal err %v", err)
	}
	msg4 := new(pb.Login)
	msg4.Userid = ws.User.GetUserid()
	msg4.Data = string(result4)
	res4, err4 := ws.rolePid.RequestFuture(msg4, timeout).Result()
	if err4 != nil {
		//TODO 断开
		glog.Errorf("Login err: %v", err4)
	}
	response4 := res4.(*pb.Logined)
	glog.Debugf("response4: %#v", response4)
	//日志
	ws.logined2(isRegist)
	//同步数据
	ws.syncUser()
	//登录成功
	ws.online = true
	//成功
	ctx.SetReceiveTimeout(0) //login Successfully, timeout off
	//启动时钟
	go ws.ticker(ctx)
}

//登录处理
func (ws *WSConn) logined2(isRegist bool) {
	ws.User.LoginIp = ws.GetIPAddr()
	if isRegist {
		//注册ip
		ws.User.RegIp = ws.GetIPAddr()
		//注册奖励发放
		var diamond int32 = config.GetEnv(data.ENV1)
		var coin int32 = config.GetEnv(data.ENV2)
		ws.addCurrency(diamond, coin, data.LogType1)
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
