package handler

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"
	"niu/betting"
)

func BettingBet(ctos *pb.CBetting, p *data.User) {
	stoc := new(pb.SBetting)
	seat := ctos.GetSeat()
	number := ctos.GetNumber()
	stoc.Seat = seat
	stoc.Number = number
	if p.GetDiamond() < number {
		stoc.Error = pb.NotEnoughDiamond
		p.Send(stoc)
		return
	}
	userid := p.GetUserid()
	code := betting.Bet(userid, seat, number)
	glog.Debugf("code %d, userid %s", code, userid)
	glog.Debugf("seat %d, number %d", seat, number)
	switch code {
	case 0:
		expend(p, number, data.LogType36)
		p.Send(stoc)
	case 1:
		stoc.Error = pb.GameNotStart
		p.Send(stoc)
	}
}

func BettingRecord(ctos *pb.CBettingRecord,
	userid string) (stoc *pb.SBettingRecord) {
	stoc := new(pb.SBettingRecord)
	page := ctos.GetPage()
	list, err := data.GetBettingRecords(userid, int(page))
	glog.Debugf("getRecord page %d, userid %s", page, userid)
	glog.Debugf("getRecord err %v, len list %d", err, len(list))
	if err != nil {
		glog.Errorf("GetRecords err:%v", err)
		return
	}
	for _, v := range list {
		msg := &pb.RecordBettings{
			Index:  v.Id,
			Cards:  v.Cards,
			Niu:    v.Niu,
			Seats:  v.Seats,
			Number: v.Lose[userid],
		}
		for _, v2 := range v.Ante[userid] {
			msg2 := &pb.RecordBetting{
				Seat:   v2.Seat,
				Number: v2.Number,
			}
			msg.List = append(msg.List, msg2)
		}
		stoc.List = append(stoc.List, msg)
	}
	glog.Debugf("getRecord len List %d", len(stoc.List))
	return
}
