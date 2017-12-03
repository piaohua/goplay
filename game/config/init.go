package config

//配置初始化
func ConfigInit() {
	InitNotice()  //公告服务
	InitShop()    //商城服务
	InitBox()     //宝箱服务
	InitPrize()   //抽奖服务
	InitEnv()     //变量服务
	InitClassic() //经典服务
	InitVip()     //VIP服务
	InitLottery() //刮奖服务
}

//节点变量初始化, 节点连接时同步数据
func Init2Gate(id, secret, key, machid, pattern, notifyUrl string) {
	InitNotice2()  //公告服务
	InitShop2()    //商城服务
	InitBox2()     //宝箱服务
	InitPrize2()   //抽奖服务
	InitEnv2()     //变量服务
	InitClassic2() //经典服务
	InitVip2()     //VIP服务
	InitLottery2() //刮奖服务

	WxLoginInit(id, secret) //微信登录

	WxPayInit(id, key, machid, pattern, notifyUrl) //微信支付
}
