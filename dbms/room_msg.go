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
		a.create(arg, ctx)
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
	switch deskData.Rtype {
	case data.ROOM_FREE:
		//百人
		deskData.Rid = a.uniqueid.GenID()
		deskData.Code = a.GenCodeFree()
	case data.ROOM_PRIVATE:
		//私人
		deskData.Rid = a.uniqueid.GenID()
		deskData.Code = a.GenCode()
	}
	//添加房间,TODO 如果响应超时？
	a.rooms[deskData.Rid] = deskData
	a.codes[deskData.Code] = deskData.Rid
	//响应登录
	rsp.Data = handler.DeskData2(deskData)
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
