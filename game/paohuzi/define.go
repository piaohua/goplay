package paohuzi

const (
	NumCard = 80 //牌数

	TEN  = 10 //牌值
	HAND = 20 //手牌数量

	SEAT uint32 = 3 //玩家数量
)

// suit
const (
	Big   uint32 = 0x10 //大牌
	Small uint32 = 0x00 //小牌

	SuitMask uint32 = 0xF0 //花色掩码
	RankMask uint32 = 0x0F //牌值掩码
)

// 胡>(提=跑)>（偎=碰）>吃
const (
	//手动操作
	CHOW uint32 = 1 << 0 // 吃
	PONG uint32 = 1 << 1 // 碰
	HU   uint32 = 1 << 2 // 胡(代表广义的胡)

	//自动操作
	WEI      uint32 = 1 << 3 // 偎
	CHOU_WEI uint32 = 1 << 4 // 臭偎
	PAO      uint32 = 1 << 5 // 跑
	TI       uint32 = 1 << 6 // 提

	//胡牌方式
	PING_HU uint32 = 1 << 8  // 平摸
	ZIMO_HU uint32 = 1 << 9  // 自摸
	PAO_HU  uint32 = 1 << 10 // 抢杠(破跑胡)
	CHOU_HU uint32 = 1 << 11 // 臭胡
	TIAN_HU uint32 = 1 << 12 // 天胡
	DI_HU   uint32 = 1 << 13 // 地胡
	HONG_HU uint32 = 1 << 14 // 红胡
	DIAN_HU uint32 = 1 << 15 // 点胡
	HONG_WU uint32 = 1 << 16 // 红乌
	WU_HU   uint32 = 1 << 17 // 乌胡
)
