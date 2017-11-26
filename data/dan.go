package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

/*
qualifying 排位赛

ranking 排行榜

tasklist 任务列表

dan 段
stars 星
points 积分

topDan 最高段位
curDan 当前段位
topRank 最高排名
curRank 当前排名

number 局数
wins 胜利局数
winRate 胜率

season 赛季

start 开始时间
end 截止时间
*/

//bronze
//silver
//gold
//crystal
//master
//Champion
//Emperor
// 0 青铜杯
// 1 白银杯
// 2 黄金杯
// 3 水晶杯
// 4 大师杯
// 5 冠军杯
// 6 帝王杯

type DanList struct {
	Id     string    `bson:"_id"`    //序号
	Name   string    `bson:"name"`   //名称
	Dan    int       `bson:"dan"`    //段0-6
	Ante   uint32    `bson:"ante"`   //隐藏底分
	Number uint32    `bson:"number"` //预计局数
	Coin   uint32    `bson:"coin"`   //最低限制金币
	Di     uint32    `bson:"di"`     //房间底分
	Level  []Stars   `bson:"level"`  //星级,积分
	Ctime  time.Time `bson:"ctime"`  //创建时间
	Del    int       `bson:"-"`      //是否移除
}

type Stars struct {
	Stars  int   `bson:"stars"`  //星
	Points int32 `bson:"points"` //积分
}

func GetDanList() []DanList {
	var list []DanList
	ListByQ(DanLists, nil, &list)
	return list
}

//级别
type DanLevel struct {
	Dan    int
	Stars  int
	Points int32
	Prev   int
	Curr   int
	Next   int
}

//任务列表
type TaskList struct {
	Id      string    `bson:"_id"`     //任务序号
	Name    string    `bson:"name"`    //任务名称
	Diamond uint32    `bson:"diamond"` //钻石奖励
	Coin    uint32    `bson:"coin"`    //金币奖励
	Ctime   time.Time `bson:"ctime"`   //创建时间
	Del     int       `bson:"-"`       //是否移除
}

func GetTaskList() []TaskList {
	var list []TaskList
	ListByQ(TaskLists, nil, &list)
	return list
}

//当前赛季
type Season struct {
	Start     uint32 `bson:"start"`     //开始时间(时间截)
	End       uint32 `bson:"end"`       //结束时间(时间截)
	Remaining uint32 `bson:"remaining"` //剩余时间(秒)
}

func (this *Season) Get() {
	GetByQ(Seasons, nil, this)
}

func (this *Season) Save() bool {
	return Upsert(Seasons, nil, this)
}

func (this *Season) Delete() bool {
	return Delete(Seasons, bson.M{"start": this.Start})
}

//

//个人任务列表
//TODO record
type DanTask struct {
	Userid string         `bson:"_id"`  //玩家ID
	List   map[string]int `bson:"list"` //任务列表map[taskid]status
	//Taskid string `bson:"taskid"` //任务序号
	//Status int    `bson:"status"` //状态0未达成,1可领取,2已领取
	//TODO 按赛季处理
	//Start     uint32 `bson:"start"`     //开始时间(时间截)
}

func (this *DanTask) Get() {
	Get(DanTasks, this.Userid, this)
}

func (this *DanTask) Save() bool {
	return Upsert(DanTasks, bson.M{"_id": this.Userid}, this)
}

func SaveDanTask(list map[string]*DanTask) {
	for _, v := range list {
		v.Save()
	}
}

//排行榜
//TODO record
type DanRanking struct {
	Userid   string `bson:"_id"`      //玩家ID
	Photo    string `bson:"photo"`    //玩家头像
	Nickname string `bson:"nickname"` //玩家昵称
	Rank     uint32 `bson:"rank"`     //排名
	Dan      int    `bson:"dan"`      //玩家段位
	Stars    int    `bson:"stars"`    //玩家星级
	Points   int32  `bson:"points"`   //玩家积分
}

func GetDanRanking() []DanRanking {
	var list []DanRanking
	ListByQ(DanRankings, nil, &list)
	return list
}

func (this *DanRanking) Save() bool {
	return Insert(DanRankings, this)
}

func SaveDanRanking(list []DanRanking) {
	for k, _ := range list {
		d := &list[k]
		d.Save()
	}
}

//个人战绩
//TODO record
type DanCombat struct {
	Userid   string `bson:"_id"`       // 用户id
	Dan      int    `bson:"dan"`       // 段
	Stars    int    `bson:"stars"`     // 星
	Points   int32  `bson:"points"`    // 积分
	Number   uint32 `bson:"number"`    // 总局数
	Wins     uint32 `bson:"wins"`      // 胜利数
	Rank     uint32 `bson:"rank"`      // 当前排名
	TopDan   int    `bson:"top_dan"`   // 最高段位
	TopRank  uint32 `bson:"top_rank"`  // 最高排名
	CurIndex int    `bson:"cur_index"` // 最高排名
	//TODO 按赛季处理
	//Start     uint32 `bson:"start"`     //开始时间(时间截)
}

func (this *DanCombat) Get() {
	Get(DanCombats, this.Userid, this)
}

func (this *DanCombat) Save() bool {
	return Upsert(DanCombats, bson.M{"_id": this.Userid}, this)
}

func SaveDanCombat(list map[string]*DanCombat) {
	for _, v := range list {
		v.Save()
	}
}
