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
	Id     string    `bson:"_id" json:"id"`        //购买ID
	Atype  uint32    `bson:"atype" json:"atype"`   //分包类型
	Status int       `bson:"status" json:"status"` //物品状态,1=热卖
	Propid int       `bson:"propid" json:"propid"` //兑换的物品,1=钻石,2=金币
	Payway int       `bson:"payway" json:"payway"` //支付方式,1=RMB,,2=钻石
	Number uint32    `bson:"number" json:"number"` //兑换的数量
	Price  uint32    `bson:"price" json:"price"`   //支付价格(单位元)
	Name   string    `bson:"name" json:"name"`     //物品名字
	Info   string    `bson:"info" json:"info"`     //物品信息
	Del    int       `bson:"del" json:"del"`       //是否移除
	Etime  time.Time `bson:"etime" json:"etime"`   //过期时间
	Ctime  time.Time `bson:"ctime" json:"ctime"`   //创建时间
}

func GetShopList() []Shop {
	var list []Shop
	ListByQ(Shops, bson.M{"del": 0, "etime": bson.M{"$gt": bson.Now()}}, &list)
	return list
}
