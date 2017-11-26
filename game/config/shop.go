package config

import (
	"sort"
	"sync"

	"goplay/data"
	"utils"
)

// 商城列表
var ShopList *Shop

type Shop struct {
	sync.RWMutex
	list map[string]data.Shop
}

// 启动初始化
func InitShop() {
	ShopList = new(Shop)
	ShopList.list = make(map[string]data.Shop)
	l := data.GetShopList()
	for _, v := range l {
		ShopList.list[v.Id] = v
	}
}

// 获取商品列表
func GetShops(atype uint32) []data.Shop {
	ShopList.RLock()
	defer ShopList.RUnlock()
	list := make([]data.Shop, 0)
	for _, v := range ShopList.list {
		if v.Del > 0 {
			continue
		}
		if v.Etime.Before(utils.BsonNow()) {
			continue
		}
		//不是自己包的商品
		if atype != v.Atype && v.Atype != 0 {
			continue
		}
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Price < list[j].Price
	})
	return list
}

// 添加新的商品
func AddShop(v data.Shop) {
	ShopList.Lock()
	defer ShopList.Unlock()
	if v.Del > 0 || v.Etime.Before(utils.BsonNow()) {
		delete(ShopList.list, v.Id)
	} else {
		ShopList.list[v.Id] = v
	}
}

// 获取商品
func GetShop(id string) data.Shop {
	ShopList.RLock()
	defer ShopList.RUnlock()
	return ShopList.list[id]
}
