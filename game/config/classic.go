package config

import (
	"sort"
	"sync"

	"goplay/data"
)

// 商城列表
var ClassicList *Classic

type Classic struct {
	sync.RWMutex
	list  map[string]data.Classic
	list2 []data.Classic
}

// 启动初始化
func InitClassic() {
	ClassicList = new(Classic)
	ClassicList.list = make(map[string]data.Classic)
	ClassicList.list2 = make([]data.Classic, 0)
	l := data.GetClassicList()
	for _, v := range l {
		ClassicList.list[v.Id] = v
	}
	syncList()
}

// 启动初始化
func InitClassic2() {
	ClassicList = new(Classic)
	ClassicList.list = make(map[string]data.Classic)
	ClassicList.list2 = make([]data.Classic, 0)
}

func syncList() {
	list2 := make([]data.Classic, 0)
	for _, v := range ClassicList.list {
		list2 = append(list2, v)
	}
	sort.Slice(list2, func(i, j int) bool {
		return list2[i].Minimum > list2[j].Minimum
	})
	ClassicList.list2 = list2
}

// 获取列表
func GetClassics() map[string]data.Classic {
	ClassicList.RLock()
	defer ClassicList.RUnlock()
	return ClassicList.list
}

// 获取列表
func GetClassicsList() []data.Classic {
	ClassicList.RLock()
	defer ClassicList.RUnlock()
	return ClassicList.list2
}

// 添加
func AddClassic(v data.Classic) {
	ClassicList.Lock()
	defer ClassicList.Unlock()
	ClassicList.list[v.Id] = v
	syncList()
}

// 删除
func DelClassic(v data.Classic) {
	ClassicList.Lock()
	defer ClassicList.Unlock()
	delete(ClassicList.list, v.Id)
	syncList()
}

// 获取
func GetClassic(id string) data.Classic {
	ClassicList.RLock()
	defer ClassicList.RUnlock()
	return ClassicList.list[id]
}
