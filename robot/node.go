/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:12:54
 * Filename      : node.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"goplay/glog"
	"goplay/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func (server *RobotServer) NewRemote(bind, name string) {
	if bind == "" {
		glog.Panic("bind empty")
	}
	// Start the remote server
	remote.Start(bind)
	server.remoteRecv(name) //接收远程消息
}

//接收远程消息
func (server *RobotServer) remoteRecv(name string) {
	//create the channel
	server.channel = make(chan *pb.RobotMsg, 100) //protos中定义

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*pb.RobotMsg); ok {
			server.channel <- msg
		}
	})
	actor.SpawnNamed(props, name)

	//consume the channel just like you use to
	go func() {
		for msg := range server.channel {
			//分配机器人
			for msg.Num > 0 {
				go func(code, phone string, rtype uint32) {
					server.RunRobot(code, phone, rtype, false)
				}(msg.Code, server.phone, msg.Rtype)
				//TODO msg.Rtype
				server.phone = utils.StringAdd(server.phone)
				msg.Num--
			}
			glog.Debugf("node msg -> %v", msg)
		}
	}()
}
