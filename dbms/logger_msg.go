package main

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *LoggerActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.LogRegist:
		arg := msg.(*pb.LogRegist)
		data.RegistRecord(arg.Userid, arg.Nickname, arg.Ip, arg.Atype)
	case *pb.LogLogin:
		arg := msg.(*pb.LogLogin)
		data.LoginRecord(arg.Userid, arg.Ip, arg.Atype)
	case *pb.LogLogout:
		arg := msg.(*pb.LogLogout)
		data.LogoutRecord(arg.Userid, int(arg.Event))
	case *pb.LogDiamond:
		arg := msg.(*pb.LogDiamond)
		data.DiamondRecord(arg.Userid, int(arg.Type), arg.Rest, arg.Num)
	case *pb.LogCoin:
		arg := msg.(*pb.LogCoin)
		data.CoinRecord(arg.Userid, int(arg.Type), arg.Rest, arg.Num)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}
