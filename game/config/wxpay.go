package config

import "api/wxpay"

var Apppay *wxpay.AppTrans //微信支付

func WxPayInit(id, key, machid, pattern, notifyUrl string) {
	cfg := &wxpay.WxConfig{
		AppId:         id,
		AppKey:        key,
		MchId:         machid,
		NotifyPattern: pattern,
		NotifyUrl:     notifyUrl,
		PlaceOrderUrl: "https://api.mch.weixin.qq.com/pay/unifiedorder",
		QueryOrderUrl: "https://api.mch.weixin.qq.com/pay/orderquery",
		TradeType:     "APP",
	}
	appTrans, err := wxpay.NewAppTrans(cfg)
	if err != nil {
		panic(err)
	}
	Apppay = appTrans
	//go Apppay.RecvNotify(wxRecvTrade) //goroutine
}
