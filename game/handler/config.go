package handler

import (
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"

	jsoniter "github.com/json-iterator/go"
)

//TODO 单条配置修改同步

//同步配置
func SyncConfig(arg *pb.SyncConfig) {
	switch arg.Type {
	case pb.CONFIG_BOX: //宝箱
		b := make([]data.Box, 0)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddBox(v)
		}
	case pb.CONFIG_ENV: //变量
		b := make(map[string]int32)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for k, v := range b {
			config.SetEnv2(k, v)
		}
	case pb.CONFIG_LOTTERY: //全民刮奖
		b := make(map[uint32]uint32)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for k, v := range b {
			config.SetLottery(k, v)
		}
	case pb.CONFIG_NOTICE: //公告
		b := make([]data.Notice, 0)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddNotice(v)
		}
	case pb.CONFIG_PRIZE: //抽奖
		b := make([]data.Prize, 0)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddPrize(v)
		}
	case pb.CONFIG_SHOP: //商城
		b := make(map[string]data.Shop)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddShop(v)
		}
	case pb.CONFIG_VIP: //VIP
		b := make(map[int]data.Vip)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddVip(v)
		}
	case pb.CONFIG_CLASSIC: //经典
		b := make(map[string]data.Classic)
		err := jsoniter.Unmarshal(arg.Data, &b)
		if err != nil {
			glog.Errorf("syncConfig Unmarshal err %v, data %#v", err, arg.Data)
			return
		}
		for _, v := range b {
			config.AddClassic(v)
		}
	default:
		glog.Errorf("syncConfig unknown type %d", arg.Type)
	}
}

//打包配置
func syncConfigMsg(ctype pb.ConfigType,
	d interface{}) interface{} {
	msg := new(pb.SyncConfig)
	msg.Type = ctype
	result, err := jsoniter.Marshal(d)
	if err != nil {
		glog.Errorf("syncConfig Marshal err %v, data %#v", err, d)
	}
	msg.Data = result
	return msg
}

//同步配置
func GetSyncConfig(ctype pb.ConfigType) interface{} {
	switch ctype {
	case pb.CONFIG_BOX: //宝箱
		return syncConfigMsg(pb.CONFIG_BOX, config.GetBoxs())
	case pb.CONFIG_ENV: //变量
		return syncConfigMsg(pb.CONFIG_ENV, config.GetEnvs())
	case pb.CONFIG_LOTTERY: //全民刮奖
		return syncConfigMsg(pb.CONFIG_LOTTERY, config.GetLotterys())
	case pb.CONFIG_NOTICE: //公告
		return syncConfigMsg(pb.CONFIG_NOTICE, config.GetNotices(data.NOTICE_TYPE1))
	case pb.CONFIG_PRIZE: //抽奖
		return syncConfigMsg(pb.CONFIG_PRIZE, config.GetPrizes())
	case pb.CONFIG_SHOP: //商城
		return syncConfigMsg(pb.CONFIG_SHOP, config.GetShops2())
	case pb.CONFIG_VIP: //VIP
		return syncConfigMsg(pb.CONFIG_VIP, config.GetVips())
	case pb.CONFIG_CLASSIC: //经典
		return syncConfigMsg(pb.CONFIG_CLASSIC, config.GetClassics())
	default:
	}
	return nil
}

//同步配置
func UpdateSyncConfig(ctype pb.ConfigType, msg []byte) (err error) {
	switch ctype {
	case pb.CONFIG_BOX: //宝箱
		b := new(data.Box)
		err = jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddBox(v)
	case pb.CONFIG_ENV: //变量
		b := make(map[string]int32)
		err := jsoniter.Unmarshal(msg, &b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		for k, v := range b {
			config.SetEnv2(k, v)
		}
	case pb.CONFIG_LOTTERY: //全民刮奖
		b := make(map[uint32]uint32)
		err := jsoniter.Unmarshal(msg, &b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		for k, v := range b {
			config.SetLottery(k, v)
		}
	case pb.CONFIG_NOTICE: //公告
		b := new(data.Notice)
		err := jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddNotice(b)
	case pb.CONFIG_PRIZE: //抽奖
		b := new(data.Prize)
		err := jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddPrize(b)
	case pb.CONFIG_SHOP: //商城
		b := new(data.Shop)
		err := jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddShop(b)
	case pb.CONFIG_VIP: //VIP
		b := new(data.Vip)
		err := jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddVip(b)
	case pb.CONFIG_CLASSIC: //经典
		b := new(data.Classic)
		err := jsoniter.Unmarshal(msg, b)
		if err != nil {
			glog.Errorf("update syncConfig Unmarshal err %v", err)
			return
		}
		config.AddClassic(b)
	default:
		glog.Errorf("syncConfig unknown type %d", ctype)
	}
	return
}
