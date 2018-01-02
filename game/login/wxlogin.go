package login

import (
	"api/wxapi"
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
	"utils"

	jsoniter "github.com/json-iterator/go"
)

//微信登录
func WxLogin(ctos *pb.WxLogin, user *data.User,
	genid *data.IDGen) (stoc *pb.WxLogined) {
	stoc = new(pb.WxLogined)
	var wxuid string = ctos.GetWxuid()
	var nickname string = ctos.GetNickname()
	var photo string = ctos.GetPhoto()
	var sex uint32 = ctos.GetSex()
	var atype uint32 = ctos.GetType()
	//已经在线
	if user.GetUserid() != "" {
		user.Wxuid = wxuid
		user.Nickname = nickname
		user.Photo = photo
		user.Sex = sex
		user.Atype = atype
		result, err := jsoniter.Marshal(user)
		if err != nil {
			glog.Errorf("user Marshal err %v", err)
			stoc.Error = pb.GetWechatUserInfoFail
			return
		}
		stoc.Data = string(result)
		return
	}
	//离线
	user.Wxuid = wxuid
	user.GetByWechat()
	if user.Userid != "" {
		//登录
		result, err := jsoniter.Marshal(user)
		if err != nil {
			glog.Errorf("user Marshal err %v", err)
			stoc.Error = pb.GetWechatUserInfoFail
			return
		}
		stoc.Data = string(result)
		return
	}
	//注册
	userid := genid.GenID()
	user.Userid = userid
	user.Wxuid = wxuid
	user.Nickname = nickname
	user.Photo = photo
	user.Sex = sex
	user.Atype = atype
	user.Ctime = utils.BsonNow()
	if !user.Save() {
		glog.Errorf("WxLogin failed : %s", userid)
		stoc.Error = pb.GetWechatUserInfoFail
		return
	}
	result, err := jsoniter.Marshal(user)
	if err != nil {
		glog.Errorf("user Marshal err %v", err)
		stoc.Error = pb.GetWechatUserInfoFail
		return
	}
	stoc.Data = string(result)
	stoc.IsRegist = true
	return
}

//微信登录验证
func WxLoginCheck(ctos *pb.CWxLogin) (stoc *pb.SWxLogin,
	wxdata *data.WxLoginData) {
	var wxcode string = ctos.GetWxcode()
	var token string = ctos.GetToken()
	//glog.Infof("weixinLogin wxcode:%s, token:%s", wxcode, token)
	wxdata = new(data.WxLoginData)
	//token登录
	if token != "" {
		err := loginByToken(token, wxdata)
		if err != nil {
			glog.Errorf("weixinLogin err : %v", err)
			stoc.Error = pb.WechatLoingFailReAuth
			token = "" //重置为空，重新授权
		} else {
			token = wxdata.RefreshToken
		}
	} else if wxcode != "" { //wxcode登录
		err := loginByCode(wxcode, wxdata)
		if err != nil {
			glog.Errorf("weixinLogin err : %v", err)
			stoc.Error = pb.WechatLoingFailReAuth
		} else {
			token = wxdata.RefreshToken
		}
	} else {
		stoc.Error = pb.WechatLoingFailReAuth
	}
	if stoc.Error != pb.OK {
		return
	}
	stoc.Token = token
	return
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
