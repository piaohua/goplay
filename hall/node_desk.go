package main

import (
	"strings"

	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家桌子请求处理
func (a *HallActor) HandlerDesk(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.MatchDesk:
		arg := msg.(*pb.MatchDesk)
		glog.Debugf("MatchDesk: %v", arg)
		//匹配房间
		if v, ok := a.serve[arg.Name]; ok {
			//TODO
			v.Tell(arg)
		}
	case *pb.JoinDesk:
		arg := msg.(*pb.JoinDesk)
		glog.Debugf("JoinDesk %#v", arg)
		//房间数据变更
		a.router[arg.Userid] = arg.Roomid
		a.rnums[arg.Roomid] += 1
		//响应
		//rsp := new(pb.JoinedDesk)
		//ctx.Respond(rsp)
	case *pb.AddDesk:
		arg := msg.(*pb.AddDesk)
		glog.Debugf("AddDesk %#v", arg)
		//房间数据变更
		a.desks[arg.Roomid] = arg.Desk
		a.rtype[arg.Roomid] = arg.Rtype
		//响应
		//rsp := new(pb.AddedDesk)
		//ctx.Respond(rsp)
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//移除
		delete(a.rtype, arg.Roomid)
		delete(a.rnums, arg.Roomid)
		delete(a.desks, arg.Roomid)
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		delete(a.router, arg.Userid)
		if n, ok := a.rnums[arg.Roomid]; ok && n > 0 {
			a.rnums[arg.Roomid] = n - 1
		}
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.SyncConfig:
		//主动通知同步配置
		arg := msg.(*pb.SyncConfig)
		//name := cfg.Section("game.free").Name()
		for k, v := range a.serve {
			if strings.Contains(k, "game.") {
				v.Tell(arg)
			}
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
