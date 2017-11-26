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
	}
	return 0
}
