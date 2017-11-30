// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//解包消息
func Unpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	case 2008:
		msg := new(CNotice)
		msg.Code = 2008
		err := msg.Unmarshal(b)
		return msg, err
	case 1020:
		msg := new(CConfig)
		msg.Code = 1020
		err := msg.Unmarshal(b)
		return msg, err
	case 1022:
		msg := new(CUserData)
		msg.Code = 1022
		err := msg.Unmarshal(b)
		return msg, err
	case 1030:
		msg := new(CBank)
		msg.Code = 1030
		err := msg.Unmarshal(b)
		return msg, err
	case 4005:
		msg := new(CKick)
		msg.Code = 4005
		err := msg.Unmarshal(b)
		return msg, err
	case 4051:
		msg := new(CFreeTrend)
		msg.Code = 4051
		err := msg.Unmarshal(b)
		return msg, err
	case 2003:
		msg := new(CChatText)
		msg.Code = 2003
		err := msg.Unmarshal(b)
		return msg, err
	case 4000:
		msg := new(CEnterRoom)
		msg.Code = 4000
		err := msg.Unmarshal(b)
		return msg, err
	case 4010:
		msg := new(CBet)
		msg.Code = 4010
		err := msg.Unmarshal(b)
		return msg, err
	case 4082:
		msg := new(CPrizeCards)
		msg.Code = 4082
		err := msg.Unmarshal(b)
		return msg, err
	case 1026:
		msg := new(CBuildAgent)
		msg.Code = 1026
		err := msg.Unmarshal(b)
		return msg, err
	case 4100:
		msg := new(CEnterZiRoom)
		msg.Code = 4100
		err := msg.Unmarshal(b)
		return msg, err
	case 1055:
		msg := new(CPrizeBox)
		msg.Code = 1055
		err := msg.Unmarshal(b)
		return msg, err
	case 1082:
		msg := new(CDanRanking)
		msg.Code = 1082
		err := msg.Unmarshal(b)
		return msg, err
	case 4205:
		msg := new(CBettingRecord)
		msg.Code = 4205
		err := msg.Unmarshal(b)
		return msg, err
	case 4220:
		msg := new(CLotteryInfo)
		msg.Code = 4220
		err := msg.Unmarshal(b)
		return msg, err
	case 3002:
		msg := new(CWxpayOrder)
		msg.Code = 3002
		err := msg.Unmarshal(b)
		return msg, err
	case 3010:
		msg := new(CShop)
		msg.Code = 3010
		err := msg.Unmarshal(b)
		return msg, err
	case 1073:
		msg := new(CGetMailItem)
		msg.Code = 1073
		err := msg.Unmarshal(b)
		return msg, err
	case 4001:
		msg := new(CCreateRoom)
		msg.Code = 4001
		err := msg.Unmarshal(b)
		return msg, err
	case 4200:
		msg := new(CBettingInfo)
		msg.Code = 4200
		err := msg.Unmarshal(b)
		return msg, err
	case 1054:
		msg := new(CPrizeDraw)
		msg.Code = 1054
		err := msg.Unmarshal(b)
		return msg, err
	case 1062:
		msg := new(CVipList)
		msg.Code = 1062
		err := msg.Unmarshal(b)
		return msg, err
	case 4016:
		msg := new(CLaunchVote)
		msg.Code = 4016
		err := msg.Unmarshal(b)
		return msg, err
	case 1081:
		msg := new(CQualifying)
		msg.Code = 1081
		err := msg.Unmarshal(b)
		return msg, err
	case 4011:
		msg := new(CNiu)
		msg.Code = 4011
		err := msg.Unmarshal(b)
		return msg, err
	case 4106:
		msg := new(CPushDiscard)
		msg.Code = 4106
		err := msg.Unmarshal(b)
		return msg, err
	case 1053:
		msg := new(CPrizeList)
		msg.Code = 1053
		err := msg.Unmarshal(b)
		return msg, err
	case 1060:
		msg := new(CClassicList)
		msg.Code = 1060
		err := msg.Unmarshal(b)
		return msg, err
	case 1080:
		msg := new(CDanInfo)
		msg.Code = 1080
		err := msg.Unmarshal(b)
		return msg, err
	case 3004:
		msg := new(CWxpayQuery)
		msg.Code = 3004
		err := msg.Unmarshal(b)
		return msg, err
	case 1052:
		msg := new(CBankrupts)
		msg.Code = 1052
		err := msg.Unmarshal(b)
		return msg, err
	case 3006:
		msg := new(CApplePay)
		msg.Code = 3006
		err := msg.Unmarshal(b)
		return msg, err
	case 1000:
		msg := new(CLogin)
		msg.Code = 1000
		err := msg.Unmarshal(b)
		return msg, err
	case 1004:
		msg := new(CWxLogin)
		msg.Code = 1004
		err := msg.Unmarshal(b)
		return msg, err
	case 1072:
		msg := new(CDeleteMail)
		msg.Code = 1072
		err := msg.Unmarshal(b)
		return msg, err
	case 4042:
		msg := new(CFreeDealer)
		msg.Code = 4042
		err := msg.Unmarshal(b)
		return msg, err
	case 1002:
		msg := new(CRegist)
		msg.Code = 1002
		err := msg.Unmarshal(b)
		return msg, err
	case 4008:
		msg := new(CDealer)
		msg.Code = 4008
		err := msg.Unmarshal(b)
		return msg, err
	case 4103:
		msg := new(CZiGameRecord)
		msg.Code = 4103
		err := msg.Unmarshal(b)
		return msg, err
	case 4201:
		msg := new(CBetting)
		msg.Code = 4201
		err := msg.Unmarshal(b)
		return msg, err
	case 1050:
		msg := new(CPing)
		msg.Code = 1050
		err := msg.Unmarshal(b)
		return msg, err
	case 4081:
		msg := new(CGetPrize)
		msg.Code = 4081
		err := msg.Unmarshal(b)
		return msg, err
	case 4108:
		msg := new(COperate)
		msg.Code = 4108
		err := msg.Unmarshal(b)
		return msg, err
	case 2004:
		msg := new(CChatVoice)
		msg.Code = 2004
		err := msg.Unmarshal(b)
		return msg, err
	case 4017:
		msg := new(CVote)
		msg.Code = 4017
		err := msg.Unmarshal(b)
		return msg, err
	case 4044:
		msg := new(CFreeSit)
		msg.Code = 4044
		err := msg.Unmarshal(b)
		return msg, err
	case 4020:
		msg := new(CGameRecord)
		msg.Code = 4020
		err := msg.Unmarshal(b)
		return msg, err
	case 4043:
		msg := new(CDealerList)
		msg.Code = 4043
		err := msg.Unmarshal(b)
		return msg, err
	case 4101:
		msg := new(CCreateZiRoom)
		msg.Code = 4101
		err := msg.Unmarshal(b)
		return msg, err
	case 4004:
		msg := new(CLeave)
		msg.Code = 4004
		err := msg.Unmarshal(b)
		return msg, err
	case 4006:
		msg := new(CReady)
		msg.Code = 4006
		err := msg.Unmarshal(b)
		return msg, err
	case 4040:
		msg := new(CEnterFreeRoom)
		msg.Code = 4040
		err := msg.Unmarshal(b)
		return msg, err
	case 4221:
		msg := new(CLottery)
		msg.Code = 4221
		err := msg.Unmarshal(b)
		return msg, err
	case 3000:
		msg := new(CBuy)
		msg.Code = 3000
		err := msg.Unmarshal(b)
		return msg, err
	case 1024:
		msg := new(CGetCurrency)
		msg.Code = 1024
		err := msg.Unmarshal(b)
		return msg, err
	case 1071:
		msg := new(CMailList)
		msg.Code = 1071
		err := msg.Unmarshal(b)
		return msg, err
	case 4046:
		msg := new(CFreeBet)
		msg.Code = 4046
		err := msg.Unmarshal(b)
		return msg, err
	case 4060:
		msg := new(CEnterClassicRoom)
		msg.Code = 4060
		err := msg.Unmarshal(b)
		return msg, err
	case 4115:
		msg := new(CRoomList)
		msg.Code = 4115
		err := msg.Unmarshal(b)
		return msg, err
	default:
		return nil, errors.New("unknown message")
	}
}