//微信登录
package data

import (
	"time"
	"utils"

	"gopkg.in/mgo.v2/bson"
)

// 微信登录数据
type WXLogin struct {
	OpenId       string    `bson:"_id"`           //OpenId
	UnionId      string    `bson:"unionid"`       //unionid
	AccessToken  string    `bson:"access_token"`  //access_token
	RefreshToken string    `bson:"refresh_token"` //refresh_token
	ATExpiresIn  int64     `bson:"at_expires_in"` //access_token过期时间(2小时)
	RTExpiresIn  int64     `bson:"rt_expires_in"` //refresh_token过期时间(30天)
	Utime        time.Time `bson:"update_time"`   //更新本条记录unix时间戳
	Ctime        time.Time `bson:"create_time"`   //生成本条记录unix时间戳
}

// 生成id,(时间截+随机字符串)
func GenWXLoginID() string {
	return utils.Base62encode(uint64(utils.TimestampNano())) +
		utils.Base62encode(uint64(utils.RandUint32()))
}

// 交易结果记录
func (this *WXLogin) Get() {
	Get(WxLogins, this.OpenId, this)
}

func (this *WXLogin) GetByToken() {
	GetByQ(WxLogins, bson.M{"access_token": this.AccessToken}, this)
}

func (this *WXLogin) Update() bool {
	this.Utime = bson.Now()
	this.Ctime = bson.Now()
	return Upsert(WxLogins, bson.M{"_id": this.OpenId}, this)
}

func (this *WXLogin) Set(access_token string, expires_in int64) bool {
	return Update(WxLogins, bson.M{"_id": this.OpenId},
		bson.M{"$set": bson.M{
			"AccessToken": access_token,
			"ATExpiresIn": expires_in,
			"Utime":       bson.Now(),
		}})
}

func (this *WXLogin) Save() bool {
	return Insert(WxLogins, this)
}

// TODO:定时刷新refresh_token

type WxLoginData struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	OpenId       string   `json:"openid"`
	Nickname     string   `json:"nickname"`
	Sex          int      `json:"sex"`
	Province     string   `json:"province"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
	HeadImagUrl  string   `json:"headimgurl"`
	Privilege    []string `json:"privilege"`
	UnionId      string   `json:"unionid"`
}
