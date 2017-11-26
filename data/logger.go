package data

import (
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

//1注册赠送,2开房消耗,3房间解散返还,
//4充值购买,5下注,6游戏收益,7上庄，8下庄
//8下庄, 9后台操作,10玩家赠送,11破产补助
//12存款,13取款,14取款抽成,15赠送,16赠送抽成
//17中奖,18商城购买,19绑定赠送,20首充赠送
//21转盘抽奖,22宝箱,23进入房间消耗
//24通比牛牛,25看牌抢庄,26牛牛抢庄
//27代理发放,28vip赠送,29天杠
//30刮奖,31跑胡子开房,32跑胡子进入房间消耗
//33跑胡子返还,34跑胡子,35跑胡子排位赛
//36疯狂投注,37全民刮奖
const (
	LogType1  = 1
	LogType2  = 2
	LogType3  = 3
	LogType4  = 4
	LogType5  = 5
	LogType6  = 6
	LogType7  = 7
	LogType8  = 8
	LogType9  = 9
	LogType10 = 10
	LogType11 = 11
	LogType12 = 12
	LogType13 = 13
	LogType14 = 14
	LogType15 = 15
	LogType16 = 16
	LogType17 = 17
	LogType18 = 18
	LogType19 = 19
	LogType20 = 20
	LogType21 = 21
	LogType22 = 22
	LogType23 = 23
	LogType24 = 24
	LogType25 = 25
	LogType26 = 26
	LogType27 = 27
	LogType28 = 28
	LogType29 = 29
	LogType30 = 30
	LogType31 = 31
	LogType32 = 32
	LogType33 = 33
	LogType34 = 34
	LogType35 = 35
	LogType36 = 36
	LogType37 = 37
)

//注册日志
type LogRegist struct {
	//Id       string    `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Nickname string    `bson:"nickname"`  //账户名称
	Ip       string    `bson:"ip"`        //注册IP
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	DayDate  int       `bson:"day_date"`  //regist day date
	Ctime    time.Time `bson:"ctime"`     //create Time
	Atype    uint32    `bson:"atype"`     //regist type
}

func (this *LogRegist) Save() bool {
	//this.Id = bson.NewObjectId().String()
	this.DayStamp = utils.TimestampTodayTime()
	this.DayDate = utils.DayDate()
	this.Ctime = bson.Now()
	return Insert(LogRegists, this)
}

//注册记录
func RegistRecord(userid, nickname, ip string, atype uint32) {
	record := &LogRegist{
		Userid:   userid,
		Nickname: nickname,
		Ip:       ip,
		Atype:    atype,
	}
	record.Save()
}

//登录日志
type LogLogin struct {
	//Id         string `bson:"_id"`
	Userid     string    `bson:"userid"`      //账户ID
	Event      int       `bson:"event"`       //事件：0=登录,1=正常退出,2＝系统关闭时被迫退出,3＝被动退出,4＝其它情况导致的退出
	Ip         string    `bson:"ip"`          //登录IP
	DayStamp   time.Time `bson:"day_stamp"`   //login Time Today
	LoginTime  time.Time `bson:"login_time"`  //login Time
	LogoutTime time.Time `bson:"logout_time"` //logout Time
	Atype      uint32    `bson:"atype"`       //login type
}

func (this *LogLogin) Save() bool {
	//this.Id = bson.NewObjectId().String()
	this.DayStamp = utils.TimestampTodayTime()
	this.LoginTime = bson.Now()
	return Insert(LogLogins, this)
}

func (this *LogLogin) Update(event int) bool {
	this.LogoutTime = bson.Now()
	return Update(LogLogins, bson.M{"userid": this.Userid, "event": 0},
		bson.M{"$set": bson.M{"event": event, "logout_time": this.LogoutTime}})
}

//登录记录
func LoginRecord(userid, ip string, atype uint32) {
	record := &LogLogin{
		Userid: userid,
		Event:  0,
		Ip:     ip,
		Atype:  atype,
	}
	record.Save()
}

//登录记录
func LogoutRecord(userid string, event int) {
	record := &LogLogin{
		Userid: userid,
	}
	record.Update(event)
}

//钻石日志
type LogDiamond struct {
	//Id     string `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

func (this *LogDiamond) Save() bool {
	//this.Id = bson.NewObjectId().String()
	this.Ctime = bson.Now()
	return Insert(LogDiamonds, this)
}

//钻石记录
func DiamondRecord(userid string, rtype int, rest uint32, num int32) {
	record := &LogDiamond{
		Userid: userid,
		Type:   rtype,
		Num:    num,
		Rest:   rest,
	}
	record.Save()
}

//金币日志
type LogCoin struct {
	//Id     string `bson:"_id"`
	Userid string    `bson:"userid"` //账户ID
	Type   int       `bson:"type"`   //类型
	Num    int32     `bson:"num"`    //数量
	Rest   uint32    `bson:"rest"`   //剩余数量
	Ctime  time.Time `bson:"ctime"`  //create Time
}

func (this *LogCoin) Save() bool {
	//this.Id = bson.NewObjectId().String()
	this.Ctime = bson.Now()
	return Insert(LogCoins, this)
}

//金币记录
func CoinRecord(userid string, rtype int, rest uint32, num int32) {
	record := &LogCoin{
		Userid: userid,
		Type:   rtype,
		Num:    num,
		Rest:   rest,
	}
	record.Save()
}

//绑定日志
type LogBuildAgency struct {
	//Id       string `bson:"_id"`
	Userid   string    `bson:"userid"`    //账户ID
	Agent    string    `bson:"agent"`     //绑定ID
	DayStamp time.Time `bson:"day_stamp"` //regist Time Today
	Day      int       `bson:"day"`       //regist day
	Month    int       `bson:"month"`     //regist month
	Ctime    time.Time `bson:"ctime"`     //create Time
	//DayDate   int    `bson:"day_date"`   //regist day date
	//MonthDate int    `bson:"month_date"` //regist month date
}

func (this *LogBuildAgency) Save() bool {
	//this.Id = bson.NewObjectId().Hex()
	this.DayStamp = utils.TimestampTodayTime()
	//this.DayDate = utils.DayDate()
	//this.MonthDate = utils.MonthDate()
	this.Day = utils.Day()
	this.Month = int(utils.Month())
	this.Ctime = bson.Now()
	return Insert(LogBuildAgencys, this)
}

//绑定记录
func BuildRecord(userid, agent string) {
	record := &LogBuildAgency{
		Userid: userid,
		Agent:  agent,
	}
	record.Save()
}

//在线日志
type LogOnline struct {
	//Id       string `bson:"_id"`
	Num      int       `bson:"num"`       //online count
	DayStamp time.Time `bson:"day_stamp"` //Time Today
	Ctime    time.Time `bson:"ctime"`     //create Time
}

func (this *LogOnline) Save() bool {
	//this.Id = bson.NewObjectId().Hex()
	this.DayStamp = utils.TimestampTodayTime()
	this.Ctime = bson.Now()
	return Insert(LogOnlines, this)
}

//绑定记录
func OnlineRecord(num int) {
	record := &LogOnline{
		Num: num,
	}
	record.Save()
}

//做牌日志
type LogSetHand struct {
	//Id       string `bson:"_id"`
	Rid      string    `bson:"rid"`      //房间
	Round    int       `bson:"round"`    //
	Userid   string    `bson:"userid"`   //
	Nickname string    `bson:"nickname"` //昵称
	SetHands []uint32  `bson:"sethands"` //设置手牌
	Hands    []uint32  `bson:"hands"`    //手牌
	Niu      int       `bson:"niu"`      //牌力
	Score    int32     `bson:"score"`    //得分
	Ctime    time.Time `bson:"ctime"`    //create Time
}

func (this *LogSetHand) Save() bool {
	//this.Id = bson.NewObjectId().Hex()
	this.Ctime = bson.Now()
	return Insert(LogSetHands, this)
}
