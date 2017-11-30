package data

import (
	"utils"

	"gopkg.in/mgo.v2/bson"
)

//在同一个collection中
const (
	ROOMID_KEY  = "LastRoomID" //房间唯一id
	USERID_KEY  = "LastUserID" //玩家唯一id
	BETTING_KEY = "betting"    //疯狂投注期号 + today
	MAILID_KEY  = "LastMailID" //邮件
)

type IDGen struct {
	Id    string `bson:"_id"`
	Value string `bson:"value"`
}

func (this *IDGen) Save() bool {
	return Upsert(IDGens, bson.M{"_id": this.Id}, this)
}

func (this *IDGen) Get() {
	Get(IDGens, this.Id, this)
}

func (this *IDGen) GenID() string {
	this.Value = utils.StringAdd(this.Value)
	this.Save()
	return this.Value
}

//初始化
func InitIDGen(key string) (r *IDGen) {
	r = new(IDGen)
	r.Id = key
	r.Get()
	if r.Value == "" {
		switch key {
		case ROOMID_KEY:
			r.Value = "1"
		case USERID_KEY:
			r.Value = "10000"
		case BETTING_KEY:
			today := utils.String(utils.DayDate())
			r.Value = key + today + "000"
		}
	}
	return
}
