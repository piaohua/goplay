package main

import (
	"encoding/json"
	"time"

	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (a *DBMSActor) Handler(msg interface{}, ctx actor.Context) {
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
		//同步配置
		a.syncConfig(arg.Sender.String())
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
		glog.Infof("dbms init: %v", ctx.Self().String())
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
	default:
		if a.logger == nil {
			glog.Errorf("unknown message %v", msg)
		} else {
			a.logger.Tell(msg)
		}
	}
}

func (a *DBMSActor) HandlerStop(ctx actor.Context) {
	glog.Debugf("HandlerStop: %s", a.Name)
	//回存数据
	for k, _ := range a.gates {
		glog.Debugf("Stop gate: %s", k)
	}
	if a.logger != nil {
		a.logger.Stop()
	}
}

//打包配置
func (a *DBMSActor) syncConfigMsg(ctype pb.ConfigType,
	d interface{}) interface{} {
	msg := new(pb.SyncConfig)
	msg.Type = ctype
	result, err := json.Marshal(d)
	if err != nil {
		glog.Errorf("syncConfig Marshal err %v, data %#v", err, d)
	}
	msg.Data = string(result)
	return msg
}

//同步配置
func (a *DBMSActor) syncConfig(key string) {
	if _, ok := a.gates[key]; !ok {
		glog.Errorf("gate not exists: %s", key)
		return
	}
	pid := a.gates[key]
	msg1 := a.syncConfigMsg(pb.CONFIG_BOX, config.GetBoxs())
	pid.Tell(msg1)
	msg2 := a.syncConfigMsg(pb.CONFIG_ENV, config.GetEnvs())
	pid.Tell(msg2)
	msg3 := a.syncConfigMsg(pb.CONFIG_LOTTERY, config.GetLotterys())
	pid.Tell(msg3)
	msg4 := a.syncConfigMsg(pb.CONFIG_NOTICE,
		config.GetNotices(data.NOTICE_TYPE1))
	pid.Tell(msg4)
	msg5 := a.syncConfigMsg(pb.CONFIG_PRIZE, config.GetLotterys())
	pid.Tell(msg5)
	msg6 := a.syncConfigMsg(pb.CONFIG_SHOP, config.GetShops2())
	pid.Tell(msg6)
	msg7 := a.syncConfigMsg(pb.CONFIG_VIP, config.GetVips())
	pid.Tell(msg7)
	msg8 := a.syncConfigMsg(pb.CONFIG_CLASSIC, config.GetClassics())
	pid.Tell(msg8)
}
