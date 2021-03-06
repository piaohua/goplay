package config

import (
	"goplay/data"
	"utils"
)

// 环境变量
var EnvMap *utils.BeeMap

// 启动初始化
func InitEnv() {
	EnvMap = utils.NewBeeMap()
	list := data.GetEnvList()
	for _, v := range list {
		SetEnv(v.Key, v.Value)
	}
}

// 启动初始化
func InitEnv2() {
	EnvMap = utils.NewBeeMap()
}

// 获取元素
func GetEnv2(k interface{}) interface{} {
	return EnvMap.Get(k)
}

// 设置元素
func SetEnv(k interface{}, v interface{}) bool {
	e := data.Env{Key: k.(string), Value: v.(int32)}
	e.SetEnv()
	return EnvMap.Set(k, v)
}

// 设置元素,不操作数据库
func SetEnv2(k interface{}, v interface{}) bool {
	return EnvMap.Set(k, v)
}

// 删除元素
func DelEnv(k interface{}) {
	e := data.Env{Key: k.(string)}
	e.DelEnv()
	EnvMap.Delete(k)
}

// 存在元素
func CheckEnv(k interface{}) bool {
	return EnvMap.Check(k)
}

// 全部元素
func ItemsEnv() map[interface{}]interface{} {
	return EnvMap.Items()
}

// 全部元素
func GetEnvs() map[string]int32 {
	m := make(map[string]int32)
	for k, v := range ItemsEnv() {
		m[k.(string)] = v.(int32)
	}
	return m
}

//获取变量,变量默认值设置

func GetEnv(k interface{}) int32 {
	key, ok := k.(string)
	if !ok {
		return 0
	}
	if CheckEnv(k) {
		return GetEnv2(k).(int32)
	}
	//默认值
	switch key {
	case data.ENV1:
		return 128 //注册赠送钻石
	case data.ENV2:
		return 88888 //注册赠送金币
	case data.ENV3:
		return 30 //绑定赠送钻石
	case data.ENV11:
		return 100 //百人场人数
	case data.ENV15:
		return 5 //全民刮奖
	case data.ENV16:
		return 100 //全民刮奖最大注
	case data.ENV17:
		return 50000000 //系统上庄限额
	case data.ENV18:
		return 100000 //庄家上庄限额
	case data.ENV19:
		return 100000 //庄家下庄限额
	case data.ENV20:
		return 8 //做庄次数限制
	case data.ENV21:
		return 300000 //玩家坐下限额
	case data.ENV22:
		return 5000000 //中奖下注限额
	}
	return 0
}
