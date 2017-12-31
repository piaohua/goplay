package main

import (
	"time"

	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *RoomActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.CreateDesk:
		arg := msg.(*pb.CreateDesk)
		//按规则创建房间
		a.create(arg, ctx)
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
		a.count[arg.Roomid] += 1
		//响应
		//rsp := new(pb.JoinedDesk)
		//ctx.Respond(rsp)
	case *pb.CloseDesk:
		arg := msg.(*pb.CloseDesk)
		glog.Debugf("CloseDesk %#v", arg)
		//TODO 私人房间
		//移除
		delete(a.count, arg.Roomid)
		delete(a.codes, arg.Code)
		delete(a.rooms, arg.Roomid)
		//响应
		//rsp := new(pb.ClosedDesk)
		//ctx.Respond(rsp)
	case *pb.LeaveDesk:
		arg := msg.(*pb.LeaveDesk)
		glog.Debugf("LeaveDesk %#v", arg)
		//移除
		delete(a.router, arg.Userid)
		if n, ok := a.count[arg.Roomid]; ok && n > 0 {
			a.count[arg.Roomid] = n - 1
		}
		//响应
		//rsp := new(pb.LeftDesk)
		//ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *RoomActor) start(ctx actor.Context) {
	glog.Infof("room start: %v", ctx.Self().String())
	//初始化建立连接
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	//timeout := 3 * time.Second
	//hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
	//if err != nil {
	//	glog.Fatalf("remote hall err %v", err)
	//}
	//a.hallPid = hallPid.Pid
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *RoomActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("room ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("room ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *RoomActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *RoomActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *RoomActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, _ := range a.rooms {
		glog.Debugf("Stop room: %s", k)
		//TODO
		//v.Save()
	}
}

//创建房间
func (a *RoomActor) create(arg *pb.CreateDesk, ctx actor.Context) {
	glog.Debugf("CreateDesk %#v", arg)
	rsp := new(pb.CreatedDesk)
	deskData := handler.Data2Desk(arg.Data)
	if deskData == nil {
		rsp.Error = pb.RoomNotExist
		ctx.Respond(rsp)
	}
	//创建
	a.createRoom(deskData)
	//响应登录
	rsp.Data = handler.Desk2Data(deskData)
	if rsp.Data == "" {
		rsp.Error = pb.RoomNotExist
		ctx.Respond(rsp)
	}
	//添加房间,TODO 如果响应超时？
	a.rooms[deskData.Rid] = deskData
	a.codes[deskData.Code] = deskData.Rid
	ctx.Respond(rsp)
}

//生成一个牌桌邀请码,全列表中唯一
func (a *RoomActor) GenCode() (s string) {
	s = utils.RandStr(6)
	//是否已经存在
	if _, ok := a.codes[s]; ok {
		return a.GenCode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (a *RoomActor) GenCodeFree() (s string) {
	s = utils.RandStr(7) //区别于私人房间
	//是否已经存在
	if _, ok := a.codes[s]; ok {
		return a.GenCode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//匹配房间
func (a *RoomActor) matchDesk(arg *pb.MatchDesk, ctx actor.Context) {
	rsp := new(pb.MatchedDesk)
	msg1 := new(pb.MatchDesk)
	//TODO 匹配成功后查找
	//msg1.Roomid = roomid
	//匹配失败后新建
	deskData := handler.NewDeskData(arg.Rtype)
	if deskData == nil {
		glog.Errorf("matchRoom err rtype: %d", arg.Rtype)
		ctx.Respond(rsp)
		return
	}
	//创建
	a.createRoom(deskData)
	msg1.Name = arg.Name
	msg1.Data = handler.Desk2Data(deskData)
	if msg1.Data == "" {
		glog.Errorf("matchRoom err rtype: %d", arg.Rtype)
		ctx.Respond(rsp)
		return
	}
	timeout := 3 * time.Second
	res1, err1 := a.hallPid.RequestFuture(msg1, timeout).Result()
	if err1 != nil {
		glog.Errorf("matchRoom err: %v", err1)
		//响应
		ctx.Respond(rsp)
		return
	}
	response1 := res1.(*pb.MatchedDesk)
	glog.Debugf("response1: %#v", response1)
	rsp.Desk = response1.Desk
	//响应
	ctx.Respond(rsp)
}

//创建房间
func (a *RoomActor) createRoom(deskData *data.DeskData) {
	switch deskData.Rtype {
	case data.ROOM_FREE:
		//百人
		deskData.Code = a.GenCodeFree()
	case data.ROOM_PRIVATE:
		//私人
		deskData.Code = a.GenCode()
	}
	deskData.Rid = a.uniqueid.GenID()
	deskData.Ctime = uint32(utils.Timestamp())
}
