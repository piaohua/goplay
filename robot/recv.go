/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-11-19 12:48:03
 * Filename      : recv.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"goplay/glog"
	"goplay/pb"
)

func (r *Robot) receive(msg interface{}) {
	switch msg.(type) {
	case *pb.SLogin:
		r.recvLogin(msg.(*pb.SLogin))
	case *pb.SEnterRoom:
		r.recvComein(msg.(*pb.SEnterRoom))
	case *pb.SLaunchVote:
		r.recvVote(msg.(*pb.SLaunchVote))
	case *pb.SCreateRoom:
		r.recvCreate(msg.(*pb.SCreateRoom))
	case *pb.SPushDealer:
		//r.recvDealer(msg.(*pb.SPushDealer))
	case *pb.SDraw:
		//r.recvDraw(msg.(*pb.SDraw))
	case *pb.SGameover:
		r.recvGameover(msg.(*pb.SGameover))
	case *pb.SLeave:
		r.recvLeave(msg.(*pb.SLeave))
	case *pb.SUserData:
		r.recvdata(msg.(*pb.SUserData))
	case *pb.SKick:
		r.recvKick(msg.(*pb.SKick))
	case *pb.SRegist:
		r.recvRegist(msg.(*pb.SRegist))
	case *pb.SEnterFreeRoom:
		r.recvEntryFree(msg.(*pb.SEnterFreeRoom))
	case *pb.SFreeSit:
		r.recvFreeSit(msg.(*pb.SFreeSit))
	case *pb.SFreeBet:
		r.recvFreeBet(msg.(*pb.SFreeBet))
	case *pb.SFreeGamestart:
		r.recvGamestartFree(msg.(*pb.SFreeGamestart))
	case *pb.SFreeGameover:
		r.recvGameoverFree(msg.(*pb.SFreeGameover))
	case *pb.SPushCurrency:
		r.recvPushCurrency(msg.(*pb.SPushCurrency))
	case *pb.SEnterClassicRoom:
		r.recvEntryClassic(msg.(*pb.SEnterClassicRoom))
	case *pb.SClassicList:
		r.recvClassicList(msg.(*pb.SClassicList))
	case *pb.SClassicGameover:
		r.recvGameover2(msg.(*pb.SClassicGameover))
	case *pb.SEnterZiRoom:
		r.recvEntryPhz(msg.(*pb.SEnterZiRoom))
	case *pb.SZiCamein:
		r.recvComeinPhz(msg.(*pb.SZiCamein))
	case *pb.SZiGameover:
		r.recvPhzOver(msg.(*pb.SZiGameover))
	case *pb.SPushDeal:
		r.recvPhzDeal(msg.(*pb.SPushDeal))
	case *pb.SPushDealerDeal:
		r.recvPhzDealerDeal(msg.(*pb.SPushDealerDeal))
	case *pb.SPushDraw:
		r.recvPhzDraw(msg.(*pb.SPushDraw))
	case *pb.SPushDiscard:
		r.recvPhzDiscard(msg.(*pb.SPushDiscard))
	case *pb.SPushAuto:
		r.recvPhzAuto(msg.(*pb.SPushAuto))
	case *pb.SOperate:
		r.recvPhzOperate(msg.(*pb.SOperate))
	case *pb.SPushStatus:
		r.recvPhzStatus(msg.(*pb.SPushStatus))
	default:
		glog.Errorf("unknow message: %#v", msg)
	}
}

//' 接收到服务器登录返回
func (r *Robot) recvRegist(stoc *pb.SRegist) {
	var errcode = stoc.GetError()
	switch {
	case errcode == 0:
		Logined(r.data.Phone) //登录成功
		r.regist = true       //注册成功
		r.data.Userid = stoc.GetUserid()
		glog.Infof("regist successful -> %s", r.data.Userid)
		r.SendGetClassic()
		r.SendUserData() // 获取玩家数据
		return
	default:
		glog.Infof("regist err -> %d", errcode)
	}
	//重新尝试登录
	go ReLogined(r.data.Phone, r.code, r.rtype)
	r.Close()
}

//.

//' 接收到服务器登录返回
func (r *Robot) recvLogin(stoc *pb.SLogin) {
	var errcode = stoc.GetError()
	switch {
	case errcode == 0:
		Logined(r.data.Phone) //登录成功
		r.data.Userid = stoc.GetUserid()
		glog.Infof("login successful -> %s", r.data.Userid)
		r.SendGetClassic()
		r.SendUserData() // 获取玩家数据
		return
	default:
		glog.Infof("login err -> %d", errcode)
	}
	r.Close()
}

//.

//' 接收到玩家数据
func (r *Robot) recvdata(stoc *pb.SUserData) {
	var errcode = stoc.GetError()
	if errcode != 0 {
		glog.Infof("get data err -> %d", errcode)
	}
	userdata := stoc.GetData()
	// 设置数据
	r.data.Userid = userdata.GetUserid()     // 用户id
	r.data.Nickname = userdata.GetNickname() // 用户昵称
	r.data.Sex = userdata.GetSex()           // 用户性别,男1 女2 非男非女3
	r.data.Coin = userdata.GetCoin()         // 金币
	r.data.Diamond = userdata.GetDiamond()   // 钻石
	//r.rtype = userdata.GetRoomtype()         // 房间类型
	//TODO
	//if r.data.Coin < 200000 {
	//	go addCoin(r.data.Userid, 10000000)
	//} else if r.data.Coin < 500000 {
	//	go addCoin(r.data.Userid, 5000000)
	//} else if r.data.Coin < 1000000 {
	//	go addCoin(r.data.Userid, 15000000)
	//}
	//查找房间-进入房间
	switch r.code {
	case "free":
		r.SendEntryFree() //进入百人场
	case "classic":
		r.SendEntryClassic()
	case "create":
		switch r.rtype {
		case 6, 7: //跑胡子
			r.SendEntryPhz()
		default:
			r.SendEntry()
		}
	default:
		glog.Infof("enter phz room -> %s, %d", r.code, r.rtype)
		if len(r.code) == 6 {
			switch r.rtype {
			case 6, 7: //跑胡子
				r.SendEntryPhz()
			default:
				r.SendEntry()
			}
		} else {
			r.Close()
		}
	}
}

//.

//' 离开房间
func (r *Robot) recvLeave(stoc *pb.SLeave) {
	var seat uint32 = stoc.GetSeat()
	if seat == r.seat {
		r.Close() //下线
	}
	if seat >= 1 && seat <= 4 && seat != r.seat {
		r.SendLeave() //离开
	}
}

//.

//' 创建房间
func (r *Robot) recvCreate(stoc *pb.SCreateRoom) {
	var errcode = stoc.GetError()
	switch {
	case errcode == 0:
		var code string = stoc.GetRdata().GetInvitecode()
		if code != "" && len(code) == 6 {
			glog.Infof("create room code -> %s", code)
			r.code = code       //设置邀请码
			r.SendEntry()       //进入房间
			Msg2Robots(code, 3) //创建房间成功,邀请3个人进入
		} else {
			glog.Errorf("create room code empty -> %s", code)
		}
	default:
		glog.Infof("create room err -> %d", errcode)
		r.Close() //进入出错,关闭
	}
}

//.

//' 进入房间
func (r *Robot) recvComein(stoc *pb.SEnterRoom) {
	var errcode = stoc.GetError()
	switch {
	case errcode == 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == r.data.Userid {
				r.seat = v.GetSeat()
				glog.Infof("enter room -> %s, %d", r.data.Userid, r.seat)
				r.SendReady() //准备
				break
			}
		}
		info := stoc.GetRoominfo()
		r.rtype = info.GetRtype() // 房间类型
	default:
		glog.Infof("comein err -> %d", errcode)
		r.Close() //进入出错,关闭
	}
}

//.

//' 解散
func (r *Robot) recvVote(stoc *pb.SLaunchVote) {
	var seat uint32 = stoc.GetSeat()
	glog.Infof("vote seat -> %d", seat)
	r.SendVote()
}

// 解散
func (r *Robot) recvKick(stoc *pb.SKick) {
	userid := stoc.GetUserid()
	if userid == r.data.Userid {
		r.Close() //下线
	}
}

//.

//' 结束
func (r *Robot) recvGameover(stoc *pb.SGameover) {
	var round uint32 = stoc.GetRound()
	r.cards = []uint32{} //清除牌
	if round == 0 {
		r.Close() //结束下线
	} else {
		r.SendReady() //准备
	}
}

//.

//' 抓牌
//STATE_DEALER = 1  //抢庄状态
//STATE_BET    = 2  //下注状态
//STATE_NIU    = 3  //选牛状态
//func (r *Robot) recvDraw(stoc *pb.SDraw) {
//	var state uint32 = stoc.GetState()
//	var seat uint32 = stoc.GetSeat()
//	var cards []uint32 = stoc.GetCards()
//	if seat != r.seat { //自己摸牌
//		return
//	}
//	r.cards = append(r.cards, cards...)
//	switch r.rtype {
//	case data.ROOM_PRIVATE:
//		switch state {
//		case 1:
//			//抢庄
//			r.SendDealer()
//		case 2:
//			//提交组合
//			r.SendNiu()
//		}
//	case data.ROOM_PRIVATE4:
//		switch state {
//		case 2:
//			//提交组合
//			r.SendNiu()
//		}
//	case data.ROOM_PRIVATE3:
//		switch state {
//		case 2:
//			//提交组合
//			r.SendNiu()
//		}
//	}
//}

//.

//' 打庄
//func (r *Robot) recvDealer(stoc *pb.SPushDealer) {
//	var dealer uint32 = stoc.GetDealer()
//	if r.seat == dealer { //做庄不下注
//		return
//	}
//	switch r.rtype {
//	case data.ROOM_PRIVATE:
//		//下注
//		r.SendBet()
//	case data.ROOM_PRIVATE3, data.ROOM_PRIVATE4:
//		//下注
//		r.SendClassicBet()
//	}
//}

//.

//' free

//进入
func (r *Robot) recvEntryFree(stoc *pb.SEnterFreeRoom) {
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
		if r.data.Coin >= 300000 {
			r.SendFreeSit()
		}
		info := stoc.GetRoominfo()
		r.rtype = info.GetRtype() // 房间类型
	default:
		r.SendFreeStandup()
	}
}

//坐下
func (r *Robot) recvFreeSit(stoc *pb.SFreeSit) {
	var errcode = stoc.GetError()
	var seat uint32 = stoc.GetSeat()
	var userid string = stoc.GetUserid()
	if userid != r.data.Userid {
		return
	}
	switch errcode {
	case 0:
		r.seat = seat //坐下位置
	default:
		//if r.sits > 7 { //尝试次数过多
		//	r.SendFreeStandup()
		//} else {
		//	r.SendFreeSit()
		//}
	}
}

//下注
func (r *Robot) recvFreeBet(stoc *pb.SFreeBet) {
	var errcode = stoc.GetError()
	var userid string = stoc.GetUserid()
	if userid != r.data.Userid {
		return
	}
	switch errcode {
	case 0:
		//if r.bits < 3 {
		//	r.SendFreeBet()
		//}
	default:
		r.SendFreeStandup()
	}
}

const (
	STATE_FREE_READY  = 0 //准备状态
	STATE_FREE_DEALER = 1 //休息中状态
	STATE_FREE_BET    = 2 //下注中状态
)

//开始
func (r *Robot) recvGamestartFree(stoc *pb.SFreeGamestart) {
	var state uint32 = stoc.GetState()
	switch state {
	case STATE_FREE_READY:
		r.SendFreeStandup()
	case STATE_FREE_DEALER:
		if r.data.Coin < 20000 {
			r.SendFreeStandup()
		}
	case STATE_FREE_BET:
		r.SendFreeBet() //下注
	default:
		r.SendFreeStandup()
	}
}

//结束
func (r *Robot) recvGameoverFree(stoc *pb.SFreeGameover) {
	r.round++
	if r.round >= 16 { //打10局下线
		r.SendFreeStandup()
	}
}

//更新金币
func (r *Robot) recvPushCurrency(stoc *pb.SPushCurrency) {
	var coin int32 = stoc.GetCoin()
	if coin > 100000 {
		r.SendChat()
	}
	newcoin := int32(r.data.Coin) + coin
	if newcoin < 20000 { //金币少于一定金额时下线
		r.SendFreeStandup()
		r.SendLeave()
		r.Close() //下线
	} else {
		r.data.Coin = uint32(newcoin)
	}
}

//.

//' classic

//进入
func (r *Robot) recvEntryClassic(stoc *pb.SEnterClassicRoom) {
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == r.data.Userid {
				r.seat = v.GetSeat()
				glog.Infof("enter classic room -> %s, %d", r.data.Userid, r.seat)
				r.SendReady() //准备
				break
			}
		}
		info := stoc.GetRoominfo()
		r.rtype = info.GetRtype() // 房间类型
	default:
		r.SendClassicClose()
	}
}

// 结束
func (r *Robot) recvClassicList(stoc *pb.SClassicList) {
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
		r.classicId = ""
		r.classic = stoc.GetList()
	default:
		r.SendClassicClose()
	}
}

// 结束
func (r *Robot) recvGameover2(stoc *pb.SClassicGameover) {
	userdata := stoc.GetData()
	r.cards = []uint32{} //清除牌
	for _, v := range userdata {
		if v.GetSeat() != r.seat {
			continue
		}
		for _, m := range r.classic {
			if m.GetId() != r.classicId {
				continue
			}
			if m.GetMinimum() > v.GetCoin() {
				r.SendClassicClose()
				return
			}
		}
		if v.GetCoin() <= 20000 {
			r.SendClassicClose()
			return
		} else {
			r.SendReady() //准备
			return
		}
	}
	r.SendReady() //准备
}

//.

//' 跑胡子

//进入
func (r *Robot) recvEntryPhz(stoc *pb.SEnterZiRoom) {
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == r.data.Userid {
				r.seat = v.GetSeat()
				glog.Infof("enter phz room -> %s, %d", r.data.Userid, r.seat)
				r.SendReady() //准备
				break
			}
		}
		info := stoc.GetRoominfo()
		r.rtype = info.GetRtype() // 房间类型
	default:
		r.SendPhzClose()
	}
}

//进入
func (r *Robot) recvComeinPhz(stoc *pb.SZiCamein) {
	v := stoc.GetUserinfo()
	if v.GetUserid() == r.data.Userid {
		r.seat = v.GetSeat()
		r.SendReady() //准备
	}
}

// 结算广播接口，游戏结束
func (r *Robot) recvPhzOver(stoc *pb.SZiGameover) {
	var round uint32 = stoc.GetRound()
	r.cards = []uint32{} //清除牌
	if round == 0 {
		r.SendPhzClose() //结束下线
	} else {
		r.SendReady() //准备
	}
}

// 发牌
func (r *Robot) recvPhzDeal(stoc *pb.SPushDeal) {
	if stoc.GetSeat() == r.seat {
		r.cards = stoc.GetCards()
	}
}

// 庄家发牌
func (r *Robot) recvPhzDealerDeal(stoc *pb.SPushDealerDeal) {
	if stoc.GetSeat() == r.seat {
		card := stoc.GetCard()
		r.cards = append(r.cards, card)
		value := stoc.GetValue()
		//胡
		r.phzOperate(card, value)
	}
}

// 摸牌
func (r *Robot) recvPhzDraw(stoc *pb.SPushDraw) {
	card := stoc.GetCard()
	r.cards = stoc.GetCards() //同步手牌
	value := stoc.GetValue()
	//胡碰吃
	r.phzOperate(card, value)
}

// 出牌
func (r *Robot) recvPhzDiscard(stoc *pb.SPushDiscard) {
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
	default:
		glog.Errorf("discard error %d", errcode)
		return
	}
	if r.seat == stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	//吃碰
	r.phzOperate(card, value)
}

//自动操作(提,跑,偎)
func (r *Robot) recvPhzAuto(stoc *pb.SPushAuto) {
	if r.seat == stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	//从手牌中去掉
	r.phzAuto(card, value)
}

// 玩家操作(吃碰胡)
func (r *Robot) recvPhzOperate(stoc *pb.SOperate) {
	if r.seat != stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	cards := stoc.GetCards()
	bione := stoc.GetBione()
	bitwo := stoc.GetBitwo()
	var errcode = stoc.GetError()
	switch errcode {
	case 0:
	default:
		glog.Errorf("operate error %d", errcode)
		glog.Errorf("operate hs %v", r.cards)
		glog.Errorf("operate value: %d, cards: %v, bione: %v, bitwo: %v",
			value, cards, bione, bitwo)
		return
	}
	//移除牌
	r.phzChow(card, value, cards, bione, bitwo)
}

//房间状态
func (r *Robot) recvPhzStatus(stoc *pb.SPushStatus) {
	if r.seat != stoc.GetSeat() {
		return
	}
	if stoc.GetStatus() == 2 { //出牌
		r.phzDiscard()
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
