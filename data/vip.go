package data

import "time"

//vip
type Vip struct {
	Id     string    `bson:"_id" json:"id"`        //ID
	Level  int       `bson:"level" json:"level"`   //等级
	Number uint32    `bson:"number" json:"number"` //等级充值金额数量限制(分)
	Pay    uint32    `bson:"pay" json:"pay"`       //充值赠送百分比5=赠送充值的5%
	Prize  uint32    `bson:"prize" json:"prize"`   //赠送抽奖次数
	Kick   int       `bson:"kick" json:"kick"`     //经典场可踢人次数
	Ctime  time.Time `bson:"ctime" json:"ctime"`   //创建时间
}

func GetVipList() []Vip {
	var list []Vip
	ListByQ(Vips, nil, &list)
	return list
}
