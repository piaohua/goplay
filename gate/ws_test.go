package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"goplay/pb"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	proto "github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

// 启动一个服务
func TestProto(t *testing.T) {
	packet := &pb.CLogin{}
	t.Logf("----:%+v", packet.String())
	s := proto.MessageName(&pb.CRegist{})
	t.Logf("s -> %s", s)
	//p := proto.MessageType(s)
	p := reflect.TypeOf(pb.CRegist{})
	t.Logf("p -> %v", p) //reflect.Type
	v := reflect.New(p)
	t.Logf("v -> %v", v)
	t.Logf("v -> %v", v.Interface())
	reg := &pb.CRegist{
		Phone: "1111",
	}
	b, err := proto.Marshal(reg)
	t.Log("err -> ", err)
	t.Logf("b -> %#v", b)
	err = proto.Unmarshal(b, v.Interface().(proto.Message))
	t.Log("err -> ", err)
	t.Logf("b -> %#v", b)
}

// 启动一个服务
func TestRunServer(t *testing.T) {
	closeSig := make(chan bool, 1)
	Run(closeSig)
}

func Run(closeSig chan bool) {
	wsServer := new(WSServer)
	//wsServer
	wsServer.Addr = "127.0.0.1:8880"
	if wsServer != nil {
		wsServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}

//const (
//	HeaderLen uint32 = 1 //包头长度
//	PROTOLen  uint32 = 4 //协议头长度
//	DataLen   uint32 = 4 //数据长度
//	HANDDLen  uint32 = 9 //消息头总长度
//)

// 发送一条消息
func TestClient(t *testing.T) {
	var HeaderLen uint32 = 1 //包头长度
	var PROTOLen uint32 = 4
	var DataLen uint32 = 4 //包信息数据长度占位长度
	var HANDDLen uint32 = 9
	var count uint32 = 1
	var p uint32 = 1022
	packet := &pb.CRegist{
		Nickname: "wwww",
		Phone:    "1111",
	}
	message, _ := proto.Marshal((proto.Message)(packet))
	t.Logf("message -> %+v", message)
	msglen := uint32(len(message))
	buff := make([]byte, int(HANDDLen)+len(message))
	t.Logf("buff -> %+v", buff)
	buff[0] = byte(count)
	t.Logf("count -> %d, buff -> %+v", count, buff)
	t.Logf("p -> %d, %+v", p, encodeUint32(p))
	t.Logf("msglen -> %d, %+v", msglen, encodeUint32(msglen))
	copy(buff[HeaderLen:HeaderLen+PROTOLen], encodeUint32(p))
	copy(buff[HeaderLen+PROTOLen:HeaderLen+PROTOLen+DataLen], encodeUint32(msglen))
	copy(buff[HANDDLen:HANDDLen+msglen], message)
	t.Logf("buff -> %+v", buff)
	client(buff)
}

// 发送一条消息
func client(buff []byte) {
	var addr string = "127.0.0.1:8880"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {""}})
	//fmt.Printf("c -> %+v\n", c)
	fmt.Printf("err -> %+v\n", err)
	if err != nil {
		fmt.Printf("err -> %+v\n", err)
	}
	if c != nil {
		c.WriteMessage(websocket.TextMessage, buff)
		c.Close()
	}
}

// 发送多条消息,粘包测试
func TestRunClient(t *testing.T) {
	// 注册协议请求消息
	packet := &pb.CRegist{
		Nickname: "wwww",
		Phone:    "1111",
	}
	// 打包protobuf协议消息
	message, _ := proto.Marshal((proto.Message)(packet))
	t.Logf("message -> %+v", message)
	// 打包完整协议消息
	buff := pack(packet.GetCode(), message, 0)
	// 消息长度
	blen := len(buff)
	t.Logf("buff -> %+v, len(buff) -> %d", buff, blen)
	// 链接服务器
	conn, err := client_conn()
	if err != nil {
		t.Logf("conn error: %s\n", err)
		return
	}
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 写入5条消息
	for i := 0; i < 5; i++ {
		msgbuf.Write(buff)
	}
	t.Logf("msgbuf len -> %d", msgbuf.Len())
	t.Logf("msgbuf -> %v", msgbuf)
	// 发送一条完整消息
	client_write(conn, msgbuf.Next(blen))
	time.Sleep(time.Second)
	// 发送一条不完整的消息头
	client_write(conn, msgbuf.Next(2))
	time.Sleep(time.Second)
	// 发送消息剩下部分
	client_write(conn, msgbuf.Next(blen-2))
	time.Sleep(time.Second)
	// 三条消息分多段发送
	client_write(conn, msgbuf.Next(blen+2))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(-2+blen-6))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(6+3))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(-3+blen))
	time.Sleep(time.Second)
	// 关闭连接
	conn.Close()
}

//召唤机器人
func client_write(c *websocket.Conn, buff []byte) {
	c.WriteMessage(websocket.TextMessage, buff)
}

//召唤机器人
func client_conn() (*websocket.Conn, error) {
	var addr string = "127.0.0.1:8880"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {""}})
	return c, err
}

// 协议路由测试
func TestRouter(t *testing.T) {
	reg := &pb.Request{
		UserName: "xxxx",
	}
	b, err := reg.Marshal()
	t.Log("err -> ", err)
	registTest(1000, pb.Request{})
	routerTest(1000, b)
}

var protoTypes map[uint32]reflect.Type = make(map[uint32]reflect.Type)

func registTest(id uint32, x interface{}) {
	if _, ok := protoTypes[id]; ok {
		// TODO: Some day, make this a panic.
		log.Printf("proto: duplicate proto type registered: %d", id)
		return
	}
	t := reflect.TypeOf(x)
	protoTypes[id] = t
}

func routerTest(id uint32, b []byte) {
	if t, ok := protoTypes[id]; ok {
		v := reflect.New(t)
		err := proto.Unmarshal(b, v.Interface().(proto.Message))
		if err != nil {
			fmt.Printf("err -> %#v\n", err)
		}
		fmt.Printf("b -> %#v\n", b)
	}
}

// actor测试
func TestActor(t *testing.T) {
	wsConn := &WSConn{}
	wsConn.pid = wsConn.initWs()
	wsConn.pid.Tell(&testWs{Who: "testpiao"})
	//---
	props4 := actor.FromInstance(&WSConn{})
	pid4 := actor.Spawn(props4)
	pid4.Tell(&testWs{Who: "test"})
	<-time.After(time.Duration(6) * time.Second)
	console.ReadLine()
}

// 打包，解包测试
func TestMarshal(t *testing.T) {
	reg := &pb.Request{
		UserName: "xxxx",
	}
	b, err := reg.Marshal()
	t.Log("err -> ", err)
	t.Log("b -> ", b)
	b, err = proto.Marshal(reg)
	t.Log("err -> ", err)
	t.Log("b -> ", b)
	//---
	msg := &pb.SRegist{
		Userid: "xxxx",
	}
	code, msg_b, err := pb.Packet(msg)
	t.Log(code, msg_b, err)
	//---
	msg1 := &pb.CRegist{
		Nickname: "xxxx",
		Phone:    "xxxx",
		Password: "xxxx",
	}
	code1, b1, err := pb.Packet(msg1)
	t.Log(code1, b1, err)
	//---
	msg2, err := pb.Unpack(code1, b1)
	t.Log(msg2, err)
}

// 打包，解包文件模板
func TestTpl(t *testing.T) {
	var str string
	var err error
	str = tpl_unpack()
	err = ioutil.WriteFile("unpack.go", []byte(str), 0644)
	t.Log(err)
	str = tpl_packet()
	err = ioutil.WriteFile("packet.go", []byte(str), 0644)
	t.Log(err)
}

func tpl_packet() string {
	str := `
// Code generated by protoc-gen-main.
// source: regist/msg.go
// DO NOT EDIT!

package socket

import (
	"errors"
	"messages"
)

//打包消息
func packet(msg interface{}) (uint32, []byte, error) {
	switch msg.(type) {
	case *messages.CRegist:
		code := msg.(*messages.CRegist).GetCode()
		b, err := msg.(*messages.CRegist).Marshal()
		return code, b, err
	case *messages.SRegist:
		code := msg.(*messages.SRegist).GetCode()
		b, err := msg.(*messages.SRegist).Marshal()
		return code, b, err
	case *messages.SLogin:
		code := msg.(*messages.SLogin).GetCode()
		b, err := msg.(*messages.SLogin).Marshal()
		return code, b, err
	default:
		return 0, []byte{}, errors.New("msg wrong")
	}
}
	`
	return str
}

func tpl_unpack() string {
	str := `
// Code generated by protoc-gen-main.
// source: regist/msg.go
// DO NOT EDIT!

package socket

import (
	"errors"
	"messages"
)

//解包消息
func unpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	case 1000:
		msg := &messages.CLogin{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1002:
		msg := &messages.CRegist{}
		err := msg.Unmarshal(b)
		return msg, err
	default:
		return nil, errors.New("msg id wrong")
	}
}
	`
	return str
}

func TestCode(t *testing.T) {
	b := encodeUint32(2014)
	t.Logf("%v\n", b)
	c := decodeUint32(b)
	t.Logf("%d\n", c)
}
