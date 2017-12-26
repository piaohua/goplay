/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 12:50:31
 * Filename      : sender.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"crypto/md5"
	"encoding/hex"

	"goplay/game/algo"
	"goplay/game/kong"
	"goplay/game/paohuzi"
	"goplay/glog"
	"goplay/pb"
	"utils"
)

//' 登录
// 发送注册请求
func (c *Robot) SendRegist() {
	ctos := &pb.CRegist{}
	ctos.Phone = c.data.Phone
	ctos.Nickname = c.data.Nickname
	h := md5.New()
	passwd := cfg.Section("robot").Key("passwd").Value()
	h.Write([]byte(passwd)) // 需要加密的字符串为
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = pwd
	c.Sender(ctos)
}

// 发送登录请求
func (c *Robot) SendLogin() {
	ctos := &pb.CLogin{}
	ctos.Phone = c.data.Phone
	h := md5.New()
	passwd := cfg.Section("robot").Key("passwd").Value()
	h.Write([]byte(passwd)) // 需要加密的字符串为
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = pwd
	//utils.Sleep(4)
	glog.Infof("ctos -> %#v", ctos)
	c.Sender(ctos)
}

// 获取玩家数据
func (c *Robot) SendUserData() {
	ctos := &pb.CUserData{}
	ctos.Userid = c.data.Userid
	c.Sender(ctos)
}

// 玩家创建房间
func (c *Robot) SendCreate() {
	ctos := &pb.CCreateRoom{}
	ctos.Round = 8
	ctos.Rtype = 1
	ctos.Ante = 1
	ctos.Count = 4
	ctos.Payment = 0
	ctos.Rname = "ddd"
	utils.Sleep(3)
	c.Sender(ctos)
}

// 玩家进入房间
func (c *Robot) SendEntry() {
	if c.code == "create" { //表示创建房间
		c.SendCreate() //创建一个房间
	} else { //表示进入房间
		ctos := &pb.CEnterRoom{}
		ctos.Invitecode = c.code
		c.Sender(ctos)
	}
}

// 玩家准备
func (c *Robot) SendReady() {
	ctos := &pb.CReady{}
	ctos.Ready = true
	utils.Sleep(5)
	c.Sender(ctos)
}

// 离开
func (c *Robot) SendLeave() {
	ctos := &pb.CLeave{}
	c.Sender(ctos)
}

// 解散
func (c *Robot) SendVote() {
	ctos := &pb.CVote{}
	ctos.Vote = 0
	c.Sender(ctos)
}

//.

//' 牛牛
// 抢庄
func (c *Robot) SendDealer() {
	ctos := &pb.CDealer{}
	ctos.Dealer = true
	num := uint32(utils.RandInt32N(3) + 1) //随机
	ctos.Num = num
	utils.Sleep(3)
	c.Sender(ctos)
}

// 下注
func (c *Robot) SendBet() {
	ctos := &pb.CBet{}
	ctos.Seatbet = c.seat
	val := uint32(utils.RandInt32N(3)) + 1
	ctos.Value = val
	utils.Sleep(3)
	c.Sender(ctos)
}

// 提交组合
func (c *Robot) SendNiu() {
	ctos := &pb.CNiu{}
	ctos.Cards = c.cards
	val := algo.Algo(c.cards)
	ctos.Value = val
	utils.Sleep(3)
	c.Sender(ctos)
}

// 提交组合
func (c *Robot) SendNiu2() {
	ctos := &pb.CNiu{}
	ctos.Cards = c.cards
	ctos.Cards = append(ctos.Cards, c.card)
	ctos.Cards = append(ctos.Cards, 0x0)
	ctos.Cards = append(ctos.Cards, 0x0)
	//val := algo.Algo(c.cards)
	val, _, _ := kong.Algo(c.cards, c.card)
	ctos.Value = val
	utils.Sleep(3)
	c.Sender(ctos)
}

// 提交组合
func (c *Robot) SendRecord() {
	ctos := &pb.CGameRecord{
		Page: 0,
	}
	utils.Sleep(3)
	c.Sender(ctos)
}

//.

//' free

// 进入房间
func (c *Robot) SendEntryFree() {
	ctos := &pb.CEnterFreeRoom{}
	utils.Sleep(2)
	c.Sender(ctos)
}

// 玩家入坐
func (c *Robot) SendFreeSit() {
	seat := uint32(utils.RandInt32N(7) + 1) //随机
	ctos := &pb.CFreeSit{
		State: true,
		Seat:  seat,
	}
	//c.sits++ //尝试次数
	utils.Sleep(2)
	c.Sender(ctos)
}

// 玩家离坐
func (c *Robot) SendFreeStandup() {
	ctos := &pb.CFreeSit{
		State: false,
		Seat:  c.seat,
	}
	utils.Sleep(2)
	c.Sender(ctos)
	utils.Sleep(2)
	c.SendLeave()
	utils.Sleep(2)
	c.Close() //下线
}

// 玩家下注
func (c *Robot) SendFreeBet() {
	var a1 []uint32 = []uint32{2, 3, 4, 5}
	var c1 []uint32 = []uint32{100, 1000, 10000, 50000, 100000, 200000}
	var coin uint32 = c.data.Coin / 4
	var n1 int32 = utils.RandInt32N(15) //随机
	if coin > 10000000 {
		n1 += 10
	}
	for {
		if n1 <= 0 {
			break
		}
		if coin <= 0 {
			break
		}
		var i2 int
		for i := 5; i >= 0; i-- {
			if coin >= c1[i] {
				i2 = i
				break
			}
		}
		var val int
		switch i2 {
		case 0, 4, 5:
			val = i2
		default:
			val = int(utils.RandInt32N(int32(i2))) + 1 //随机
		}
		var i1 int32 = utils.RandInt32N(4) //随机
		ctos := &pb.CFreeBet{
			Value: c1[val],
			Seat:  a1[i1],
		}
		//c.bits++
		utils.Sleep(1)
		c.Sender(ctos)
		//
		n1--
		coin -= c1[val]
	}
}

// 玩家聊天
func (c *Robot) SendChat() {
	if c.rtype == 0 {
		return
	}
	if utils.RandInt32N(10) > 4 {
		return
	}
	content := chat[utils.RandInt32N(int32(len(chat)))]
	ctos := &pb.CChatText{
		Content: content,
	}
	utils.Sleep(3)
	c.Sender(ctos)
}

var chat []string = []string{
	"收金币，200W金币150，微伈：13699755455",
	"卖金币，200W金币160，微伈：13699755455",
	"不要走，决战到天亮",
	"又赢了",
	"哈哈",
	"哈哈哈",
	"nnd",
	"hahahaha",
	"搞什么飞机",
	"这运气",
	"这都行",
	"也是没谁了",
	"怎么搞的，又输，没金币了",
	"输惨了，有没有土豪送点金币",
	"还好，只输了一点",
	"这次一定是要赢",
	"连输2把，下次一定要赢回来",
	"又赢了，手气不错",
	"哈哈，这样都能赢",
	"连胜4把，快超神了",
	"今天要连赢19把",
	"赢得豪爽",
	"这把牌简直逆天了",
	"这游戏赢得豪爽啊",
	"哈哈哈哈",
	"你们继续，我先撤了",
	"晚上继续搞起",
	"我靠，这牌，这么大！",
	"我靠，好大！",
	"牌小也能赢，牛逼",
}

//.

//' classic

// 进入房间
func (c *Robot) SendEntryClassic() {
	ctos := &pb.CEnterClassicRoom{}
	for _, v := range c.classic {
		min := v.GetMinimum()
		max := v.GetMaximum()
		if c.data.Coin > min && (c.data.Coin < max || max == 0) {
			c.classicId = v.GetId()
			break
		}
	}
	ctos.Id = c.classicId
	utils.Sleep(3)
	c.Sender(ctos)
}

// 进入房间
func (c *Robot) SendGetClassic() {
	ctos := &pb.CClassicList{}
	c.Sender(ctos)
}

// 玩家下注
func (c *Robot) SendClassicBet() {
	val := uint32(utils.RandInt32N(3)) + 1 //随机
	ctos := &pb.CBet{
		Value:   val,
		Seatbet: c.seat,
	}
	//c.bits++
	utils.Sleep(1)
	c.Sender(ctos)
}

func (c *Robot) SendClassicClose() {
	c.classicId = ""
	c.SendLeave()
	c.Close() //下线
}

//.

//' 跑胡子

// 玩家进入房间
func (c *Robot) SendEntryPhz() {
	if c.code == "create" { //表示创建房间
		c.SendCreatePhz() //创建一个房间
	} else { //表示进入房间
		glog.Infof("enter phz room -> %s, %d", c.code, c.rtype)
		ctos := &pb.CEnterZiRoom{}
		ctos.Invitecode = c.code
		c.Sender(ctos)
	}
}

// 玩家创建房间
func (c *Robot) SendCreatePhz() {
	ctos := &pb.CCreateZiRoom{}
	ctos.Round = 8
	ctos.Rtype = 6
	ctos.Ante = 1
	ctos.Count = 3
	ctos.Payment = 0
	ctos.Rname = "phz"
	ctos.Xi = 10
	utils.Sleep(3)
	c.Sender(ctos)
}

// 玩家操作(吃碰胡)
func (c *Robot) SendOperatePhz(value uint32, cards, bione, bitwo []uint32) {
	ctos := &pb.COperate{}
	ctos.Value = value //掩码值
	ctos.Cards = cards //吃牌
	ctos.Bione = bione //比牌
	ctos.Bitwo = bitwo //比牌
	utils.Sleep(3)
	c.Sender(ctos)
}

// 玩家出牌
func (c *Robot) SendDiscardPhz(card uint32) {
	ctos := &pb.CPushDiscard{}
	ctos.Card = card
	utils.Sleep(3)
	c.Sender(ctos)
}

func (c *Robot) SendPhzClose() {
	c.classicId = ""
	c.SendLeave()
	c.Close() //下线
}

//胡碰吃
func (c *Robot) phzOperate(card, value uint32) {
	if value&paohuzi.HU > 0 {
		cs := []uint32{}
		c.SendOperatePhz(paohuzi.HU, cs, cs, cs)
	} else if value&paohuzi.PONG > 0 {
		cs := []uint32{}
		c.SendOperatePhz(paohuzi.PONG, cs, cs, cs)
	} else if value&paohuzi.CHOW > 0 {
		//TODO BUG 选出吃牌
		v, cs := paohuzi.RobotHasChow(card, c.cards)
		if v > 0 {
			var cs1, cs2, cs3 []uint32
			if len(cs) >= 3 {
				cs1 = cs[:3]
			}
			if len(cs) >= 6 {
				cs2 = cs[3:6]
			}
			if len(cs) >= 9 {
				cs3 = cs[6:]
			}
			//c.SendOperatePhz(paohuzi.CHOW, cs1, cs2, cs3)
			//TODO BUG FIXME
			c.SendOperatePhz(0, cs1, cs2, cs3)
		}
	}
}

//自动操作(提,跑,偎)
func (c *Robot) phzAuto(card, value uint32) {
	if value&paohuzi.TI > 0 {
		c.cards = paohuzi.Remove(card, c.cards, 4)
	} else if value&paohuzi.PAO > 0 {
		c.cards = paohuzi.Remove(card, c.cards, 3)
	} else if value&paohuzi.WEI > 0 ||
		value&paohuzi.CHOU_WEI > 0 {
		c.cards = paohuzi.Remove(card, c.cards, 2)
	}
}

//移除吃碰组合
func (c *Robot) phzChow(card, value uint32, cards, bione, bitwo []uint32) {
	c.cards = paohuzi.Remove(card, c.cards, 4)
	for _, v := range cards {
		c.cards = paohuzi.Remove(v, c.cards, 1)
	}
	for _, v := range bione {
		c.cards = paohuzi.Remove(v, c.cards, 1)
	}
	for _, v := range bitwo {
		c.cards = paohuzi.Remove(v, c.cards, 1)
	}
}

//出牌
func (c *Robot) phzDiscard() {
	//TODO 选择最优出牌
	c.SendDiscardPhz(c.cards[0])
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
