package handler

/*
import (
	"math"
	"niu/data"
	"niu/images"
	"niu/inter"
	"niu/pb"
	"niu/players"
	"niu/protocol"
	"utils"

	"github.com/golang/protobuf/proto"
)

func vipList(ctos *protocol.CVipList, p inter.IPlayer) {
	stoc := &protocol.SVipList{Error: proto.Uint32(0)}
	list := images.GetVipsList()
	for _, v := range list {
		s := &protocol.Vip{
			Level:  proto.Uint32(uint32(v.Level)),
			Number: proto.Uint32(uint32(v.Number) / 100),
			Pay:    proto.Uint32(v.Pay),
			Prize:  proto.Uint32(v.Prize),
			Kick:   proto.Uint32(uint32(v.Kick)),
		}
		stoc.List = append(stoc.List, s)
	}
	p.Send(stoc)
}

func classicList(ctos *protocol.CClassicList, p inter.IPlayer) {
	stoc := &protocol.SClassicList{Error: proto.Uint32(0)}
	list := images.GetClassics()
	for _, v := range list {
		s := &protocol.Classic{
			Id:      proto.String(v.Id),
			Ptype:   proto.Uint32(uint32(v.Ptype)),
			Rtype:   proto.Uint32(uint32(v.Rtype)),
			Ante:    proto.Uint32(v.Ante),
			Minimum: proto.Uint32(v.Minimum),
			Maximum: proto.Uint32(v.Maximum),
		}
		stoc.List = append(stoc.List, s)
	}
	p.Send(stoc)
}

func prizeBox(ctos *protocol.CPrizeBox, p inter.IPlayer) {
	stoc := &protocol.SPrizeBox{}
	state := ctos.GetState()
	id := p.GetBox()
	num := p.GetDuration()
	d := images.GetBox(id)
	boxstate := p.GetBoxState()
	if d.Id == "" { //不存在或全部完成
		stoc.Error = proto.Uint32(pb.NotBox)
		p.Send(stoc)
		return
	}
	if boxstate == 1 {
		stoc.State = proto.Uint32(2)
		p.Send(stoc)
		return
	}
	if state == 2 && d.Duration <= num { //领取
		l := &protocol.Prize{
			Id:     proto.String(d.Id),
			Rtype:  proto.Uint32(uint32(d.Rtype)),
			Number: proto.Uint32(uint32(d.Amount)),
		}
		stoc.List = append(stoc.List, l)
		addPrize(p, d.Rtype, data.LogType22, d.Amount)
		//下一个宝箱
		d = images.NextBox(d.Duration)
		if d.Id == "" {
			p.SetBox(d.Id, 1)
		} else {
			p.SetBox(d.Id, 0)
		}
		p.SetDuration()
	}
	if d.Id == "" {
		stoc.State = proto.Uint32(2)
		p.Send(stoc)
		return
	}
	n := &protocol.Prize{
		Id:     proto.String(d.Id),
		Rtype:  proto.Uint32(uint32(d.Rtype)),
		Number: proto.Uint32(uint32(d.Amount)),
	}
	stoc.Next = append(stoc.Next, n)
	num = p.GetDuration()
	stoc.Duration = proto.Uint32(d.Duration)
	stoc.Time = proto.Uint32(num)
	if d.Duration <= num {
		stoc.State = proto.Uint32(1)
	}
	p.Send(stoc)
}

func prizeDraw(ctos *protocol.CPrizeDraw, p inter.IPlayer) {
	stoc := &protocol.SPrizeDraw{}
	var num int32
	if images.CheckEnv("prizedraw") {
		num = images.GetEnv("prizedraw").(int32)
	}
	//vip
	vip := images.GetVip(p.GetVipLevel())
	num += int32(vip.Prize) //vip赠送
	draw := p.GetPrizeDraw()
	if int32(draw) >= num {
		stoc.Error = proto.Uint32(pb.NotPrizeDraw)
		p.Send(stoc)
		return
	}
	rate := images.GetPrizeRate()
	var n uint32
	if rate > 0 {
		n = uint32(utils.RandInt64N(int64(rate)))
	}
	list := images.GetPrizes()
	for _, v := range list {
		if n > v.Rate {
			l := &protocol.Prize{
				Id:     proto.String(v.Id),
				Rtype:  proto.Uint32(uint32(v.Rtype)),
				Number: proto.Uint32(uint32(v.Amount)),
			}
			left := num - int32(draw) - 1
			if left < 0 {
				left = 0
			}
			stoc.Leftdraw = proto.Uint32(uint32(left))
			stoc.Prizedraw = proto.Uint32(draw + 1)
			stoc.List = append(stoc.List, l)
			p.SetPrizeDraw()
			p.Send(stoc)
			addPrize(p, v.Rtype, data.LogType21, v.Amount)
			return
		}
	}
	stoc.Error = proto.Uint32(pb.NotGotPrizeDraw)
	p.Send(stoc)
}

func addPrize(p inter.IPlayer, rtype, ltype int, amount int32) {
	switch uint32(rtype) {
	case data.DIAMOND:
		p.AddDiamond(amount)
		//日志
		data.DiamondRecord(p.GetUserid(), ltype, p.GetDiamond(), amount)
		pushCurrency(p, ltype, 0, amount)
	case data.COIN:
		p.AddCoin(amount)
		//日志
		data.CoinRecord(p.GetUserid(), ltype, p.GetCoin(), amount)
		pushCurrency(p, ltype, amount, 0)
	}
}

func pushCurrency(p inter.IPlayer, rtype int, coin, diamond int32) {
	msg := &protocol.SPushCurrency{
		Rtype:   proto.Uint32(uint32(rtype)),
		Diamond: proto.Int32(int32(diamond)),
		Coin:    proto.Int32(int32(coin)),
	}
	p.Send(msg)
}

func prizeList(ctos *protocol.CPrizeList, p inter.IPlayer) {
	stoc := &protocol.SPrizeList{}
	list := images.GetPrizes()
	for _, v := range list {
		l := &protocol.Prize{
			Id:     proto.String(v.Id),
			Rtype:  proto.Uint32(uint32(v.Rtype)),
			Number: proto.Uint32(uint32(v.Amount)),
		}
		stoc.List = append(stoc.List, l)
	}
	p.Send(stoc)
}

func bankrupt(ctos *protocol.CBankrupts, p inter.IPlayer) {
	stoc := &protocol.SBankrupts{}
	var coin1 int32
	if images.CheckEnv("bankrupt_coin") {
		coin1 = images.GetEnv("bankrupt_coin").(int32)
	}
	if int32(p.GetCoin()) >= coin1 {
		stoc.Error = proto.Uint32(pb.NotBankrupt)
		p.Send(stoc)
		return
	}
	var num int32
	if images.CheckEnv("relieve") {
		num = images.GetEnv("relieve").(int32)
	}
	num2 := p.GetBankrupts()
	if int32(num2) > num {
		stoc.Error = proto.Uint32(pb.NotRelieves)
		p.Send(stoc)
		return
	}
	var coin int32
	if images.CheckEnv("relieve_coin") {
		coin = images.GetEnv("relieve_coin").(int32)
	}
	if coin > 0 {
		l := &protocol.Prize{
			Rtype:  proto.Uint32(uint32(data.COIN)),
			Number: proto.Uint32(uint32(coin)),
		}
		left := num - (int32(num2) + 1)
		if left < 0 {
			left = 0
		}
		stoc.Relieve = proto.Uint32(uint32(left))
		stoc.Bankrupt = proto.Uint32(num2 + 1)
		stoc.List = append(stoc.List, l)
		p.AddCoin(coin)
		//日志
		data.CoinRecord(p.GetUserid(), data.LogType11, p.GetCoin(), coin)
		p.SetBankrupts()
		pushCurrency(p, data.LogType11, coin, 0)
	}
	p.Send(stoc)
}

func ping(ctos *protocol.CPing, p inter.IPlayer) {
	stoc := &protocol.SPing{}
	stoc.Time = proto.Uint32(ctos.GetTime())
	p.Send(stoc)
}

func config(ctos *protocol.CConfig, c inter.IConn) {
	stoc := &protocol.SConfig{}
	url := "https://" + data.Conf.ServerHost + ":" + data.Conf.ServerPort + data.Conf.ImagePattern
	stoc.Imageurl = proto.String(url)
	stoc.Version = proto.String(data.Conf.Version)
	c.Send(stoc)
}

func getUserDataHdr(ctos *protocol.CUserData, p inter.IPlayer) {
	stoc := &protocol.SUserData{}
	stoc.Data = &protocol.UserData{}
	userid := ctos.GetUserid()
	if userid == "" {
		stoc.Error = proto.Uint32(pb.UsernameEmpty)
		p.Send(stoc)
		return
	}
	user := &data.User{Userid: userid}
	// 获取玩家自己的详细资料
	if userid == p.GetUserid() {
		stoc.Data.Give = proto.Uint32(p.GetGive())
		stoc.Data.Bank = proto.Uint32(p.GetBank())
		user = p.GetUser().(*data.User)
		iroom := p.GetRoom()
		if iroom != nil {
			rdata := iroom.GetData().(*data.DeskData)
			stoc.Data.Roomtype = proto.Uint32(rdata.Rtype)
			stoc.Data.Invitecode = proto.String(rdata.Code)
			stoc.Data.Roomid = proto.String(rdata.Rid)
		}
		first, relieve, bankrupt, prizedraw, leftdraw, kicktimes := getActivity(p)
		stoc.Data.Data = &protocol.Activity{
			Firstpay:  proto.Uint32(first),
			Relieve:   proto.Uint32(relieve),
			Bankrupt:  proto.Uint32(bankrupt),
			Prizedraw: proto.Uint32(prizedraw),
			Leftdraw:  proto.Uint32(leftdraw),
			Kicktimes: proto.Uint32(kicktimes),
		}
		stoc.Data.Vip = &protocol.VipInfo{
			Level:  proto.Uint32(uint32(p.GetVipLevel())),
			Number: proto.Uint32(p.GetVip() / 100),
		}
	} else {
		player := players.Get(userid) //在线列表中取
		if player != nil {
			user = player.GetUser().(*data.User)
			stoc.Data.Give = proto.Uint32(player.GetGive())
		} else {
			user.Get() //数据库中取
		}
	}
	stoc.Data.Agent = proto.String(user.Agent)
	stoc.Data.Userid = proto.String(userid)
	stoc.Data.Photo = proto.String(user.Photo)
	stoc.Data.Nickname = proto.String(user.Nickname)
	stoc.Data.Sex = proto.Uint32(user.Sex)
	stoc.Data.Phone = proto.String(user.Phone)
	stoc.Data.Coin = proto.Uint32(user.Coin)
	stoc.Data.Diamond = proto.Uint32(user.Diamond)
	p.Send(stoc)
}

func getActivity(p inter.IPlayer) (first, relieve, bankrupt, prizedraw, leftdraw, kicktimes uint32) {
	if p.GetMoney() == 0 {
		first = 1
	}
	//vip
	vip := images.GetVip(p.GetVipLevel())
	//
	bankrupt = p.GetBankrupts()
	prizedraw = p.GetPrizeDraw()
	var num1 int32
	if images.CheckEnv("prizedraw") {
		num1 = images.GetEnv("prizedraw").(int32)
	}
	num1 += int32(vip.Prize) //vip赠送
	var num2 int32
	if images.CheckEnv("relieve") {
		num2 = images.GetEnv("relieve").(int32)
	}
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
	//vip := images.GetVip(p.GetVipLevel())
	kick_times := vip.Kick - p.GetKickTimes()
	if kick_times < 0 {
		kicktimes = 0
	} else {
		kicktimes = uint32(kick_times)
	}
	return
}

func buildAgent(ctos *protocol.CBuildAgent, p inter.IPlayer) {
	stoc := &protocol.SBuildAgent{}
	userid := ctos.GetUserid()
	agent := p.GetAgent()
	if agent == userid {
		stoc.Result = proto.Uint32(1)
	} else if agent != "" {
		stoc.Result = proto.Uint32(2)
	} else {
		agency := new(data.Agency)
		agency.Get(userid)
		if agency.Agent == "" || agency.Status != 0 || userid == "" {
			stoc.Result = proto.Uint32(5)
		} else {
			agencySelf := new(data.Agency)
			agencySelf.Get(p.GetUserid())
			if agencySelf.Agent != "" {
				stoc.Result = proto.Uint32(4) //已经是代理商不能绑定
			} else {
				p.SetAgent(userid)
				stoc.Result = proto.Uint32(0)
				//日志
				data.BuildRecord(p.GetUserid(), userid)
				//赠送
				var num int32 = 30
				if images.CheckEnv("build") {
					num = images.GetEnv("build").(int32)
				}
				p.AddDiamond(num)
				//消息
				msg := &protocol.SPushCurrency{
					Rtype:   proto.Uint32(uint32(data.LogType19)),
					Diamond: proto.Int32(int32(num)),
				}
				p.Send(msg)
				//日志
				data.DiamondRecord(p.GetUserid(), data.LogType19, p.GetDiamond(), num)
			}
		}
	}
	p.Send(stoc)
}

func getCurrency(ctos *protocol.CGetCurrency, p inter.IPlayer) {
	stoc := &protocol.SGetCurrency{}
	stoc.Coin = proto.Uint32(p.GetCoin())
	stoc.Diamond = proto.Uint32(p.GetDiamond())
	stoc.Roomcard = proto.Uint32(p.GetRoomCard())
	p.Send(stoc)
}

//1存入,2取出,3赠送
func bank(ctos *protocol.CBank, p inter.IPlayer) {
	stoc := &protocol.SBank{}
	rtype := ctos.GetRtype()
	amount := ctos.GetAmount()
	userid := ctos.GetUserid()
	coin := p.GetCoin()
	//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
	switch rtype {
	case 1: //存入
		if int32(coin-amount) < int32(data.BANKRUPT) {
			//glog.Infof("coin %d, userid %s, rtype %d, amount %d", coin, userid, rtype, amount)
			stoc.Error = proto.Uint32(pb.NotEnoughCoin)
		} else if int32(amount) <= 0 {
			stoc.Error = proto.Uint32(pb.DepositNumberError)
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
			stoc.Error = proto.Uint32(pb.DrawMoneyNumberError)
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
			stoc.Error = proto.Uint32(pb.GiveTooMuch)
		} else if amount < data.DRAW_MONEY || amount > p.GetBank() {
			stoc.Error = proto.Uint32(pb.GiveNumberError)
		} else if userid == "" {
			stoc.Error = proto.Uint32(pb.GiveUseridError)
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
					stoc.Error = proto.Uint32(pb.GiveUseridError)
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
						stoc.Error = proto.Uint32(pb.GiveUseridError)
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
	stoc.Rtype = proto.Uint32(rtype)
	stoc.Amount = proto.Uint32(amount)
	stoc.Userid = proto.String(userid)
	stoc.Balance = proto.Uint32(p.GetBank())
	p.Send(stoc)
}
*/
