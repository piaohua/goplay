package main

import (
	"log"
	"testing"
	"time"

	"goplay/pb"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func TestRun(t *testing.T) {
	//召唤机器人
	code := "199280"
	go callRobot(code)
}

//创建远程连接消息通道
func callRobot(code string) {
	message := &pb.RobotMsg{
		Code: code,
		Num:  3,
	}
	remote := actor.NewPID("127.0.0.1:7070", "RobotMsg")
	remote.Tell(msg)
	remote.Stop(msg)
}

//创建远程连接消息通道
func newMyMessageSenderChannel() chan<- *pb.RobotMsg {
	channel := make(chan *pb.RobotMsg, 10)
	remote := actor.NewPID("127.0.0.1:7070", "RobotMsg")
	go func() {
		for msg := range channel {
			remote.Tell(msg)
		}
	}()

	return channel
}

//向远程连接通道写入消息
func remote_send() {
	remote.Start("127.0.0.1:7071")
	channel := newMyMessageSenderChannel()

	for i := 0; i < 1; i++ {
		message := &pb.RobotMsg{
			Code: "create",
			Num:  1,
		}
		channel <- message
	}
}

//远程连接客户端
func TestRobot(t *testing.T) {
	remote_send()
	//go remote_note1()
	//go remote_note2()
	<-time.After(time.Duration(6) * time.Second)
	console.ReadLine()
}

func remote_note2() {
	remote.Start("127.0.0.1:8088")
	channel := newMyMessageSenderChannel2()

	for i := 0; i < 3; i++ {
		message := &pb.RobotMsg{
			Code: "199999",
			Num:  3,
		}
		channel <- message
	}

	console.ReadLine()
}

func remote_note1() {
	remote.Start("127.0.0.1:8080")
	//create the channel
	channel := make(chan *pb.RobotMsg)

	//create an actor receiving messages and pushing them onto the channel
	props := actor.FromFunc(func(context actor.Context) {
		if msg, ok := context.Message().(*pb.RobotMsg); ok {
			channel <- msg
		}
	})
	actor.SpawnNamed(props, "RobotMsg")

	//consume the channel just like you use to
	go func() {
		for msg := range channel {
			log.Println("node2 msg -> ", msg)
		}
	}()

	console.ReadLine()
}

func newMyMessageSenderChannel2() chan<- *pb.RobotMsg {
	channel := make(chan *pb.RobotMsg)
	remote := actor.NewPID("127.0.0.1:8080", "RobotMsg")
	go func() {
		for msg := range channel {
			remote.Tell(msg)
		}
	}()

	return channel
}
