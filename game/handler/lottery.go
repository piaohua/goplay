package handle

import (
	"math/rand"
	"time"

	"goplay/data"
	"goplay/game/algo"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
)

func LotteryInfo(ctos *pb.CLotteryInfo) (stoc *pb.SLotteryInfo) {
	stoc = new(pb.SLotteryInfo)
	l := config.GetLotterys()
	for k, v := range l {
		msg := &pb.Lottery{
			Niu:    k,
			Number: v,
		}
		stoc.List = append(stoc.List, msg)
	}
	num, max := get_lettery()
	stoc.Single = uint32(num)
	stoc.Maxnumber = uint32(max)
	return
}

func Lottery(ctos *pb.CLottery, p *data.User) (stoc *pb.SLottery,
	number, prize uint32, ok bool) {
	stoc = new(pb.SLottery)
	times := ctos.GetTimes()
	num, max := get_lettery()
	if times <= 0 || int32(times) > max {
		stoc.Error = pb.NotTimes
		return
	}
	number = times * uint32(num) //消耗
	glog.Info("number ", times, num, max, number)
	if p.GetDiamond() < number {
		stoc.Error = pb.NotEnoughDiamond
		return
	}
	cards := shuffle()
	key, niu := algo.Lottery(cards)
	stoc.Cards = cards
	stoc.Niu = niu
	prize, ok = config.GetLottery(key)
	prize = prize * times //注数
	glog.Info("prize ", prize)
	stoc.Number = prize
	return
}

//洗牌
func shuffle() (cards []uint32) {
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

func get_lettery() (diamond, max int32) {
	diamond = config.GetEnv(data.ENV15)
	max = config.GetEnv(data.ENV16)
	return
}
