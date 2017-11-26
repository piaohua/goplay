package paohuzi

import (
	"fmt"
	"testing"
)

// 测试
func TestHu(t *testing.T) {
	//cards := []uint32{0x01, 0x01, 0x11, 0x03, 0x04, 0x05,
	//	0x11, 0x12, 0x13, 0x12, 0x13}
	cards := []uint32{0x14, 0x14, 0x0a, 0x0a, 0x1a, 0x06, 0x06, 0x16}
	t.Log(cards)
	status, n := Hu(cards)
	t.Log(status, n)
}

// 测试
func TestHasHu(t *testing.T) {
	//cards := []uint32{0x01, 0x01, 0x11, 0x03, 0x04, 0x05,
	//	0x11, 0x12, 0x13, 0x12, 0x13}
	//cards := []uint32{0x14, 0x14, 0x0a, 0x0a, 0x1a, 0x06, 0x06, 0x16}
	//cards := []uint32{0x12, 0x13, 0x19, 0x19}
	cards := []uint32{0x01, 0x02, 0x03, 0x18}
	t.Log(cards)
	status, n := HasHu(0x18, cards)
	t.Log(status, n)
}

func TestDraw(t *testing.T) {
	cards := []uint32{0x01, 0x01, 0x02, 0x03, 0x02, 0x03,
		0x04, 0x05, 0x11, 0x11, 0x12, 0x13, 0x13}
	t.Log(cards)
}

func TestRemove(t *testing.T) {
	cards := []uint32{0x01, 0x01, 0x02, 0x03, 0x02, 0x03,
		0x04, 0x05, 0x11, 0x11, 0x12, 0x13, 0x13}
	t.Log(cards)
	cards = Remove(0x01, cards, 2)
	t.Log(cards)
}

func TestAlgo(t *testing.T) {
	//cards := []uint32{0x01, 0x03, 0x04, 0x05,
	//	0x04, 0x05, 0x11, 0x01, 0x22, 0x13, 0x14}
	//t.Log(Drawcard(true, 0x11, cards))
	//cards := []uint32{0x04, 0x04, 0x14,
	//	0x0a, 0x1a, 0x1a, 0x17, 0x17}
	cards := []uint32{0x01, 0x11, 0x11,
		0x12, 0x12, 0x0a, 0x0a}
	t.Log(Drawcard(false, 0x12, cards))
}

func TestDiscard(t *testing.T) {
	cards := []uint32{0x01, 0x11, 0x11,
		0x12, 0x12, 0x0a, 0x0a}
	t.Log(Discard(true, 0x12, cards))
}

func Test2Drawcard(t *testing.T) {
	cards := []uint32{0x09, 0x19, 0x19,
		0x02, 0x02}
	t.Log(Drawcard2(true, 0x12, cards))
}

func TestWei(t *testing.T) {
	cards := []uint32{0x19, 0x19}
	t.Log(WeiCard(0x19, cards))
}

func TestEncode(t *testing.T) {
	v := Encode(1, 0x01, PONG)
	t.Log(v)
	s, c, val := Decode(v)
	t.Log(s, c, val)
}

func TestDecode(t *testing.T) {
	s, c, val := Decode(0x80107)
	t.Log(s, c, val)
}

func TestExistChow(t *testing.T) {
	hs := []uint32{19, 24, 2, 23, 26, 7, 6, 6, 4, 25, 26, 18, 25, 24, 22, 4, 18, 7, 8, 20}
	cs := []uint32{21, 20, 19}
	t.Log(Shunza(cs))
	os := []uint32{}
	ts := []uint32{}
	t.Log(ExistChow(21, cs, os, ts, hs))
}

// 测试
func TestMask(t *testing.T) {
	var v uint32 = 16388
	//胡牌方式
	if v&CHOW > 0 {
		t.Log("CHOW -> ", v)
	}
	if v&PONG > 0 {
		t.Log("PONG -> ", v)
	}
	if v&HU > 0 {
		t.Log("HU -> ", v)
	}
	if v&WEI > 0 {
		t.Log("WEI -> ", v)
	}
	if v&CHOU_HU > 0 {
		t.Log("CHOU_HU -> ", v)
	}
	if v&PAO > 0 {
		t.Log("PAO -> ", v)
	}
	if v&TI > 0 {
		t.Log("TI -> ", v)
	}
	if v&PING_HU > 0 {
		t.Log("PING_HU -> ", v)
	}
	if v&ZIMO_HU > 0 {
		t.Log("ZIMO_HU -> ", v)
	}
	if v&PAO_HU > 0 {
		t.Log("PAO_HU -> ", v)
	}
	if v&CHOU_HU > 0 {
		t.Log("CHOU_HU -> ", v)
	}
	if v&TIAN_HU > 0 {
		t.Log("TIAN_HU -> ", v)
	}
	if v&DI_HU > 0 {
		t.Log("DI_HU -> ", v)
	}
	if v&HONG_HU > 0 {
		t.Log("HONG_HU -> ", v)
	}
	if v&DIAN_HU > 0 {
		t.Log("DIAN_HU -> ", v)
	}
	if v&HONG_WU > 0 {
		t.Log("HONG_WU -> ", v)
	}
	if v&WU_HU > 0 {
		t.Log("WU_HU -> ", v)
	}
}

func TestChow(t *testing.T) {
	//hs := []uint32{19, 7, 1, 26, 8, 2, 21, 9, 20, 3, 18, 2, 7, 10, 4, 4, 1, 26, 22}
	//cs := []uint32{3, 3, 19}
	//hs := []uint32{5, 22, 10, 25, 18, 21, 17, 26, 7, 1, 1, 18}
	//cs := []uint32{18, 18, 2}
	hs := []uint32{6, 21, 4, 4, 20, 18, 8, 21, 1, 6, 5, 19, 23, 9, 20, 10, 22, 1, 8, 17}
	//cs := []uint32{21, 21, 5}
	cs := []uint32{21, 5, 5}
	os := []uint32{6, 5, 4}
	ts := []uint32{}
	t.Log(chowOperate2(0x05, hs, cs, os, ts))
	//n := ExistN(9, cs)
	//t.Log(n)
}

func TestHasChow(t *testing.T) {
	//cards := []uint32{0x06, 0x16, 0x16, 0x01, 0x11, 0x11, 0x05,
	//	0x06, 0x07, 0x12, 0x12, 0x18, 0x1a, 0x04, 0x14, 0x09, 0x03}
	//cards := []uint32{0x09, 0x19, 0x19, 0x06, 0x16, 0x16, 0x04,
	//	0x14, 0x14, 0x01, 0x01, 0x11, 0x17, 0x18, 0x12, 0x02}
	cards := []uint32{0x15, 0x15, 0x19, 0x19}
	t.Log(HasChow(0x05, cards))
}

//吃牌,比牌操作
func chowOperate2(card uint32, hs, cards, bione, bitwo []uint32) bool {
	//没有牌
	if len(cards) != 3 {
		fmt.Println("1 ", 1)
		return false
	}
	//手牌中是否存在
	if !ExistChow(card, cards, bione, bitwo, hs) {
		fmt.Println("2 ", 1)
		return false
	}
	//比牌正常
	//顺子
	if !Shunza(cards) {
		fmt.Println("4 ", 1)
		return false
	}
	if !Exist(card, cards, 1) {
		return false
	}
	if len(bione) == 3 {
		if !Shunza(bione) {
			return false
		}
		if !Exist(card, bione, 1) {
			return false
		}
	}
	if len(bitwo) == 3 {
		if !Shunza(bitwo) {
			return false
		}
		if !Exist(card, bitwo, 1) {
			return false
		}
	}
	return true
}
