package config

import "api/wxapi"

var WxLogin *wxapi.Wxapi //微信登录

func WxLoginInit(id, secret string) {
	cfg := &wxapi.WxapiConfig{
		AppId:           id,
		AppSecret:       secret,
		AccessUrl:       "https://api.weixin.qq.com/sns/oauth2/access_token",
		RefreshUrl:      "https://api.weixin.qq.com/sns/oauth2/refresh_token",
		UserinfoUrl:     "https://api.weixin.qq.com/sns/userinfo",
		VerifyAccessUrl: "https://api.weixin.qq.com/sns/auth",
	}
	wx, err := wxapi.NewWxapi(cfg)
	if err != nil {
		panic(err)
	}
	WxLogin = wx
}
