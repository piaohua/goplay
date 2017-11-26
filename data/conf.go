package data

import (
	"encoding/json"
	"os"
)

var Conf Config

type Config struct {
	Timeout        int    `json:"timeout"`          //超时时间
	ServerId       string `json:"server_id"`        //服务器ID
	DbHost         string `json:"db_host"`          //数据库地址
	DbPort         string `json:"db_port"`          //数据库端口
	DbUser         string `json:"db_user"`          //数据库用户
	DbPassword     string `json:"db_password"`      //数据库密码
	DbName         string `json:"db_name"`          //数据库名
	ServerHost     string `json:"server_host"`      //服务器地址
	ServerPort     string `json:"server_port"`      //服务器端口
	RobotHost      string `json:"robot_host"`       //机器人端口
	RobotPort      string `json:"robot_port"`       //机器人端口
	RobotPhone     string `json:"robot_phone"`      //机器人注册账号
	WebPort        string `json:"web_port"`         //web服务器端口
	WebPattern     string `json:"web_pattern"`      //后台地址
	ImagePattern   string `json:"image_pattern"`    //头像地址
	PayWxPattern   string `json:"pay_wx_pattern"`   //微信支付回调路径
	PayIappPattern string `json:"pay_iapp_pattern"` //爱贝支付回调路径
	ShareAddr      string `json:"share_addr"`       //分享地址
	GmKey          string `json:"gm_key"`           //gm密钥
	Version        string `json:"version"`          //版本号
}

func LoadConf(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Conf)
	if err != nil {
		panic(err)
	}
}
