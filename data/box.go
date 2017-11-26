package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//宝箱
type Box struct {
	Id       string    `bson:"_id"`      //id
	Duration uint32    `bson:"duration"` //时间(秒)
	Rtype    int       `bson:"rtype"`    //类型,1钻石,2金币
	Amount   int32     `bson:"amount"`   //数量
	Del      int       `bson:"del"`      //是否移除
	Ctime    time.Time `bson:"ctime"`    //创建时间
}

func GetBoxList() []Box {
	var list []Box
	ListByQ(Boxs, bson.M{"del": 0}, &list)
	return list
}
