package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// 全民刮奖
//times 1有牛,2牛牛,3五花牛,4五小牛,5四炸
type LotteryInfo struct {
	Id      string    `bson:"_id" json:"id"`          //id
	Name    string    `bson:"name" json:"name"`       //名称
	Times   uint32    `bson:"times" json:"times"`     //次数
	Diamond uint32    `bson:"diamond" json:"diamond"` //钻石
	Del     int       `bson:"del" json:"del"`         //是否移除
	Ctime   time.Time `bson:"ctime" json:"ctime"`     //创建时间
}

func GetLotteryList() []LotteryInfo {
	var list []LotteryInfo
	ListByQ(Lotterys, bson.M{"del": 0}, &list)
	return list
}
