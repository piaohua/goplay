package main

import (
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
	case *pb.EnterRoom:
		arg := msg.(*pb.EnterRoom)
		glog.Debugf("EnterRoom %#v", arg)
		//房间数据变更
		a.router[arg.Userid] = arg.Roomid
		a.rnums[arg.Roomid] += 1
		//响应
		//rsp := new(pb.EnteredRoom)
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
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
