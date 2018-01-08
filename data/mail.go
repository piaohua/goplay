package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Items struct {
	Rtype  int    `bson:"rtype" json:"rtype"`   //类型,1钻石,2金币
	Number uint32 `bson:"number" json:"number"` //数量
}

//TODO 发一份给所有人的邮件
//TODO 离线邮件
type Mail struct {
	Id         string    `bson:"_id" json:"id"`                //邮件id
	From       string    `bson:"from" json:"from"`             //发件人
	To         string    `bson:"to" json:"to"`                 //收件人id
	Attachment []Items   `bson:"attachment" json:"attachment"` //附件
	Status     int       `bson:"status" json:"status"`         //状态0未读,1已读,2已领取,3已经过期,4已经删除
	Content    string    `bson:"content" json:"content"`       //内容
	Etime      time.Time `bson:"etime" json:"etime"`           //过期时间
	Ctime      time.Time `bson:"ctime" json:"ctime"`           //创建时间
}

func (this *Mail) Get() {
	Get(Mails, this.Id, this)
}

func (this *Mail) Save() bool {
	this.Ctime = bson.Now()
	return Insert(Mails, this)
}

func GetNewMailList(maxid, userid string) []Mail {
	var list []Mail
	ListByQLimit(Mails, bson.M{"status": 0, "to": userid,
		"id":    bson.M{"$gt": maxid},
		"etime": bson.M{"$gt": bson.Now()}}, &list, 10)
	return list
}

func (this *Mail) Del() {
	Delete(Mails, bson.M{"_id": this.Id, "to": this.To})
}

func (this *Mail) Update() {
	if Update(Mails, bson.M{"_id": this.Id},
		bson.M{"$set": bson.M{"status": this.Status}}) {
	}
}
