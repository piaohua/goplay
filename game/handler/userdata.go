package handler

import (
	"math"

	"goplay/data"
	"goplay/game/config"
	"goplay/pb"
	"niu/players"
	"utils"
)

//网关直接调用
func GetCurrency(ctos *pb.CGetCurrency, p *data.User) (stoc *pb.SGetCurrency) {
	stoc = new(pb.SGetCurrency)
	if p == nil {
		return
	}
	stoc.Coin = p.GetCoin()
	stoc.Diamond = p.GetDiamond()
	stoc.Roomcard = p.GetRoomCard()
	return
}

func Ping(ctos *pb.CPing) (stoc *pb.SPing) {
	stoc = new(pb.SPing)
	stoc.Time = ctos.GetTime()
	return
}

func Config(ctos *pb.CConfig) (stoc *pb.SConfig) {
	stoc = new(pb.SConfig)
	return
}

func VipList(ctos *pb.CVipList) (stoc *pb.SVipList) {
	stoc = new(pb.SVipList)
	list := config.GetVipsList()
	for _, v := range list {
		s := &pb.Vip{
			Level:  uint32(v.Level),
			Number: uint32(v.Number) / 100,
			Pay:    v.Pay,
			Prize:  v.Prize,
			Kick:   uint32(v.Kick),
		}
		stoc.List = append(stoc.List, s)
	}
	return
}

func ClassicList(ctos *pb.CClassicList) (stoc *pb.SClassicList) {
	stoc = new(pb.SClassicList)
	list := config.GetClassics()
	for _, v := range list {
		s := &pb.Classic{
			Id:      v.Id,
			Ptype:   uint32(v.Ptype),
			Rtype:   uint32(v.Rtype),
			Ante:    v.Ante,
			Minimum: v.Minimum,
			Maximum: v.Maximum,
		}
		stoc.List = append(stoc.List, s)
	}
	return
}

func PrizeBox(ctos *pb.CPrizeBox, p *data.User) (stoc *pb.SPrizeBox) {
	stoc = new(pb.SPrizeBox)
	state := ctos.GetState()
	id := p.GetBox()
	num := p.GetDuration()
	d := config.GetBox(id)
	boxstate := p.GetBoxState()
	if d.Id == "" { //不存在或全部完成
		stoc.Error = pb.NotBox
		return
	}
	if boxstate == 1 {
		stoc.State = 2
		return
	}
	if state == 2 && d.Duration <= num { //领取
		l := &pb.Prize{
			Id:     d.Id,
			Rtype:  uint32(d.Rtype),
			Number: uint32(d.Amount),
		}
		stoc.List = append(stoc.List, l)
		addPrize(p, d.Rtype, data.LogType22, d.Amount)
		//下一个宝箱
		d = config.NextBox(d.Duration)
		if d.Id == "" {
			p.SetBox(d.Id, 1)
		} else {
			p.SetBox(d.Id, 0)
		}
		p.SetDuration()
	}
	if d.Id == "" {
		stoc.State = 2
		return
	}
	n := &pb.Prize{
		Id:     d.Id,
		Rtype:  uint32(d.Rtype),
		Number: uint32(d.Amount),
	}
	stoc.Next = append(stoc.Next, n)
	num = p.GetDuration()
	stoc.Duration = d.Duration
	stoc.Time = num
	if d.Duration <= num {
		stoc.State = 1
	}
	return
}

func PrizeDraw(ctos *pb.CPrizeDraw, p *data.User) {
	stoc := new(pb.SPrizeDraw)
	var num int32 = config.GetEnv(data.ENV7)
	//vip
	vip := config.GetVip(p.GetVipLevel())
	num += int32(vip.Prize) //vip赠送
	draw := p.GetPrizeDraw()
	if int32(draw) >= num {
		stoc.Error = pb.NotPrizeDraw
		return
	}
	rate := config.GetPrizeRate()
	var n uint32
	if rate > 0 {
		n = uint32(utils.RandInt64N(int64(rate)))
	}
	list := config.GetPrizes()
	for _, v := range list {
		if n > v.Rate {
			l := &pb.Prize{
				Id:     v.Id,
				Rtype:  uint32(v.Rtype),
				Number: uint32(v.Amount),
			}
			left := num - int32(draw) - 1
			if left < 0 {
				left = 0
			}
			stoc.Leftdraw = uint32(left)
			stoc.Prizedraw = draw + 1
			stoc.List = append(stoc.List, l)
			p.SetPrizeDraw()
			addPrize(p, v.Rtype, data.LogType21, v.Amount)
			return
		}
	}
	stoc.Error = pb.NotGotPrizeDraw
	return
}

func PrizeList(ctos *pb.CPrizeList) (stoc *pb.SPrizeList) {
	stoc = new(pb.SPrizeList)
	list := config.GetPrizes()
	for _, v := range list {
		l := &pb.Prize{
			Id:     v.Id,
			Rtype:  uint32(v.Rtype),
			Number: uint32(v.Amount),
		}
		stoc.List = append(stoc.List, l)
	}
	return
}

func bankrupt(ctos *pb.CBankrupts, p *data.User) {
	stoc := &pb.SBankrupts{}
	var coin1 int32 = config.GetEnv(data.ENV8)
	if int32(p.GetCoin()) >= coin1 {
		stoc.Error = pb.NotBankrupt
		return
	}
	var num int32 = config.GetEnv(data.ENV6)
	num2 := p.GetBankrupts()
	if int32(num2) > num {
		stoc.Error = pb.NotRelieves
		return
	}
	var coin int32 = config.GetEnv(data.ENV9)
	if coin > 0 {
		l := &pb.Prize{
			Rtype:  uint32(data.COIN),
			Number: uint32(coin),
		}
		left := num - (int32(num2) + 1)
		if left < 0 {
			left = 0
		}
		stoc.Relieve = uint32(left)
		stoc.Bankrupt = num2 + 1
		stoc.List = append(stoc.List, l)
		p.AddCoin(coin)
		//日志
		data.CoinRecord(p.GetUserid(), data.LogType11, p.GetCoin(), coin)
		p.SetBankrupts()
		pushCurrency(p, data.LogType11, coin, 0)
	}
	return
}

func GetUserData(ctos *pb.CUserData, p *data.User) (stoc *pb.SUserData) {
	stoc = new(pb.SUserData)
	stoc.Data = new(pb.UserData)
	userid := ctos.GetUserid()
	if userid == "" {
		stoc.Error = pb.UsernameEmpty
		return
	}
	// 获取玩家自己的详细资料
	if userid == p.GetUserid() {
		stoc.Data.Bank = p.GetBank()
		first, relieve, bankrupt, prizedraw,
			leftdraw, kicktimes := getActivity(p)
		stoc.Data.Data = &pb.Activity{
			Firstpay:  first,
			Relieve:   relieve,
			Bankrupt:  bankrupt,
			Prizedraw: prizedraw,
			Leftdraw:  leftdraw,
			Kicktimes: kicktimes,
		}
		stoc.Data.Vip = &pb.VipInfo{
			Level:  uint32(p.GetVipLevel()),
			Number: p.GetVip() / 100,
		}
	}
	stoc.Data.Give = p.GetGive()
	stoc.Data.Agent = user.Agent
	stoc.Data.Userid = userid
	stoc.Data.Photo = user.Photo
	stoc.Data.Nickname = user.Nickname
	stoc.Data.Sex = user.Sex
	stoc.Data.Phone = user.Phone
	stoc.Data.Coin = user.Coin
	stoc.Data.Diamond = user.Diamond
	return
}

func getActivity(p *data.User) (first, relieve, bankrupt,
	prizedraw, leftdraw, kicktimes uint32) {
	if p.GetMoney() == 0 {
		first = 1
	}
	//vip
	vip := config.GetVip(p.GetVipLevel())
	//
	bankrupt = p.GetBankrupts()
	prizedraw = p.GetPrizeDraw()
	var num1 int32 = config.GetEnv(data.ENV7)
	num1 += int32(vip.Prize) //vip赠送
	var num2 int32 = config.GetEnv(data.ENV6)
	num2 = num2 - int32(bankrupt)
	if num2 < 0 {
		num2 = 0
	}
	num1 = num1 - int32(prizedraw)
	if num1 < 0 {
		num1 = 0
	}
	relieve = uint32(num2)
	leftdraw = uint32(num1)
	//vip
	//vip := config.GetVip(p.GetVipLevel())
	kick_times := vip.Kick - p.GetKickTimes()
	if kick_times < 0 {
		kicktimes = 0
	} else {
		kicktimes = uint32(kick_times)
	}
	return
}

func buildAgent(ctos *pb.CBuildAgent, p *data.User) (stoc *pb.SBuildAgent) {
	stoc = new(pb.SBuildAgent)
	userid := ctos.GetUserid()
	agent := p.GetAgent()
	if agent == userid {
		stoc.Result = 1
		return
	}
	if agent != "" {
		stoc.Result = 2
		return
	}
	agency := new(data.Agency)
	agency.Get(userid)
	if agency.Agent == "" || agency.Status != 0 || userid == "" {
		stoc.Result = 5
		return
	}
	agencySelf := new(data.Agency)
	agencySelf.Get(p.GetUserid())
	if agencySelf.Agent != "" {
		stoc.Result = 4 //已经是代理商不能绑定
		return
	}
	p.SetAgent(userid)
	stoc.Result = 0
	//日志
	data.BuildRecord(p.GetUserid(), userid)
	//赠送
	var num2 int32 = config.GetEnv(data.ENV3)
	p.AddDiamond(num)
	//消息
	msg := &pb.SPushCurrency{
		Rtype:   uint32(data.LogType19),
		Diamond: int32(num),
	}
	//日志
	data.DiamondRecord(p.GetUserid(), data.LogType19, p.GetDiamond(), num)
	return
}

//1存入,2取出,3赠送
func bank(ctos *pb.CBank, p *data.User) {
	stoc := &pb.SBank{}
	rtype := ctos.GetRtype()
	amount := ctos.GetAmount()
	userid := ctos.GetUserid()
	coin := p.GetCoin()
	//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
	switch rtype {
	case 1: //存入
		if int32(coin-amount) < int32(data.BANKRUPT) {
			//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
			stoc.Error = pb.NotEnoughCoin
		} else if int32(amount) <= 0 {
			stoc.Error = pb.DepositNumberError
			//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
		} else {
			num := -1 * int32(amount)
			p.AddCoin(num)
			p.AddBank(int32(amount))
			//日志
			data.CoinRecord(p.GetUserid(), data.LogType12, p.GetCoin(), num)
			//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
		}
	case 2: //取出
		if amount < data.DRAW_MONEY || amount > p.GetBank() {
			stoc.Error = pb.DrawMoneyNumberError
		} else {
			var tax uint32
			if amount < data.TAX_NUMBER {
				tax = 1
			} else {
				tax = uint32(math.Trunc(float64(amount) * data.GIVE_PERCENT))
			}
			num := -1 * int32(amount)
			coin := int32(amount - tax)
			p.AddCoin(coin)
			p.AddBank(num)
			//日志
			data.CoinRecord(p.GetUserid(), data.LogType14, p.GetCoin(), int32(tax))
			//日志
			data.CoinRecord(p.GetUserid(), data.LogType13, p.GetCoin(), coin)
		}
	case 3: //赠送
		if amount+p.GetGive() > data.GIVE_LIMIT {
			stoc.Error = pb.GiveTooMuch
		} else if amount < data.DRAW_MONEY || amount > p.GetBank() {
			stoc.Error = pb.GiveNumberError
		} else if userid == "" {
			stoc.Error = pb.GiveUseridError
		} else {
			var tax uint32
			if amount < data.TAX_NUMBER {
				tax = 1
			} else {
				tax = uint32(math.Trunc(float64(amount) * data.GIVE_PERCENT))
			}
			num := -1 * int32(amount)
			coin := int32(amount - tax) //实际获得
			//glog.Infof("bank %d", p.GetBank())
			//glog.Infof("userid %s, amount %d, num %d, tax %d, coin %d", userid, amount, num, tax, coin)
			player := players.Get(userid)
			if player == nil {
				user := new(data.User)
				user.GetById(userid)
				if user.Userid == "" || coin <= 0 {
					stoc.Error = pb.GiveUseridError
				} else {
					if user.UpdateCoin(uint32(coin)) {
						p.AddGive(amount)
						p.AddBank(num)
						//日志
						data.CoinRecord(userid, data.LogType15, user.Coin+uint32(coin), coin)
						data.CoinRecord(p.GetUserid(), data.LogType15, p.GetCoin(), (-1 * coin))
						//日志
						data.CoinRecord(p.GetUserid(), data.LogType16, p.GetBank(), int32(tax))
					} else {
						stoc.Error = pb.GiveUseridError
					}
				}
			} else {
				player.AddCoin(coin)
				p.AddGive(amount)
				p.AddBank(num)
				//glog.Infof("userid %s, amount %d, num %d, tax %d, coin %d", userid, amount, num, tax, coin)
				//日志
				data.CoinRecord(userid, data.LogType15, player.GetCoin(), coin)
				data.CoinRecord(p.GetUserid(), data.LogType15, p.GetCoin(), (-1 * coin))
				//日志
				data.CoinRecord(p.GetUserid(), data.LogType16, p.GetBank(), int32(tax))
				//glog.Infof("bank %d", p.GetBank())
			}
		}
	case 4: //查询
	}
	//glog.Infof("rtype %d, bank %d", rtype, p.GetBank())
	stoc.Rtype = rtype
	stoc.Amount = amount
	stoc.Userid = userid
	stoc.Balance = p.GetBank()
	return
}
