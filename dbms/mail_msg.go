package main

import (
	"time"

	"goplay/data"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (a *MailActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.HallConnect:
		//初始化建立连接
		glog.Infof("mail init: %v", ctx.Self().String())
		//连接
		bind := cfg.Section("hall").Key("bind").Value()
		name := cfg.Section("cookie").Key("name").Value()
		timeout := 3 * time.Second
		hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
		glog.Infof("a.hallPid: %s", a.hallPid.String())
		if err != nil {
			glog.Fatalf("remote hall err %v", err)
		}
		a.hallPid = hallPid.Pid
		connect := &pb.HallConnect{
			Sender: ctx.Self(),
			Name:   a.Name,
		}
		a.hallPid.Tell(connect)
	case *pb.ServeStop:
		//关闭服务
		a.HandlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.CMailList:
		arg := msg.(*pb.CMailList)
		userid := a.getUserid(ctx)
		//响应
		rsp := handler.GetMailList(arg, userid)
		ctx.Respond(rsp)
	case *pb.CDeleteMail:
		arg := msg.(*pb.CDeleteMail)
		userid := a.getUserid(ctx)
		//响应
		rsp := handler.DelMail(arg, userid)
		ctx.Respond(rsp)
	case *pb.CGetMailItem:
		arg := msg.(*pb.CGetMailItem)
		userid := a.getUserid(ctx)
		//响应
		rsp, list := handler.GetMailItem(arg, userid)
		for _, v := range list {
			a.addCurrency(userid, int(v.Rtype),
				data.LogType22, int32(v.Number), ctx)
		}
		ctx.Respond(rsp)
	case *pb.NewMail:
		arg := msg.(*pb.NewMail)
		id := a.uniqueid.GenID()
		glog.Debugf("new mail %#v", arg)
		glog.Debugf("new mail From %s", arg.From)
		glog.Debugf("new mail To %s", arg.To)
		glog.Debugf("new mail Content %s", arg.Content)
		glog.Debugf("new mail id %s", id)
		//test
		userid := a.getUserid(ctx)
		glog.Debugf("new mail userid %s", userid)
		handler.NewMail(userid, id)
		//
		msg := new(pb.SMailNotice)
		ctx.Respond(msg)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *MailActor) HandlerStop(ctx actor.Context) {
	glog.Debugf("HandlerStop: %s", a.Name)
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	for k, _ := range a.mails {
		glog.Debugf("Stop Mail: %s", k)
		//TODO
		//v.Save()
	}
}

func (a *MailActor) getUserid(ctx actor.Context) (userid string) {
	timeout := 1 * time.Second
	req := new(pb.GetUserid)
	req.Sender = ctx.Sender()
	res, _ := rolePid.RequestFuture(req, timeout).Result()
	res1 := res.(*pb.GotUserid)
	if res1 != nil {
		return
	}
	userid = res1.Userid
	return
}

func (a *MailActor) addCurrency(userid string,
	rtype, ltype int, amount int32, ctx actor.Context) {
	switch uint32(rtype) {
	case data.DIAMOND:
		a.sendCurrency(userid, amount, 0, ltype, ctx)
	case data.COIN:
		a.sendCurrency(userid, 0, amount, ltype, ctx)
	}
}

func (a *MailActor) sendCurrency(userid string,
	diamond, coin int32, ltype int, ctx actor.Context) {
	if userid == "" {
		glog.Errorf("sendCurrency user empty: %d", ltype)
	}
	msg3 := new(pb.ChangeCurrency)
	msg3.Userid = userid
	msg3.Type = int32(ltype)
	msg3.Coin = coin
	msg3.Diamond = diamond
	ctx.Sender().Tell(msg3)
}
