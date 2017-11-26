package config

import (
	"sort"
	"sync"

	"goplay/data"
)

// 列表
var PrizeList *Prize

type Prize struct {
	sync.RWMutex
	rate uint32
	list []data.Prize
}

// 启动初始化
func InitPrize() {
	PrizeList = new(Prize)
	PrizeList.list = make([]data.Prize, 0)
	l := data.GetPrizeList()
	PrizeList.list = append(PrizeList.list, l...)
	sortPrize()
	ratePrize()
}

// 获取列表
func GetPrizes() []data.Prize {
	PrizeList.RLock()
	defer PrizeList.RUnlock()
	return PrizeList.list
}

// 添加新的
func AddPrize(v data.Prize) {
	PrizeList.Lock()
	defer PrizeList.Unlock()
	for i, p := range PrizeList.list {
		if p.Id == v.Id {
			if v.Del > 0 {
				PrizeList.list = append(PrizeList.list[:i], PrizeList.list[i+1:]...)
				ratePrize()
				return
			}
			PrizeList.list[i] = v
			ratePrize()
			return
		}
	}
	PrizeList.list = append(PrizeList.list, v)
	sortPrize()
	ratePrize()
}

func sortPrize() {
	sort.Slice(PrizeList.list, func(i, j int) bool {
		return PrizeList.list[i].Rate > PrizeList.list[j].Rate
	})
}

func ratePrize() {
	PrizeList.rate = 0
	for _, v := range PrizeList.list {
		PrizeList.rate += v.Rate
	}
}

func GetPrizeRate() uint32 {
	PrizeList.RLock()
	defer PrizeList.RUnlock()
	return PrizeList.rate
}

// 获取
func GetPrize(id string) data.Prize {
	PrizeList.RLock()
	defer PrizeList.RUnlock()
	for _, p := range PrizeList.list {
		if p.Id == id {
			return p
		}
	}
	return data.Prize{}
}
