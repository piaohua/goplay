package main

import (
	"net"
	"net/http"
	"sync"
	"time"

	"goplay/glog"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	Addr            string        //地址
	MaxConnNum      int           //最大连接数
	PendingWriteNum int           //等待写入消息长度
	MaxMsgLen       uint32        //最大消息长度
	HTTPTimeout     time.Duration //超时时间
	ln              net.Listener  //监听
	handler         *WSHandler    //处理
}

type WSHandler struct {
	maxConnNum      int                //最大连接数
	pendingWriteNum int                //等待写入消息长度
	maxMsgLen       uint32             //最大消息长
	upgrader        websocket.Upgrader //升级http连接
	conns           WebsocketConnSet   //连接集合
	mutexConns      sync.Mutex         //互斥锁
	wg              sync.WaitGroup     //同步机制
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorf("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		glog.Errorf("too many connections: %d", len(handler.conns))
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := newWSConn(conn, handler.pendingWriteNum, handler.maxMsgLen)
	wsConn.pid = wsConn.initWs()
	go wsConn.writePump()
	wsConn.readPump()

	// cleanup
	wsConn.Close()
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
	// pid stop
	glog.Infof("wsConn.pid: %s", wsConn.pid.String())
	wsConn.pid.Stop()
}

func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		glog.Fatal("%v", err)
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 10000
		glog.Infof("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		glog.Infof("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 40960
		glog.Infof("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		glog.Infof("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}

	server.ln = ln
	server.handler = &WSHandler{
		maxConnNum:      server.MaxConnNum,
		pendingWriteNum: server.PendingWriteNum,
		maxMsgLen:       server.MaxMsgLen,
		conns:           make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			ReadBufferSize:   1024, //default 4096
			WriteBufferSize:  1024, //default 4096
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
