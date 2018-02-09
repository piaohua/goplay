/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-12-17 18:22:30
 * Filename      : free.go
 * Description   : 自由场协议消息请求
 * *******************************************************/
package handler

import (
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"

	jsoniter "github.com/json-iterator/go"
)

//自由场数据
func NewDeskData(rtype uint32) *data.DeskData {
	switch rtype {
	case data.ROOM_FREE:
		//自由场数据
		count := uint32(config.GetEnv(data.ENV11))
		return data.NewDeskData(0, 0, data.ROOM_FREE, 1,
			0, 0, 0, 0, count, 0, 0, 0, 0, "", "", "", "")
	}
	return nil
}

//自由场数据
func FreeData() string {
	//now := uint32(utils.Timestamp())
	count := uint32(config.GetEnv(data.ENV11))
	deskData := data.NewDeskData(0, 0, data.ROOM_FREE, 1,
		0, 0, 0, 0, count, 0, 0, 0, 0, "", "", "", "")
	return Desk2Data(deskData)
}

//打包
func Desk2Data(deskData *data.DeskData) string {
	result, err := jsoniter.Marshal(deskData)
	if err != nil {
		glog.Errorf("Desk2Data Marshal err %v", err)
		return ""
	}
	return string(result)
}

//解析
func Data2Desk(deskDataStr string) *data.DeskData {
	deskData := new(data.DeskData)
	err := jsoniter.Unmarshal([]byte(deskDataStr), deskData)
	if err != nil {
		glog.Errorf("Data2Desk Unmarshal err %v", err)
		return nil
	}
	return deskData
}

/*
import (
	"goplay/data"
	"goplay/game/config"
	"goplay/pb"
	"niu/desk"
	"niu/rooms"
	"utils"
)

// 进入自由场
func EntryFreeRoom(ctos *pb.CEnterFreeRoom, p *data.User) (stoc *pb.SEnterFreeRoom) {
	stoc = new(pb.SEnterFreeRoom)
	//匹配可以进入的房间
	rdata = rooms.MatchFree()
	if rdata == nil {
		now := uint32(utils.Timestamp())
		rid := rooms.GenID()
		code := rooms.GenInvitecodeFree()
		count := uint32(config.GetEnv(data.ENV11))
		deskData := data.NewDeskData(0, 0, 2, 1, 0, 0, 0, 0, count, now, 0, 0, 0, rid, "", "", code)
		rdata = desk.NewDeskFree(deskData)
		rooms.AddFree(code, rdata)
	}
	if rdata == nil {
		stoc.Error = pb.RoomNotExist
		return
	}
	d := rdata.GetData()
	if d == nil {
		stoc.Error = pb.RoomNotExist
		return
	}
	var code int = rdata.Enter(p)
	switch code {
	case 1:
		stoc.Error = pb.RoomFull
	}
}
*/

/*
//操作
func freesit(ctos *pb.CFreeSit, p inter.IPlayer) {
	stoc := &pb.SFreeSit{}
	rdata := p.GetRoom()
	if rdata == nil {
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
		return
	}
	var seat uint32 = ctos.GetSeat()
	var state bool = ctos.GetState()
	//glog.Infof("sit free room -> %s, %d, %v", p.GetUserid(), seat, state)
	if !(seat >= 1 && seat <= 8) {
		stoc.Error = proto.Uint32(pb.OperateError)
		p.Send(stoc)
		return
	}
	if p.GetCoin() < desk.LIMIT_SIT {
		stoc.Error = proto.Uint32(pb.NotEnoughCoin)
		p.Send(stoc)
		return
	}
	err := rdata.SitDown(p.GetUserid(), seat, state)
	//glog.Infof("sit free room err -> %d", err)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(pb.SitDownFailed)
		p.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(pb.StandUpFailed)
		p.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(pb.DealerSitFailed)
		p.Send(stoc)
	}
}

//操作
func freedealer(ctos *pb.CFreeDealer, p inter.IPlayer) {
	stoc := &pb.SFreeDealer{}
	rdata := p.GetRoom()
	if rdata == nil {
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
		return
	}
	var state uint32 = ctos.GetState()
	var num uint32 = ctos.GetCoin()
	//glog.Infof("dealer free room -> %s, %d, %d", p.GetUserid(), num, state)
	if p.GetCoin() < num {
		stoc.Error = proto.Uint32(pb.NotEnoughCoin)
		p.Send(stoc)
		return
	}
	code := rdata.BeDealer(p.GetUserid(), state, num)
	//glog.Infof("dealer free room err -> %d", code)
	switch code {
	case 1:
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(pb.BeDealerAlready)
		p.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(pb.BeDealerAlreadySit)
		p.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(pb.BeDealerNotEnough)
		p.Send(stoc)
	case 5:
		stoc.Error = proto.Uint32(pb.GameStartedCannotLeave)
		p.Send(stoc)
	}
}

//操作
func freedealerlist(ctos *pb.CDealerList, p inter.IPlayer) {
	stoc := &pb.SDealerList{}
	rdata := p.GetRoom()
	if rdata == nil {
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
		return
	}
	//glog.Infof("dealer list free room -> %s", p.GetUserid())
	rdata.BeDealerList(p.GetUserid())
}

//操作
func freebet(ctos *pb.CFreeBet, p inter.IPlayer) {
	stoc := &pb.SFreeBet{}
	rdata := p.GetRoom()
	if rdata == nil {
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
		return
	}
	value := ctos.GetValue()
	seatBet := ctos.GetSeat()
	//glog.Infof("bet free room -> %s, %d, %d", p.GetUserid(), value, seatBet)
	if !(seatBet >= 2 && seatBet <= 5) {
		stoc.Error = proto.Uint32(pb.OperateError)
		p.Send(stoc)
		return
	}
	if p.GetCoin() < value {
		stoc.Error = proto.Uint32(pb.NotEnoughCoin)
		p.Send(stoc)
		return
	}
	err := rdata.ChoiceBet(p.GetUserid(), seatBet, value)
	//glog.Infof("bet free room err -> %d", err)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(pb.GameNotStart)
		p.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(pb.BetDealerFailed)
		p.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(pb.BetTopLimit)
		p.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
	case 5:
		stoc.Error = proto.Uint32(pb.BetNotSeat)
		p.Send(stoc)
	}
}

//房间趋势
func getTrend(ctos *pb.CFreeTrend, p inter.IPlayer) {
	stoc := &pb.SFreeTrend{}
	rdata := p.GetRoom()
	if rdata == nil {
		//stoc.Error = proto.Uint32(pb.NotInRoom)
		p.Send(stoc)
		return
	}
	d := rdata.GetData()
	if d == nil {
		//stoc.Error = proto.Uint32(pb.RoomNotExist)
		p.Send(stoc)
		return
	}
	rid := d.(*data.DeskData).Rid
	list, err := data.GetTrends(rid, 1)
	if err != nil {
		p.Send(stoc)
		return
	}
	var i uint32
	for k, v := range list {
		l := &pb.FreeTrendList{
			//Round: proto.Uint32(v.Round),
			Round: proto.Uint32(uint32(k + 1)),
		}
		for i = 2; i <= 5; i++ {
			s := &pb.FreeTrend{
				Seat: proto.Uint32(i),
			}
			switch i {
			case 2:
				s.Win = proto.Bool(v.Seat2)
			case 3:
				s.Win = proto.Bool(v.Seat3)
			case 4:
				s.Win = proto.Bool(v.Seat4)
			case 5:
				s.Win = proto.Bool(v.Seat5)
			}
			l.List = append(l.List, s)
		}
		stoc.List = append(stoc.List, l)
	}
	p.Send(stoc)
}
*/
