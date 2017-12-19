package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *HallActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connect:
		//初始化建立连接
		arg := msg.(*pb.Connect)
		a.serve[arg.Name] = ctx.Sender()
		//响应
		connected := &pb.Connected{
			Name: a.Name,
		}
		ctx.Respond(connected)
		glog.Debugf("connected name %s", arg.Name)
		glog.Debugf("serve len %d", len(a.serve))
	case *pb.Disconnect:
		arg := msg.(*pb.Disconnect)
		delete(a.serve, arg.Name)
		//响应
		disconnected := &pb.Disconnected{
			Name: a.Name,
		}
		ctx.Respond(disconnected)
		glog.Debugf("disconnected name %s", arg.Name)
		glog.Debugf("serve len %d", len(a.serve))
	case *pb.LoginHall:
		arg := msg.(*pb.LoginHall)
		userid := arg.GetUserid()
		name := arg.GetNodeName()
		//断开旧连接
		k := a.roles[userid]
		if p, ok := a.serve[k]; ok {
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
		rsp.Message = ctx.Self().String()
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
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		glog.Debugf("WxpayCallback: %v", arg)
		//支付回调
		for k, v := range a.serve {
			//TODO 优化
			s := utils.Split(k, "game.node")
			if len(s) == 2 && s[0] == "" && s[1] != "" {
				v.Tell(arg)
				//TODO 优化
				//选择一个验证即可,
				//暂时不知道哪个节点的订单
				break
			}
		}
	case *pb.WxpayGoods:
		arg := msg.(*pb.WxpayGoods)
		glog.Debugf("WxpayGoods: %v", arg)
		userid := arg.Userid
		gate := a.roles[userid]
		if v, ok := a.serve[gate]; ok {
			v.Tell(arg)
		} else {
			glog.Errorf("WxpayGoods: %v", arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *HallActor) start(ctx actor.Context) {
	glog.Infof("hall start: %v", ctx.Self().String())
	//启动
	go a.ticker(ctx)
}

//时钟
func (a *HallActor) ticker(ctx actor.Context) {
	tick := time.Tick(30 * time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("hall ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("hall ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (a *HallActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	//TODO
}

//关闭时钟
func (a *HallActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

func (a *HallActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	//msg := new(pb.ServeClose)
	//for k, v := range a.roles {
	//	glog.Debugf("Stop role: %s", k)
	//	v.Tell(msg)
	//}
	dbmsName := cfg.Section("dbms").Name()
	roomName := cfg.Section("room").Name()
	roleName := cfg.Section("role").Name()
	loginName := cfg.Section("login").Name()
	mailName := cfg.Section("mail").Name()
	betsName := cfg.Section("bets").Name()
	if v, ok := a.serve[loginName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[roleName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[roomName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[mailName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[betsName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[dbmsName]; ok {
		v.Stop()
	}
	//for k, v := range a.serve {
	//	v.Stop()
	//}
}
