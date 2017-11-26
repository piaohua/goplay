package config

import (
	"sort"
	"sync"

	"goplay/data"
)

// 商城列表
var VipList *Vip

type Vip struct {
	sync.RWMutex
	list  map[int]data.Vip
	list2 []data.Vip
}

// 启动初始化
func InitVip() {
	VipList = new(Vip)
	VipList.list = make(map[int]data.Vip)
	VipList.list2 = make([]data.Vip, 0)
	l := data.GetVipList()
	for _, v := range l {
		VipList.list[v.Level] = v
	}
	syncVipList()
}

func syncVipList() {
	list2 := make([]data.Vip, 0)
	for _, v := range VipList.list {
		list2 = append(list2, v)
	}
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].Level > list2[j].Level
	})
	VipList.list2 = list2
}

// 获取列表
func GetVips() map[int]data.Vip {
	VipList.RLock()
	defer VipList.RUnlock()
	return VipList.list
}

// 获取列表
func GetVipsList() []data.Vip {
	VipList.RLock()
	defer VipList.RUnlock()
	return VipList.list2
}

// 添加
func AddVip(v data.Vip) {
	VipList.Lock()
	defer VipList.Unlock()
	VipList.list[v.Level] = v
	syncVipList()
}

// 删除
func DelVip(v data.Vip) {
	VipList.Lock()
	defer VipList.Unlock()
	delete(VipList.list, v.Level)
	syncVipList()
}

// 获取
func GetVip(level int) data.Vip {
	VipList.RLock()
	defer VipList.RUnlock()
	return VipList.list[level]
}

// 获取
func GetVipLevel(num uint32) int {
	VipList.RLock()
	defer VipList.RUnlock()
	for _, v := range VipList.list2 {
		if num >= v.Number {
			return v.Level
		}
	}
	return 0
}

// 获取
func GetVipPay(level int) uint32 {
	VipList.RLock()
	defer VipList.RUnlock()
	if v, ok := VipList.list[level]; ok {
		return v.Pay
	}
	return 0
}
