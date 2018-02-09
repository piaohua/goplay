package main

import (
	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *DeskActor) HandlerMsg(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		//连接成功
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		//成功断开
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.SpawnDesk:
		arg := msg.(*pb.SpawnDesk)
		glog.Debugf("SpawnDesk %#v", arg)
		a.spawnDesk(arg, ctx)
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//移除
		delete(a.desks, arg.Roomid)
		delete(a.count, arg.Roomid)
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		delete(a.roles, arg.Userid)
		if n, ok := a.count[arg.Roomid]; ok && n > 0 {
			a.count[arg.Roomid] = n - 1
		}
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.JoinDesk:
		arg := msg.(*pb.JoinDesk)
		glog.Debugf("JoinDesk %#v", arg)
		//房间数据变更
		a.roles[arg.Userid] = arg.Sender
		a.count[arg.Roomid] += 1
		a.roomPid.Request(arg, ctx.Self())
		a.hallPid.Request(arg, ctx.Self())
		//响应
		//rsp := new(pb.EnteredRoom)
		//ctx.Respond(rsp)
	case *pb.SyncConfig:
		//同步配置
		arg := msg.(*pb.SyncConfig)
		handler.SyncConfig(arg)
	case *pb.ChangeCurrency:
		arg := msg.(*pb.ChangeCurrency)
		//TODO
		//离线用户结算问题通知到dbms,
		//应该检测是否已经登录其它节点，
		//如果有要同步更新到其它节点
		//如果没有dbms
		//a.changeCurrency(arg)
		//a.hallPid
		//a.dbmsPid
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动新服务
func (a *DeskActor) spawnDesk(arg *pb.SpawnDesk, ctx actor.Context) {
	rsp := new(pb.SpawnedDesk)
	//解析
	deskData := handler.Data2Desk(arg.Data)
	if deskData == nil {
		glog.Error("createRoom failed")
		ctx.Respond(rsp)
		return
	}
	switch deskData.Rtype {
	case data.ROOM_FREE:
		//新桌子
		deskFree := NewDeskFree(deskData)
		deskPid := deskFree.newDesk()
		//添加新桌子
		a.desks[deskData.Rid] = deskPid
		//响应
		rsp.Desk = deskPid
		//启动
		deskPid.Tell(new(pb.ServeStart))
		//添加桌子
		msg2 := new(pb.AddDesk)
		msg2.Desk = deskPid
		msg2.Roomid = deskData.Rid
		msg2.Rtype = deskData.Rtype
		a.hallPid.Request(msg2, ctx.Self())
	default:
		glog.Errorf("deskData.Rtype error %d", deskData.Rtype)
	}
	ctx.Respond(rsp)
}

////更新货币
//func (a *DeskActor) changeCurrency(arg *pb.ChangeCurrency) {
//	if v, ok := a.roles[arg.Userid]; ok {
//		v.Tell(arg)
//	}
//}
