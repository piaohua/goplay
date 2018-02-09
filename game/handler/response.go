/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-12-27 23:16:38
 * Filename      : response.go
 * Description   : 响应协议消息
 * *******************************************************/
package handler

import (
	"goplay/data"
	"goplay/pb"
)

//坐下,起来离开消息
func FreeSitMsg(userid string, seat uint32,
	state bool, p *data.User) interface{} {
	return &pb.SFreeSit{
		Seat:     seat,
		State:    state,
		Userid:   userid,
		Nickname: p.GetName(),
		Photo:    p.GetPhoto(),
		Coin:     p.GetCoin(),
	}
}

//下注消息
func FreeBetMsg(seat, beseat, val, coin,
	bets uint32, userid string) interface{} {
	return &pb.SFreeBet{
		Seat:   seat,
		Beseat: beseat,
		Value:  val,
		Coin:   coin,
		Bets:   bets,
		Userid: userid,
	}
}

//上下庄消息
func FreeBeDealerMsg(state, coin uint32, dealer,
	userid, name string) interface{} {
	return &pb.SFreeDealer{
		State:    state,
		Coin:     coin,
		Userid:   userid,
		Dealer:   dealer,
		Nickname: name,
	}
}

//开始消息
func FreeStartMsg(dealer, photo string, state, carry,
	dealerNum, left uint32) interface{} {
	return &pb.SFreeGamestart{
		Dealer:        dealer,
		Photo:         photo,
		State:         state,
		Coin:          carry,
		DealerNum:     dealerNum,
		LeftDealerNum: left,
	}
}

//文本聊天消息
func ChatMsg(seat uint32, userid string, msg []byte) interface{} {
	return &pb.SChatText{
		Seat:    seat,
		Userid:  userid,
		Content: msg,
	}
}

//语音聊天消息
func ChatMsg2(seat uint32, userid string, msg []byte) interface{} {
	return &pb.SChatVoice{
		Seat:    seat,
		Userid:  userid,
		Content: msg,
	}
}

//进入消息
func CameinMsg(v *data.User, ready bool,
	seat uint32, score int32) interface{} {
	user := new(pb.RoomUser)
	UserInfoMsg(v, user)
	user.Score = score
	user.Ready = ready
	user.Seat = seat
	user.Vip = &pb.VipInfo{
		Level:  uint32(v.GetVipLevel()),
		Number: v.GetVip(),
	}
	return &pb.SCamein{Userinfo: user}
}

//离开消息
func LeaveMsg(userid string, seat uint32) interface{} {
	return &pb.SLeave{
		Seat:   seat,
		Userid: userid,
	}
}

func res_ready(seat uint32, ready bool) interface{} {
	return &pb.SReady{
		Seat:  seat,
		Ready: ready,
	}
}

func DrawMsg(seat, state uint32, cards []uint32) interface{} {
	return &pb.SDraw{
		State: state,
		Seat:  seat,
		Cards: cards,
	}
}

func res_dealer(seat uint32, dealer bool, num uint32) interface{} {
	return &pb.SDealer{
		Dealer: dealer,
		Seat:   seat,
		Num:    num,
	}
}

func res_pushdealer(dealer uint32) interface{} {
	return &pb.SPushDealer{
		Dealer: dealer,
	}
}

func res_bet(seat, seatBet, value uint32) interface{} {
	return &pb.SBet{
		Value:   value,
		Seat:    seat,
		Seatbet: seatBet,
	}
}

func res_niu(seat, value uint32, cards []uint32) interface{} {
	return &pb.SNiu{
		Value: value,
		Seat:  seat,
		Cards: cards,
	}
}

func UserInfoMsg2(v *data.User, r *pb.FreeUser) {
	r.Userid = v.GetUserid()
	r.Nickname = v.GetName()
	r.Phone = v.GetPhone()
	r.Photo = v.GetPhoto()
	r.Sex = v.GetSex()
	r.Coin = v.GetCoin()
	r.Diamond = v.GetDiamond()
}

func UserInfoMsg(v *data.User, r *pb.RoomUser) {
	r.Userid = v.GetUserid()
	r.Nickname = v.GetName()
	r.Phone = v.GetPhone()
	r.Photo = v.GetPhoto()
	r.Sex = v.GetSex()
	r.Coin = v.GetCoin()
	r.Diamond = v.GetDiamond()
}

func RoomInfoMsg2(id string, d *data.DeskData, r *pb.FreeRoom) {
	r.Roomid = id
	r.Rtype = d.Rtype
	r.Rname = d.Rname
	r.Count = d.Count
	r.Ante = d.Ante
	r.Chat = d.Chat
}
