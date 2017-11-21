package main

import (
	"goplay/glog"
	"goplay/pb"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *HallActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.GateConnect:
		arg := msg.(*pb.GateConnect)
		connected := &pb.GateConnected{
			Message: ctx.Self().String(),
		}
		arg.Sender.Tell(connected)
		glog.Infof("GateConnect %s", arg.Sender.String())
		//网关注册
		a.gates[arg.Sender.String()] = arg.Sender
	case *pb.RegistDbms:
		a.dbmsPid = arg.Sender
		rsp := new(pb.RegistedDbms)
		rsp.Message = a.Self().String()
		ctx.Respond(rsp)
	case *pb.RegistRoom:
		a.roomPid = arg.Sender
		rsp := new(pb.RegistedRoom)
		rsp.Message = a.Self().String()
		ctx.Respond(rsp)
	case *pb.RegistRole:
		a.rolePid = arg.Sender
		rsp := new(pb.RegistedRole)
		rsp.Message = a.Self().String()
		ctx.Respond(rsp)
	case *pb.LoginHall:
		arg := msg.(*pb.LoginHall)
		userid := arg.GetUserid()
		name := arg.GetNodeName()
		//断开旧连接
		k := a.roles[userid]
		if p, ok := a.gates[k]; ok {
			msg1 := new(pb.LoginElse)
			msg1.Userid = userid
			//p.Tell(msg1)
			timeout := 2 * time.Second
			res1, err1 := p.RequestFuture(msg1, timeout).Result()
			if err1 != nil {
				glog.Errorf("LoginHall err: %v", err1)
			}
			response1 := res1.(*pb.LoginedHall)
			glog.Debugf("response1: %#v", response1)
		} else {
			//增加
			a.count[name] += 1
		}
		//添加
		a.roles[userid] = name
		//响应登录
		rsp := new(pb.LoginedHall)
		rsp.Message = a.Self().String()
		ctx.Respond(rsp)
	case *pb.Logout:
		//登出成功
		arg := msg.(*pb.Logout)
		glog.Debugf("Logout userid: %s", arg.Userid)
		//减少
		userid := arg.GetUserid()
		name := a.roles[userid]
		a.count[name] -= 1
		//移除
		delete(a.roles, userid)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
