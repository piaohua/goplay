package config

//配置初始化
func ConfigInit(id, secret, key, machid, pattern, notifyUrl string) {
	InitNotice()  //公告服务
	InitShop()    //商城服务
	InitBox()     //宝箱服务
	InitPrize()   //抽奖服务
	InitEnv()     //变量服务
	InitClassic() //经典服务
	InitVip()     //VIP服务
	InitLottery() //刮奖服务

	WxLoginInit(id, secret) //微信登录

	WxPayInit(id, key, machid, pattern, notifyUrl) //微信支付
}
