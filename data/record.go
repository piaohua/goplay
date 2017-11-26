package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//TODO 优化，优化写入，优化返回给客户端的格式

//整局记录
type GameRecord struct {
	Roomid     string    `bson:"_id"`
	RecordJson []byte    `bson:"record_json"`
	Ctime      time.Time `bson:"ctime"`
}

func (this *GameRecord) Save() bool {
	this.Ctime = bson.Now()
	return Insert(GameRecords, this)
}

func (this *GameRecord) SaveZi() bool {
	this.Ctime = bson.Now()
	return Insert(GameRecordsZi, this)
}

func (this *GameRecord) Get() {
	Get(GameRecords, this.Roomid, this)
}

//个人记录
type UseridRecord struct {
	Userid   string    `bson:"userid"`
	Roomtype uint32    `bson:"roomtype"`
	Roomid   string    `bson:"roomid"`
	Ctime    time.Time `bson:"ctime"`
}

func (this *UseridRecord) Save() bool {
	this.Ctime = bson.Now()
	return Insert(UserRecords, this)
}

func (this *UseridRecord) SaveZi() bool {
	this.Ctime = bson.Now()
	return Insert(UserRecordsZi, this)
}

func (this *UseridRecord) Get() {
	Get(UserRecords, this.Roomid, this)
}

/*
//获取记录
func GetRecords(userid string, page int) ([]*GameRecord, error) {
	pageSize := 30
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	var list []*UseridRecord
	err := UserRecords.
		Find(bson.M{"userid": userid}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err != nil {
		return nil, err
	}
	rs := make([]*GameRecord, 0)
	for _, v := range list {
		r := new(GameRecord)
		r.Roomid = v.Roomid
		r.Get()
		if len(r.RecordJson) == 0 {
			continue
		}
		rs = append(rs, r)
	}
	return rs, nil
}
*/

//获取记录
func GetRecords(userid string, page int) ([]*GameRecord, error) {
	pageSize := 4
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	//fmt.Println(userid, skipNum, sortFieldR)
	var list []*UseridRecord
	err := UserRecords.
		Find(bson.M{"userid": userid}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err != nil {
		return nil, err
	}
	//fmt.Println("len ", len(list))
	var rs []*GameRecord
	in := make([]string, 0)
	for _, v := range list {
		in = append(in, v.Roomid)
	}
	ListByQ(GameRecords, bson.M{"_id": bson.M{"$in": in}}, &rs)
	return rs, nil
}

//获取记录
func GetRecordsZi(userid string, page int) ([]*GameRecord, error) {
	pageSize := 4
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	//fmt.Println(userid, skipNum, sortFieldR)
	var list []*UseridRecord
	err := UserRecordsZi.
		Find(bson.M{"userid": userid}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err != nil {
		return nil, err
	}
	//fmt.Println("len ", len(list))
	var rs []*GameRecord
	in := make([]string, 0)
	for _, v := range list {
		in = append(in, v.Roomid)
	}
	ListByQ(GameRecordsZi, bson.M{"_id": bson.M{"$in": in}}, &rs)
	return rs, nil
}

//输赢趋势，天 地 玄 黄
type Trend struct {
	Roomid string    `bson:"roomid"`
	Round  uint32    `bson:"round"`
	Dealer string    `bson:"dealer"`
	Seat2  bool      `bson:"seat2"`
	Seat3  bool      `bson:"seat3"`
	Seat4  bool      `bson:"seat4"`
	Seat5  bool      `bson:"seat5"`
	Ctime  time.Time `bson:"ctime"`
}

func (this *Trend) Save() bool {
	this.Ctime = bson.Now()
	return Insert(Trends, this)
}

func (this *Trend) Get(id string) {
	Get(Trends, id, this)
}

//获取记录
func GetTrends(roomid string, page int) ([]*Trend, error) {
	pageSize := 10
	skipNum, sortFieldR := parsePageAndSort(page, pageSize, "ctime", false)
	var list []*Trend
	err := Trends.
		Find(bson.M{"roomid": roomid}).
		Sort(sortFieldR).
		Skip(skipNum).
		Limit(pageSize).
		All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
