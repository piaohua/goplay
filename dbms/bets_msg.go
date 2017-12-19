package main

import (
	"math"
	"math/rand"
	"time"

	"goplay/data"
	"goplay/game/algo"
	"goplay/game/handler"
	"goplay/glog"
	"goplay/pb"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func (a *BetsActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *pb.Connected:
		arg := msg.(*pb.Connected)
		glog.Infof("Connected %s", arg.Name)
	case *pb.Disconnected:
		arg := msg.(*pb.Disconnected)
		glog.Infof("Disconnected %s", arg.Name)
	case *pb.ServeStop:
		//关闭服务
		a.handlerStop(ctx)
		//响应登录
		rsp := new(pb.ServeStoped)
		ctx.Respond(rsp)
	case *pb.ServeStart:
		a.start(ctx)
		//响应
		//rsp := new(pb.ServeStarted)
		//ctx.Respond(rsp)
	case *pb.Tick:
		a.ding(ctx)
	case *pb.CBettingInfo:
		arg := msg.(*pb.CBettingInfo)
		glog.Debugf("CBettingInfo %#v", arg)
		//响应
		rsp := a.getInfo()
		ctx.Respond(rsp)
	case *pb.BetsOn:
		arg := msg.(*pb.BetsOn)
		glog.Debugf("BetsOn %#v", arg)
		seat := arg.GetSeat()
		number := arg.GetNumber()
		userid := arg.GetUserid()
		//响应
		rsp := a.bet(userid, seat, number)
		ctx.Respond(rsp)
	case *pb.BetsRecord:
		arg := msg.(*pb.BetsRecord)
		glog.Debugf("BetsRecord %#v", arg)
		page := arg.GetPage()
		userid := arg.GetUserid()
		//响应
		rsp := handler.BettingRecord(page, userid)
		ctx.Respond(rsp)
	default:
		glog.Errorf("unknown message %v", msg)
	}
}

//启动服务
func (a *BetsActor) start(ctx actor.Context) {
	glog.Infof("bets start: %v", ctx.Self().String())
	//初始化建立连接
	bind := cfg.Section("hall").Key("bind").Value()
	name := cfg.Section("cookie").Key("name").Value()
	//timeout := 3 * time.Second
	//hallPid, err := remote.SpawnNamed(bind, a.Name, name, timeout)
	//if err != nil {
	//	glog.Fatalf("remote hall err %v", err)
	//}
	//a.hallPid = hallPid.Pid
	a.hallPid = actor.NewPID(bind, name)
	glog.Infof("a.hallPid: %s", a.hallPid.String())
	connect := &pb.Connect{
		Name: a.Name,
	}
	a.hallPid.Request(connect, ctx.Self())
	//初始化
	a.initBetting(ctx)
	go a.ticker(ctx) //goroutine,计时
}

//时钟
func (a *BetsActor) ticker(ctx actor.Context) {
	glog.Info("Betting ticker started")
	tick := time.Tick(time.Second)
	msg := new(pb.Tick)
	for {
		select {
		case <-a.stopCh:
			glog.Info("bets ticker closed")
			return
		default: //防止阻塞
		}
		select {
		case <-a.stopCh:
			glog.Info("bets ticker closed")
			return
		case <-tick:
			ctx.Self().Tell(msg)
		}
	}
}

//钟声
func (t *BetsActor) ding(ctx actor.Context) {
	glog.Debugf("ding: %v", ctx.Self().String())
	switch t.timer {
	case t.betTime:
		//投注结束
		t.timer = 0
		t.state = 1
		//结算
		t.over()
	case t.betRest:
		if t.state == 0 {
			t.timer++
			return
		}
		//休息结束
		t.timer = 0
		t.state = 0
		//广播开始消息
		t.gamestart()
	default:
		t.timer++
	}
}

//关闭时钟
func (a *BetsActor) closeTick() {
	select {
	case <-a.stopCh:
		return
	default:
		//停止发送消息
		close(a.stopCh)
	}
}

//停止服务
func (a *BetsActor) handlerStop(ctx actor.Context) {
	glog.Debugf("handlerStop: %s", a.Name)
	//关闭
	a.closeTick()
	//回存数据
	if a.uniqueid != nil {
		a.uniqueid.Save()
	}
	//结算
	if a.state == 0 && len(a.jackpot) != 0 {
		a.over()
	}
}

func (a *BetsActor) initBetting(ctx actor.Context) {
	a.betTime = 180 //最少90秒,index "000" 开始
	a.betRest = 25
	a.timer = 0
	a.state = 1 //等待开始
	a.stopCh = make(chan struct{})
	a.odds = make(map[uint32]float32)
	a.jackpot = make(map[uint32]uint32)
	a.ante = make(map[string]map[uint32]uint32)
	a.result = make(map[string]map[uint32]int32)
	a.today = utils.String(utils.DayDate())
	//初始化
	a.initOdds()
}

func (t *BetsActor) initOdds() {
	list := data.GetBettingList()
	for _, v := range list {
		t.odds[v.Seat] = v.Odds
	}
}

func (a *BetsActor) print() {
	glog.Infof("timer -> %d", a.timer)
	glog.Infof("state -> %d", a.state)
	glog.Infof("today -> %s", a.today)
	glog.Infof("index -> %s", a.uniqueid.Index())
	glog.Infof("odds -> %v", a.odds)
	glog.Infof("jackpot -> %v", a.jackpot)
	glog.Infof("ante -> %#v", a.ante)
	glog.Infof("result -> %v", a.result)
	glog.Infof("prize -> %v", a.prize)
	glog.Infof("winner -> %v", a.winner)
	glog.Infof("lose -> %v", a.lose)
}

func (a *BetsActor) reset() {
	a.ante = make(map[string]map[uint32]uint32)
	//个人赢
	a.prize = make(map[string]int32)
	//个人输赢
	a.lose = make(map[string]int32)
	//个人位置输赢
	a.result = make(map[string]map[uint32]int32)
	//
	a.winner = make(map[uint32]bool)
	//
	a.jackpot = make(map[uint32]uint32)
	//
	a.niu = 0
}

//洗牌
func (a *BetsActor) shuffle() (cards []uint32) {
	rand.Seed(time.Now().UnixNano())
	d := make([]uint32, algo.NumCard, algo.NumCard)
	copy(d, algo.NiuCARDS)
	//测试暂时去掉洗牌
	for i := range d {
		j := rand.Intn(i + 1)
		d[i], d[j] = d[j], d[i]
	}
	cards = d[:5]
	return
}

func (a *BetsActor) gamestart() {
	today := utils.String(utils.DayDate())
	if today != a.today { //第二天重置
		a.uniqueid.Save()
		a.uniqueid.Id = data.BETTING_KEY + today
		a.uniqueid.Value = today + "000"
		a.today = today
	}
	//新期号
	a.uniqueid.GenID()
	//消息广播
	msg := new(pb.SPushNewBetting)
	msg.Index = a.uniqueid.Index()
	msg.Status = &pb.Betting{
		Times: a.betTime,
		State: a.state,
	}
	glog.Infof("gamestart today %s, state %d", a.today, a.state)
	//Broadcast(msg)
	//a.hallPid.Tell(msg)
	nodePid.Tell(msg)
	//重置
	a.reset()
	//打印
	//a.print()
}

//发奖
func (t *BetsActor) sendPrize() {
	for k, v := range t.prize {
		if v <= 0 {
			continue
		}
		msg := new(pb.ChangeCurrency)
		msg.Userid = k
		msg.Diamond = v
		msg.Type = int32(data.LogType36)
		msg.Upsert = true
		nodePid.Tell(msg)
	}
}

//结算
func (t *BetsActor) jiesuan() {
	//个人赢
	t.prize = make(map[string]int32)
	//个人输赢
	t.lose = make(map[string]int32)
	//个人位置输赢
	t.result = make(map[string]map[uint32]int32)
	//结算
	for k, m := range t.ante {
		if _, ok := t.result[k]; !ok {
			t.result[k] = make(map[uint32]int32)
		}
		for s, v := range m {
			if w, ok := t.winner[s]; ok && w {
				rate := t.odds[s]
				val := int32(math.Trunc(float64(float32(v) * rate)))
				t.prize[k] += val
				t.result[k][s] += val
				t.lose[k] += val
			} else {
				t.lose[k] -= int32(v)
				t.result[k][s] -= int32(v)
			}
		}
	}
}

func (t *BetsActor) over() {
	t.cards = t.shuffle()
	//中奖位置
	t.winner, t.niu = algo.Betting(t.cards)
	//结算
	t.jiesuan()
	//发奖
	t.sendPrize()
	//消息广播
	status := &pb.Betting{
		Times: t.betRest,
		State: t.state,
	}
	var i uint32
	info := make([]*pb.JackpotOver, 0)
	for i = 1; i <= 11; i++ {
		list := &pb.JackpotOver{
			Seat:  i,
			Win:   t.winner[i],
			Odds:  t.odds[i],
			Count: t.jackpot[i],
		}
		info = append(info, list)
	}
	for k, m := range t.ante {
		msg := new(pb.SPushBetting)
		msg.Status = status
		msg.Info = info
		msg.Cards = t.cards
		msg.Niu = t.niu
		msg.Number = t.lose[k]
		for i = 1; i <= 11; i++ {
			list := &pb.JackpotSelf{
				Seat: i,
				Ante: m[i],
				//Number: t.result[k][i],
			}
			if t.result[k][i] > 0 {
				list.Number = t.result[k][i]
			}
			msg.List = append(msg.List, list)
		}
		msg1 := new(pb.BetsResult)
		b, err1 := msg.Marshal()
		if err1 != nil {
			glog.Errorf("BetsResult k %s err1 %v", k, err1)
		}
		msg1.Userid = k
		msg1.Message = string(b)
		//TODO 个人消息推送
		nodePid.Tell(msg1)
	}
	//日志
	t.record()
	//
	//t.print()
	//
	t.reset()
}

//记录
func (t *BetsActor) record() {
	for k, v := range t.lose {
		d := &data.BettingUser{
			Userid: k,
			Index:  t.uniqueid.Index(),
			Lose:   v,
		}
		if !d.Save() {
			glog.Errorf("save record failed d:%#v", d)
		}
	}
	seats := make([]uint32, 0)
	for k, _ := range t.winner {
		seats = append(seats, k)
	}
	d := &data.BettingRecord{
		Id:    t.uniqueid.Index(),
		Cards: t.cards,
		Niu:   t.niu,
		Seats: seats,
		Lose:  t.lose,
	}
	d.Ante = make(map[string][]data.SeatBetting)
	for u, m := range t.ante {
		s := make([]data.SeatBetting, 0)
		for k, v := range m {
			s1 := data.SeatBetting{
				Seat:   k,
				Number: v,
			}
			s = append(s, s1)
		}
		d.Ante[u] = s
	}
	if !d.Save() {
		glog.Errorf("save record failed d:%#v", d)
	}
}

//获取信息
func (t *BetsActor) getInfo() (msg *pb.SBettingInfo) {
	status := &pb.Betting{
		State: t.state,
	}
	if t.state == 0 {
		status.Times = t.betTime - t.timer
	} else {
		status.Times = t.betRest - t.timer
	}
	msg = &pb.SBettingInfo{
		Index:  t.uniqueid.Index(),
		Status: status,
	}
	var i uint32
	for i = 1; i <= 11; i++ {
		list := &pb.JackpotInfo{
			Seat:   i,
			Odds:   t.odds[i],
			Number: t.jackpot[i],
		}
		msg.List = append(msg.List, list)
	}
	return msg
}

//投注
func (t *BetsActor) bet(userid string, seat,
	number uint32) (stoc *pb.SBetting) {
	stoc = new(pb.SBetting)
	if t.state == 1 {
		stoc.Error = pb.GameNotStart
		return
	}
	t.jackpot[seat] += number
	if m, ok := t.ante[userid]; ok {
		m[seat] += number
		t.ante[userid] = m
	} else {
		m := make(map[uint32]uint32)
		m[seat] += number
		t.ante[userid] = m
	}
	//广播消息,TODO 只广播参与人员
	msg := new(pb.SPushJackpot)
	for k, v := range t.jackpot {
		msg1 := &pb.Jackpot{
			Seat:   k,
			Number: v,
		}
		msg.List = append(msg.List, msg1)
	}
	//Broadcast(msg)
	//a.hallPid.Tell(msg)
	nodePid.Tell(msg)
	//响应
	stoc.Seat = seat
	stoc.Number = number
	return
}

//后台设置
func (t *BetsActor) setOdds(seat uint32, number float32) {
	t.odds[seat] = number
}
