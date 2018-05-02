/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:12:54
 * Filename      : node.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"errors"

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
	server.channel = make(chan interface{}, 100) //protos中定义
	server.closeCh = make(chan struct{})

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		server.remoteSend(context.Message())
	})
	nodePid, err = actor.SpawnNamed(props, name)
	server.remoteHall()

	//consume the channel just like you use to
	go func() {
		for msg := range server.channel {
			err := server.remoteHandler(msg)
			if err != nil {
				//停止发送消息
				close(server.closeCh)
				break
			}
		}
	}()
}

func (server *RobotServer) remoteSend(message interface{}) {
	if server.channel == nil {
		glog.Errorf("server channel closed %#v", message)
		return
	}
	if len(server.channel) == cap(server.channel) {
		glog.Errorf("send msg channel full -> %d", len(server.channel))
		return
	}
	select {
	case <-server.closeCh:
		return
	default:
	}
	select {
	case <-server.closeCh:
		return
	case server.channel <- message:
	}
}

//处理
func (server *RobotServer) remoteHandler(message interface{}) error {
	switch message.(type) {
	case *pb.RobotMsg:
		msg := message.(*pb.RobotMsg)
		//分配机器人
		for msg.Num > 0 {
			go func(code, phone string, rtype uint32) {
				server.RunRobot(code, phone, rtype, false)
			}(msg.Code, server.phone, msg.Rtype)
			//TODO msg.Rtype
			server.phone = utils.StringAdd(server.phone)
			msg.Num--
		}
		glog.Debugf("node msg -> %#v", msg)
	case closeFlag:
		return errors.New("msg channel closed")
	}
	return nil
}

func (server *RobotServer) remoteHall() {
	//hall
	name := cfg.Section("cookie").Key("name").Value()
	bind := cfg.Section("hall").Key("bind").Value()
	hallPid = actor.NewPID(bind, name)
	//name
	server.Name = cfg.Section("robot").Name()
	connect := &pb.Connect{
		Name: server.Name,
	}
	hallPid.Request(connect, nodePid)
}
