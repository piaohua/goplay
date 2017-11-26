package login

import (
	"api/wxapi"
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
	"utils"
)

func WxLogin(ctos *pb.CWxLogin, genid *data.IDGen) (stoc *pb.SWxLogin, user *data.User) {
	var wxcode string = ctos.GetWxcode()
	var token string = ctos.GetToken()
	var ipaddr string = ctos.GetIpaddr()
	var atype uint32 = ctos.GetType()
	//glog.Infof("weixinLogin wxcode:%s, token:%s", wxcode, token)
	var isRegist bool
	wxdata := new(data.WxLoginData)
	//token登录
	if token != "" {
		err := loginByToken(token, wxdata)
		if err != nil {
			glog.Errorf("weixinLogin err : %v", err)
			stoc.Error = pb.WechatLoingFailReAuth
			token = "" //重置为空，重新授权
		} else {
			user, isRegist = wxLogin2(wxdata, ipaddr, genid)
			if user == nil {
				glog.Errorf("weixinLogin err : %v", err)
				stoc.Error = pb.GetWechatUserInfoFail
			}
			token = wxdata.RefreshToken
		}
	} else if wxcode != "" { //wxcode登录
		err := loginByCode(wxcode, wxdata)
		if err != nil {
			glog.Errorf("weixinLogin err : %v", err)
			stoc.Error = pb.WechatLoingFailReAuth
		} else {
			user, isRegist = wxLogin2(wxdata, ipaddr, genid)
			if user == nil {
				glog.Errorf("weixinLogin err : %v", err)
				stoc.Error = pb.GetWechatUserInfoFail
			}
			token = wxdata.RefreshToken
		}
	} else {
		stoc.Error = pb.WechatLoingFailReAuth
	}
	if stoc.Error != pb.OK {
		return
	}
	user.Atype = atype
	stoc.Userid = user.Userid
	stoc.Token = token
	stoc.Isreg = isRegist
	return
}

func wxLogin2(wxdata *data.WxLoginData,
	ipaddr string, genid *data.IDGen) (*data.User, bool) {
	//TODO 在线表中查找
	user := &data.User{Wxuid: wxdata.OpenId}
	user.GetByWechat()
	if user.Userid != "" {
		return user, false
	}
	//isregist
	userid := genid.GenID()
	user.Wxuid = wxdata.OpenId
	user.Nickname = wxdata.Nickname
	user.Photo = wxdata.HeadImagUrl
	user.Sex = uint32(wxdata.Sex)
	user.Userid = userid
	user.RegIp = ipaddr
	user.Ctime = utils.BsonNow()
	if !user.Save() {
		return nil, false
	}
	return user, true
}

//直接使用refresh_token

//refresh_token登录
func loginByToken(refresh_token string, wxdata *data.WxLoginData) error {
	//刷新refresh_token
	refreshResult, err := config.WxLogin.Refresh(refresh_token)
	if err != nil {
		return err
	}
	//获取个人信息
	userinfoResult, err := config.WxLogin.UserInfo(refreshResult.Openid, refreshResult.AccessToken)
	if err != nil {
		return err
	}
	wxdata.AccessToken = refreshResult.AccessToken
	wxdata.ExpiresIn = refreshResult.ExpiresIn
	wxdata.RefreshToken = refreshResult.RefreshToken
	loginData(wxdata, userinfoResult)
	return nil
}

//wxcode登录
func loginByCode(wxcode string, wxdata *data.WxLoginData) error {
	//获取access_token
	accessResult, err := config.WxLogin.Auth(wxcode)
	if err != nil {
		return err
	}
	//获取个人信息
	userinfoResult, err := config.WxLogin.UserInfo(accessResult.OpenId, accessResult.AccessToken)
	if err != nil {
		return err
	}
	wxdata.AccessToken = accessResult.AccessToken
	wxdata.ExpiresIn = accessResult.ExpiresIn
	wxdata.RefreshToken = accessResult.RefreshToken
	loginData(wxdata, userinfoResult)
	return nil
}

func loginData(wxdata *data.WxLoginData,
	userinfo wxapi.UserInfoResult) {
	wxdata.OpenId = userinfo.OpenId
	wxdata.Nickname = userinfo.Nickname
	wxdata.Sex = userinfo.Sex
	wxdata.Province = userinfo.Province
	wxdata.City = userinfo.City
	wxdata.Country = userinfo.Country
	wxdata.HeadImagUrl = userinfo.HeadImagUrl
	wxdata.Privilege = userinfo.Privilege
	wxdata.UnionId = userinfo.UnionId
}
