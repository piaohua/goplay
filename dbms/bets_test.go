package main

import (
	"fmt"
	"goplay/pb"
	"testing"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func TestBets(t *testing.T) {
	remote.Start("127.0.0.1:0")
	timeout := 1 * time.Second
	pid, err := remote.SpawnNamed("127.0.0.1:7001", "remote1", "hello", timeout)
	if err != nil {
		//fmt.Println(err)
		//return
	}
	res, _ := pid.RequestFuture(new(pb.Request), timeout).Result()
	fmt.Println("res ", res)
	response := res.(*pb.Response)
	fmt.Println(response)
	pid.Stop()
	console.ReadLine()
}
