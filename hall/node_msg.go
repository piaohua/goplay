package main

import (
	"goplay/glog"
	"goplay/pb"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
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
	case *pb.GateDisconnect:
		arg := msg.(*pb.GateDisconnect)
		connected := &pb.GateDisconnected{
			Message: ctx.Self().String(),
		}
		arg.Sender.Tell(connected)
		glog.Infof("GateDisconnect %s", arg.Sender.String())
		//网关注销
		delete(a.gates, arg.Sender.String())
	case *pb.HallConnect:
		//初始化建立连接
		arg := msg.(*pb.HallConnect)
		a.serve[arg.Name] = arg.Sender
		//connected := &pb.HallConnected{
		//	Message: ctx.Self().String(),
		//	Name:    arg.Name,
		//}
		//ctx.Respond(rsp)
		glog.Debugf("name %s", arg.Name)
		glog.Debugf("serve len %d", len(a.serve))
	case *pb.LoginHall:
		arg := msg.(*pb.LoginHall)
		userid := arg.GetUserid()
		name := arg.GetNodeName()
		//断开旧连接
		k := a.roles[userid]
		if p, ok := a.gates[k]; ok {
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
		a.HandlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.WxpayCallback:
		arg := msg.(*pb.WxpayCallback)
		glog.Debugf("WxpayCallback: %v", arg)
		//支付回调
		for _, v := range a.gates {
			v.Tell(arg)
			//TODO 优化
			//选择一个验证即可,
			//暂时不知道哪个节点的订单
			break
		}
	case *pb.WxpayGoods:
		arg := msg.(*pb.WxpayGoods)
		glog.Debugf("WxpayGoods: %v", arg)
		userid := arg.Userid
		gate := a.roles[userid]
		if v, ok := a.gates[gate]; ok {
			v.Tell(arg)
		} else {
			glog.Errorf("WxpayGoods: %v", arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

func (a *HallActor) HandlerStop(ctx actor.Context) {
	glog.Debugf("HandlerStop: %s", a.Name)
	//回存数据
	//msg := new(pb.ServeStop)
	//for k, v := range a.roles {
	//	glog.Debugf("Stop role: %s", k)
	//	v.Tell(msg)
	//}
	dbmsName := cfg.Section("dbms").Name()
	roomName := cfg.Section("room").Name()
	roleName := cfg.Section("role").Name()
	loginName := cfg.Section("login").Name()
	mailName := cfg.Section("mail").Name()
	bettingName := cfg.Section("betting").Name()
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
	if v, ok := a.serve[bettingName]; ok {
		v.Stop()
	}
	if v, ok := a.serve[dbmsName]; ok {
		v.Stop()
	}
	//for k, v := range a.serve {
	//	v.Stop()
	//}
}
