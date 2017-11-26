/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-23 21:54:36
 * Filename      : niu.go
 * Description   : 玩牌算法
 * *******************************************************/
package algo

import (
	"testing"
	"utils"
)

/*
type Card uint32

//const NilCard = Card(0)

func (card Card) Rank() int {
	return int(card & RankMask)
}

func (card Card) Suit() int {
	return int(card & SuitMask)
}
*/

// 测试
func TestAlgo(t *testing.T) {
	//var c Card = 0x1b
	//t.Logf("%x, %x", c.Suit(), c.Rank())
	t.Log(Compare([]uint32{0x32, 0x4a, 0x44, 0x22, 0x21},
		[]uint32{0x4a, 0x44, 0x32, 0x22, 0x11}))
	t.Log(0x19&0xF0, 0x19)
	t.Logf("%x, %x", 0x19>>4, 0x19&0xF0)
	t.Logf("%x", 0x40|4)
	t.Logf("%x", 0x49&0xF0)
	t.Logf("%x, %x", 0x1d, King)
	t.Logf("%d, %d", 0x1d, King)
	t.Log(0x4b & RankMask)
	t.Log(0x4b + 0x21 + 0x48)
	t.Log((0x4a + 0x22 + 0x48) % 10)
	//t.Log((0x4a + 0x22 + 0x48) & RemMask)
	//cards := []int{0x4b, 0x4b, 0x4d, 0x4d, 0x4b}
	cards := []uint32{0x1c, 0x19, 0x14, 0x15, 0x17}
	t.Log(Algo(cards))
	t.Log([]uint32(cards))
	cs := []uint32{0x13, 0x23, 0x1a, 0x2a, 0x1d}
	t.Log(Algo(cs))
	t.Log(AlgoVerify(cs, 0x06))
	cs = []uint32{0x4d, 0x32, 0x36, 0x16, 0x18}
	t.Log(AlgoVerify(cs, 0x02))
	cs = []uint32{0x31, 0x25, 0x21, 0x12, 0x11}
	t.Log(Algo(cs))
	cs = []uint32{19, 67, 51, 35, 23}
	cs = []uint32{0x13, 0x43, 0x33, 0x23, 0x17}
	//t.Log(Algo(cs))
	cs = []uint32{0x13, 0x13, 0x13, 0x13, 0x17}
	t.Log(Algo(cs))
}

func TestNius(t *testing.T) {
	//所有有牛组合
	var nius [][]int = [][]int{{1, 1, 8}, {1, 2, 7}, {1, 3, 6}, {1, 4, 5},
		{2, 3, 5}, {2, 2, 6}, {2, 4, 4}, {3, 3, 4}, {10, 9, 1}, {10, 8, 2},
		{10, 7, 3}, {10, 6, 4}, {10, 5, 5}, {9, 8, 3}, {9, 7, 4}, {9, 6, 5},
		{8, 7, 5}, {8, 8, 4}, {8, 6, 6}, {7, 7, 6}, {10, 10, 10}}
	t.Log(nius)
	//5选3组合
	var niu_zuhe [][]int = [][]int{{1, 1, 1, 0, 0}, {1, 1, 0, 1, 0},
		{1, 0, 1, 1, 0}, {0, 1, 1, 1, 0}, {1, 1, 0, 0, 1}, {1, 0, 1, 0, 1},
		{0, 1, 1, 0, 1}, {1, 0, 0, 1, 1}, {0, 1, 0, 1, 1}, {0, 0, 1, 1, 1}}
	t.Log(niu_zuhe)
}

// 测试
func TestCom(t *testing.T) {
	t.Log(Compare2([]uint32{0x32, 0x4a, 0x44, 0x22, 0x22},
		[]uint32{0x4a, 0x44, 0x32, 0x22, 0x31}))
	m := map[int]int{1: 1, 2: 2, 3: 3}
	t.Log(m)
	for k, _ := range m {
		delete(m, k)
	}
	t.Log(m)
}

// 测试
func TestPare(t *testing.T) {
	t.Log(Compare([]uint32{0x3c, 0x33, 0x1c, 0x1a, 0x17},
		[]uint32{0x4a, 0x35, 0x2c, 0x23, 0x22}))
	t.Log([]uint32{0x3c, 0x33, 0x1c, 0x1a, 0x17})
	t.Log([]uint32{0x4a, 0x35, 0x2c, 0x23, 0x22})
}

// 测试
func TestBomb(t *testing.T) {
	cs := []uint32{19, 67, 51, 35, 23}
	cs = []uint32{0x13, 0x43, 0x33, 0x23, 0x17}
	//cs = []uint32{0x13, 0x13, 0x13, 0x13, 0x17}
	cs = []uint32{0x3d, 0x2d, 0x1d, 0x4b, 0x4d}
	cs = []uint32{0x4d, 0x2c, 0x28, 0x39, 0x31}
	cs = []uint32{0x43, 0x3d, 0x39, 0x32, 0x4c}
	t.Log(Algo(cs))
	gs := []uint32{0x4d, 0x2c, 0x28, 0x39, 0x31}
	gs = []uint32{0x43, 0x3d, 0x39, 0x32, 0x4c}
	t.Log(AlgoVerify(gs, 0))
	endTime := utils.LocalTime().Unix()
	t.Log(endTime)
	t.Log(utils.Timestamp())
	t.Log(utils.LocalTime())
	t.Log(utils.Minute())
	t.Log(utils.Hour())
}
