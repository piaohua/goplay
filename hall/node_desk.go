package main

import (
	"time"

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
		//按规则匹配房间
		a.matchDesk(arg, ctx)
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
		//TODO 添加桌子匹配规则
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
	case *pb.WebRequest:
		arg := msg.(*pb.WebRequest)
		rsp := new(pb.WebResponse)
		rsp.Code = arg.Code
		a.HandlerWeb(arg, rsp, ctx)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//匹配房间
func (a *HallActor) matchDesk(arg *pb.MatchDesk, ctx actor.Context) {
	rsp := new(pb.MatchedDesk)
	//已经存在
	if v, ok := a.desks[arg.Roomid]; ok {
		rsp.Desk = v
		//响应
		ctx.Respond(rsp)
		return
	}
	//新建
	rsp.Desk = a.spawnRoom(arg, ctx)
	//响应
	ctx.Respond(rsp)
}

//创建新桌子
func (a *HallActor) spawnRoom(arg *pb.MatchDesk, ctx actor.Context) *actor.PID {
	if arg.Data == "" {
		return nil
	}
	//逻辑节点存在
	if v, ok := a.serve[arg.Name]; ok {
		msg2 := new(pb.SpawnDesk)
		msg2.Data = arg.Data
		timeout := 3 * time.Second
		res2, err2 := v.RequestFuture(msg2, timeout).Result()
		if err2 != nil {
			glog.Errorf("spawnRoom err: %v", err2)
			return nil
		}
		response2 := res2.(*pb.SpawnedDesk)
		glog.Debugf("response2: %#v", response2)
		return response2.Desk
	}
	return nil
}
