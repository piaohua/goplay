// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//解包消息
func Runpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	case 1004:
		msg := new(SWxLogin)
		msg.Code = 1004
		err := msg.Unmarshal(b)
		return msg, err
	case 1062:
		msg := new(SVipList)
		msg.Code = 1062
		err := msg.Unmarshal(b)
		return msg, err
	case 1082:
		msg := new(SDanRanking)
		msg.Code = 1082
		err := msg.Unmarshal(b)
		return msg, err
	case 4000:
		msg := new(SEnterRoom)
		msg.Code = 4000
		err := msg.Unmarshal(b)
		return msg, err
	case 4113:
		msg := new(SPushPaoHu)
		msg.Code = 4113
		err := msg.Unmarshal(b)
		return msg, err
	case 1002:
		msg := new(SRegist)
		msg.Code = 1002
		err := msg.Unmarshal(b)
		return msg, err
	case 4004:
		msg := new(SLeave)
		msg.Code = 4004
		err := msg.Unmarshal(b)
		return msg, err
	case 4017:
		msg := new(SVote)
		msg.Code = 4017
		err := msg.Unmarshal(b)
		return msg, err
	case 4100:
		msg := new(SEnterZiRoom)
		msg.Code = 4100
		err := msg.Unmarshal(b)
		return msg, err
	case 4007:
		msg := new(SDraw)
		msg.Code = 4007
		err := msg.Unmarshal(b)
		return msg, err
	case 4041:
		msg := new(SFreeCamein)
		msg.Code = 4041
		err := msg.Unmarshal(b)
		return msg, err
	case 4044:
		msg := new(SFreeSit)
		msg.Code = 4044
		err := msg.Unmarshal(b)
		return msg, err
	case 4115:
		msg := new(SRoomList)
		msg.Code = 4115
		err := msg.Unmarshal(b)
		return msg, err
	case 2006:
		msg := new(SBroadcast)
		msg.Code = 2006
		err := msg.Unmarshal(b)
		return msg, err
	case 4008:
		msg := new(SDealer)
		msg.Code = 4008
		err := msg.Unmarshal(b)
		return msg, err
	case 4105:
		msg := new(SPushDraw)
		msg.Code = 4105
		err := msg.Unmarshal(b)
		return msg, err
	case 4204:
		msg := new(SPushNewBetting)
		msg.Code = 4204
		err := msg.Unmarshal(b)
		return msg, err
	case 4202:
		msg := new(SPushJackpot)
		msg.Code = 4202
		err := msg.Unmarshal(b)
		return msg, err
	case 2008:
		msg := new(SNotice)
		msg.Code = 2008
		err := msg.Unmarshal(b)
		return msg, err
	case 1020:
		msg := new(SConfig)
		msg.Code = 1020
		err := msg.Unmarshal(b)
		return msg, err
	case 1071:
		msg := new(SMailList)
		msg.Code = 1071
		err := msg.Unmarshal(b)
		return msg, err
	case 1073:
		msg := new(SGetMailItem)
		msg.Code = 1073
		err := msg.Unmarshal(b)
		return msg, err
	case 4062:
		msg := new(SClassicGameover)
		msg.Code = 4062
		err := msg.Unmarshal(b)
		return msg, err
	case 4080:
		msg := new(SPubDraw)
		msg.Code = 4080
		err := msg.Unmarshal(b)
		return msg, err
	case 1028:
		msg := new(SPushCurrency)
		msg.Code = 1028
		err := msg.Unmarshal(b)
		return msg, err
	case 4001:
		msg := new(SCreateRoom)
		msg.Code = 4001
		err := msg.Unmarshal(b)
		return msg, err
	case 4011:
		msg := new(SNiu)
		msg.Code = 4011
		err := msg.Unmarshal(b)
		return msg, err
	case 4018:
		msg := new(SVoteResult)
		msg.Code = 4018
		err := msg.Unmarshal(b)
		return msg, err
	case 4110:
		msg := new(SPushDealerDeal)
		msg.Code = 4110
		err := msg.Unmarshal(b)
		return msg, err
	case 4221:
		msg := new(SLottery)
		msg.Code = 4221
		err := msg.Unmarshal(b)
		return msg, err
	case 3000:
		msg := new(SBuy)
		msg.Code = 3000
		err := msg.Unmarshal(b)
		return msg, err
	case 1080:
		msg := new(SDanInfo)
		msg.Code = 1080
		err := msg.Unmarshal(b)
		return msg, err
	case 4040:
		msg := new(SEnterFreeRoom)
		msg.Code = 4040
		err := msg.Unmarshal(b)
		return msg, err
	case 4060:
		msg := new(SEnterClassicRoom)
		msg.Code = 4060
		err := msg.Unmarshal(b)
		return msg, err
	case 4106:
		msg := new(SPushDiscard)
		msg.Code = 4106
		err := msg.Unmarshal(b)
		return msg, err
	case 2004:
		msg := new(SChatVoice)
		msg.Code = 2004
		err := msg.Unmarshal(b)
		return msg, err
	case 4012:
		msg := new(SGameover)
		msg.Code = 4012
		err := msg.Unmarshal(b)
		return msg, err
	case 4043:
		msg := new(SDealerList)
		msg.Code = 4043
		err := msg.Unmarshal(b)
		return msg, err
	case 4109:
		msg := new(SPushDeal)
		msg.Code = 4109
		err := msg.Unmarshal(b)
		return msg, err
	case 4112:
		msg := new(SPushDealerBu)
		msg.Code = 4112
		err := msg.Unmarshal(b)
		return msg, err
	case 1030:
		msg := new(SBank)
		msg.Code = 1030
		err := msg.Unmarshal(b)
		return msg, err
	case 1072:
		msg := new(SDeleteMail)
		msg.Code = 1072
		err := msg.Unmarshal(b)
		return msg, err
	case 4020:
		msg := new(SGameRecord)
		msg.Code = 4020
		err := msg.Unmarshal(b)
		return msg, err
	case 3004:
		msg := new(SWxpayQuery)
		msg.Code = 3004
		err := msg.Unmarshal(b)
		return msg, err
	case 4050:
		msg := new(SFreeGameover)
		msg.Code = 4050
		err := msg.Unmarshal(b)
		return msg, err
	case 4102:
		msg := new(SZiGameover)
		msg.Code = 4102
		err := msg.Unmarshal(b)
		return msg, err
	case 4046:
		msg := new(SFreeBet)
		msg.Code = 4046
		err := msg.Unmarshal(b)
		return msg, err
	case 4048:
		msg := new(SFreeGamestart)
		msg.Code = 4048
		err := msg.Unmarshal(b)
		return msg, err
	case 3002:
		msg := new(SWxpayOrder)
		msg.Code = 3002
		err := msg.Unmarshal(b)
		return msg, err
	case 1053:
		msg := new(SPrizeList)
		msg.Code = 1053
		err := msg.Unmarshal(b)
		return msg, err
	case 1055:
		msg := new(SPrizeBox)
		msg.Code = 1055
		err := msg.Unmarshal(b)
		return msg, err
	case 1060:
		msg := new(SClassicList)
		msg.Code = 1060
		err := msg.Unmarshal(b)
		return msg, err
	case 1070:
		msg := new(SMailNotice)
		msg.Code = 1070
		err := msg.Unmarshal(b)
		return msg, err
	case 4009:
		msg := new(SPushDealer)
		msg.Code = 4009
		err := msg.Unmarshal(b)
		return msg, err
	case 4201:
		msg := new(SBetting)
		msg.Code = 4201
		err := msg.Unmarshal(b)
		return msg, err
	case 1063:
		msg := new(SPushVip)
		msg.Code = 1063
		err := msg.Unmarshal(b)
		return msg, err
	case 4005:
		msg := new(SKick)
		msg.Code = 4005
		err := msg.Unmarshal(b)
		return msg, err
	case 4010:
		msg := new(SBet)
		msg.Code = 4010
		err := msg.Unmarshal(b)
		return msg, err
	case 4081:
		msg := new(SGetPrize)
		msg.Code = 4081
		err := msg.Unmarshal(b)
		return msg, err
	case 4108:
		msg := new(SOperate)
		msg.Code = 4108
		err := msg.Unmarshal(b)
		return msg, err
	case 4203:
		msg := new(SPushBetting)
		msg.Code = 4203
		err := msg.Unmarshal(b)
		return msg, err
	case 3010:
		msg := new(SShop)
		msg.Code = 3010
		err := msg.Unmarshal(b)
		return msg, err
	case 2003:
		msg := new(SChatText)
		msg.Code = 2003
		err := msg.Unmarshal(b)
		return msg, err
	case 1022:
		msg := new(SUserData)
		msg.Code = 1022
		err := msg.Unmarshal(b)
		return msg, err
	case 1052:
		msg := new(SBankrupts)
		msg.Code = 1052
		err := msg.Unmarshal(b)
		return msg, err
	case 4016:
		msg := new(SLaunchVote)
		msg.Code = 4016
		err := msg.Unmarshal(b)
		return msg, err
	case 4103:
		msg := new(SZiGameRecord)
		msg.Code = 4103
		err := msg.Unmarshal(b)
		return msg, err
	case 4220:
		msg := new(SLotteryInfo)
		msg.Code = 4220
		err := msg.Unmarshal(b)
		return msg, err
	case 4101:
		msg := new(SCreateZiRoom)
		msg.Code = 4101
		err := msg.Unmarshal(b)
		return msg, err
	case 4107:
		msg := new(SPushAuto)
		msg.Code = 4107
		err := msg.Unmarshal(b)
		return msg, err
	case 3006:
		msg := new(SApplePay)
		msg.Code = 3006
		err := msg.Unmarshal(b)
		return msg, err
	case 1026:
		msg := new(SBuildAgent)
		msg.Code = 1026
		err := msg.Unmarshal(b)
		return msg, err
	case 1050:
		msg := new(SPing)
		msg.Code = 1050
		err := msg.Unmarshal(b)
		return msg, err
	case 1054:
		msg := new(SPrizeDraw)
		msg.Code = 1054
		err := msg.Unmarshal(b)
		return msg, err
	case 4006:
		msg := new(SReady)
		msg.Code = 4006
		err := msg.Unmarshal(b)
		return msg, err
	case 4051:
		msg := new(SFreeTrend)
		msg.Code = 4051
		err := msg.Unmarshal(b)
		return msg, err
	case 4200:
		msg := new(SBettingInfo)
		msg.Code = 4200
		err := msg.Unmarshal(b)
		return msg, err
	case 4205:
		msg := new(SBettingRecord)
		msg.Code = 4205
		err := msg.Unmarshal(b)
		return msg, err
	case 1000:
		msg := new(SLogin)
		msg.Code = 1000
		err := msg.Unmarshal(b)
		return msg, err
	case 1024:
		msg := new(SGetCurrency)
		msg.Code = 1024
		err := msg.Unmarshal(b)
		return msg, err
	case 1084:
		msg := new(SDanNotice)
		msg.Code = 1084
		err := msg.Unmarshal(b)
		return msg, err
	case 4003:
		msg := new(SCamein)
		msg.Code = 4003
		err := msg.Unmarshal(b)
		return msg, err
	case 4082:
		msg := new(SPrizeCards)
		msg.Code = 4082
		err := msg.Unmarshal(b)
		return msg, err
	case 1006:
		msg := new(SLoginOut)
		msg.Code = 1006
		err := msg.Unmarshal(b)
		return msg, err
	case 1081:
		msg := new(SQualifying)
		msg.Code = 1081
		err := msg.Unmarshal(b)
		return msg, err
	case 4042:
		msg := new(SFreeDealer)
		msg.Code = 4042
		err := msg.Unmarshal(b)
		return msg, err
	case 4104:
		msg := new(SZiCamein)
		msg.Code = 4104
		err := msg.Unmarshal(b)
		return msg, err
	case 4111:
		msg := new(SPushStatus)
		msg.Code = 4111
		err := msg.Unmarshal(b)
		return msg, err
	default:
		return nil, errors.New("unknown message")
	}
}