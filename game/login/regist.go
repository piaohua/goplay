package login

import (
	"encoding/json"

	"goplay/data"
	"goplay/glog"
	"goplay/pb"
	"utils"
)

//注册验证
func RegistCheck(ctos *pb.CRegist) (stoc *pb.SRegist) {
	stoc = new(pb.SRegist)
	var nickname string = ctos.GetNickname()
	var phone string = ctos.GetPhone()
	var passwd string = ctos.GetPassword()
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
	return
}

//注册
func Regist(arg *pb.RoleRegist, genid *data.IDGen) (stoc *pb.RoleRegisted) {
	var nickname string = arg.GetNickname()
	var phone string = arg.GetPhone()
	var passwd string = arg.GetPassword()
	var atype uint32 = ctos.GetType()
	stoc = new(pb.RoleRegisted)
	user := new(data.User)
	user.Phone = phone
	if user.ExistsPhone() {
		stoc.Error = pb.PhoneRegisted
		return
	}
	userid := genid.GenID()
	if userid == "" {
		stoc.Error = pb.RegistError
		return
	}
	auth := string(utils.GetAuth())
	user = &data.User{
		Userid:   userid,
		Nickname: nickname,
		Auth:     auth,
		Pwd:      utils.Md5(passwd + auth),
		Phone:    phone,
		Atype:    atype,
		Ctime:    utils.BsonNow(),
	}
	if !user.Save() {
		stoc.Error = pb.RegistError
		return
	}
	result, err := json.Marshal(user)
	if err != nil {
		glog.Errorf("user Marshal err %v", err)
		stoc.Error = pb.RegistError
		return
	}
	stoc.Data = string(result)
	return
}
