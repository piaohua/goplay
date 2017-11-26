/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-07-04 18:50:50
 * Filename      : kong.go
 * Description   : 玩牌算法
 * *******************************************************/
package kong

import "sort"

/*

选庄时间显示10秒，实际时间为20秒
下注时间显示10秒，实际时间为20秒
组合手牌时间显示12秒，实际时间为30秒

流程：
开始牌局 发1张牌(正面显示) 抢庄阶段(多人抢或无人抢时随机分配)

闲家下注(超时默认1倍) 发第2张牌(正面显示) 组合手牌 推开手牌

比拼大小 积分结算 结束牌局

【基本规则】
无鬼牌和10;
散牌情况下，K、Q、J的点数按10算；
最少需要2名玩家参与，最多可支持6名玩家进行游戏；
允许玩家中途加入游戏（房间人数未满时可中途加入）；
游戏中，桌面上会出现一张“公众牌”，而每位玩家会有2张手牌，并按手牌组合后，各闲家与庄家进行比牌

【比拼规则】
散牌：两个牌数值相加取个位数算点(花牌算10，1+9算0点)
豹子：一对相同牌组成
天杠：固定2，8组成

花色比较：
黑桃>红心>梅花>方块

单张比较：
K>Q>J>9>8>7>6>5>4>3>2>A

同牌型比较:
天杠：按数字8的花色进行比较；
如果数字8是公众牌，则按数字2的花色进行比较

豹子：先按牌值比较（由于是对子，所以跟单张的比较方法一样）；
如果牌值也是相同的，则按花色进行比较

散牌：先按点数比较（9最大，0最小，别把天杠当0了！）；
如果点数相同，则按最大牌的牌值比较；
如果最大牌牌值相同，则按最大牌花色比较

积分计算
散牌：        1倍
豹子：        2倍
天杠：        3倍

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
	ONEPAIRS
	TWOEIGHT
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
	Spade   uint32 = 0x40 //黑桃
	Heart   uint32 = 0x30 //红桃
	Club    uint32 = 0x20 //梅花
	Diamond uint32 = 0x10 //方块

	SuitMask uint32 = 0xF0
)

const (
	NumCard = 48
)

func Rank(card uint32) uint32 {
	return card & RankMask
}

func Suit(card uint32) uint32 {
	return card & SuitMask
}

var NiuCARDS = []uint32{
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1b, 0x1c, 0x1d,
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2b, 0x2c, 0x2d,
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3b, 0x3c, 0x3d,
	0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4b, 0x4c, 0x4d,
}

var NIUS [][]int = [][]int{{0, 1}, {0, 2}, {1, 2}}

func Algo(hs []uint32, card uint32) (uint32, uint32, uint32) {
	cs := make([]uint32, len(hs), len(hs))
	copy(cs, hs)
	cs = append(cs, card)
	return algo2(cs)
}

//[]int{1, 5, K}
func algo2(cs []uint32) (niu, c, s uint32) {
	if len(cs) != 3 {
		return 0, 0, 0
	}
	descSort(cs)
	for _, v := range NIUS {
		if Rank(cs[v[0]]) == Deuce && Rank(cs[v[1]]) == Eight {
			return TWOEIGHT, cs[v[0]], cs[v[1]]
		}
		if Rank(cs[v[1]]) == Deuce && Rank(cs[v[0]]) == Eight {
			return TWOEIGHT, cs[v[1]], cs[v[0]]
		}
	}
	for _, v := range NIUS {
		if Rank(cs[v[0]]) == Rank(cs[v[1]]) {
			return ONEPAIRS, cs[v[0]], cs[v[1]]
		}
	}
	//var niu uint32 = HgihCard
	for _, v := range NIUS {
		n := (Trunc(cs[v[0]]) + Trunc(cs[v[1]])) % 10
		if n > niu {
			niu = n
			c = cs[v[0]]
			s = cs[v[1]]
		}
	}
	return niu, cs[0], cs[1]
}

func Trunc(n uint32) uint32 {
	if Rank(n) > Ten {
		return Ten
	}
	return Rank(n)
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

//比较{a, b} > {c, d}
func Compare(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) == 0 || len(b) == 0 {
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

//积分倍数
func Multiple(n uint32) uint32 {
	switch n {
	case ONEPAIRS:
		return 3
	case TWOEIGHT:
		return 2
	default:
		return 1
	}
	return 1
}

//[]int{1, 5, K}
func AlgoVerify(cs []uint32, val uint32) bool {
	if len(cs) != 3 {
		return false
	}
	descSort(cs)
	for _, v := range NIUS {
		if Rank(cs[v[0]]) == Deuce && Rank(cs[v[1]]) == Eight &&
			val == TWOEIGHT {
			return true
		}
		if Rank(cs[v[1]]) == Deuce && Rank(cs[v[0]]) == Eight &&
			val == TWOEIGHT {
			return true
		}
	}
	for _, v := range NIUS {
		if Rank(cs[v[0]]) == Rank(cs[v[1]]) && val == ONEPAIRS {
			return true
		}
	}
	var niu uint32 = HgihCard
	for _, v := range NIUS {
		n := (Trunc(cs[v[0]]) + Trunc(cs[v[1]])) % 10
		if n > niu {
			niu = n
		}
	}
	if niu == val {
		return true
	}
	return false
}
