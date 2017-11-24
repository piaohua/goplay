package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"goplay/glog"

	ini "gopkg.in/ini.v1"
)

var (
	cfg *ini.File
	sec *ini.Section
	err error
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	defer glog.Flush()
	//日志定义
	glog.Init()
	//加载配置
	cfg, err = ini.Load("conf.ini")
	if err != nil {
		panic(err)
	}
	cfg.BlockMode = false //只读
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	NewRemote(bind, name)
	signalListen()
	//TODO 关闭处理
	//延迟等待
	<-time.After(5 * time.Second) //延迟关闭
}

func signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c)
	//signal.Stop(c)
	for {
		s := <-c
		glog.Error("get signal:", s)
	}
}
