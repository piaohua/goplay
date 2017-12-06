package data

import (
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

const (
	BANKRUPT      uint32  = 2000      //破产补助限制和金额
	BANKRUPT_TIME uint32  = 3         //每天破产补助次数
	DRAW_MONEY    uint32  = 10        //提现和赠送最低金额限制
	GIVE_PERCENT  float64 = 0.1       //赠送抽成
	GIVE_LIMIT    uint32  = 100000000 //赠送上限
	TAX_NUMBER    uint32  = 100       //小于这个数抽成为1
)

type User struct {
	Userid       string     `bson:"_id" json:"userid"`                  // 用户id
	Nickname     string     `bson:"nickname" json:"nickname"`           // 用户昵称
	Sex          uint32     `bson:"sex" json:"sex"`                     // 用户性别,男1 女2 非男非女3
	Phone        string     `bson:"phone" json:"phone"`                 // 绑定的手机号码
	Auth         string     `bson:"auth" json:"auth"`                   // 密码验证码
	Pwd          string     `bson:"pwd" json:"pwd"`                     // MD5密码
	RegIp        string     `bson:"regist_ip" json:"reg_ip"`            // 注册账户时的IP地址
	Coin         uint32     `bson:"coin" json:"coin"`                   // 金币
	Diamond      uint32     `bson:"diamond" json:"diamond"`             // 钻石
	RoomCard     uint32     `bson:"room_card" json:"room_card"`         // 房卡
	Status       uint32     `bson:"status" json:"status"`               // 正常1  锁定2  黑名单3
	Address      string     `bson:"address" json:"address"`             // 物理地址
	Photo        string     `bson:"photo" json:"photo"`                 // 头像
	Wxuid        string     `bson:"wxuid" json:"wxuid"`                 // 微信uid
	Win          uint32     `bson:"win" json:"win"`                     // 赢
	Lost         uint32     `bson:"lost" json:"lost"`                   // 输
	Ping         uint32     `bson:"ping" json:"ping"`                   // 平
	Robot        bool       `bson:"robot" json:"robot"`                 // 是否是机器人
	Money        uint32     `bson:"money" json:"money"`                 // 充值总金额(分)
	TopDia       uint32     `bson:"top_diamond" json:"top_dia"`         // 钻石总金额
	Bank         uint32     `bson:"bank" json:"bank"`                   // 个人银行
	GiveCoin     uint32     `bson:"give_coin" json:"give_coin"`         // 当天赠送金额
	Bankrupts    uint32     `bson:"bankrupts" json:"bankrupts"`         // 破产次数
	Agent        string     `bson:"agent" json:"agent"`                 // 代理ID
	Atype        uint32     `bson:"atype" json:"atype"`                 // 代理登录包类型
	Atime        time.Time  `bson:"agent_time" json:"atime"`            // 绑定代理时间
	Ctime        time.Time  `bson:"create_time" json:"ctime"`           // 注册时间
	GiveTime     time.Time  `bson:"give_time" json:"give_time"`         // 赠送时间
	BankruptTime time.Time  `bson:"bankrupt_time" json:"bankrupt_time"` // 破产补助时间
	PrizeDraw    uint32     `bson:"prize_draw" json:"prize_draw"`       // 抽奖次数
	PrizeTime    time.Time  `bson:"prize_time" json:"prize_time"`       // 抽奖时间
	Duration     uint32     `bson:"duration" json:"duration"`           // 在线时长
	Box          string     `bson:"box" json:"box"`                     // 领取的id
	BoxState     int        `bson:"box_state" json:"box_state"`         // 领取完成
	BoxTime      time.Time  `bson:"box_time" json:"box_time"`           // 在线时间
	Vip          uint32     `bson:"vip" json:"vip"`                     // vip充值金额
	VipLevel     int        `bson:"vip_level" json:"vip_level"`         // vip等级,充值时改变
	KickTimes    int        `bson:"kick_times" json:"kick_times"`       // vip已经踢人次数
	TodayTime    time.Time  `bson:"today_time" json:"today_time"`       // 凌晨重置时间标识
	DanCombat    *DanCombat `bson:"-" json:"dan_combat"`                // 竞技场
	DanTask      *DanTask   `bson:"-" json:"dan_task"`                  // 竞技场
}

func (this *User) Save() bool {
	return Upsert(PlayerUsers, bson.M{"_id": this.Userid}, this)
}

func (this *User) Get() {
	Get(PlayerUsers, this.Userid, this)
}

func (this *User) GetById(userid string) {
	GetByQ(PlayerUsers, bson.M{"_id": userid}, this)
}

func (this *User) UpdatePhoto() bool {
	return UpdateByQField(PlayerUsers, bson.M{"_id": this.Userid}, "photo", this.Photo)
}

func (this *User) UpdateAgent(agent string) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$set": bson.M{"agent": agent}})
}

func (this *User) UpdateDiamond(num int32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$inc": bson.M{"diamond": num}})
}

func (this *User) UpdateCoin(num int32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$inc": bson.M{"coin": num}})
}

func (this *User) UpdateBank(num int32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$inc": bson.M{"bank": num}})
}

func (this *User) UpdatePay(diamond, money, coin uint32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid},
		bson.M{"$inc": bson.M{"diamond": diamond, "top_diamond": diamond, "money": money, "coin": coin}})
}

func (this *User) UpdateVip(level int, money uint32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid},
		bson.M{"$set": bson.M{"vip_level": level}, "$inc": bson.M{"vip": money}})
}

func (this *User) UpsertVip(level int, money uint32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid},
		bson.M{"$set": bson.M{"vip_level": level, "vip": money}})
}

func (this *User) SetDiamond(num uint32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$set": bson.M{"diamond": num}})
}

func (this *User) SetCoin(num uint32) bool {
	return Update(PlayerUsers, bson.M{"_id": this.Userid}, bson.M{"$set": bson.M{"coin": num}})
}

func (this *User) ExistsPhone() bool {
	return Has(PlayerUsers, bson.M{"phone": this.Phone})
}

func (this *User) GetByPhone() {
	GetByQ(PlayerUsers, bson.M{"phone": this.Phone}, this)
}

func (this *User) GetByWechat() {
	GetByQ(PlayerUsers, bson.M{"wxuid": this.Wxuid}, this)
}

func (this *User) VerifyPwdByPhone(pwd string) bool {
	if this.Phone == "" {
		return false
	}
	this.GetByPhone()
	if this.Userid == "" {
		return false
	}
	return this.VerifyPwd(pwd)
}

//用户登陆密码验证
func (this *User) PWDIsOK(userid, pwd string) bool {
	if userid == "" {
		return false
	}
	this.Get()
	if this.Userid == "" {
		return false
	}
	return this.VerifyPwd(pwd)
}

//密码验证
func (this *User) VerifyPwd(pwd string) bool {
	return utils.Md5(pwd+this.Auth) == this.Pwd
}

//

func (this *User) UserSave() {
	this.Save()
}

func (this *User) GetUserid() string {
	return this.Userid
}

func (this *User) GetDiamond() uint32 {
	return this.Diamond
}

func (this *User) AddDiamond(num int32) {
	diamond := int32(this.Diamond) + num
	if diamond < 0 {
		diamond = 0
	}
	this.Diamond = uint32(diamond)
	if num > 0 {
		this.TopDia += uint32(num)
	}
}

func (this *User) AddMoney(num uint32) {
	this.Money += num
}

func (this *User) GetMoney() uint32 {
	return this.Money
}

func (this *User) GetCoin() uint32 {
	return this.Coin
}

func (this *User) AddCoin(num int32) {
	coin := int32(this.Coin) + num
	if coin < 0 {
		coin = 0
	}
	this.Coin = uint32(coin)
}

func (this *User) GetBank() uint32 {
	return this.Bank
}

func (this *User) AddBank(num int32) {
	coin := int32(this.Bank) + num
	if coin < 0 {
		coin = 0
	}
	this.Bank = uint32(coin)
}

func (this *User) GetGive() uint32 {
	dayStamp := utils.TimestampTodayTime()
	if this.GiveTime.Before(dayStamp) { //是否昨天
		this.GiveCoin = 0
	}
	return this.GiveCoin
}

func (this *User) AddGive(num uint32) {
	dayStamp := utils.TimestampTodayTime()
	if this.GiveTime.Before(dayStamp) { //是否昨天
		this.GiveCoin = 0
	}
	this.GiveTime = utils.LocalTime()
	this.GiveCoin += num
}

func (this *User) GetName() string {
	return this.Nickname
}

func (this *User) GetPhoto() string {
	return this.Photo
}

func (this *User) GetPhone() string {
	return this.Phone
}

func (this *User) SetAgent(agent string) {
	this.Agent = agent
	this.Atime = utils.LocalTime()
}

func (this *User) GetAgent() string {
	return this.Agent
}

func (this *User) GetRoomCard() uint32 {
	return this.RoomCard
}

func (this *User) GetSex() uint32 {
	return this.Sex
}

func (this *User) GetAtype() uint32 {
	return this.Atype
}

func (this *User) SetAtype(atype uint32) {
	this.Atype = atype
}

func (this *User) GetBankrupts() uint32 {
	dayStamp := utils.TimestampTodayTime()
	if this.BankruptTime.Before(dayStamp) { //是否昨天
		this.Bankrupts = 0
	}
	return this.Bankrupts
}

func (this *User) SetBankrupts() {
	this.Bankrupts += 1
	this.BankruptTime = utils.LocalTime()
}

func (this *User) GetPrizeDraw() uint32 {
	dayStamp := utils.TimestampTodayTime()
	if this.PrizeTime.Before(dayStamp) { //是否昨天
		this.PrizeDraw = 0
	}
	return this.PrizeDraw
}

func (this *User) SetPrizeDraw() {
	this.PrizeDraw += 1
	this.PrizeTime = utils.LocalTime()
}

func (this *User) GetBox() string {
	return this.Box
}

func (this *User) SetBox(box string, state int) {
	this.Box = box
	this.BoxState = state
}

func (this *User) GetBoxState() int {
	return this.BoxState
}

func (this *User) SetDuration() {
	this.BoxTime = utils.LocalTime()
	this.Duration = 0
}

func (this *User) GetDuration() uint32 {
	if this.BoxTime.IsZero() {
		this.BoxTime = utils.LocalTime()
	}
	n := utils.LocalTime().Unix() - this.BoxTime.Unix()
	if n > 0 {
		this.Duration += uint32(n)
	}
	return this.Duration
}

func (this *User) GetRegIp() string {
	return this.RegIp
}

func (this *User) SetRecord(value int) {
	if value > 0 {
		this.Win += 1
	} else if value < 0 {
		this.Lost += 1
	} else {
		this.Ping += 1
	}
}

func (this *User) GetVipLevel() int {
	return this.VipLevel
}

func (this *User) GetVip() uint32 {
	return this.Vip
}

func (this *User) SetVip(level int, num uint32) {
	this.VipLevel = level
	this.Vip += num
}

func (this *User) GetKickTimes() int {
	if this.TodayTime != utils.TimestampTodayTime() {
		this.KickTimes = 0
	}
	return this.KickTimes
}

//TODO TodayTime 优化
func (this *User) SetKickTimes(num int) {
	this.TodayTime = utils.TimestampTodayTime()
	this.KickTimes += num
}

func (this *User) setLoginTime(in bool) {
	if in {
		dayStamp := utils.TimestampTodayTime()
		if this.BoxTime.Before(dayStamp) { //是否昨天
			this.Duration = 0
			this.Box = ""
			this.BoxState = 0
			this.BoxTime = utils.LocalTime()
		}
	} else {
		this.GetDuration()
	}
}

func (this *User) GetDan() int {
	return this.DanCombat.Dan
}

func (this *User) GetDanCombat() interface{} {
	return this.DanCombat
}

func (this *User) GetDanTask() interface{} {
	return this.DanTask
}
