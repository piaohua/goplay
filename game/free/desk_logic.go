/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-30 16:21:11
 * Filename      : desk.go
 * Description   : 玩牌逻辑
 * *******************************************************/
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"goplay/data"
	"goplay/game/algo"
	"goplay/game/config"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//房间牌桌数据结构
type DeskFree struct {
	id    string         //房间id
	data  *data.DeskData //房间类型基础数据
	cards []uint32       //没摸起的海底牌
	//-------------
	players map[string]*data.User //房间无座玩家
	seats   map[string]uint32     //userid:seat (seat:1~8)
	//-------------
	state uint32 //房间状态,0准备中,1游戏中
	timer int    //计时
	//closeCh chan bool //关闭通道
	//-------------
	round     uint32                      //局数
	pond      uint32                      //奖池
	dealerNum uint32                      //做庄次数
	dealer    string                      //庄家
	carry     uint32                      //庄家的携带,小于一定值时下庄,字段只做记录,真实数据直接写入玩家数据
	num       uint32                      //当前局下注总数
	dealers   []map[string]uint32         //上庄列表,userid: carry
	bets      map[string]uint32           //userid:num, 玩家下注金额
	seatBets  map[uint32]uint32           //userid:num, 玩家下注金额
	tian      map[string]uint32           //天,seat:value
	di        map[string]uint32           //地
	xuan      map[string]uint32           //玄
	huang     map[string]uint32           //黄
	handCards map[uint32][]uint32         //手牌 seat:cards,seat=(1,2,3,4,5)
	power     map[uint32]uint32           //位置(1-5)对应牌力
	multiple  map[uint32]int32            //结果 seat:num,seat=(1,2,3,4,5),倍数
	score     map[uint32]int32            //位置(1-5)输赢总量
	score2    map[string]int32            //每个闲家输赢总量
	score3    map[uint32]map[string]int32 //位置(1-5)上每个玩家输赢
	score4    map[uint32]uint32           //位置(1-5)分到奖池金额
	score5    map[string]uint32           //玩家分到奖池金额
	//-------------
	//userid-playerPid
	pids map[string]*actor.PID
	//playerPid-userid
	router map[string]string
	//关闭通道
	stopCh chan struct{}
	//name
	Name string
}

/*
得到奖金的条件:
1.五花牛 炸弹 五小牛才能分一定比例的奖池(5%,10%,20%)
还需要有庄、天、地、玄、黄中有五花牛、炸弹、五小牛的那一家下注达到一定数量（比如300万，庄家需携带300万）
奖池暂定每盘抽取5%
坐位中的人可否上庄
*/

const (
	STATE_FREE_READY  = 0  //准备状态
	STATE_FREE_DEALER = 1  //休息中状态
	STATE_FREE_BET    = 2  //下注中状态
	FREE_DT           = 20 //下注超时时间
	FREE_RT           = 19 //休息超时时间
	//
	SYS_CARRY    uint32 = 50000000 //系统上庄限额
	FREE_CARRY   uint32 = 100000   //庄家上庄限额
	LIMIT_CARRY  uint32 = 100000   //庄家下庄限额
	DEALER_TIMES uint32 = 8        //做庄次数限制
	LIMIT_SIT    uint32 = 300000   //玩家坐下限额
	PRIZE_LIMIT  uint32 = 5000000  //中奖下注限额
	//
	SEAT1 uint32 = 1 //庄
	SEAT2 uint32 = 2 //天
	SEAT3 uint32 = 3 //地
	SEAT4 uint32 = 4 //玄
	SEAT5 uint32 = 5 //黄
)

//// external function

//新建一张牌桌
func NewDeskFree(deskData *data.DeskData) *DeskFree {
	desk := &DeskFree{
		id:      deskData.Rid,
		data:    deskData,
		players: make(map[string]*data.User),
		//------
		seats:   make(map[string]uint32),
		dealers: make([]map[string]uint32, 0),
		//closeCh: make(chan bool, 1),
		//
		pids:   make(map[string]*actor.PID),
		router: make(map[string]string),
		stopCh: make(chan struct{}),
		Name:   deskData.Rid,
	}
	desk.gameInit() //初始化
	return desk
}

//初始化
func (t *DeskFree) gameInit() {
	t.num = 0                                    //
	t.bets = make(map[string]uint32)             //
	t.seatBets = make(map[uint32]uint32)         //
	t.tian = make(map[string]uint32)             //
	t.di = make(map[string]uint32)               //
	t.xuan = make(map[string]uint32)             //
	t.huang = make(map[string]uint32)            //
	t.handCards = make(map[uint32][]uint32)      //手牌
	t.multiple = make(map[uint32]int32)          //倍数
	t.score = make(map[uint32]int32)             //位置(1-5)输赢总量
	t.score2 = make(map[string]int32)            //个人输赢结果userid: value
	t.score3 = make(map[uint32]map[string]int32) //个人输赢结果userid: value
	t.score4 = make(map[uint32]uint32)           //位置(1-5)分到奖池金额
	t.score5 = make(map[string]uint32)           //分到奖池金额
	t.power = make(map[uint32]uint32)
}

//初始化
func (t *DeskFree) gameOverInit() {
	t.gameInit()
	t.gameStart()
}

//初始化
func (t *DeskFree) gameStart() {
	if len(t.players) > 0 {
		t.state = STATE_FREE_DEALER //休息停顿
	} else {
		t.state = STATE_FREE_READY //准备
	}
	t.gameStartMsg()
}

func (t *DeskFree) gameStartMsg() {
	photo := ""
	p := t.getPlayer(t.dealer)
	if p != nil {
		photo = p.GetPhoto()
	}
	var left uint32 = t.leftDealerTimes()
	//if t.dealer == "" {
	//	//left = 1
	//	//dealerNum = 1
	//} else {
	//	left = t.leftDealerTimes()
	//	//dealerNum = DEALER_TIMES
	//}
	msg := handler.FreeStartMsg(t.dealer, photo,
		t.state, t.carry, DEALER_TIMES, left)
	t.broadcast(msg)
}

func (t *DeskFree) gameStartMsgL() {
	t.checkDealer()          //不足做庄
	t.state = STATE_FREE_BET //下注
	t.gameStartMsg()
}

func (t *DeskFree) leftDealerTimes() uint32 {
	return DEALER_TIMES - t.dealerNum
}

//不足做庄,或者检测是否有人上庄
func (t *DeskFree) checkDealer() {
	if t.dealer != "" {
		if t.carry < LIMIT_CARRY || t.leftDealerTimes() == 0 {
			p := t.players[t.dealer]
			t.delBeDealer(t.dealer, p)
		}
	}
	t.beComeDealer()    //成为庄家
	if t.dealer == "" { //无人坐庄
		//if t.carry < LIMIT_CARRY {
		//	t.carry = SYS_CARRY //庄家3百万
		//}
		t.carry = SYS_CARRY //庄家3百万,庄家每次都补庄
		//重置次数
		if t.leftDealerTimes() == 0 {
			t.dealerNum = 0
		}
	}
}

//成为庄家
func (t *DeskFree) beComeDealer() {
	if t.dealer != "" {
		return
	}
	i, userid, num := t.findBeDealer()
	if userid == "" || num < FREE_CARRY {
		return
	}
	p := t.players[userid]
	if p == nil {
		return
	}
	//成为庄家
	t.dealer = userid
	t.carry = num
	t.dealerNum = 0
	t.dealers = append(t.dealers[:i], t.dealers[i+1:]...)
	//t.leaveSeat(userid, false, p)
	//msg := handler.FreeBeDealerMsg(1, num, t.dealer, userid, p.GetName())
	//t.broadcast(msg)
}

//携带最大的优先做庄
func (t *DeskFree) findBeDealer() (int, string, uint32) {
	var index int
	var userid string
	var maxNum uint32
	for i, m := range t.dealers {
		for k, v := range m {
			if v > maxNum {
				index = i
				userid = k
				maxNum = v
			}
		}
	}
	return index, userid, maxNum
}

////获取牌桌数据
//func (t *DeskFree) GetData() interface{} {
//	return t.data
//}
//
////获取牌桌数据
//func (t *DeskFree) GetType() uint32 {
//	return t.data.Rtype
//}
//
////获取牌桌数据
//func (t *DeskFree) GetLtype() uint32 {
//	return t.data.Ltype
//}
//
////房间满人
//func (t *DeskFree) GetFull() bool {
//	return uint32(len(t.players)) < t.data.Count
//}
//
////房间满人
//func (t *DeskFree) GetPlayers() uint32 {
//	return uint32(len(t.players))
//}
//
////房间状态
//func (t *DeskFree) GetState() bool {
//	return t.state == STATE_FREE_READY
//}

////关闭房间,停服or过期清除等等.TODO:玩牌中是否关闭?
//func (t *DeskFree) Closed(ok bool) {
//	t.close() //ok=true强制解散,=false清理
//}

//房间消息广播,聊天
//func (t *DeskFree) Broadcasts(mtype int, userid, msg string) {
//	seat := t.seats[userid]
//	if mtype == 1 {
//		t.broadcast(handler.ChatMsg2(seat, userid, msg))
//	} else {
//		t.broadcast(handler.ChatMsg(seat, userid, msg))
//	}
//}

//进入
func (t *DeskFree) Enter(p *data.User) pb.ErrCode {
	if _, ok := t.players[p.GetUserid()]; !ok {
		if uint32(len(t.players)) >= t.data.Count {
			return pb.RoomFull //人数已满
		}
	}
	t.players[p.GetUserid()] = p
	//t.reEnter(p) //检测重复进入
	msg2 := t.freecamein(p)
	t.broadcast(msg2)
	//TODO //p.SetRoom(t)
	return pb.OK
}

//玩家离开,下注也可以离开
func (t *DeskFree) Leave(userid string) pb.ErrCode {
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	//庄家下注时不能离开
	if t.state == STATE_FREE_BET && userid == t.dealer {
		return pb.GameStartedCannotLeave
	}
	p := t.players[userid]
	if t.dealer == userid && userid != "" {
		t.delBeDealer(userid, p)
	}
	t.leaveBeDealer(userid, p) //清空上庄列表
	//广播消息
	msg := handler.LeaveMsg(userid, t.seats[userid])
	t.broadcast(msg)
	//清除数据
	delete(t.seats, userid)
	delete(t.players, userid)
	//p.SetRoom(nil) //清除玩家房间数据
	return pb.OK
}

//离开位置
func (t *DeskFree) leaveSeat(userid string, st bool, p *data.User) {
	if seat, ok := t.seats[userid]; ok {
		delete(t.seats, userid)
		//广播消息
		msg := handler.FreeSitMsg(userid, seat, st, p) //离开坐位
		t.broadcast(msg)
	}
}

//0下庄 1上庄 2补庄
func (t *DeskFree) addBeDealer(userid string, st, num uint32, p *data.User) {
	if userid == t.dealer && st == 2 { //庄家补庄
		t.carry += num
	} else {
		if t.state != STATE_FREE_BET && t.dealer == "" && st == 2 { //系统做庄
			//成为庄家
			t.dealer = userid
			t.carry = num
			//t.leaveSeat(userid, false, p)
		} else if st == 1 {
			m := map[string]uint32{userid: num}
			t.dealers = append(t.dealers, m)
		}
	}
	t.sendFreeCoin(p, (-1 * int32(num)), data.LogType7)
	msg := handler.FreeBeDealerMsg(st, num, t.dealer, userid, p.GetName())
	t.broadcast(msg)
}

func (t *DeskFree) delBeDealer(userid string, p *data.User) {
	var num uint32
	if t.dealer == userid && t.dealer != "" {
		num = t.carry
		t.sendFreeCoin(p, int32(num), data.LogType8)
		t.dealer = ""
		t.carry = 0
	}
	t.dealerNum = 0
	msg := handler.FreeBeDealerMsg(0, num, t.dealer, userid, p.GetName())
	t.broadcast(msg)
}

//离开房间返还上庄列表
func (t *DeskFree) leaveBeDealer(userid string, p *data.User) {
	for {
		had := true
		for i, m := range t.dealers {
			if num, ok := m[userid]; ok {
				t.sendFreeCoin(p, int32(num), data.LogType8)
				delete(m, userid)
				t.dealers = append(t.dealers[:i], t.dealers[i+1:]...)
				had = false
				break
			}
		}
		if had {
			break
		}
	}
}

//玩家离开
func (t *DeskFree) SitDown(userid string, seat uint32, st bool) pb.ErrCode {
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	if _, ok := t.seats[userid]; ok && st { //坐下
		return pb.SitDownFailed
	}
	if _, ok := t.seats[userid]; !ok && !st { //站起
		return pb.StandUpFailed
	}
	//if userid == t.dealer { //庄家不能坐
	//	return pb.DealerSitFailed
	//}
	if st {
		for _, s := range t.seats {
			if s == seat {
				return pb.SitDownFailed
			}
		}
		t.seats[userid] = seat
	} else {
		delete(t.seats, userid)
	}
	//glog.Infof("SitDown -> %s, %d, %v", userid, seat, st)
	//广播消息
	p := t.players[userid]
	msg := handler.FreeSitMsg(userid, seat, st, p)
	t.broadcast(msg)
	return pb.OK
}

//抢庄,没人上庄时都可以选择上庄,可以多次上庄，已经上庄的人可以补庄
//0下庄 1上庄 2补庄
func (t *DeskFree) BeDealer(userid string, st, num uint32) pb.ErrCode {
	if num < FREE_CARRY && st == 1 {
		return pb.BeDealerNotEnough
	}
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	//庄家下注时不能下庄
	if t.state == STATE_FREE_BET && userid == t.dealer && st == 0 {
		return pb.GameStartedCannotLeave
	}
	p := t.players[userid]
	if st == 1 || st == 2 {
		t.addBeDealer(userid, st, num, p)
	} else {
		t.delBeDealer(userid, p)
	}
	return pb.OK
}

//上庄列表
//func (t *DeskFree) BeDealerList(userid string) {
//	msg := t.res_bedealerlist()
//	p := t.players[userid]
//	//glog.Infof("dealer list free room -> %s", userid)
//	p.Send(msg)
//}

//下注
func (t *DeskFree) ChoiceBet(userid string, seatBet, num uint32) pb.ErrCode {
	if userid == t.dealer { //庄家不用下注
		return pb.BetDealerFailed
	}
	if t.state != STATE_FREE_BET {
		return pb.GameNotStart
	}
	if _, ok := t.players[userid]; !ok {
		return pb.NotInRoom
	}
	if (t.num + num) > (t.carry / 4) { //下注不能大于庄家携带1/4
		return pb.BetTopLimit //下注限制
	}
	p := t.getPlayer(userid)
	coin := p.GetCoin()                         //剩余金额
	a_bets := t.bets[userid]                    //已经下注额
	if (num + a_bets) > ((coin + a_bets) / 4) { //本轮下注不能超过1/4
		return pb.BetTopLimit //下注限制
	}
	//glog.Infof("ChoiceBet -> %s, %d, %d", userid, seatBet, num)
	seat := t.seats[userid]
	t.bets[userid] += num      //个人总下注额
	t.seatBets[seatBet] += num //当前位置总下注额
	t.num += num               //当局总下注额
	var betsNum uint32
	switch seatBet {
	case SEAT2:
		t.tian[userid] += num
		betsNum = t.tian[userid]
	case SEAT3:
		t.di[userid] += num
		betsNum = t.di[userid]
	case SEAT4:
		t.xuan[userid] += num
		betsNum = t.xuan[userid]
	case SEAT5:
		t.huang[userid] += num
		betsNum = t.huang[userid]
	}
	t.sendFreeCoin(p, (-1 * int32(num)), data.LogType5)
	msg := handler.FreeBetMsg(seat, seatBet, num,
		t.seatBets[seatBet], betsNum, userid)
	t.broadcast(msg)
	return pb.OK
}

//结束牌局,ok=true投票解散
func (t *DeskFree) close() {
	//下庄
	if t.dealer != "" {
		p := t.players[t.dealer]
		t.delBeDealer(t.dealer, p)
	}
	for k, p := range t.players {
		t.leaveBeDealer(k, p) //清空上庄列表
		msg := handler.LeaveMsg(k, t.seats[k])
		t.broadcast(msg)
		//TODO p.SetRoom(nil) //清除玩家房间数据
		//glog.Infof("desk close userid : %s", k)
	}
	//if t.closeCh != nil {
	//	t.closeCh <- true
	//	close(t.closeCh) //关闭计时器
	//	t.closeCh = nil  //消除计时器
	//}
	//TODO
	//go func(key string) {
	//	rooms.DelFree(key) //从房间列表中清除
	//}(t.data.Code)
}

//计时器
func (t *DeskFree) tickerHandler() {
	//超时判断
	if t.state != STATE_FREE_READY {
		t.timerL() //逻辑处理
	} else {
		t.timerLStart()
	}
}

//超时处理
func (t *DeskFree) timerLStart() {
	//开始
	t.gameStart()
}

//超时处理
func (t *DeskFree) timerL() {
	if t.timer == FREE_DT && t.state == STATE_FREE_BET { //结束下注
		//出牌超时处理
		t.timer = 0
		t.gameOver()
	} else if t.timer == FREE_RT && t.state == STATE_FREE_DEALER { //开始下注
		t.timer = 0
		t.gameStartMsgL()
	} else {
		t.timer++
	}
}

//洗牌
func (t *DeskFree) shuffle() {
	rand.Seed(time.Now().UnixNano())
	d := make([]uint32, algo.NumCard, algo.NumCard)
	copy(d, algo.NiuCARDS)
	//测试暂时去掉洗牌
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	t.cards = d
}

//发牌,直接发5张,百人场发固定位置
func (t *DeskFree) deal() {
	var hand int = 5
	for i := SEAT1; i <= SEAT5; i++ {
		cards := make([]uint32, hand, hand)
		tmp := t.cards[:hand]
		copy(cards, tmp)
		t.handCards[i] = cards
		t.cards = t.cards[hand:]
		msg := handler.DrawMsg(i, t.state, cards)
		t.broadcast(msg)
	}
}

//牌力计算
func (t *DeskFree) paili() {
	if t.dealer != "" {
		return
	}
	//配置概率
	r := config.GetEnv(data.ENV10)
	if r <= 0 || r >= 100 {
		return
	}
	if (utils.RandInt32N(100) + 1) > r { //百分比
		return
	}
	cs1 := t.getHandCards(SEAT1)   //庄家牌
	cs2 := t.getHandCards(SEAT2)   //庄家牌
	cs3 := t.getHandCards(SEAT3)   //庄家牌
	cs4 := t.getHandCards(SEAT4)   //庄家牌
	cs5 := t.getHandCards(SEAT5)   //庄家牌
	var a1 uint32 = algo.Algo(cs1) //1位置庄家牌力
	var a2 uint32 = algo.Algo(cs2) //2位置庄家牌力
	var a3 uint32 = algo.Algo(cs3) //3位置庄家牌力
	var a4 uint32 = algo.Algo(cs4) //4位置庄家牌力
	var a5 uint32 = algo.Algo(cs5) //5位置庄家牌力
	power := make(map[uint32]uint32)
	power[SEAT1] = a1
	power[SEAT2] = a2
	power[SEAT3] = a3
	power[SEAT4] = a4
	power[SEAT5] = a5
	var seat, val uint32
	for k, v := range power {
		if v > val {
			seat = k
			val = v
		}
	}
	if seat == SEAT1 || seat == 0 {
		return
	}
	t.handCards[SEAT1] = t.getHandCards(seat)
	t.handCards[seat] = cs1
}

//结束游戏
func (t *DeskFree) gameOver() {
	//glog.Infof("game Over -> %s, round -> %d", t.id, t.round)
	t.shuffle() //洗牌
	t.deal()    //发牌
	t.paili()   //发牌
	// 结算
	cs1 := t.getHandCards(SEAT1)   //庄家牌
	cs2 := t.getHandCards(SEAT2)   //庄家牌
	cs3 := t.getHandCards(SEAT3)   //庄家牌
	cs4 := t.getHandCards(SEAT4)   //庄家牌
	cs5 := t.getHandCards(SEAT5)   //庄家牌
	var a1 uint32 = algo.Algo(cs1) //1位置庄家牌力
	var a2 uint32 = algo.Algo(cs2) //2位置庄家牌力
	var a3 uint32 = algo.Algo(cs3) //3位置庄家牌力
	var a4 uint32 = algo.Algo(cs4) //4位置庄家牌力
	var a5 uint32 = algo.Algo(cs5) //5位置庄家牌力
	t.power[SEAT1] = a1
	t.power[SEAT2] = a2
	t.power[SEAT3] = a3
	t.power[SEAT4] = a4
	t.power[SEAT5] = a5
	//奖池中奖
	t.winPrize()
	//发奖
	t.sendPrize()
	//各位置和庄家对比的赔付倍数
	//正表示赢,负表示赔
	t.multiple[SEAT2] = muliti(a1, a2, cs1, cs2)
	t.multiple[SEAT3] = muliti(a1, a3, cs1, cs3)
	t.multiple[SEAT4] = muliti(a1, a4, cs1, cs4)
	t.multiple[SEAT5] = muliti(a1, a5, cs1, cs5)
	//牌局数累加一次
	t.round++
	//玩家坐庄次数累加
	//if t.dealer != "" {
	//	t.dealerNum++
	//}
	t.dealerNum++
	t.xianjia_jiesuan() //结算,闲家赔付
	//庄家收钱,奖池抽成
	t.dealer_win()
	//庄家赔付,闲家收钱,奖池抽成
	t.dealer_jiesuan()
	t.print()               //打印信息
	msg := t.res_overFree() //结束消息
	t.broadcast(msg)        //广播
	//TODO t.saveTrend()           //记录房间趋势
	t.setRecord()    //个人记录
	t.gameOverInit() //重置状态
	t.checkOffline() //踢除离线玩家
	//t.sendOver()            //发放给玩家金币
	//t.saveRecord()          //日志
}

////打印牌局状态信息,test
//func (t *DeskFree) Print() {
//	t.print()
//}

// 打印
func (t *DeskFree) print() {
	glog.Infof("game over players -> %d", len(t.players))
	glog.Infof("game over dealer -> %s, dealerNum -> %d", t.dealer, t.dealerNum)
	glog.Infof("game over num -> %d, carry -> %d", t.num, t.carry)
	glog.Infof("game over score -> %#v, seats -> %#v", t.score, t.seats)
	glog.Infof("game over power -> %#v", t.power)
	glog.Infof("game over bets -> %#v", t.bets)
	glog.Infof("game over seatBets -> %#v", t.seatBets)
	glog.Infof("game over handCards -> %+x", t.handCards)
	glog.Infof("game over score2 -> %#v", t.score2)
	glog.Infof("game over score3 -> %#v", t.score3)
}

// 踢除离线玩家
func (t *DeskFree) checkOffline() {
	for userid, p := range t.players {
		//TODO
		//if p.GetConn() != nil {
		//	continue
		//}
		if t.dealer == userid && userid != "" {
			t.delBeDealer(userid, p)
		}
		t.leaveBeDealer(userid, p)    //清空上庄列表
		t.leaveSeat(userid, false, p) //离开位置
		//广播消息
		msg := handler.LeaveMsg(userid, t.seats[userid])
		t.broadcast(msg)
		//清除数据
		delete(t.players, userid)
		//TODO p.SetRoom(nil) //清除玩家房间数据
		//glog.Infof("userid logout : %s", userid)
	}
}

/*
//闲家结算金币发放
func (t *DeskFree) sendOver() {
	for k, v := range t.score2 {
		if v <= 0 {
			continue
		}
		p := t.getPlayer(k)
		if p == nil {
			continue
		}
		t.sendFreeCoin(p, v, data.LogType6)
	}
}
*/

//个人记录
func (t *DeskFree) setRecord() {
	for k, v := range t.score2 {
		//p := t.getPlayer(k)
		//if p == nil {
		//	continue
		//}
		//if v > 0 {
		//	p.SetRecord(1) //胜利
		//} else if v < 0 {
		//	p.SetRecord(-1) //输了
		//} else {
		//	p.SetRecord(0) //荒庄
		//}
		if p, ok := t.pids[k]; ok {
			if v > 0 {
				msg2 := &pb.SetRecord{
					Rtype: 1,
				}
				p.Tell(msg2)
			} else if v < 0 {
				msg2 := &pb.SetRecord{
					Rtype: -1,
				}
				p.Tell(msg2)
			} else {
				msg2 := &pb.SetRecord{
					Rtype: 0,
				}
				p.Tell(msg2)
			}
		}
	}
}

//中奖金币发放
func (t *DeskFree) sendPrize() {
	for k, v := range t.score4 {
		tmp := t.getBets(k)
		switch k {
		case SEAT1:
			if t.dealer != "" {
				t.carry += v
			}
		}
		num := uint32(len(tmp))
		if num > 0 {
			v2 := v / num //平分奖金
			for k1, _ := range tmp {
				p := t.getPlayer(k1)
				if p == nil {
					continue
				}
				t.score5[k1] += v2
				t.sendFreeCoin(p, int32(v2), data.LogType17)
			}
		}
	}
}

//中奖
//1.五花牛 炸弹 五小牛才能分一定比例的奖池(5%,10%,20%)
func (t *DeskFree) winPrize() {
	for k, v := range t.power {
		if k == SEAT1 { //庄家限制
			if t.dealer == "" || t.carry < PRIZE_LIMIT {
				continue
			}
		} else if t.seatBets[k] < PRIZE_LIMIT { //闲家限制
			continue
		}
		switch v {
		case algo.FiveTiny:
			t.score4[k] = uint32(math.Trunc(float64(t.pond) * 0.2))
		case algo.Bomb:
			t.score4[k] = uint32(math.Trunc(float64(t.pond) * 0.1))
		case algo.FiveFlower:
			t.score4[k] = uint32(math.Trunc(float64(t.pond) * 0.05))
		}
	}
	var val uint32
	for _, v := range t.score4 {
		val += v
	}
	if val > 0 {
		if t.pond >= val {
			t.pond -= val
		} else {
			t.pond = 0
		}
	}
}

//庄家收钱,奖池抽成
func (t *DeskFree) dealer_win() {
	//庄家收钱总额
	var val int32
	for seat, m := range t.score3 {
		for _, v := range m {
			t.score[seat] += v
			val += v
		}
	}
	//庄家收钱转为正数
	if val < 0 {
		val *= -1
	}
	//奖池抽成
	val1 := int32(math.Trunc(float64(val) * 0.05))
	t.pond += uint32(val1) //更新奖池
	val2 := val - val1
	if val2 < 0 {
		val2 = 0
	}
	t.score[SEAT1] = val2
	t.carry += uint32(val2) //更新庄家收入携带
}

//闲家赔付
func (t *DeskFree) xianjia_jiesuan() {
	for k, v := range t.multiple {
		if v < 0 { //表示庄家输
			continue
		}
		tmp := t.getBets(k)
		t.score3[k] = make(map[string]int32)
		if v > 1 { //表示庄家赢,且大于1倍从玩家身上扣赔付倍数
			val := uint32(v - 1)
			for userid, betNum := range tmp {
				p := t.getPlayer(userid)
				if p == nil {
					continue
				}
				coin := p.GetCoin()
				num := val * betNum
				if num > coin {
					num = coin
				}
				//扣除位置数
				t.sendFreeCoin(p, -1*int32(num), data.LogType6)
				t.score3[k][userid] = -1 * int32((betNum + num))
				t.score2[userid] += -1 * int32((betNum + num))
			}
		} else if v > 0 {
			for userid, betNum := range tmp {
				t.score3[k][userid] = -1 * int32(betNum)
				t.score2[userid] += -1 * int32(betNum)
			}
		}
	}
}

//庄家赔付
func (t *DeskFree) dealer_jiesuan() {
	var num uint32
	for k, v := range t.multiple {
		if v > 0 { //表示庄家赢
			continue
		}
		if v < 0 { //表示庄家输
			v2 := -1 * v
			if v2 < 0 {
				v2 = 1 //1倍
			}
			num += t.seatBets[k] * uint32(v2)
		}
	}
	if t.carry >= num { //足够赔付
		t.dealer_jiesuan1()
	} else { //不足赔付
		t.dealer_jiesuan2(num)
	}
}

//足够赔付
func (t *DeskFree) dealer_jiesuan1() {
	for k, v := range t.multiple {
		if v > 0 { //表示庄家赢
			continue
		}
		tmp := t.getBets(k)
		t.score3[k] = make(map[string]int32)
		if v < 0 { //表示庄家输
			for userid, betNum := range tmp {
				p := t.getPlayer(userid)
				if p == nil {
					continue
				}
				val := int32(v * -1)
				if val < 0 {
					val = 0
				}
				num := uint32(val) * betNum
				if num > t.carry {
					num = t.carry
				}
				t.carry -= num
				n := int32(num + betNum)
				//奖池抽成
				val1 := int32(math.Trunc(float64(n) * 0.05))
				t.pond += uint32(val1) //更新奖池
				val2 := n - val1
				if val2 < 0 {
					val2 = 0
				}
				//扣除位置数
				t.sendFreeCoin(p, val2, data.LogType6)
				t.score3[k][userid] = val2
				t.score[k] += val2
				t.score2[userid] += val2
				t.score[SEAT1] -= int32(num)
			}
		}
	}
}

//不足赔付
func (t *DeskFree) dealer_jiesuan2(num uint32) {
	m := make(map[uint32]uint32)
	for k, v := range t.multiple {
		if v > 0 { //表示庄家赢
			continue
		}
		if v < 0 { //表示庄家输
			v2 := -1 * v
			if v2 < 0 {
				v2 = 1 //1倍
			}
			num1 := t.seatBets[k] * uint32(v2) //当前位置的总金额
			num2 := (num1 / num) * t.carry
			m[k] = num2 //位置分到金额
		}
	}
	for seat, val := range m {
		tmp := t.getBets(seat)
		betsNum := t.seatBets[seat]
		t.score3[seat] = make(map[string]int32)
		for userid, betNum := range tmp {
			p := t.getPlayer(userid)
			if p == nil {
				continue
			}
			num2 := (betNum / betsNum) * val //分到金额
			num3 := num2 + betNum            //加上下注额
			//奖池抽成
			val1 := int32(math.Trunc(float64(num3) * 0.05))
			t.pond += uint32(val1) //更新奖池
			val2 := int32(num3) - val1
			if val2 < 0 {
				val2 = 0
			}
			//扣除位置数
			t.sendFreeCoin(p, val2, data.LogType6)
			t.score3[seat][userid] = val2
			t.score[seat] += val2
			t.score2[userid] += val2
			t.score[SEAT1] -= int32(num2)
		}
	}
}

//返回庄家赢倍数,a1庄家牌力,an闲家牌力,庄家赢返回正数,输返回负数
func muliti(a1, an uint32, cs1, csn []uint32) int32 {
	switch {
	case a1 > an:
		return int32(algo.Multiple(a1))
	case a1 < an:
		return -1 * int32(algo.Multiple(an))
	case a1 == an:
		if algo.Compare(cs1, csn) {
			return int32(algo.Multiple(a1))
		}
		return -1 * int32(algo.Multiple(an))
	}
	return 1
}

//房间消息广播
func (t *DeskFree) broadcast(msg interface{}) {
	//for _, p := range t.players {
	//	p.Send(msg)
	//}
	for _, p := range t.pids {
		p.Tell(msg)
	}
}

//房间消息广播(除seat外)
func (t *DeskFree) broadcast_(userid string, msg interface{}) {
	//for i, p := range t.players {
	//	if i != userid {
	//		p.Send(msg)
	//	}
	//}
	for i, p := range t.pids {
		if i != userid {
			p.Tell(msg)
		}
	}
}

//获取对应位置下注列表
func (t *DeskFree) getBets(seat uint32) map[string]uint32 {
	switch seat {
	case SEAT2:
		return t.tian
	case SEAT3:
		return t.di
	case SEAT4:
		return t.xuan
	case SEAT5:
		return t.huang
	}
	tmp := make(map[string]uint32)
	return tmp
}

//获取
func (t *DeskFree) getPlayer(userid string) *data.User {
	if v, ok := t.players[userid]; ok && v != nil {
		return v
	}
	//panic(fmt.Sprintf("getPlayer error:%s", userid))
	return nil
}

//获取手牌
func (t *DeskFree) getHandCards(seat uint32) []uint32 {
	if v, ok := t.handCards[seat]; ok && v != nil {
		return v
	}
	//return []uint32{}
	panic(fmt.Sprintf("getHandCards error:%d", seat))
}

////检测重复进入
//func (t *DeskFree) reEnter(p *data.User) bool {
//	//userid := p.GetUserid()
//	//if seat, ok := t.seats[userid]; !ok {
//	//	return false
//	//}
//	msg1 := t.res_reEnter(p)
//	p.Send(msg1)
//	return true
//}

//进入房间响应消息
func (t *DeskFree) res_reEnter() interface{} {
	stoc := new(pb.SEnterFreeRoom)
	var timer int
	var dealerNum, left uint32
	var dealerPhoto string
	if t.state == STATE_FREE_BET {
		timer = FREE_DT - t.timer
	} else if t.state == STATE_FREE_DEALER {
		timer = FREE_RT - t.timer
	}
	if timer < 0 {
		timer = 0
	}
	if t.dealer == "" {
		left = 1
		dealerNum = 1
	} else {
		left = t.leftDealerTimes()
		dealerNum = DEALER_TIMES
		d_p := t.getPlayer(t.dealer)
		if d_p != nil {
			dealerPhoto = d_p.GetPhoto()
		}
	}
	roomdata := &pb.FreeRoom{
		Dealer:        SEAT1,
		Userid:        t.dealer,
		Coin:          t.carry,
		Pond:          t.pond,
		State:         t.state,
		Timer:         uint32(timer),
		DealerNum:     dealerNum,
		LeftDealerNum: left,
		Photo:         dealerPhoto,
	}
	handler.RoomInfoMsg2(t.id, t.data, roomdata)
	stoc.Roominfo = roomdata
	//
	for k, v := range t.seatBets {
		betsinfo := &pb.RoomBets{
			Seat: k,
			Bets: v,
		}
		stoc.Betsinfo = append(stoc.Betsinfo, betsinfo)
	}
	for k, v := range t.players { //为坐下玩家
		user := new(pb.FreeUser)
		handler.UserInfoMsg2(v, user)
		user.Seat = t.seats[k]
		user.Ready = false
		user.Bet = t.bets[k]
		user.Vip = &pb.VipInfo{
			Level:  uint32(v.GetVipLevel()),
			Number: v.GetVip(),
		}
		for i := SEAT2; i <= SEAT5; i++ {
			bets2 := &pb.RoomBets{
				Seat: i,
			}
			switch i {
			case SEAT2:
				bets2.Bets = t.tian[k]
			case SEAT3:
				bets2.Bets = t.di[k]
			case SEAT4:
				bets2.Bets = t.xuan[k]
			case SEAT5:
				bets2.Bets = t.huang[k]
			}
			user.Bets = append(user.Bets, bets2)
		}
		stoc.Userinfo = append(stoc.Userinfo, user)
	}
	return stoc
}

func (t *DeskFree) freecamein(v *data.User) interface{} {
	k := v.GetUserid()
	seat := t.seats[k]
	bet := t.bets[k]
	user := new(pb.FreeUser)
	handler.UserInfoMsg2(v, user)
	user.Seat = seat
	user.Ready = false
	user.Bet = bet
	for i := SEAT2; i <= SEAT5; i++ {
		bets2 := &pb.RoomBets{
			Seat: i,
		}
		switch i {
		case SEAT2:
			bets2.Bets = t.tian[k]
		case SEAT3:
			bets2.Bets = t.di[k]
		case SEAT4:
			bets2.Bets = t.xuan[k]
		case SEAT5:
			bets2.Bets = t.huang[k]
		}
		user.Bets = append(user.Bets, bets2)
	}
	return &pb.SFreeCamein{Userinfo: user}
}

//上庄列表消息
func (t *DeskFree) res_bedealerlist() interface{} {
	stoc := new(pb.SDealerList)
	for _, m := range t.dealers {
		for k, v := range m {
			p := t.getPlayer(k)
			if p == nil {
				continue
			}
			list := &pb.DealerList{
				Userid:   k,
				Nickname: p.GetName(),
				Photo:    p.GetPhoto(),
				Coin:     v,
			}
			stoc.List = append(stoc.List, list)
		}
	}
	return stoc
}

//结算消息
func (t *DeskFree) res_overFree() interface{} {
	var left uint32
	if t.dealer == "" {
		left = 1
	} else {
		left = t.leftDealerTimes()
	}
	stoc := &pb.SFreeGameover{
		State:         t.state,
		Dealer:        t.dealer,
		Coin:          t.carry,
		Pond:          t.pond,
		DealerNum:     DEALER_TIMES,
		LeftDealerNum: left,
	}
	for k, v := range t.power {
		data := &pb.FreeRoomOver{
			Seat:  k,
			Cards: t.handCards[k],
			Value: v,
			Total: int32(t.seatBets[k]),
			Score: t.score[k],
		}
		for k1, v1 := range t.score3[k] {
			list := &pb.RoomScore{
				Seat:   t.seats[k1],
				Userid: k1,
				Score:  v1,
				Pond:   t.score5[k1],
			}
			data.List = append(data.List, list)
		}
		stoc.Data = append(stoc.Data, data)
	}
	for k2, v2 := range t.score2 {
		p := t.getPlayer(k2)
		if p == nil {
			continue
		}
		list := &pb.RoomScore{
			Seat:   t.seats[k2],
			Userid: k2,
			Score:  v2,
			Pond:   t.score5[k2],
			Coin:   p.GetCoin(),
		}
		stoc.List = append(stoc.List, list)
	}
	return stoc
}
