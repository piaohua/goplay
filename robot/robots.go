/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 11:32:23
 * Filename      : robots.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"time"

	"goplay/glog"
	"goplay/pb"
	"utils"
)

//消息通知
func Msg2Robots(msg interface{}, num uint32) {
	for num > 0 {
		rbs.Send2rbs(msg)
		num--
	}
}

//登录成功
func Logined(phone string) {
	msg := &pb.RobotLogin{
		Phone: phone,
	}
	Msg2Robots(msg, 1)
}

//登出成功
func Logout(phone, code string) {
	msg := &pb.RobotLogout{
		Phone: phone,
		Code:  code,
	}
	Msg2Robots(msg, 1)
}

// 已经注册,重新登录
func ReLogined(phone, code string, rtype uint32) {
	msg := &pb.RobotReLogin{
		Phone: phone,
		Code:  code,
		Rtype: rtype,
	}
	Msg2Robots(msg, 1)
}

//发送消息
func (r *RobotServer) Send2rbs(msg interface{}) {
	if r.msgCh == nil {
		glog.Errorf("server msg channel closed %#v", msg)
		return
	}
	if len(r.msgCh) == cap(r.msgCh) {
		glog.Errorf("send msg channel full -> %d", len(r.msgCh))
		return
	}
	select {
	case <-r.stopCh:
		return
	default:
	}
	select {
	case <-r.stopCh:
		return
	default:
		r.msgCh <- msg
	}
}

func (r *RobotServer) runTest1() {
	msg := &pb.RobotMsg{
		Code: "free",
	}
	//TODO 数量添加到配置
	go Msg2Robots(msg, 1)
}

//机器人测试
func (r *RobotServer) runTest() {
	glog.Infof("runTest started phone -> %s", r.phone)
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			if len(r.online) < 3 {
				msg := &pb.RobotMsg{
					Code: "create",
				}
				go Msg2Robots(msg, 3)
			}
		}
	}
}

//机器人测试
func (r *RobotServer) runFree() {
	glog.Infof("runFree started phone -> %s", r.phone)
	tick := time.Tick(5 * time.Minute)
	num := 0
	msg1 := &pb.RobotMsg{
		Code: "free",
	}
	msg2 := &pb.RobotMsg{
		Code: "classic",
	}
	for {
		select {
		case <-tick:
			now := utils.Timestamp()
			today := utils.TimestampToday()
			if (now - today) < (8 * 3600) {
				continue
			}
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化,按时间段运行
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			go Msg2Robots(msg1, 5)
			if len(r.online) < 60 {
				go Msg2Robots(msg2, 10)
			}
			if num > 1000 {
				break
			}
			num++
		}
	}
}

//处理
func (r *RobotServer) run() {
	defer func() {
		glog.Infof("Robots closed online -> %d\n", len(r.online))
		glog.Infof("Robots closed offline -> %d\n", len(r.offline))
		glog.Infof("Robots closed phone -> %s\n", r.phone)
	}()
	glog.Infof("Robots started -> %s", r.phone)
	tick := time.Tick(time.Minute)
	for {
		select {
		case m, ok := <-r.msgCh:
			if !ok {
				glog.Errorf("Robots msgCh closed phone -> %s\n", r.phone)
				return
			}
			switch m.(type) {
			case *pb.RobotMsg:
				//启动机器人
				msg := m.(*pb.RobotMsg)
				var code string = msg.Code
				var rtype uint32 = msg.Rtype
				var phone string
				for k, v := range r.offline {
					if v {
						phone = k
						r.offline[k] = false
						go r.RunRobot(phone, code, rtype, false)
						break
					}
				}
				if len(phone) == 0 {
					phone = r.phone
					r.phone = utils.StringAdd(r.phone)
					//go r.RunRobot(phone, code, rtype, true)
					//新机器人不用注册
					go r.RunRobot(phone, code, rtype, false)
				}
				//phone = r.phone
				//r.phone = utils.StringAdd(r.phone)
				//go r.RunRobot(phone, code, rypte, r.msgCh)
				glog.Infof("phone -> %s", phone)
			case *pb.RobotReLogin:
				//重新尝试登录
				msg := m.(*pb.RobotReLogin)
				glog.Infof("ReLogin -> %#v", msg)
				go r.RunRobot(msg.Phone, msg.Code, msg.Rtype, false)
			case *pb.RobotLogin:
				//登录成功
				msg := m.(*pb.RobotLogin)
				glog.Infof("login -> %#v", msg)
				delete(r.offline, msg.Phone)
				r.online[msg.Phone] = true
			case *pb.RobotLogout:
				//登出成功
				msg := m.(*pb.RobotLogout)
				glog.Infof("logout -> %#v", msg)
				if _, ok := r.online[msg.Phone]; ok {
					delete(r.online, msg.Phone)
					r.offline[msg.Phone] = true
				}
			case closeFlag:
				//停止发送消息
				close(r.stopCh)
				return
			}
		case <-tick:
			//逻辑处理
		}
	}
}
