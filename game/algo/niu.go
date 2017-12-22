/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-23 21:54:36
 * Filename      : niu.go
 * Description   : 玩牌算法
 * *******************************************************/
package algo

import "sort"

/*

选庄时间显示10秒，实际时间为20秒
下注时间显示10秒，实际时间为20秒
组合手牌时间显示12秒，实际时间为30秒

流程：
开始牌局 发4张牌(正面显示) 抢庄阶段(多人抢或无人抢时随机分配)

闲家下注(超时默认1倍) 发第5张牌(正面显示) 组合手牌 推开手牌

比拼大小 积分结算 结束牌局

【基本规则】
房主在创建房间时，需要选择进入人数，最少需要3名玩家后，才能开始游戏，在第一局游戏开始后，无法进入房间
总共52张牌（无大小王）
定庄后，正面展示玩家手牌，玩家可以自己根据5张牌进行排列组合，组合后，与庄家进行比拼，决定胜负
下注时，只允许闲家下注，庄家无法下注

【比拼规则】
玩家从手牌中任意选择3张进行组合（确定有牛），
组合后，点击摊牌按钮，即打开手牌
当所有玩家摊开手牌后，则与庄家比拼大小
如果是五小或炸弹牌型，则无需组合，即可选择摊牌

牌型:(牌面分数：10、J、Q、K都为10，其他按牌面数字计算)
无牛：5张中任意3张相加不能成为10的倍数
有牛：5张中任意3张相加能成为10的倍数(牛几等于另2张相加后的个位数)
牛牛：5张中任意3张相加能成为10的倍数,另外2张相加也是10的倍数
四花：5张中1张10,另外4张为花牌(大于10的牌)
五花：5张全部为花牌
五小：5张全部小于5且相加小于10
炸弹：5张中4张相同的牌

大小比较
牌型：炸弹>五小>五花>四花>牛牛>有牛>无牛
花色：黑桃>红桃>草花>方块
单张：K>Q>J>10>9>8>7>6>5>4>3>2>A
无牛：比单张大小
有牛：比牌型大小，牛九>牛八>牛七>牛六>牛五>牛四>牛三>牛二>牛一
牛牛：比单张+花色的大小
四花：比单张+花色的大小
五花：比单张+花色的大小
五小：比点数+单张+花色的大小
炸弹：大牌吃小牌，K最大，A最小

积分计算
无牛,牛1~6:   1倍
牛7~9:        2倍
牛牛：        3倍
四花：        4倍
五花：        4倍
五小：        6倍
炸弹：        5倍

玩家胜负积分=倍数*抢庄下注额度
如：玩家A抢庄成为庄家，玩家B下注4，以牛牛牌型获得胜利，
则，玩家A需支付给玩家B的分数为： 4*3=12

100人的是只有五家
庄家，天，地，玄，黄

100人场下注达到一定额度后，超时开始

私人场，进入人数全部准备后开始(1人以上)

游戏中玩家可以进入房间(进入玩家等待)

*/

const (
	HgihCard uint32 = iota + 0x00
	Niu1
	Niu2
	Niu3
	Niu4
	Niu5
	Niu6
	Niu7
	Niu8
	Niu9
	NiuNiu
	FourFlower
	FiveFlower
	Bomb
	FiveTiny
)

// rank
const (
	Ace uint32 = iota + 0x01
	Deuce
	Trey
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King

	RankMask uint32 = 0x0F
)

// suit
const (
	Spade   uint32 = 0x40 //黑桃(♠)
	Heart   uint32 = 0x30 //红桃(♥)
	Club    uint32 = 0x20 //梅花(♣)
	Diamond uint32 = 0x10 //方块(♦)

	SuitMask uint32 = 0xF0
)

const (
	NumCard = 52
)

func Rank(card uint32) uint32 {
	return card & RankMask
}

func Suit(card uint32) uint32 {
	return card & SuitMask
}

var NiuCARDS = []uint32{
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d,
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d,
	0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d,
}

var NIUS [][]int = [][]int{{0, 1, 2}, {0, 1, 3}, {0, 2, 3}, {1, 2, 3}, {0, 1, 4}, {0, 2, 4}, {1, 2, 4}, {0, 3, 4}, {1, 3, 4}, {2, 3, 4}}
var NIUL [][]int = [][]int{{3, 4}, {2, 4}, {1, 4}, {0, 4}, {2, 3}, {1, 3}, {0, 3}, {1, 2}, {0, 2}, {0, 1}}

//[]uint32{1, 5, 8, 9, K}
func Algo(cs []uint32) uint32 {
	if len(cs) != 5 {
		return 0
	}
	descSort(cs)
	bomb_n := make(map[uint32]int)
	var tiny_n int
	var tiny_v uint32
	var flower int
	var ten int
	for _, v := range cs {
		bomb_n[Rank(v)] += 1
		switch Rank(v) {
		case Jack, Queen, King:
			flower++
		case Ten:
			ten++
		case Ace, Deuce, Trey, Four, Five, Six:
			tiny_n++
			tiny_v += Rank(v)
		}
	}
	if tiny_n == 5 && tiny_v <= Ten {
		return FiveTiny
	}
	for _, v := range bomb_n {
		if v == 4 {
			return Bomb
		}
	}
	if flower == 5 {
		return FiveFlower
	}
	//if ten == 1 && flower == 4 {
	//	return FourFlower
	//}
	var niu uint32 = HgihCard
	for k, v := range NIUS {
		if ((Trunc(cs[v[0]]) + Trunc(cs[v[1]]) + Trunc(cs[v[2]])) % 10) != 0 {
			continue
		}
		switch (Trunc(cs[NIUL[k][0]]) + Trunc(cs[NIUL[k][1]])) % 10 {
		case 0:
			return NiuNiu
		case 1:
			niu = max(niu, Niu1)
		case 2:
			niu = max(niu, Niu2)
		case 3:
			niu = max(niu, Niu3)
		case 4:
			niu = max(niu, Niu4)
		case 5:
			niu = max(niu, Niu5)
		case 6:
			niu = max(niu, Niu6)
		case 7:
			niu = max(niu, Niu7)
		case 8:
			niu = max(niu, Niu8)
		case 9:
			niu = max(niu, Niu9)
		}
	}
	return niu
}

func Trunc(n uint32) uint32 {
	if Rank(n) > Ten {
		return Ten
	}
	return Rank(n)
}

//取大值
func max(n, m uint32) uint32 {
	if n > m {
		return n
	}
	return m
}

//降序排序
func descSort(cards []uint32) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] >= cards[j]
	})
}

//比较 a == b
func Equal(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	descSort(a)
	descSort(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

//比较 a >= b
func Compare2(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	descSort(a)
	descSort(b)
	for i, v := range a {
		if v < b[i] {
			return false
		}
	}

	return true
}

type hands struct {
	Suit uint32
	Rank uint32
}

//降序排序
func descSortHands(cards []hands) {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Rank == cards[j].Rank {
			return cards[i].Suit > cards[j].Suit
		}
		return cards[i].Rank > cards[j].Rank
	})
}

//比较 a >= b (先比较牌值,牌值相同再比较花)
//同等牛的比其中牌值最大的一个,如果最大的一个牌值一样,则比花色(相同牌花色永远不同)
func Compare(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	a1 := make([]hands, len(a))
	b1 := make([]hands, len(b))
	for i, v := range a {
		a1[i].Suit = Suit(v)
		a1[i].Rank = Rank(v)
	}
	for i, v := range b {
		b1[i].Suit = Suit(v)
		b1[i].Rank = Rank(v)
	}

	descSortHands(a1)
	descSortHands(b1)
	//牌值比较
	if a1[0].Rank == b1[0].Rank {
		return a1[0].Suit > b1[0].Suit
	}
	return a1[0].Rank > b1[0].Rank
}

//比较 a >= b (先比较牌值,牌值相同再比较花)
func Compare3(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	a1 := make([]hands, len(a))
	b1 := make([]hands, len(b))
	for i, v := range a {
		a1[i].Suit = Suit(v)
		a1[i].Rank = Rank(v)
	}
	for i, v := range b {
		b1[i].Suit = Suit(v)
		b1[i].Rank = Rank(v)
	}

	descSortHands(a1)
	descSortHands(b1)
	//牌值比较
	for i, v := range a1 {
		if v.Rank < b1[i].Rank {
			return false
		} else if v.Rank > b1[i].Rank {
			return true
		}
	}

	//牌值相同,比较花色
	for i, v := range a1 {
		if v.Suit < b1[i].Suit {
			return false
		} else if v.Suit > b1[i].Suit {
			return true
		}
	}

	return true
}

//积分倍数
func Multiple(n uint32) uint32 {
	switch n {
	case HgihCard, Niu1, Niu2, Niu3, Niu4, Niu5, Niu6:
		return 1
	case Niu7, Niu8, Niu9:
		return 2
	case NiuNiu:
		return 3
	case FourFlower:
		return 4
	case FiveFlower:
		return 4
	case FiveTiny:
		return 6
	case Bomb:
		return 5
	}
	return 1
}

//[]uint32{1, 5, 8, 9, K}
func AlgoVerify(cs []uint32, val uint32) bool {
	if len(cs) != 5 {
		return false
	}
	descSort(cs)
	bomb_n := make(map[uint32]int)
	var tiny_n int
	var tiny_v uint32
	var flower int
	var ten int
	for _, v := range cs {
		bomb_n[Rank(v)] += 1
		switch Rank(v) {
		case Jack, Queen, King:
			flower++
		case Ten:
			ten++
		case Ace, Deuce, Trey, Four, Five, Six:
			tiny_n++
			tiny_v += Rank(v)
		}
	}
	if tiny_n == 5 && tiny_v <= Ten && val == FiveTiny {
		return true
	}
	for _, v := range bomb_n {
		if v == 4 && val == Bomb {
			return true
		}
	}
	if flower == 5 && val == FiveFlower {
		return true
	}
	//if ten == 1 && flower == 4 && val == FourFlower {
	//	return true
	//}
	for k, v := range NIUS {
		if ((Trunc(cs[v[0]]) + Trunc(cs[v[1]]) + Trunc(cs[v[2]])) % 10) != 0 {
			continue
		}
		switch (Trunc(cs[NIUL[k][0]]) + Trunc(cs[NIUL[k][1]])) % 10 {
		case 0:
			if val == NiuNiu {
				return true
			}
		case 1:
			if val == Niu1 {
				return true
			}
		case 2:
			if val == Niu2 {
				return true
			}
		case 3:
			if val == Niu3 {
				return true
			}
		case 4:
			if val == Niu4 {
				return true
			}
		case 5:
			if val == Niu5 {
				return true
			}
		case 6:
			if val == Niu6 {
				return true
			}
		case 7:
			if val == Niu7 {
				return true
			}
		case 8:
			if val == Niu8 {
				return true
			}
		case 9:
			if val == Niu9 {
				return true
			}
		}
	}
	return false
}

// 移除一个牌
func Remove(c uint32, cs []uint32) []uint32 {
	for i, v := range cs {
		if c == v {
			cs = append(cs[:i], cs[i+1:]...)
			break
		}
	}
	return cs
}

// 移除一个牌
func SameCard(cs, hs []uint32) bool {
	for _, c := range cs {
		for _, h := range hs {
			if c == h {
				return true
			}
		}
	}
	return false
}

// 验证手牌(设置时存在0)
func VerifyCard(cs []uint32) bool {
	for _, c := range cs {
		if c == 0 {
			continue
		}
		switch Suit(c) {
		case Club, Heart, Spade, Diamond:
		default:
			return false
		}
		switch Rank(c) {
		case Ace, Deuce, Trey, Four, Five, Six:
		case Seven, Eight, Nine, Ten, Jack, Queen, King:
		default:
			return false
		}
	}
	return true
}

//选牌
/*
无牛：5张中任意3张相加不能成为10的倍数
有牛：5张中任意3张相加能成为10的倍数(牛几等于另2张相加后的个位数)
牛牛：5张中任意3张相加能成为10的倍数,另外2张相加也是10的倍数
四花：5张中1张10,另外4张为花牌(大于10的牌)
五花：5张全部为花牌
五小：5张全部小于5且相加小于10
炸弹：5张中4张相同的牌

	HgihCard uint32 = iota + 0x00
	Niu1
	Niu2
	Niu3
	Niu4
	Niu5
	Niu6
	Niu7
	Niu8
	Niu9
	NiuNiu
	FourFlower
	FiveFlower
	Bomb
	FiveTiny

	{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	{0x1a, 0x1b, 0x1c, 0x1d}

//选一手指定牛的牌
func NiuAlgo(niu uint32) []uint32 {
	var cards []uint32 = make([]uint32, 0)
	if niu < HgihCard || niu > FiveTiny {
		return []uint32{}
	}
	var finger uint32 = 5 //选5个
	//FiveTiny
	for finger > 0 {
		var n uint32 = utils.RandUint32(5) + 1
		m = (10 - n) -
		finger--
	}
}
*/
