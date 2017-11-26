/*
1.战绩里需要给玩家nickname,photo,coin,diamond四个字段（战绩里详细记录需要显示昵称和头像，要么别人不知道谁是谁，点击头像会有玩家的信息，那里面需要显示钻石和金币）
2.百人场的休息时间还是改为原来的时间，现在我玩着发现时间太长，很容易就不想玩了

3.福利：
1）转盘，每天只有一次机会，每天凌晨重置（转盘表文祥发在群里）
2）救济金，救济金不足2000才可领取，每天次数2次（这个最好后台能够配置）
3）首充礼包：首次充值送双倍钻石和2万金币（没人仅有一次）
4）宝箱（在线时长奖励，后台配置在线时长及对应的奖励）：
统计每天用户累计登录在线时长多久，
达到指定时间后点击宝箱获取奖励
（假10分钟到但是玩家没有领取3分钟奖励和5分钟奖励，
则用户可以一次领取3/5/10分钟的奖励），每天零点更新
*/

/*
type 1 金币 2钻石
id	duration	type	number
1	30	1	1888
2	60	2	3
3	120	1	4766
4	300	2	8
5	900	1	13888
6	1800	2	38


type 2 金币 1钻石
id rate type amount
1  50   1    1000
2  50   1    1000
3  50   1    1000
4  50   1    1000
5  50   1    1000
*/
package config

import (
	"sort"
	"sync"

	"goplay/data"
	"utils"
)

// 公告列表
var NoticeList *Notice

type Notice struct {
	sync.RWMutex
	list map[string]data.Notice
}

// 启动初始化
func InitNotice() {
	NoticeList = new(Notice)
	NoticeList.list = make(map[string]data.Notice)
	l := data.GetNoticeList(data.NOTICE_TYPE1)
	for _, v := range l {
		NoticeList.list[v.Id] = v
	}
}

// 获取消息列表
func GetNotices(atype uint32) []data.Notice {
	NoticeList.RLock()
	defer NoticeList.RUnlock()
	tops := make([]data.Notice, 0)
	list := make([]data.Notice, 0)
	for _, v := range NoticeList.list {
		if v.Del > 0 {
			continue
		}
		if v.Etime.Before(utils.BsonNow()) {
			continue
		}
		//不是自己包的消息
		if atype != v.Atype && v.Atype != 0 {
			continue
		}
		if v.Top > 0 {
			tops = append(tops, v)
			continue
		}
		list = append(list, v)
	}
	sort.Slice(tops, func(i, j int) bool {
		return tops[i].Ctime.After(tops[j].Ctime)
	})
	sort.Slice(list, func(i, j int) bool {
		return list[i].Ctime.After(list[j].Ctime)
	})
	return append(tops, list...)
}

// 添加新的公告
func AddNotice(v data.Notice) {
	NoticeList.Lock()
	defer NoticeList.Unlock()
	if v.Del > 0 || v.Etime.Before(utils.BsonNow()) {
		delete(NoticeList.list, v.Id)
	} else {
		NoticeList.list[v.Id] = v
	}
}
