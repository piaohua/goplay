package config

import (
	"sync"

	"goplay/data"
)

// 列表
var LotteryList *Lottery

type Lottery struct {
	sync.RWMutex
	//1有牛,2牛牛,3五花牛,4五小牛,5四炸
	prize map[uint32]uint32 //奖励
}

// 启动初始化
func InitLottery() {
	LotteryList = new(Lottery)
	LotteryList.prize = make(map[uint32]uint32)
	l := data.GetLotteryList()
	for _, v := range l {
		LotteryList.prize[v.Times] = v.Diamond
	}
}

// 启动初始化
func InitLottery2() {
	LotteryList = new(Lottery)
	LotteryList.prize = make(map[uint32]uint32)
}

// 获取列表
func GetLotterys() map[uint32]uint32 {
	LotteryList.RLock()
	defer LotteryList.RUnlock()
	return LotteryList.prize
}

// 获取
func GetLottery(times uint32) (uint32, bool) {
	LotteryList.RLock()
	defer LotteryList.RUnlock()
	if v, ok := LotteryList.prize[times]; ok {
		return v, true
	}
	return 0, false
}

// 设置
func SetLottery(times, diamond uint32) {
	LotteryList.Lock()
	defer LotteryList.Unlock()
	LotteryList.prize[times] = diamond
}
