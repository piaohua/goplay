package main

import (
	"fmt"
	"strings"

	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//web请求处理
func (a *HallActor) HandlerWeb(arg *pb.WebRequest,
	rsp *pb.WebResponse, ctx actor.Context) {
	switch arg.Code {
	case pb.WebOnline:
		msg1 := make([]string, 0)
		err1 := json.Unmarshal(arg.Data, &msg1)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//响应
		resp := make(map[string]int)
		for _, v := range msg1 {
			if _, ok := a.roles[v]; ok {
				resp[v] = 1
			} else {
				resp[v] = 0
			}
		}
		result, err2 := json.Marshal(resp)
		if err2 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err2)
			return
		}
		rsp.Result = result
	case pb.WebShop:
		//更新配置
		err1 := handler.UpdateSyncConfig(pb.CONFIG_SHOP, arg.Data)
		if err1 != nil {
			rsp.ErrMsg = fmt.Sprintf("msg err: %v", err1)
			return
		}
		//广播所有节点,主动通知同步配置,TODO 只同步修改数据
		msg2 := handler.GetSyncConfig(pb.CONFIG_SHOP)
		a.broadcast(msg2, ctx)
	case pb.WebSetEnv:
	default:
		glog.Errorf("unknown message %v", arg)
	}
}

//广播所有节点,逻辑服,db
func (a *HallActor) broadcast(msg interface{}, ctx actor.Context) {
	//name := cfg.Section("game.free").Name()
	for k, v := range a.serve {
		if strings.Contains(k, "dbms.") ||
			strings.Contains(k, "gate.") ||
			strings.Contains(k, "game.") {
			v.Tell(msg)
		}
	}
}

func (a *HallActor) handling() {
}

//name := cfg.Section("dbms").Name()
//if v, ok := a.serve[name]; ok {
//	timeout := 3 * time.Second
//	res2, err2 := v.RequestFuture(arg, timeout).Result()
//	if err2 != nil {
//		glog.Errorf("WebRequest err: %v", err2)
//		//TODO ctx.Respond(rsp3)
//		return
//	}
//	ctx.Respond(rsp2)
//}
////TODO ctx.Respond(rsp3)
