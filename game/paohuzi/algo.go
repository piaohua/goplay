/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-07-14 14:39:20
 * Filename      : algo.go
 * Description   : 玩牌算法
 * *******************************************************/
package paohuzi

import "sort"

func Rank(card uint32) uint32 {
	return card & RankMask
}

func Suit(card uint32) uint32 {
	return card & SuitMask
}

func Card(suit uint32, rank int) uint32 {
	return uint32(int(suit) | rank)
}

//正常流程走牌令牌移到下一家
func NextSeat(seat uint32) uint32 {
	if seat == SEAT {
		return 1
	}
	return seat + 1
}

// 移除n个牌
func Remove(c uint32, cs []uint32, n int) []uint32 {
	for n > 0 {
		for i, v := range cs {
			if c == v {
				cs = append(cs[:i], cs[i+1:]...)
				break
			}
		}
		n--
	}
	return cs
}

// 是否存在n个牌
func Exist(c uint32, cs []uint32, n int) bool {
	for _, v := range cs {
		if n == 0 {
			return true
		}
		if c == v {
			n--
		}
	}
	return n == 0
}

// 是否存在n个牌
func ExistN(c uint32, cs []uint32) (n int) {
	for _, v := range cs {
		if c == v {
			n++
		}
	}
	return
}

// 吃牌验证
func ExistChow(c uint32, cs, os, ts, hs []uint32) bool {
	m := make(map[uint32]int)
	n := make(map[uint32]int)
	for _, v := range cs {
		if v == 0 {
			return false
		}
		m[v] += 1
	}
	for _, v := range os {
		if v == 0 {
			return false
		}
		m[v] += 1
	}
	for _, v := range ts {
		if v == 0 {
			return false
		}
		m[v] += 1
	}
	m[c] -= 1
	for _, v := range hs {
		n[v] += 1
	}
	for k, v := range m {
		if v <= 0 {
			continue
		}
		if n[k] < v {
			return false
		}
	}
	return true
}

// n个红牌
func Hong(cs []uint32) (n int) {
	for _, v := range cs {
		switch v {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 1
		}
	}
	return
}

// n个红牌
func HongPong(cs []uint32) (n int) {
	for _, v := range cs {
		_, c, _ := Decode(v)
		switch c {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 3
		}
	}
	return
}

// n个红牌
func HongTi(cs []uint32) (n int) {
	for _, v := range cs {
		_, c, _ := Decode(v)
		switch c {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 4
		}
	}
	return
}

// n个红牌
func HongChow(cs []uint32) (n int) {
	for _, v := range cs {
		c1, c2, c3 := Decode(v)
		switch c1 {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 1
		}
		switch c2 {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 1
		}
		switch c3 {
		case 0x02, 0x07, 0x0a, 0x12, 0x17, 0x1a:
			n += 1
		}
	}
	return
}

// 顺子
func Shunza(cs []uint32) bool {
	if len(cs) != 3 {
		return false
	}
	AscSort(cs)
	//fmt.Println(cs)
	//一二三
	if cs[0] == (cs[1]-1) && cs[1] == (cs[2]-1) {
		return true
	}
	//二七十
	if cs[0] == 0x02 && cs[1] == 0x07 && cs[2] == 0x0a {
		return true
	}
	if cs[0] == 0x12 && cs[1] == 0x17 && cs[2] == 0x1a {
		return true
	}
	//二贰贰
	if Rank(cs[0]) == Rank(cs[1]) && Rank(cs[1]) == Rank(cs[2]) {
		if Suit(cs[0]) == Suit(cs[1]) && Suit(cs[1]) != Suit(cs[2]) {
			return true
		}
		if Suit(cs[0]) != Suit(cs[1]) && Suit(cs[1]) == Suit(cs[2]) {
			return true
		}
	}
	return false
}

//降序排序
func DescSort(cards []uint32) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] > cards[j]
	})
}

//升序排序
func AscSort(cards []uint32) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})
}

//seat 位置, card 牌值, value 掩码值
func Encode(seat, card, value uint32) uint32 {
	value = value << 16
	value |= (seat << 8)
	value |= card
	return value
}

func Decode(value uint32) (seat, card, v uint32) {
	v = value >> 16
	seat = (value >> 8) & 0xFF
	card = value & 0xFF
	return
}

//自动

//提牌
func TiCard(cards []uint32) map[uint32]int {
	m := make(map[uint32]int)
	for _, v := range cards {
		m[v] += 1
	}
	for k, v := range m {
		if v != 4 {
			delete(m, k)
		}
	}
	return m
}

//坎牌
func KanCard(cards []uint32) map[uint32]int {
	m := make(map[uint32]int)
	for _, v := range cards {
		m[v] += 1
	}
	for k, v := range m {
		if v != 3 {
			delete(m, k)
		}
	}
	return m
}

//跑牌
func PaoCard(card uint32, cards []uint32) uint32 {
	for _, v := range cards {
		if card == v {
			return PAO
		}
	}
	return 0
}

//偎牌
func WeiCard(card uint32, cards []uint32) uint32 {
	var i int
	l := len(cards)
	for _, v := range cards {
		if card == v {
			i++
		}
		if i == 2 && l >= 2 {
			return WEI
		}
	}
	return 0
}

//手动

//吃
func HasChow(card uint32, cards []uint32) (status uint32) {
	s, b := analysis(cards)
	i := int(Rank(card) - 1)
	t := Suit(card)
	addCard(i, t, s, b)
	if len(cards) > 2 {
		status |= chowCard(i, t, s, b)
	}
	return
}

//吃碰
func Discard(chow bool, card uint32, cards []uint32) (status uint32) {
	s, b := analysis(cards)
	i := int(Rank(card) - 1)
	t := Suit(card)
	if len(cards) > 2 {
		status |= pongCard(i, t, s, b)
	}
	addCard(i, t, s, b)
	if chow {
		if len(cards) > 2 {
			status |= chowCard(i, t, s, b)
		}
	}
	return
}

//只检测吃碰
func Drawcard2(chow bool, card uint32, cards []uint32) (status uint32) {
	return Discard(chow, card, cards)
}

//吃碰胡
func Drawcard(chow bool, card uint32, cards []uint32) (status, n uint32) {
	s, b := analysis(cards)
	i := int(Rank(card) - 1)
	t := Suit(card)
	if len(cards) > 2 {
		status |= pongCard(i, t, s, b)
	}
	addCard(i, t, s, b)
	//fmt.Printf("%#v\n", s)
	//fmt.Printf("%#v\n", b)
	if chow {
		//有吃
		var c []int = make([]int, TEN, TEN)
		var d []int = make([]int, TEN, TEN)
		copy(c, s)
		copy(d, b)
		//吃
		if len(cards) > 2 {
			status |= chowCard(i, t, s, b)
		}
		//fmt.Printf("%#v\n", c)
		//fmt.Printf("%#v\n", d)
		status2, n2 := huCard(len(cards)+1, c, d)
		status |= status2
		n = n2
	} else {
		//无吃
		status2, n2 := huCard(len(cards)+1, s, b)
		status |= status2
		n = n2
	}
	return
}

//胡
func HasHu(card uint32, cards []uint32) (status, n uint32) {
	s, b := analysis(cards)
	i := int(Rank(card) - 1)
	t := Suit(card)
	addCard(i, t, s, b)
	status, n = huCard(len(cards)+1, s, b)
	return
}

//胡
func Hu(cards []uint32) (status, n uint32) {
	s, b := analysis(cards)
	status, n = huCard(len(cards), s, b)
	return
}

//碰
func pongCard(i int, t uint32, s, b []int) (status uint32) {
	switch t {
	case Big:
		switch b[i] {
		case 2:
			status |= PONG
		}
	case Small:
		switch s[i] {
		case 2:
			status |= PONG
		}
	}
	return
}

//吃
func chowCard(i int, t uint32, s, b []int) (status uint32) {
	switch t {
	case Big:
		status |= chow(i, b[i], b, s)
	case Small:
		//fmt.Printf("s %v\n", s)
		//fmt.Printf("b %v\n", b)
		status |= chow(i, s[i], s, b)
	}
	if (status & CHOW) > 0 {
		if isChow(s, b) {
			status = 0 //没有可出的牌
		}
	}
	return
}

func isChow(s, b []int) bool {
	for i := 0; i < TEN; i++ {
		if b[i] > 0 {
			return false
		}
		if s[i] > 0 {
			return false
		}
	}
	return true
}

//加上打出的牌
func addCard(i int, t uint32, s, b []int) {
	switch t {
	case Big:
		b[i]++
	case Small:
		s[i]++
	}
}

//吃,i 吃牌, n 比牌, p,r 手牌
func chow(i, n int, p, r []int) (status uint32) {
	for {
		//绞(二贰贰)
		if p[i] > 1 && r[i] > 0 {
			p[i] -= 2
			r[i] -= 1
		}
		if p[i] > 0 && r[i] > 1 {
			p[i] -= 1
			r[i] -= 2
		}
		//一二三
		switch i {
		case 0:
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
			}
		case 9:
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
			}
		case 1:
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
			}
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
			}
		case 8:
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
			}
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
			}
		default:
			//fmt.Printf("p %v\n", p)
			//fmt.Printf("r %v\n", r)
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
			}
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
			}
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
			}
		}
		switch i {
		case 1:
			//二七十
			if p[i] > 0 && p[6] > 0 && p[9] > 0 {
				p[i]--
				p[6]--
				p[9]--
			}
		case 6:
			//二七十
			if p[i] > 0 && p[1] > 0 && p[9] > 0 {
				p[i]--
				p[1]--
				p[9]--
			}
		case 9:
			//二七十
			if p[i] > 0 && p[1] > 0 && p[6] > 0 {
				p[i]--
				p[1]--
				p[6]--
			}
		}
		//比牌
		if p[i] == 0 {
			status = CHOW
			return
		}
		if n == 0 {
			break
		}
		n--
	}
	return
}

func copyCards(s, b []int) (c, d []int) {
	c = make([]int, TEN, TEN)
	d = make([]int, TEN, TEN)
	copy(c, s)
	copy(d, b)
	return
}

func analysis(cards []uint32) (s, b []int) {
	s = make([]int, TEN, TEN)
	b = make([]int, TEN, TEN)
	analyse(s, b, cards)
	return
}

func analyse(s, b []int, cards []uint32) {
	for _, v := range cards {
		if Suit(v) == Big {
			b[(Rank(v) - 1)] += 1
		} else {
			s[(Rank(v) - 1)] += 1
		}
	}
}

//胡,二七十,一二三,二贰贰
func huCard(l int, s, b []int) (status, n uint32) {
	switch l % 3 {
	case 0:
		//不存在提跑
		status, n = notJiang(s, b)
		return
	case 2:
		//存在提或跑
		status, n = isJiang(s, b)
		if status > 0 {
			return
		}
	}
	return
}

func notJiang(s, b []int) (status, n uint32) {
	n = algoHu(s, b)
	//不存在提跑
	if isHu(s, b) {
		status = HU
		return
	}
	return
}

func isJiang(s, b []int) (status, n uint32) {
	for i := 0; i < TEN; i++ {
		if s[i] < 2 {
			continue
		}
		c, d := copyCards(s, b)
		c[i] -= 2
		n = algoHu(c, d)
		if isHu(c, d) {
			status = HU
			return
		}
	}
	for i := 0; i < TEN; i++ {
		if b[i] < 2 {
			continue
		}
		c, d := copyCards(s, b)
		d[i] -= 2
		n = algoHu(c, d)
		if isHu(c, d) {
			status = HU
			return
		}
	}
	return
}

//TODO loop
func algoHu(s, b []int) (n uint32) {
	//一二三
	n += shun1(s, b)
	//坎, 绞(二贰贰)
	n += jiao(s, b)
	//二七十
	n += shun2(s, b)
	return
}

func isHu(s, b []int) bool {
	for i := 0; i < TEN; i++ {
		if b[i] > 0 {
			return false
		}
		if s[i] > 0 {
			return false
		}
	}
	return true
}

//绞(二贰贰)
func jiao(s, b []int) (n uint32) {
	//坎
	for i := 0; i < TEN; i++ {
		//坎
		if b[i] == 3 {
			b[i] = 0
			n = 3 //碰胡
		}
		if s[i] == 3 {
			s[i] = 0
			n = 1 //碰胡
		}
		//绞(二贰贰)
		if b[i] > 0 && s[i] == 2 {
			b[i] -= 1
			s[i] -= 2
		}
		if s[i] > 0 && b[i] == 2 {
			s[i] -= 1
			b[i] -= 2
		}
	}
	return
}

//一二三
func shun1(s, b []int) (n uint32) {
	for j := 0; j < 2; j++ {
		//一二三
		for i := 0; i < (TEN - 2); i++ {
			if b[i] > 0 && b[i+1] > 0 && b[i+2] > 0 {
				b[i]--
				b[i+1]--
				b[i+2]--
				if i == 0 {
					n += 6
				}
			}
			if s[i] > 0 && s[i+1] > 0 && s[i+2] > 0 {
				s[i]--
				s[i+1]--
				s[i+2]--
				if i == 0 {
					n += 3
				}
			}
		}
	}
	return
}

//二七十
func shun2(s, b []int) (n uint32) {
	//二七十
	for i := 0; i < 2; i++ {
		if b[1] > 0 && b[6] > 0 && b[9] > 0 {
			b[1]--
			b[6]--
			b[9]--
			n += 6
		}
		if s[1] > 0 && s[6] > 0 && s[9] > 0 {
			s[1]--
			s[6]--
			s[9]--
			n += 3
		}
	}
	return
}

//robot

//吃
func RobotHasChow(card uint32, cards []uint32) (status uint32, cs []uint32) {
	s, b := analysis(cards)
	i := int(Rank(card) - 1)
	t := Suit(card)
	addCard(i, t, s, b)
	if len(cards) > 2 {
		status, cs = RobotChowCard(i, t, s, b)
	}
	return
}

//吃
func RobotChowCard(i int, t uint32, s, b []int) (status uint32, cs []uint32) {
	switch t {
	case Big:
		status, cs = robotChow(i, b[i], b, s, Big, Small)
	case Small:
		//fmt.Printf("s %v\n", s)
		//fmt.Printf("b %v\n", b)
		status, cs = robotChow(i, s[i], s, b, Small, Big)
	}
	return
}

//吃,i 吃牌, n 比牌, p,r 手牌
func robotChow(i, n int, p, r []int, b, s uint32) (status uint32, cs []uint32) {
	for {
		//绞(二贰贰)
		if p[i] > 1 && r[i] > 0 {
			p[i] -= 2
			r[i] -= 1
			cs = append(cs, Card(b, p[i]), Card(b, p[i]), Card(s, r[i]))
		}
		if p[i] > 0 && r[i] > 1 {
			p[i] -= 1
			r[i] -= 2
			cs = append(cs, Card(b, p[i]), Card(s, r[i]), Card(s, r[i]))
		}
		//一二三
		switch i {
		case 0:
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i+1]), Card(b, p[i+2]))
			}
		case 9:
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i-2]))
			}
		case 1:
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i+1]))
			}
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i+1]), Card(b, p[i+2]))
			}
		case 8:
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i+1]))
			}
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i-2]))
			}
		default:
			//fmt.Printf("p %v\n", p)
			//fmt.Printf("r %v\n", r)
			if p[i] > 0 && p[i-1] > 0 && p[i+1] > 0 {
				p[i]--
				p[i-1]--
				p[i+1]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i+1]))
			}
			if p[i] > 0 && p[i+1] > 0 && p[i+2] > 0 {
				p[i]--
				p[i+1]--
				p[i+2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i+1]), Card(b, p[i+2]))
			}
			if p[i] > 0 && p[i-1] > 0 && p[i-2] > 0 {
				p[i]--
				p[i-1]--
				p[i-2]--
				cs = append(cs, Card(b, p[i]), Card(b, p[i-1]), Card(b, p[i-2]))
			}
		}
		switch i {
		case 1:
			//二七十
			if p[i] > 0 && p[6] > 0 && p[9] > 0 {
				p[i]--
				p[6]--
				p[9]--
				cs = append(cs, Card(b, p[1]), Card(b, p[6]), Card(b, p[9]))
			}
		case 6:
			//二七十
			if p[i] > 0 && p[1] > 0 && p[9] > 0 {
				p[i]--
				p[1]--
				p[9]--
				cs = append(cs, Card(b, p[1]), Card(b, p[6]), Card(b, p[9]))
			}
		case 9:
			//二七十
			if p[i] > 0 && p[1] > 0 && p[6] > 0 {
				p[i]--
				p[1]--
				p[6]--
				cs = append(cs, Card(b, p[1]), Card(b, p[6]), Card(b, p[9]))
			}
		}
		//比牌
		if p[i] == 0 {
			status = CHOW
			return
		}
		if n == 0 {
			break
		}
		n--
	}
	return
}
