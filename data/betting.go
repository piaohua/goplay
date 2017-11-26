package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//seat
//(1有牛，2无牛)
//(3牛1-3， 4牛4-6，5牛7-9)
//(6牛牛，7两对，8同花，9五花牛，10四炸，11五小牛)

// 疯狂投注
type BettingInfo struct {
	Id    string    `bson:"_id" json:"id"`      //id
	Seat  uint32    `bson:"seat" json:"seat"`   //位置
	Odds  float32   `bson:"odds" json:"odds"`   //赔率
	Del   int       `bson:"del" json:"del"`     //是否移除
	Ctime time.Time `bson:"ctime" json:"ctime"` //创建时间
}

func GetBettingList() []BettingInfo {
	var list []BettingInfo
	ListByQ(Bettings, bson.M{"del": 0}, &list)
	return list
}

// 疯狂投注
type BettingRecord struct {
	Id    string                   `bson:"_id" json:"id"`
	Cards []uint32                 `bson:"cards" json:"cards"`
	Niu   uint32                   `bson:"niu" json:"niu"`
	Seats []uint32                 `bson:"seats" json:"seats"`
	Lose  map[string]int32         `bson:"lose" json:"lose"`
	Ante  map[string][]SeatBetting `bson:"ante" json:"ante"`
	Ctime time.Time                `bson:"ctime" json:"ctime"`
}

type SeatBetting struct {
	Seat   uint32 `json:"seat"`
	Number uint32 `json:"number"`
}

func (this *BettingRecord) Save() bool {
	this.Ctime = bson.Now()
	return Insert(BettingRecords, this)
}

//个人记录
type BettingUser struct {
	Userid string    `bson:"userid" json:"userid"`
	Index  string    `bson:"index" json:"index"`
	Lose   int32     `bson:"lose" json:"lose"`
	Ctime  time.Time `bson:"ctime" json:"ctime"`
}

func (this *BettingUser) Save() bool {
	this.Ctime = bson.Now()
	return Insert(BettingUsers, this)
}

//获取记录
func GetBettingRecords(userid string, page int) ([]*BettingRecord, error) {
	pageSize := 5
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	//fmt.Println(userid, skipNum, sortFieldR)
	var list []*BettingUser
	err := BettingUsers.
		Find(bson.M{"userid": userid}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err != nil {
		return nil, err
	}
	//fmt.Println("len ", len(list))
	var rs []*BettingRecord
	in := make([]string, 0)
	for _, v := range list {
		in = append(in, v.Index)
	}
	ListByQ(BettingRecords, bson.M{"_id": bson.M{"$in": in}}, &rs)
	return rs, nil
}
