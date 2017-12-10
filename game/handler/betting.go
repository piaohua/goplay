package handler

import (
	"goplay/data"
	"goplay/glog"
	"goplay/pb"
)

func BettingRecord(page uint32, userid string) (stoc *pb.SBettingRecord) {
	stoc = new(pb.SBettingRecord)
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
