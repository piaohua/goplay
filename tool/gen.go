package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	Init()
	Gen()
}

//TODO 内部协议通信

var (
	protoPacket = make(map[string]uint32) //响应协议
	protoUnpack = make(map[string]uint32) //请求协议
	packetPath  = "../pb/packet.go"       //打包协议文件路径
	unpackPath  = "../pb/unpack.go"       //解包协议文件路径
	rPacketPath = "../pb/rpacket.go"      //机器人打包协议
	rUnpackPath = "../pb/runpack.go"      //机器人解包协议
	luaPath     = "./MsgID.lua"           //lua文件
	packetFunc  = "Packet"                //打包协议函数名字
	unpackFunc  = "Unpack"                //解包协议函数名字
	rPacketFunc = "Rpacket"               //机器人打包协议函数名字
	rUnpackFunc = "Runpack"               //机器人解包协议函数名字
)

type proto struct {
	name string
	code uint32
}

var protosUnpack = []proto{
	//buy
	{name: "CBuy", code: 3000},
	{name: "CWxpayOrder", code: 3002},
	{name: "CWxpayQuery", code: 3004},
	{name: "CApplePay", code: 3006},
	{name: "CShop", code: 3010},
	//chat
	{name: "CChatText", code: 2003},
	{name: "CChatVoice", code: 2004},
	{name: "CNotice", code: 2008},
	//login
	{name: "CLogin", code: 1000},
	{name: "CRegist", code: 1002},
	{name: "CWxLogin", code: 1004},
	//user
	{name: "CConfig", code: 1020},
	{name: "CUserData", code: 1022},
	{name: "CGetCurrency", code: 1024},
	{name: "CBuildAgent", code: 1026},
	{name: "CBank", code: 1030},
	{name: "CPing", code: 1050},
	{name: "CBankrupts", code: 1052},
	{name: "CPrizeList", code: 1053},
	{name: "CPrizeDraw", code: 1054},
	{name: "CPrizeBox", code: 1055},
	{name: "CClassicList", code: 1060},
	{name: "CVipList", code: 1062},
	{name: "CMailList", code: 1071},
	{name: "CDeleteMail", code: 1072},
	{name: "CGetMailItem", code: 1073},
	{name: "CDanInfo", code: 1080},
	{name: "CQualifying", code: 1081},
	{name: "CDanRanking", code: 1082},
	//room
	{name: "CEnterRoom", code: 4000},
	{name: "CCreateRoom", code: 4001},
	{name: "CLeave", code: 4004},
	{name: "CKick", code: 4005},
	{name: "CReady", code: 4006},
	{name: "CDealer", code: 4008},
	{name: "CBet", code: 4010},
	{name: "CNiu", code: 4011},
	{name: "CLaunchVote", code: 4016},
	{name: "CVote", code: 4017},
	{name: "CGameRecord", code: 4020},
	{name: "CEnterFreeRoom", code: 4040},
	{name: "CFreeDealer", code: 4042},
	{name: "CDealerList", code: 4043},
	{name: "CFreeSit", code: 4044},
	{name: "CFreeBet", code: 4046},
	{name: "CFreeTrend", code: 4051},
	{name: "CEnterClassicRoom", code: 4060},
	{name: "CGetPrize", code: 4081},
	{name: "CPrizeCards", code: 4082},
	{name: "CEnterZiRoom", code: 4100},
	{name: "CCreateZiRoom", code: 4101},
	{name: "CZiGameRecord", code: 4103},
	{name: "CPushDiscard", code: 4106},
	{name: "COperate", code: 4108},
	{name: "CRoomList", code: 4115},
	//betting
	{name: "CBettingInfo", code: 4200},
	{name: "CBetting", code: 4201},
	{name: "CBettingRecord", code: 4205},
	//lottery
	{name: "CLotteryInfo", code: 4220},
	{name: "CLottery", code: 4221},
}

var protosPacket = []proto{
	//buy
	{name: "SBuy", code: 3000},
	{name: "SWxpayOrder", code: 3002},
	{name: "SWxpayQuery", code: 3004},
	{name: "SApplePay", code: 3006},
	{name: "SShop", code: 3010},
	//chat
	{name: "SChatText", code: 2003},
	{name: "SChatVoice", code: 2004},
	{name: "SBroadcast", code: 2006},
	{name: "SNotice", code: 2008},
	//login
	{name: "SLogin", code: 1000},
	{name: "SRegist", code: 1002},
	{name: "SWxLogin", code: 1004},
	{name: "SLoginOut", code: 1006},
	//user
	{name: "SConfig", code: 1020},
	{name: "SUserData", code: 1022},
	{name: "SGetCurrency", code: 1024},
	{name: "SBuildAgent", code: 1026},
	{name: "SPushCurrency", code: 1028},
	{name: "SBank", code: 1030},
	{name: "SPing", code: 1050},
	{name: "SBankrupts", code: 1052},
	{name: "SPrizeList", code: 1053},
	{name: "SPrizeDraw", code: 1054},
	{name: "SPrizeBox", code: 1055},
	{name: "SClassicList", code: 1060},
	{name: "SVipList", code: 1062},
	{name: "SPushVip", code: 1063},
	{name: "SMailNotice", code: 1070},
	{name: "SMailList", code: 1071},
	{name: "SDeleteMail", code: 1072},
	{name: "SGetMailItem", code: 1073},
	{name: "SDanInfo", code: 1080},
	{name: "SQualifying", code: 1081},
	{name: "SDanRanking", code: 1082},
	{name: "SDanNotice", code: 1084},
	//room
	{name: "SEnterRoom", code: 4000},
	{name: "SCreateRoom", code: 4001},
	{name: "SCamein", code: 4003},
	{name: "SLeave", code: 4004},
	{name: "SKick", code: 4005},
	{name: "SReady", code: 4006},
	{name: "SDraw", code: 4007},
	{name: "SDealer", code: 4008},
	{name: "SPushDealer", code: 4009},
	{name: "SBet", code: 4010},
	{name: "SNiu", code: 4011},
	{name: "SGameover", code: 4012},
	{name: "SLaunchVote", code: 4016},
	{name: "SVote", code: 4017},
	{name: "SVoteResult", code: 4018},
	{name: "SGameRecord", code: 4020},
	{name: "SEnterFreeRoom", code: 4040},
	{name: "SFreeCamein", code: 4041},
	{name: "SFreeDealer", code: 4042},
	{name: "SDealerList", code: 4043},
	{name: "SFreeSit", code: 4044},
	{name: "SFreeBet", code: 4046},
	{name: "SFreeGamestart", code: 4048},
	{name: "SFreeGameover", code: 4050},
	{name: "SFreeTrend", code: 4051},
	{name: "SEnterClassicRoom", code: 4060},
	{name: "SClassicGameover", code: 4062},
	{name: "SPubDraw", code: 4080},
	{name: "SGetPrize", code: 4081},
	{name: "SPrizeCards", code: 4082},
	{name: "SEnterZiRoom", code: 4100},
	{name: "SCreateZiRoom", code: 4101},
	{name: "SZiGameover", code: 4102},
	{name: "SZiGameRecord", code: 4103},
	{name: "SZiCamein", code: 4104},
	{name: "SPushDraw", code: 4105},
	{name: "SPushDiscard", code: 4106},
	{name: "SPushAuto", code: 4107},
	{name: "SOperate", code: 4108},
	{name: "SPushDeal", code: 4109},
	{name: "SPushDealerDeal", code: 4110},
	{name: "SPushStatus", code: 4111},
	{name: "SPushDealerBu", code: 4112},
	{name: "SPushPaoHu", code: 4113},
	{name: "SRoomList", code: 4115},
	//betting
	{name: "SBettingInfo", code: 4200},
	{name: "SBetting", code: 4201},
	{name: "SPushJackpot", code: 4202},
	{name: "SPushBetting", code: 4203},
	{name: "SPushNewBetting", code: 4204},
	{name: "SBettingRecord", code: 4205},
	//lottery
	{name: "SLotteryInfo", code: 4220},
	{name: "SLottery", code: 4221},
}

//初始化
func Init() {
	//request
	for _, v := range protosUnpack {
		registUnpack(v.name, v.code)
	}
	//response
	for _, v := range protosPacket {
		registPacket(v.name, v.code)
	}
	//最后生成MsgID.lua文件
	genMsgID()
}

func registUnpack(key string, code uint32) {
	if _, ok := protoUnpack[key]; ok {
		panic(fmt.Sprintf("%s registered: %d", key, code))
	}
	protoUnpack[key] = code
}

func registPacket(key string, code uint32) {
	if _, ok := protoPacket[key]; ok {
		panic(fmt.Sprintf("%s registered: %d", key, code))
	}
	protoPacket[key] = code
}

//生成文件
func Gen() {
	gen_packet()
	gen_unpack()
	//client
	gen_client_packet()
	gen_client_unpack()
}

//生成打包文件
func gen_packet() {
	var str string
	str += head_packet()
	str += body_packet()
	str += end_packet()
	err := ioutil.WriteFile(packetPath, []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

//生成解包文件
func gen_unpack() {
	var str string
	str += head_unpack()
	str += body_unpack()
	str += end_unpack()
	err := ioutil.WriteFile(unpackPath, []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

func body_unpack() string {
	var str string
	for k, v := range protoUnpack {
		//str += fmt.Sprintf("case %d:\n\t\tmsg := new(%s)\n\t\t%s\n\t", v, k, result_unpack())
		str += fmt.Sprintf("case %d:\n\t\tmsg := new(%s)\n\t\t%s\n\t\t%s\n\t", v, k, body_unpack_code(v), result_unpack())
	}
	return str
}

func body_packet() string {
	var str string
	for k, v := range protoPacket {
		//str += fmt.Sprintf("case *%s:\n\t\tb, err := msg.(*%s).Marshal()\n\t\t%s\n\t", k, k, result_packet(v))
		str += fmt.Sprintf("case *%s:\n\t\t%s\n\t\tb, err := msg.(*%s).Marshal()\n\t\t%s\n\t", k, body_packet_code(v, k), k, result_packet(v))
	}
	return str
}

func body_unpack_code(code uint32) (str string) {
	str = fmt.Sprintf("msg.Code = %d", code)
	return
}

func body_packet_code(code uint32, name string) (str string) {
	str = fmt.Sprintf("msg.(*%s).Code = %d", name, code)
	return
}

func head_packet() string {
	return fmt.Sprintf(`// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//打包消息
func Packet(msg interface{}) (uint32, []byte, error) {
	switch msg.(type) {
	`)
}

func head_unpack() string {
	return fmt.Sprintf(`// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//解包消息
func Unpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	`)
}

func head_rpacket() string {
	return fmt.Sprintf(`// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//打包消息
func Rpacket(msg interface{}) (uint32, []byte, error) {
	switch msg.(type) {
	`)
}

func head_runpack() string {
	return fmt.Sprintf(`// Code generated by tool/gen.go.
// DO NOT EDIT!

package pb

import (
	"errors"
)

//解包消息
func Runpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	`)
}

func result_packet(code uint32) string {
	return fmt.Sprintf("return %d, b, err", code)
}

func result_unpack() string {
	return fmt.Sprintf(`err := msg.Unmarshal(b)
		return msg, err`)
}

func end_packet() string {
	return fmt.Sprintf(`default:
		return 0, []byte{}, errors.New("unknown message")
	}
}`)
}

func end_unpack() string {
	return fmt.Sprintf(`default:
		return nil, errors.New("unknown message")
	}
}`)
}

//生成lua文件
func genMsgID() {
	var str string
	str += fmt.Sprintf("msgID = {")
	for k, v := range protoUnpack {
		str += fmt.Sprintf("\n\t%s = %d,", k, v)
	}
	str += fmt.Sprintf("\n")
	for k, v := range protoPacket {
		str += fmt.Sprintf("\n\t%s = %d,", k, v)
	}
	str += fmt.Sprintf("\n}")
	err := ioutil.WriteFile(luaPath, []byte(str), 0666)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

//生成机器人打包文件
func gen_client_packet() {
	var str string
	str += head_rpacket()
	str += body_client_packet()
	str += end_packet()
	err := ioutil.WriteFile(rPacketPath, []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

//生成机器人解包文件
func gen_client_unpack() {
	var str string
	str += head_runpack()
	str += body_client_unpack()
	str += end_unpack()
	err := ioutil.WriteFile(rUnpackPath, []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

func body_client_packet() string {
	var str string
	for k, v := range protoUnpack {
		//str += fmt.Sprintf("case *%s:\n\t\tb, err := msg.(*%s).Marshal()\n\t\t%s\n\t", k, k, result_packet(v))
		str += fmt.Sprintf("case *%s:\n\t\t%s\n\t\tb, err := msg.(*%s).Marshal()\n\t\t%s\n\t", k, body_client_packet_code(v, k), k, result_packet(v))
	}
	return str
}

func body_client_unpack() string {
	var str string
	for k, v := range protoPacket {
		//str += fmt.Sprintf("case %d:\n\t\tmsg := new(%s)\n\t\t%s\n\t", v, k, result_unpack())
		str += fmt.Sprintf("case %d:\n\t\tmsg := new(%s)\n\t\t%s\n\t\t%s\n\t", v, k, body_client_unpack_code(v), result_unpack())
	}
	return str
}

func body_client_unpack_code(code uint32) (str string) {
	str = fmt.Sprintf("msg.Code = %d", code)
	return
}

func body_client_packet_code(code uint32, name string) (str string) {
	str = fmt.Sprintf("msg.(*%s).Code = %d", name, code)
	return
}
