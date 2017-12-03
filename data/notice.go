package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	NOTICE_TYPE1 = 1 //公告消息
	NOTICE_TYPE2 = 2 //广播消息
)

const (
	NOTICE_ACT_TYPE0 = 0 //无操作消息
	NOTICE_ACT_TYPE1 = 1 //支付消息
	NOTICE_ACT_TYPE2 = 2 //活动消息
)

//公告
type Notice struct {
	Id      string    `bson:"_id" json:"id"`           //公告ID
	Rtype   int       `bson:"rtype" json:"rtype"`      //类型,1=公告消息,2=广播消息
	Atype   uint32    `bson:"atype" json:"atype"`      //分包类型
	Acttype int       `bson:"act_type" json:"acttype"` //操作类型,0=无操作,1=支付,2=活动
	Top     int       `bson:"top" json:"top"`          //置顶
	Num     int       `bson:"num" json:"num"`          //广播次数
	Del     int       `bson:"del" json:"del"`          //是否移除
	Content string    `bson:"content" json:"content"`  //广播内容
	Etime   time.Time `bson:"etime" json:"etime"`      //过期时间
	Ctime   time.Time `bson:"ctime" json:"ctime"`      //创建时间
}

func GetNoticeList(rtype int) []Notice {
	var list []Notice
	ListByQ(Notices, bson.M{"del": 0, "rtype": rtype, "etime": bson.M{"$gt": bson.Now()}}, &list)
	return list
}
