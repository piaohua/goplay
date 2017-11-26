package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//代理管理(代理ID为游戏内ID)
type Agency struct {
	Id         string    `bson:"_id"`         // AUTO_INCREMENT, PRIMARY KEY (`id`),
	UserName   string    `bson:"user_name"`   // 用户名, UNIQUE KEY `user_name` (`user_name`)
	Password   string    `bson:"password"`    // 密码
	Salt       string    `bson:"salt"`        // 密码盐
	Sex        int       `bson:"sex"`         // 性别
	Email      string    `bson:"email"`       // 邮箱
	LastLogin  time.Time `bson:"last_login"`  // 最后登录时间
	LastIp     string    `bson:"last_ip"`     // 最后登录IP
	Status     int       `bson:"status"`      // 状态，0正常 -1禁用
	CreateTime time.Time `bson:"create_time"` // 创建时间
	UpdateTime time.Time `bson:"update_time"` // 更新时间
	//RoleList   []Role    `bson:"role_list"`   // 角色列表
	//代理
	Phone    string    `bson:"phone"`     //绑定的手机号码(备用:非手机号注册时或多个手机时)
	Agent    string    `bson:"agent"`     //代理ID==Userid
	Level    int       `bson:"level"`     //代理等级ID:1级,2级...
	Weixin   string    `bson:"weixin"`    //微信ID
	Alipay   string    `bson:"alipay"`    //支付宝ID
	QQ       string    `bson:"qq"`        //qq号码
	Address  string    `bson:"address"`   //详细地址
	Number   uint32    `bson:"number"`    //当前余额
	Expend   uint32    `bson:"expend"`    //总消耗
	Cash     float32   `bson:"cash"`      //当前可提取额
	Extract  float32   `bson:"extract"`   //已经提取额
	CashTime time.Time `bson:"cash_time"` //提取指定时间前所有
}

func (this *Agency) Get(userid string) {
	GetByQ(Agencys, bson.M{"agent": userid}, this)
}
