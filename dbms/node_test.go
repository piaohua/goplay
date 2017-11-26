package main

import (
	"fmt"
	"goplay/pb"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func TestNode(t *testing.T) {
	remote.Start("127.0.0.1:0")
	activate("hello")
	//<-time.After(time.Minute)
	//activate("hello")
	server := actor.NewPID("127.0.0.1:7001", "hello")
	fmt.Println("server ", server)
	//发送Proto message
	server.Tell(new(pb.Request))
	console.ReadLine()
}

func activate(name string) {
	timeout := 1 * time.Second
	pid, err := remote.SpawnNamed("127.0.0.1:7001", "remote1", name, timeout)
	if err != nil {
		//fmt.Println(err)
		//return
	}
	res, _ := pid.RequestFuture(new(pb.Request), timeout).Result()
	fmt.Println("res ", res)
	response := res.(*pb.Response)
	fmt.Println(response)
	//pid.Stop()
	//
	//pid, _ = remote.SpawnNamed("127.0.0.1:8080", "remote2", name, timeout)
	res, _ = pid.RequestFuture(new(pb.Request), timeout).Result()
	response = res.(*pb.Response)
	fmt.Println(response)
}
