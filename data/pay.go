//充值交易记录
package data

import (
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

const (
	TradeSuccess = 0 //交易成功
	TradeFail    = 1 //交易失败
	Tradeing     = 2 //交易中(下单状态)
	TradeGoods   = 3 //发货失败
)

// 交易记录
type TradeRecord struct {
	Id        string    `bson:"_id" json:"id"`              //商户订单号(游戏内自定义订单号)
	Transid   string    `bson:"transid" json:"transid"`     //交易流水号(计费支付平台的交易流水号,微信订单号)
	Userid    string    `bson:"userid" json:"userid"`       //用户在商户应用的唯一标识(userid)
	Itemid    string    `bson:"itemid" json:"itemid"`       //购买商品ID
	Amount    string    `bson:"amount" json:"amount"`       //购买商品数量
	Diamond   uint32    `bson:"diamond" json:"diamond"`     //购买钻石数量
	Money     uint32    `bson:"money" json:"money"`         //交易总金额(单位为分)
	Transtime string    `bson:"transtime" json:"transtime"` //交易完成时间 yyyy-mm-dd hh24:mi:ss
	Result    int       `bson:"result" json:"result"`       //交易结果(0–交易成功,1–交易失败,2-交易中,3-发货中)
	Waresid   uint32    `bson:"waresid" json:"waresid"`     //商品编码(平台为应用内需计费商品分配的编码)
	Currency  string    `bson:"currency" json:"currency"`   //货币类型(RMB,CNY)
	Transtype int       `bson:"transtype" json:"transtype"` //交易类型(0–支付交易)
	Feetype   int       `bson:"feetype" json:"feetype"`     //计费方式(表示商品采用的计费方式)
	Paytype   uint32    `bson:"paytype" json:"paytype"`     //支付方式(表示用户采用的支付方式,403-微信支付)
	Clientip  string    `bson:"clientip" json:"clientip"`   //客户端ip
	Agent     string    `bson:"agent" json:"agent"`         //绑定的父级代理商游戏ID
	Atype     uint32    `bson:"atype" json:"atype"`         //代理包类型
	First     int       `bson:"first" json:"first"`         //首次充值
	Utime     time.Time `bson:"utime" json:"utime"`         //本条记录更新unix时间戳
	DayStamp  time.Time `bson:"day_stamp" json:"day_stamp"` //Time Today
	Ctime     time.Time `bson:"ctime" json:"ctime"`         //本条记录生成unix时间戳
}

// 生成订单id,(时间截+角色id)
func GenCporderid(userid string) string {
	return utils.Base62encode(uint64(utils.TimestampNano())) + userid
}

func GenOrderid() string {
	return bson.NewObjectId().Hex()
}

// 交易结果记录
func (this *TradeRecord) Get() {
	Get(TradeRecords, this.Id, this)
}

func (this *TradeRecord) Has() bool {
	return Has(TradeRecords, bson.M{"_id": this.Id})
}

func (this *TradeRecord) GetByTransid(transid string) {
	GetByQ(TradeRecords, bson.M{"transid": transid}, this)
}

func (this *TradeRecord) Update() bool {
	this.Utime = bson.Now()
	return Update(TradeRecords, bson.M{"_id": this.Id}, this)
}

func (this *TradeRecord) Save() bool {
	this.Ctime = bson.Now()
	return Insert(TradeRecords, this)
}

func (this *TradeRecord) Upsert() bool {
	this.Utime = bson.Now()
	return Upsert(TradeRecords, bson.M{"_id": this.Id}, this)
}

/*
func (this *TradeRecord) Delete() bool {
	return Delete(TradeRecords, bson.M{"_id": this.Id})
}
*/

// 获取某玩家的所有离线订单,用于上线补单
func GetTradeOff(userid string) []*TradeRecord {
	var list []*TradeRecord
	ListByQ(TradeRecords, bson.M{"userid": userid, "result": TradeFail}, &list)
	return list
}
