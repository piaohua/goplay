package main

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"goplay/glog"
	"utils"

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
	bind := cfg.Section("login").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	NewRemote(bind, name)
	//监听地址
	addr := cfg.Section("login").Key("addr").Value()
	//启动监听
	go Start(addr)
	signalListen() //监听关闭信号
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
	key := cfg.Section("login").Key("key").Value()
	aesEnc.SetKey([]byte(key))
}

//加密
func aesEn(doc string) (arrEncrypt []byte) {
	arrEncrypt, err = aesEnc.Encrypt([]byte(doc))
	if err != nil {
		glog.Errorf("arrEncrypt: %s", doc)
	}
	return
}

//解密
func aesDe(arrEncrypt []byte) (strMsg string) {
	bMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		glog.Errorf("arrEncrypt: %s", string(arrEncrypt))
	}
	strMsg = string(bMsg)
	return
}
