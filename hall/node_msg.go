package main

import (
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
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
		rsp := &pb.RegistedDbms{
			Message: a.Self().String(),
		}
		ctx.Respond(rsp)
	case *pb.RegistRoom:
		a.roomPid = arg.Sender
		rsp := &pb.RegistedRoom{
			Message: a.Self().String(),
		}
		ctx.Respond(rsp)
	case *pb.RegistRole:
		a.rolePid = arg.Sender
		rsp := &pb.RegistedRole{
			Message: a.Self().String(),
		}
		ctx.Respond(rsp)
	case *pb.LoginHall:
		arg := msg.(*pb.LoginHall)
		userid := arg.GetUserid()
		name := arg.GetNodeName()
		//断开旧连接
		k := a.roles[userid]
		if p, ok := a.gates[k]; ok {
			msg1 := &pb.LoginElse{
				Userid: proto.String(userid),
			}
			p.Tell(msg1)
		} else {
			a.count[name] += 1
		}
		a.roles[userid] = name
		//响应登录
		rsp := &pb.LoginedHall{
			Message: a.Self().String(),
		}
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
