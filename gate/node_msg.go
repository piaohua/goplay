package main

import (
	"goplay/glog"
	"goplay/pb"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (a *GateActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.GateConnected:
		//连接成功
		arg := msg.(*pb.GateConnected)
		glog.Infof("GateConnected %v", arg.Message)
	case *pb.GateDisconnected:
		//成功断开
		arg := msg.(*pb.GateDisconnected)
		glog.Infof("GateDisconnected %s", arg.Message)
	case *pb.LoginElse:
		//别处登录
		arg := msg.(*pb.LoginElse)
		userid := arg.GetUserid()
		if p, ok := a.roles[userid]; ok {
			p.Tell(arg)
		}
		glog.Debugf("LoginElse userid: %s", userid)
		//移除
		delete(a.roles, userid)
		//响应登录
		rsp := new(pb.LoginedElse)
		rsp.Userid = userid
		ctx.Respond(rsp)
	case *pb.SetLogin:
		arg := msg.(*pb.SetLogin)
		set := &pb.SetLogined{
			Message: ctx.Self().String(),
			DbmsPid: a.dbmsPid,
			RoomPid: a.roomPid,
			RolePid: a.rolePid,
			HallPid: a.hallPid,
		}
		arg.Sender.Tell(set)
		glog.Infof("SetLogin %s", arg.Sender.String())
	case *pb.LoginGate:
		//登录成功
		arg := msg.(*pb.LoginGate)
		userid := arg.GetUserid()
		//添加
		a.roles[userid] = arg.Sender
		glog.Debugf("LoginGate userid: %s", userid)
		//响应登录
		rsp := new(pb.LoginedGate)
		rsp.Message = ctx.Self().String()
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		userid := arg.GetUserid()
		glog.Debugf("Logout userid: %s", userid)
		//移除
		delete(a.roles, userid)
		a.hallPid.Tell(arg)
		a.rolePid.Tell(arg)
		a.roomPid.Tell(arg)
	case *pb.HallConnect:
		//初始化建立连接
		a.init(ctx)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *GateActor) init(ctx actor.Context) {
	glog.Infof("gate init: %v", ctx.Self().String())
	//dbms
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	role := cfg.Section("cookie").Key("role").Value()
	timeout := 3 * time.Second
	//a.dbmsPid, err = remote.SpawnNamed(bind, a.Name, name, timeout)
	//if err != nil {
	//	glog.Fatalf("remote dbms err %v", err)
	//}
	//a.roomPid, err = remote.SpawnNamed(bind, a.Name, room, timeout)
	//if err != nil {
	//	glog.Fatalf("remote room err %v", err)
	//}
	//a.rolePid, err = remote.SpawnNamed(bind, a.Name, role, timeout)
	//if err != nil {
	//	glog.Fatalf("remote role err %v", err)
	//}
	a.dbmsPid = actor.NewPID(bind, name)
	a.roomPid = actor.NewPID(bind, room)
	a.rolePid = actor.NewPID(bind, role)
	//a.dbmsPid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//a.roomPid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//a.rolePid.
	//	RequestFuture(&pb.Request{}, 2*time.Second).
	//	Wait()
	//hall
	bind = cfg.Section("hall").Key("bind").Value()
	//name
	hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
	if err != nil {
		glog.Fatalf("remote hall err %v", err)
	}
	a.hallPid = hallPid.Pid
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.GateConnect{
		Sender: ctx.Self(),
	}
	a.dbmsPid.Tell(connect)
	a.hallPid.Tell(connect)
	glog.Infof("a.dbmsPid: %s", a.dbmsPid.String())
	glog.Infof("a.roomPid: %s", a.roomPid.String())
	glog.Infof("a.rolePid: %s", a.rolePid.String())
}

func (a *GateActor) disc(ctx actor.Context) {
	//TODO 关闭消息
	//for k, v := range a.roles {
	//}
	//断开处理
	msg := &pb.GateDisconnect{
		Sender: ctx.Self(),
	}
	if a.dbmsPid != nil {
		a.dbmsPid.Tell(msg)
	}
	if a.roomPid != nil {
		a.roomPid.Tell(msg)
	}
	if a.rolePid != nil {
		a.rolePid.Tell(msg)
	}
	if a.hallPid != nil {
		a.hallPid.Tell(msg)
	}
}
