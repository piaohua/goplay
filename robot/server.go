/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:11:57
 * Filename      : server.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/gorilla/websocket"
)

type RobotServer struct {
	PendingWriteNum int    //等待写入消息长度
	MaxMsgLen       uint32 //最大消息长度

	conns      WebsocketConnSet //连接集合
	mutexConns sync.Mutex       //互斥锁
	wg         sync.WaitGroup   //同步机制

	channel chan *pb.RobotMsg //消息通道

	phone   string           //注册登录账号
	online  map[string]bool  //map[phone]状态,true=在线,
	offline map[string]bool  //map[phone]状态,true=离线,false=登录中
	msgCh   chan interface{} //消息通道
}

func (server *RobotServer) Start() {
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		glog.Infof("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 1024
		glog.Infof("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}

	server.conns = make(WebsocketConnSet)

	//初始化
	server.online = make(map[string]bool)
	server.offline = make(map[string]bool)
	server.msgCh = make(chan interface{}, 100)
	//启动管理服务
	go server.run()
	//启动测试
	go server.runTest1()
	//go server.runTest()
	//go server.runFree()
}

//关闭连接
func (server *RobotServer) Close() {
	close(server.channel)
	close(server.msgCh) //关闭

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()

	server.wg.Wait()
}

//启动一个机器人
func (server *RobotServer) RunRobot(code, phone string, rtype uint32, regist bool) {
	host := getHost()
	u := url.URL{Scheme: "ws", Host: host, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Errorf("robot run dial -> %v", err)
		return
	}

	server.wg.Add(1)
	defer server.wg.Done()

	server.mutexConns.Lock()
	if server.conns == nil {
		server.mutexConns.Unlock()
		conn.Close()
		return
	}
	server.conns[conn] = struct{}{}
	server.mutexConns.Unlock()

	//new robot
	robot := newRobot(conn, server.PendingWriteNum, server.MaxMsgLen)
	robot.code = code //设置邀请码
	robot.rtype = rtype
	robot.data.Phone = phone
	//robot.data.Nickname = RandNickName()
	glog.Infof("run robot -> %s", phone)
	glog.Infof("run robot -> code:%s, rtype:%d, regist:%v", code, rtype, regist)
	go robot.writePump()
	if regist {
		go robot.SendRegist() //发起请求,注册-登录-进入房间
	} else {
		go robot.SendLogin() //登录
	}
	robot.readPump()

	// cleanup
	robot.Close()
	server.mutexConns.Lock()
	delete(server.conns, conn)
	server.mutexConns.Unlock()
}

//获取网关
func getHost() (host string) {
	addr := cfg.Section("login").Key("addr").Value()
	addr = "http://" + addr + "/gate"
	b, err := doHttpPost(addr, []byte{})
	if err != nil {
		glog.Errorf("getHost err: %v", err)
		return
	}
	glog.Infof("getHost body: %s", string(b))
	host = aesDe(b)
	glog.Infof("getHost host: %s", host)
	return
}

// doRequest get the order in json format with a sign
func doHttpPost(targetUrl string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("GET", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}
