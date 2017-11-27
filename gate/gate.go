package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"
	"utils"

	"goplay/glog"

	ini "gopkg.in/ini.v1"
)

var (
	cfg *ini.File
	sec *ini.Section
	err error

	aesEnc *utils.AesEncrypt
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
	//初始化
	aesInit()
	//启动服务
	bind := cfg.Section("gate.node1").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	NewRemote(bind, name)
	//wsServer
	addr := cfg.Section("gate.node1").Key("addr").Value()
	wsServer := new(WSServer)
	wsServer.Addr = addr
	if wsServer != nil {
		wsServer.Start()
	}
	signalListen() //监听关闭信号
	//关闭服务
	//关闭websocket连接, 先关监听
	if wsServer != nil {
		wsServer.Close()
	}
	//关闭服务
	Stop()
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

//加密初始化
func aesInit() {
	aesEnc = new(utils.AesEncrypt)
	key := cfg.Section("gate").Key("key").Value()
	aesEnc.SetKey([]byte(key))
}

//加密
func aesEn(doc []byte) (arrEncrypt []byte) {
	arrEncrypt, err = aesEnc.Encrypt(doc)
	if err != nil {
		glog.Errorf("arrEncrypt: %s", string(doc))
	}
	return
}

//解密
func aesDe(arrEncrypt []byte) (bMsg []byte) {
	bMsg, err = aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		glog.Errorf("arrEncrypt: %s", string(arrEncrypt))
	}
	return
}
