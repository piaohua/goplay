package data

import "time"

//经典
type Classic struct {
	Id      string    `bson:"_id" json:"id"`          //id
	Ptype   int       `bson:"ptype" json:"ptype"`     //玩法类型1看牌抢庄,3通比牛牛4牛牛坐庄
	Rtype   int       `bson:"rtype" json:"rtype"`     //房间类型1初级,2中级,3高级,4大师
	Ante    uint32    `bson:"ante" json:"ante"`       //房间底分
	Minimum uint32    `bson:"minimum" json:"minimum"` //房间最低
	Maximum uint32    `bson:"maximum" json:"maximum"` //房间最高0表示没限制
	Ctime   time.Time `bson:"ctime" json:"ctime"`     //
}

func GetClassicList() []Classic {
	var list []Classic
	ListByQ(Classics, nil, &list)
	return list
}

/*
id = ptype * 100 + rtype

看牌抢庄
初级场  100   1000 ~ 20,0000
中级场  1000  10,0000 ~ 200,0000
高级场  10000 100,0000 ~

牛牛做庄
初级场  100   1000 ~ 20,0000
中级场  1000  10,0000 ~ 200,0000
高级场  10000 100,0000 ~

通比牛牛
初级场  10000 20,0000 ~
中级场  10000 500,0000 ~
高级场  20000 1000,0000 ~
*/
