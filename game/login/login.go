package login

import (
	"goplay/data"
	"goplay/pb"
	"utils"
)

func Login(ctos *pb.CLogin) (stoc *pb.SLogin, user *data.User) {
	stoc = new(pb.SLogin)
	var phone string = ctos.GetPhone()
	var passwd string = ctos.GetPassword()
	var atype uint32 = ctos.GetType()
	//glog.Infoln("login phone -> ", phone)
	if !utils.PhoneRegexp(phone) {
		stoc.Error = pb.PhoneNumberError
	}
	//TODO 在线表中查找
	user = &data.User{Phone: phone}
	if !user.VerifyPwdByPhone(passwd) {
		stoc.Error = pb.UsernameOrPwdError
	}
	if user.Userid == "" {
		stoc.Error = pb.LoginError
	}
	if stoc.Error != pb.OK {
		user = nil
		return
	}
	user.Atype = atype
	//重复登录检测
	stoc.Userid = user.Userid
	return
}
