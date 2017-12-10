package handler

import (
	"encoding/json"

	"api/apple"
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
	"utils"
)

func AppleOrder(ctos *pb.CApplePay, p *data.User) (stoc *pb.SApplePay,
	tradeRecord *data.TradeRecord, trade string) {
	stoc = new(pb.SApplePay)
	id := ctos.GetId()
	receipt := ctos.GetReceipt()
	stoc.Id = id
	result, err := apple.PayVeriy(receipt)
	if err != nil || result == nil {
		glog.Errorf("apple pay err %#v", err)
		stoc.Error = pb.AppleOrderFail
		return
	}
	if result.Status != 0 {
		glog.Errorf("apple pay err result %#v", result)
		stoc.Error = pb.AppleOrderFail
		return
	}
	glog.Debugf("apple pay result %#v", result)
	for _, v := range result.Receipt.InApp {
		tradeRecord = new(data.TradeRecord)
		tradeRecord.Id = v.Transaction_id
		tradeRecord.Transid = v.Transaction_id
		tradeRecord.Amount = v.Quantity
		tradeRecord.Transtime = v.Purchase_date
		if !tradeVerify(utils.String(id), tradeRecord, p) {
			glog.Errorf("apple pay verify err %#v", tradeRecord)
			stoc.Error = pb.AppleOrderFail
			return
		}
	}
	trade1, err := json.Marshal(tradeRecord)
	if err != nil {
		glog.Errorf("tradeRecord Marshal err %v", err)
		stoc.Error = pb.AppleOrderFail
		return
	}
	trade = string(trade1)
	return
}

func tradeVerify(product_id string, tradeRecord *data.TradeRecord,
	p *data.User) bool {
	d := config.GetShop(product_id)
	if uint32(d.Payway) != data.RMB {
		glog.Errorf("apple pay %#v, %#v", d, tradeRecord)
		return false
	}
	tradeRecord.Currency = "RMB"
	tradeRecord.Itemid = utils.String(d.Propid)
	tradeRecord.Diamond = d.Number
	tradeRecord.Money = uint32(d.Price * 100) //转换为分
	tradeRecord.Result = data.TradeSuccess
	tradeRecord.Clientip = p.LoginIp
	tradeRecord.Agent = p.GetAgent()
	tradeRecord.Atype = p.GetAtype()
	tradeRecord.Userid = p.GetUserid()
	//tradeRecord.Paytype = "ios"
	if p.GetMoney() == 0 {
		tradeRecord.First = 1
	}
	tradeRecord.Ctime = utils.BsonNow()
	return true
}

//订单是否存在和订单数据存储
func AppleVerify(arg *pb.ApplePay) (stoc *pb.ApplePaid) {
	stoc = new(pb.ApplePaid)
	tradeRecord := new(data.TradeRecord)
	err := json.Unmarshal([]byte(arg.Trade), tradeRecord)
	if err != nil {
		glog.Errorf("tradeRecord Marshal err %v", err)
		stoc.Result = false
		return
	}
	if tradeRecord.Has() {
		//重复发货
		glog.Errorf("apple pay already exist %#v", tradeRecord)
		stoc.Result = false
		return
	}
	if !tradeRecord.Save() { //
		glog.Errorf("apple pay save err %#v", tradeRecord)
		stoc.Result = false
		return
	}
	stoc.Result = true
	return
}
