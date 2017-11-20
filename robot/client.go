/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 13:11:36
 * Filename      : client.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"bytes"
	"errors"
	"time"

	"goplay/glog"
	"goplay/pb"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second // Time allowed to read the next pong message from the peer.
	pingPeriod     = 9 * time.Second  // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 1024             // Maximum message size allowed from peer.
	waitForLogin   = 5 * time.Second  // 连接建立后5秒内没有收到登陆请求,断开socket
)

type WebsocketConnSet map[*websocket.Conn]struct{}

// 机器人连接数据
type Robot struct {
	conn *websocket.Conn // websocket连接

	stopCh chan struct{}    // 关闭通道
	msgCh  chan interface{} // 消息通道

	maxMsgLen uint32 // 最大消息长度
	index     int    // 包序

	//游戏数据
	data      *user    //数据
	code      string   //邀请码
	seat      uint32   //位置
	round     uint32   //次数
	sits      uint32   //尝试坐下次数
	bits      uint32   //尝试下注次数
	cards     []uint32 //手牌
	regist    bool     //注册标识
	rtype     uint32   //房间类型
	classic   []*pb.Classic
	classicId string
	timer     uint32 //在线时间
}

// 基本数据
type user struct {
	Userid   string // 用户id
	Nickname string // 用户昵称
	Sex      uint32 // 用户性别,男1 女2 非男非女3
	Phone    string // 绑定的手机号码
	Coin     uint32 // 金币
	Diamond  uint32 // 钻石
}

//创建连接
func newRobot(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *Robot {
	return &Robot{
		maxMsgLen: maxMsgLen,

		conn: conn,
		data: new(user),

		msgCh:  make(chan interface{}, pendingWriteNum),
		stopCh: make(chan struct{}),
	}
}

//断开连接
func (ws *Robot) Close() {
	select {
	case <-ws.stopCh:
		return
	default:
		//停止发送消息
		close(ws.stopCh)
		//关闭消息通道
		close(ws.msgCh)
		//关闭连接
		ws.conn.Close()
		//Logout message
		Logout(ws.data.Phone, ws.code)
	}
}

//接收
func (ws *Robot) Router(id uint32, body []byte) {
	msg, err := pb.Runpack(id, body)
	if err != nil {
		glog.Error("protocol unpack err:", id, err)
		return
	}
	ws.receive(msg)
}

//发送消息
func (ws *Robot) Sender(msg interface{}) {
	if ws.msgCh == nil {
		glog.Errorf("WSConn msg channel closed %x", msg)
		return
	}
	if len(ws.msgCh) == cap(ws.msgCh) {
		glog.Errorf("send msg channel full -> %d", len(ws.msgCh))
		return
	}
	select {
	case <-ws.stopCh:
		return
	case ws.msgCh <- msg:
	}
}

func (ws *Robot) readPump() {
	defer ws.Close()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 消息长度
	var length int = 0
	// 包序长度
	var index int = 0
	// 协议编号
	var proto uint32 = 0
	for {
		n, message, err := ws.conn.ReadMessage()
		if err != nil {
			glog.Errorf("Read error: %s, %d\n", err, n)
			break
		}
		// 数据添加到消息缓冲
		m, err := msgbuf.Write(message)
		if err != nil {
			glog.Errorf("Buffer write error: %s, %d\n", err, m)
			return
		}
		// 消息分割循环
		for {
			// 消息头
			if length == 0 && msgbuf.Len() >= 9 {
				index = int(msgbuf.Next(1)[0])             //包序
				proto = decodeUint32(msgbuf.Next(4))       //协议号
				length = int(decodeUint32(msgbuf.Next(4))) //消息长度
				// 检查超长消息
				if length > 1024 {
					glog.Errorf("Message too length: %d\n", length)
					return
				}
			}
			//fmt.Printf("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
			// 消息体
			if length > 0 && msgbuf.Len() >= length {
				//fmt.Printf("Client messge: %s\n", string(msgbuf.Next(length)))
				//包序验证
				ws.index++
				ws.index = ws.index % 256
				//fmt.Printf("Message index error: %d, %d\n", index, ws.index)
				if ws.index != index {
					glog.Errorf("Message index error: %d, %d\n", index, ws.index)
					//return
				}
				//路由
				ws.Router(proto, msgbuf.Next(length))
				length = 0
			} else {
				break
			}
		}
	}
}

//消息写入 TODO write Buff
func (ws *Robot) writePump() {
	tick := time.Tick(pingPeriod)
	for {
		select {
		case <-tick:
			err := ws.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		default:
		}
		select {
		case message, ok := <-ws.msgCh:
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				return
			}
			err := ws.write(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

//写入
func (ws *Robot) write(mt int, msg interface{}) error {
	var message []byte
	switch msg.(type) {
	case []byte:
		message = msg.([]byte)
	default:
		code, body, err := pb.Rpacket(msg)
		if err != nil {
			glog.Errorf("write msg err %v", msg)
			return err
		}
		message = pack(code, body, ws.index)
	}
	if uint32(len(message)) > ws.maxMsgLen {
		glog.Errorf("write msg too long -> %d", len(message))
		return errors.New("write msg too long")
	}
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, message)
}

func decodeUint32(b []byte) (i uint32) {
	i = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return
}

func encodeUint32(i uint32) (b []byte) {
	b = append(b, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return
}

//封包
func pack(code uint32, msg []byte, index int) []byte {
	buff := make([]byte, 9+len(msg))
	msglen := uint32(len(msg))
	buff[0] = byte(index)
	copy(buff[1:5], encodeUint32(code))
	copy(buff[5:9], encodeUint32(msglen))
	copy(buff[9:], msg)
	return buff
}
