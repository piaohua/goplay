package login

import (
	"goplay/data"
	"goplay/pb"
	"utils"
)

func Regist(ctos *pb.CRegist, genid *data.IDGen) (stoc *pb.SRegist, user *data.User) {
	stoc = new(pb.SRegist)
	var nickname string = ctos.GetNickname()
	var phone string = ctos.GetPhone()
	var passwd string = ctos.GetPassword()
	var ipaddr string = ctos.GetIpaddr()
	var atype uint32 = ctos.GetType()
	if nickname == "" {
		stoc.Error = pb.UsernameEmpty
	}
	if !utils.LegalName(nickname, 7) {
		stoc.Error = pb.NameTooLong
	}
	if !utils.PhoneRegexp(phone) {
		stoc.Error = pb.PhoneNumberError
	}
	if len(passwd) != 32 {
		stoc.Error = pb.PwdFormatError
	}
	//TODO 在线表中查找
	user = &data.User{Phone: phone}
	if user.ExistsPhone() {
		stoc.Error = pb.PhoneRegisted
	}
	if stoc.Error != pb.OK {
		user = nil
		return
	}
	userid := genid.GenID()
	if userid == "" {
		user = nil
		stoc.Error = pb.RegistError
		return
	}
	auth := string(utils.GetAuth())
	user = &data.User{
		Userid:   userid,
		Nickname: nickname,
		RegIp:    ipaddr,
		Auth:     auth,
		Pwd:      utils.Md5(passwd + auth),
		Phone:    phone,
		Atype:    atype,
		Ctime:    utils.BsonNow(),
	}
	if !user.Save() {
		user = nil
		stoc.Error = pb.RegistError
		return
	}
	stoc.Userid = userid
	return
}
