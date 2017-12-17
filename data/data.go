package data

const (
	ROOM_PRIVATE  uint32 = 1 //私人房房间类型(看牌抢庄)
	ROOM_FREE     uint32 = 2 //自由场房间类型
	ROOM_PRIVATE3 uint32 = 3 //私人房房间类型(通比牛牛)
	ROOM_PRIVATE4 uint32 = 4 //私人房房间类型(牛牛坐庄)
	ROOM_PRIVATE5 uint32 = 5 //私人房房间类型(天杠玩法)
	ROOM_PAOHUZI  uint32 = 6 //私人房房间类型(跑胡子)
	ROOM_PAOHUZI2 uint32 = 7 //自由房房间类型(跑胡子)
	ROOM_PAOHUZI3 uint32 = 8 //竞技场房间类型(跑胡子)

	ROOM_LEVEL1 uint32 = 1 //初级场
	ROOM_LEVEL2 uint32 = 2 //中级场
	ROOM_LEVEL3 uint32 = 3 //高级场
	ROOM_LEVEL4 uint32 = 4 //大师场
)

func NewDeskData(round, expire, rtype, ante, cost, xi,
	payment, chat, count, ctime, ltype, min, max uint32,
	rid, creator, rname, invitecode string) *DeskData {
	return &DeskData{
		Rid:     rid,
		Rtype:   rtype,
		Ltype:   ltype,
		Rname:   rname,
		Ante:    ante,
		Payment: payment,
		Cost:    cost,
		Cid:     creator,
		Expire:  expire,
		Round:   round,
		Count:   count,
		Chat:    chat,
		Xi:      xi,
		Code:    invitecode,
		CTime:   ctime,
		Minimum: min,
		Maximum: max,
		Score:   make(map[string]int32),
		Record:  make([]DeskRecord, 0),
	}
}

type DeskData struct {
	Rid        string           `json:"rid"`         //房间ID
	Rtype      uint32           `json:"rtype"`       //房间类型
	Rname      string           `json:"rname"`       //房间名字
	Cid        string           `json:"cid"`         //房间创建人
	Expire     uint32           `json:"expire"`      //牌局设定的过期时间
	Code       string           `json:"code"`        //房间邀请码
	Count      uint32           `json:"count"`       //牌局人数限制
	Chat       uint32           `json:"chat"`        //1对讲机 2语音 3语音转文字
	Round      uint32           `json:"round"`       //牌局数
	Ante       uint32           `json:"ante"`        //私人房底分
	Payment    uint32           `json:"payment"`     //付费方式1=AA or 0=房主支付
	Cost       uint32           `json:"cost"`        //创建消耗
	CTime      uint32           `json:"ctime"`       //创建时间
	Ltype      uint32           `json:"ltype"`       //房间等级类型
	Minimum    uint32           `json:"minimum"`     //房间最低限制
	Maximum    uint32           `json:"maximum"`     //房间最高限制
	Xi         uint32           `json:"xi"`          //起胡息
	Score      map[string]int32 `json:"score"`       //私人局用户战绩积分
	Record     []DeskRecord     `json:"record"`      //记录
	RecordFree DeskFreeRecord   `json:"record_free"` //记录
}

//每轮记录
type DeskRecord struct {
	Seat     uint32   `json:"seat"`     //玩家座位号
	Userid   string   `json:"userid"`   //玩家ID
	Cards    []uint32 `json:"cards"`    //玩家手牌
	Card     uint32   `json:"card"`     //公共牌
	Value    uint32   `json:"value"`    //牌力
	Round    uint32   `json:"round"`    //第几轮
	Score    int32    `json:"score"`    //输赢数量
	Dealer   uint32   `json:"dealer"`   //庄家seat
	Bets     uint32   `json:"bets"`     //下注倍数
	Nickname string   `json:"nickname"` //玩家
	Photo    string   `json:"photo"`    //玩家
}

//每轮记录
type DeskFreeRecord struct {
	Round     uint32                      `json:"round"`     //局数
	Pond      uint32                      `json:"pond"`      //奖池
	Dealer    string                      `json:"dealer"`    //庄家
	Carry     uint32                      `json:"carry"`     //庄家的携带,小于一定值时下庄,字段只做记录,真实数据直接写入玩家数据
	Num       uint32                      `json:"num"`       //当前局下注总数
	Bets      map[string]uint32           `json:"bets"`      //userid:num, 玩家下注金额
	SeatBets  map[uint32]uint32           `json:"seatBets"`  //userid:num, 玩家下注金额
	Tian      map[string]uint32           `json:"tian"`      //天,seat:value
	Di        map[string]uint32           `json:"di"`        //地
	Xuan      map[string]uint32           `json:"xuan"`      //玄
	Huang     map[string]uint32           `json:"huang"`     //黄
	HandCards map[uint32][]uint32         `json:"handCards"` //手牌 seat:cards,seat=(1,2,3,4,5)
	Power     map[uint32]uint32           `json:"power"`     //牌力
	Score     map[uint32]int32            `json:"score"`     //位置(1-5)输赢总量
	Score2    map[string]int32            `json:"score2"`    //闲家输赢总量
	Score3    map[uint32]map[string]int32 `json:"score3"`    //位置上每个玩家输赢
	Seats     map[string]uint32           `json:"seats"`     //userid:seat (seat:1~8)
}

//操作
type DeskOperate struct {
	Seat  uint32   `json:"seat"`  //操作玩家座位号
	Card  uint32   `json:"card"`  //操作牌
	Value uint32   `json:"value"` //掩码
	Cards []uint32 `json:"cards"` //
	Bione []uint32 `json:"bione"` //
	Bitwo []uint32 `json:"bitwo"` //
}
