package data

import "gopkg.in/mgo.v2/bson"

//设置变量
//key             value
const (
	ENV1  = "regist_diamond"    //注册赠送钻石
	ENV2  = "regist_coin"       //注册赠送金币
	ENV3  = "build"             //绑定赠送
	ENV4  = "first_pay_multi"   //首充送n倍
	ENV5  = "first_pay_coin"    //首充送金币
	ENV6  = "relieve"           //救济金次数
	ENV7  = "prizedraw"         //转盘抽奖次数
	ENV8  = "bankrupt_coin"     //破产金额
	ENV9  = "relieve_coin"      //救济金额
	ENV10 = "free_random"       //百人场概率
	ENV11 = "free_count"        //百人场人数
	ENV12 = "prize_card"        //刮奖牌
	ENV13 = "prize_coin"        //刮奖金币
	ENV14 = "prize_diamond"     //刮奖钻石
	ENV15 = "lottery_diamond"   //全民刮奖
	ENV16 = "lottery_maxnumber" //全民刮奖最大注
	ENV17 = "sys_carry"         //系统上庄限额
	ENV18 = "free_carry"        //庄家上庄限额
	ENV19 = "limit_carry"       //庄家下庄限额
	ENV20 = "dealer_times"      //做庄次数限制
	ENV21 = "limit_sit"         //玩家坐下限额
	ENV22 = "prize_limit"       //中奖下注限额
)

type Env struct {
	Key   string `bson:"_id" json:"key"`     //key
	Value int32  `bson:"value" json:"value"` //value
}

func GetEnvList() []Env {
	var list []Env
	ListByQ(Envs, nil, &list)
	return list
}

func (this *Env) DelEnv() bool {
	return Delete(Envs, this)
}

func (this *Env) SetEnv() bool {
	return Upsert(Envs, bson.M{"_id": this.Key}, this)
}
