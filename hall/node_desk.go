package main

import (
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
		//匹配房间
		if v, ok := a.serve[arg.Name]; ok {
			//TODO
			v.Tell(arg)
		}
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
