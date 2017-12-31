package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//玩家桌子常用共有操作请求处理
func (ws *WSConn) HandlerDesk(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//TODO
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//离开房间
		ws.leaveRoom(arg, ctx)
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	case *pb.CChatText:
		arg := msg.(*pb.CChatText)
		glog.Debugf("CChatText %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SChatText)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CChatVoice:
		arg := msg.(*pb.CChatVoice)
		glog.Debugf("CChatVoice %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SChatVoice)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CLeave:
		arg := msg.(*pb.CLeave)
		glog.Debugf("CLeave %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SLeave)
			rsp.Error = pb.NotInRoomCannotLeave
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CKick:
		arg := msg.(*pb.CKick)
		glog.Debugf("CKick %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SKick)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CReady:
		arg := msg.(*pb.CReady)
		glog.Debugf("CReady %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SReady)
			rsp.Error = pb.NotInRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CLaunchVote:
		arg := msg.(*pb.CLaunchVote)
		glog.Debugf("CLaunchVote %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SLaunchVote)
			rsp.Error = pb.NotInPrivateRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CVote:
		arg := msg.(*pb.CVote)
		glog.Debugf("CVote %#v", arg)
		if ws.gamePid == nil {
			rsp := new(pb.SVote)
			rsp.Error = pb.NotInPrivateRoom
			ws.Send(rsp)
			return
		}
		ws.gamePid.Request(arg, ctx.Self())
	case *pb.CRoomList:
		arg := msg.(*pb.CRoomList)
		glog.Debugf("CRoomList %#v", arg)
	case *pb.SetRecord:
		arg := msg.(*pb.SetRecord)
		glog.Debugf("SetRecord %#v", arg)
		handler.SetRecord(ws.User, int(arg.Rtype))
	default:
		//glog.Errorf("unknown message %v", msg)
		ws.HandlerFree(msg, ctx)
	}
}

//离开房间
func (ws *WSConn) leaveRoom(arg *pb.LeaveDesk, ctx actor.Context) {
	ws.gamePid = nil
	ws.roomPid.Request(arg, ctx.Self())
	ws.hallPid.Request(arg, ctx.Self())
}

//进入房间
func (ws *WSConn) entryRoom(ctx actor.Context) {
	if ws.gamePid == nil {
		glog.Errorf("not in the room: %s", ws.User.GetUserid())
		return
	}
	result4, err4 := json.Marshal(ws.User)
	if err4 != nil {
		glog.Errorf("user Marshal err %v", err4)
		return
	}
	msg4 := new(pb.EnterDesk)
	msg4.Data = string(result4)
	//进入房间
	timeout := 3 * time.Second
	res1, err1 := ws.gamePid.RequestFuture(msg4, timeout).Result()
	if err1 != nil {
		glog.Errorf("entryRoom err: %v", err1)
		return
	}
	response1 := res1.(*pb.EnteredDesk)
	glog.Debugf("response1: %#v", response1)
}

//大厅中匹配可用房间
func (ws *WSConn) matchRoom(rtype uint32) *pb.MatchedDesk {
	//匹配可以进入的房间
	msg1 := new(pb.MatchDesk)
	switch rtype {
	case data.ROOM_FREE:
		msg1.Rtype = data.ROOM_FREE
		//节点注册名称,TODO 多节点处理
		msg1.Name = cfg.Section("game.free").Name()
	}
	timeout := 3 * time.Second
	res1, err1 := ws.hallPid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("matchRoom err: %v", err1)
		return nil
	}
	response1 := res1.(*pb.MatchedDesk)
	glog.Debugf("response1: %#v", response1)
	return response1
}

//数据中心创建房间,TODO hall 中创建添加
func (ws *WSConn) createRoom(rtype uint32) *pb.CreatedDesk {
	msg2 := new(pb.CreateDesk)
	switch rtype {
	case data.ROOM_FREE:
		msg2.Data = handler.FreeData()
	}
	if msg2.Data == "" {
		return nil
	}
	timeout := 3 * time.Second
	res2, err2 := ws.roomPid.RequestFuture(msg2, timeout).Result()
	if err2 != nil {
		glog.Errorf("createRoom err: %v", err2)
		return nil
	}
	response2 := res2.(*pb.CreatedDesk)
	glog.Debugf("response2: %#v", response2)
	if response2.Error != pb.OK {
		glog.Error("createRoom failed")
		return nil
	}
	return response2
}

//创建新桌子
func (ws *WSConn) spawnRoom(deskNode *actor.PID, rdata string) *pb.SpawnedDesk {
	if rdata == "" || deskNode == nil {
		return nil
	}
	msg2 := new(pb.SpawnDesk)
	msg2.Data = rdata
	timeout := 3 * time.Second
	res2, err2 := deskNode.RequestFuture(msg2, timeout).Result()
	if err2 != nil {
		glog.Errorf("spawnRoom err: %v", err2)
		return nil
	}
	response2 := res2.(*pb.SpawnedDesk)
	glog.Debugf("response2: %#v", response2)
	if response2.Desk == nil {
		return nil
	}
	if response2.Error != pb.OK {
		return nil
	}
	return response2
}
