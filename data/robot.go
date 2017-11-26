package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Robot struct {
	Userid    string    `bson:"_id"`         // 用户id
	Nickname  string    `bson:"nickname"`    // 用户昵称
	Sex       uint32    `bson:"sex"`         // 用户性别,男1 女2 非男非女3
	Phone     string    `bson:"phone"`       // 绑定的手机号码
	Auth      string    `bson:"auth"`        // 密码验证码
	Pwd       string    `bson:"pwd"`         // MD5密码
	RegIp     string    `bson:"regist_ip"`   // 注册账户时的IP地址
	LoginIp   string    `bson:"login_ip"`    // 登录时的IP地址
	Coin      uint32    `bson:"coin"`        // 金币
	Diamond   uint32    `bson:"diamond"`     // 钻石
	RoomCard  uint32    `bson:"room_card"`   // 房卡
	Status    uint32    `bson:"status"`      // 正常1  锁定2  黑名单3
	Address   string    `bson:"address"`     // 物理地址
	Latitude  string    `bson:"latitude"`    // 纬度
	Lontitude string    `bson:"lontitude"`   // 经度
	Photo     string    `bson:"photo"`       // 头像
	Wxuid     string    `bson:"wxuid"`       // 微信uid
	Win       uint32    `bson:"win"`         // 赢
	Lost      uint32    `bson:"lost"`        // 输
	Ping      uint32    `bson:"ping"`        // 平
	Piao      uint32    `bson:"piao"`        // 漂
	Robot     bool      `bson:"robot"`       // 是否是机器人
	Money     uint32    `bson:"money"`       // 充值总金额(分)
	TopDia    uint32    `bson:"top_diamond"` // 钻石总金额
	Agent     string    `bson:"agent"`       // 代理ID
	Xftoken   string    `bson:"xftoken"`     // xftoken
	Atime     time.Time `bson:"agent_time"`  // 绑定代理时间
	Ctime     time.Time `bson:"create_time"` // 注册时间
}

func (this *Robot) Save() bool {
	return Upsert(Robots, bson.M{"_id": this.Userid}, this)
}

func (this *Robot) Get() {
	Get(Robots, this.Userid, this)
}

//读取头像不为空的数据
func GetRobots() []Robot {
	var list []Robot
	ListByQ(Robots, bson.M{"photo": bson.M{"$ne": ""}}, &list)
	return list
}
