package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//抽奖
type Prize struct {
	Id     string    `bson:"_id"`    //id
	Rate   uint32    `bson:"rate"`   //概率
	Rtype  int       `bson:"rtype"`  //类型,1钻石,2金币
	Amount int32     `bson:"amount"` //数量
	Del    int       `bson:"del"`    //是否移除
	Ctime  time.Time `bson:"ctime"`  //创建时间
}

func GetPrizeList() []Prize {
	var list []Prize
	ListByQ(Prizes, bson.M{"del": 0}, &list)
	return list
}
