package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	COIN    uint32 = 2
	DIAMOND uint32 = 1

	RMB uint32 = 1
	DIA uint32 = 2
)

//商城
type Shop struct {
	Id     string    `bson:"_id"`    //购买ID
	Atype  uint32    `bson:"atype"`  //分包类型
	Status int       `bson:"status"` //物品状态,1=热卖
	Propid int       `bson:"propid"` //兑换的物品,1=钻石,2=金币
	Payway int       `bson:"payway"` //支付方式,1=RMB,,2=钻石
	Number uint32    `bson:"number"` //兑换的数量
	Price  uint32    `bson:"price"`  //支付价格(单位元)
	Name   string    `bson:"name"`   //物品名字
	Info   string    `bson:"info"`   //物品信息
	Del    int       `bson:"del"`    //是否移除
	Etime  time.Time `bson:"etime"`  //过期时间
	Ctime  time.Time `bson:"ctime"`  //创建时间
}

func GetShopList() []Shop {
	var list []Shop
	ListByQ(Shops, bson.M{"del": 0, "etime": bson.M{"$gt": bson.Now()}}, &list)
	return list
}
