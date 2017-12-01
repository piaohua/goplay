package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"goplay/data"
	"goplay/game/config"
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
	//数据库连接
	host := cfg.Section("mongod").Key("host").Value()
	port := cfg.Section("mongod").Key("port").Value()
	user := cfg.Section("mongod").Key("user").Value()
	passwd := cfg.Section("mongod").Key("passwd").Value()
	dbname := cfg.Section("mongod").Key("name").Value()
	data.InitMgo(host, port, user, passwd, dbname)
	//配置初始化
	appid := cfg.Section("weixin").Key("appid").Value()
	appsecret := cfg.Section("weixin").Key("appsecret").Value()
	appkey := cfg.Section("weixin").Key("appkey").Value()
	mchid := cfg.Section("weixin").Key("mchid").Value()
	pattern := cfg.Section("weixin").Key("notifyPattern").Value()
	notifyUrl := cfg.Section("weixin").Key("notifyUrl").Value()
	config.ConfigInit(appid, appsecret, appkey, mchid, pattern, notifyUrl)
	//启动服务
	bind := cfg.Section("dbms").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	room := cfg.Section("cookie").Key("room").Value()
	role := cfg.Section("cookie").Key("role").Value()
	mail := cfg.Section("cookie").Key("mail").Value()
	bets := cfg.Section("cookie").Key("bets").Value()
	NewRemote(bind, name, room, role, mail, bets)
	signalListen()
	//关闭服务
	Stop()
	data.Close() //数据库断开
	//延迟等待
	<-time.After(3 * time.Second) //延迟关闭
}

func signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c)
	//signal.Stop(c)
	for {
		s := <-c
		glog.Error("get signal:", s)
		return
	}
}
