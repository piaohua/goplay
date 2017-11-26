package config

import (
	"sort"
	"sync"

	"goplay/data"
)

// 列表
var BoxList *Box

type Box struct {
	sync.RWMutex
	list []data.Box
}

// 启动初始化
func InitBox() {
	BoxList = new(Box)
	BoxList.list = make([]data.Box, 0)
	l := data.GetBoxList()
	BoxList.list = append(BoxList.list, l...)
	sortBox()
}

// 获取列表
func GetBoxs(atype uint32) []data.Box {
	BoxList.RLock()
	defer BoxList.RUnlock()
	return BoxList.list
}

func sortBox() {
	sort.Slice(BoxList.list, func(i, j int) bool {
		return BoxList.list[i].Duration < BoxList.list[j].Duration
	})
}

// 添加新的
func AddBox(v data.Box) {
	BoxList.Lock()
	defer BoxList.Unlock()
	for i, b := range BoxList.list {
		if b.Id == v.Id {
			if v.Del > 0 {
				BoxList.list = append(BoxList.list[:i], BoxList.list[i+1:]...)
				return
			}
			BoxList.list[i] = v
			return
		}
	}
	BoxList.list = append(BoxList.list, v)
	sortBox()
}

// 获取
func GetBox(id string) data.Box {
	BoxList.RLock()
	defer BoxList.RUnlock()
	if id == "" {
		if len(BoxList.list) > 0 {
			return BoxList.list[0]
		}
	}
	for _, b := range BoxList.list {
		if b.Id == id {
			return b
		}
	}
	return data.Box{}
}

// 获取
func NextBox(duration uint32) data.Box {
	BoxList.RLock()
	defer BoxList.RUnlock()
	for _, b := range BoxList.list {
		if b.Duration > duration {
			return b
		}
	}
	return data.Box{}
}
