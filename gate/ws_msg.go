package main

import (
	"goplay/glog"
	"goplay/pb"
)

//路由
func (ws *WSConn) Router(id uint32, body []byte) {
	body = aesDe(body) //解密
	msg, err := pb.Unpack(id, body)
	if err != nil {
		glog.Error("protocol unpack err:", id, err)
		return
	}
	ws.pid.Tell(msg)
}

//发送消息
func (ws *WSConn) Send(msg interface{}) {
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

//封包
func pack(code uint32, msg []byte, index int) []byte {
	msg = aesEn(msg) //加密
	buff := make([]byte, 9+len(msg))
	msglen := uint32(len(msg))
	buff[0] = byte(index)
	copy(buff[1:5], encodeUint32(code))
	copy(buff[5:9], encodeUint32(msglen))
	copy(buff[9:], msg)
	return buff
}

func encodeUint32(i uint32) (b []byte) {
	b = append(b, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return
}

func decodeUint32(b []byte) (i uint32) {
	i = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return
}
