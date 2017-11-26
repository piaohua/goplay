package algo

//seat
//(1有牛，2无牛)
//(3牛1-3， 4牛4-6，5牛7-9)
//(6牛牛，7两对，8同花，9五花牛，10四炸，11五小牛)
//
//（有牛，两对，同花）
//（牛牛，无牛）

//押注,返回中奖位置
func Betting(cs []uint32) (m map[uint32]bool, n uint32) {
	m = make(map[uint32]bool)
	//计算牛并排序
	n = Algo(cs)
	//同花
	if isFlush(cs) {
		m[8] = true
	}
	//两对
	if isPairs(cs) {
		m[7] = true
	}
	//有牛
	if n > 0 {
		m[1] = true
	}
	//牛
	switch n {
	case HgihCard:
		m[2] = true
	case Niu1, Niu2, Niu3:
		m[3] = true
	case Niu4, Niu5, Niu6:
		m[4] = true
	case Niu7, Niu8, Niu9:
		m[5] = true
	case NiuNiu:
		m[6] = true
	case FiveFlower:
		m[9] = true
	case Bomb:
		m[10] = true
	case FiveTiny:
		m[11] = true
	}
	return
}

//同花
func isFlush(cs []uint32) bool {
	if len(cs) != 5 {
		return false
	}
	if Suit(cs[0]) == Suit(cs[1]) &&
		Suit(cs[0]) == Suit(cs[2]) &&
		Suit(cs[0]) == Suit(cs[3]) &&
		Suit(cs[0]) == Suit(cs[4]) {
		return true
	}
	return false
}

//两对(有序)
func isPairs(cs []uint32) bool {
	if len(cs) != 5 {
		return false
	}
	if cs[0] == cs[1] &&
		cs[2] == cs[3] &&
		cs[0] != cs[2] {
		return true
	}
	if cs[1] == cs[2] &&
		cs[3] == cs[4] &&
		cs[1] != cs[3] {
		return true
	}
	return false
}

// 全民刮奖
//n 1有牛,2牛牛,3五花牛,4五小牛,5四炸
func Lottery(cs []uint32) (k, n uint32) {
	//计算牛并排序
	n = Algo(cs)
	//牛
	switch n {
	case HgihCard:
		k = 0
	case Niu1, Niu2, Niu3:
		k = 1
	case Niu4, Niu5, Niu6:
		k = 1
	case Niu7, Niu8, Niu9:
		k = 1
	case NiuNiu:
		k = 2
	case FiveFlower:
		k = 3
	case Bomb:
		k = 5
	case FiveTiny:
		k = 4
	}
	return
}
