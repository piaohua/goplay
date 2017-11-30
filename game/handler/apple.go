package handler

/*
import (
	"goplay/pb"
	"niu/apple"
	"niu/data"
	"niu/images"
	"utils"

	"github.com/golang/glog"
)

func appleOrder(ctos *pb.CApplePay, p *data.User) {
	stoc := new(pb.SApplePay)
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
	glog.Infof("apple pay result %#v", result)
	for _, v := range result.Receipt.InApp {
		tradeRecord := new(data.TradeRecord)
		tradeRecord.Id = v.Transaction_id
		tradeRecord.Transid = v.Transaction_id
		tradeRecord.Amount = v.Quantity
		tradeRecord.Transtime = v.Purchase_date
		//if !appleVerify(v.Product_id, tradeRecord, p) {
		if !appleVerify(utils.String(id), tradeRecord, p) {
			glog.Errorf("apple pay verify err %#v", tradeRecord)
			stoc.Error = pb.AppleOrderFail
			return
		}
	}
}

func appleVerify(product_id string, tradeRecord *data.TradeRecord, p *data.User) bool {
	if tradeRecord.Has() {
		//重复发货
		glog.Errorf("apple pay already exist %#v", tradeRecord)
		return false
	}
	d := images.GetShop(product_id)
	if uint32(d.Payway) != data.RMB {
		glog.Errorf("apple pay %#v, %#v", d, tradeRecord)
		return false
	}
	tradeRecord.Currency = "RMB"
	tradeRecord.Itemid = utils.String(d.Propid)
	tradeRecord.Diamond = d.Number
	tradeRecord.Money = uint32(d.Price * 100) //转换为分
	tradeRecord.Result = data.TradeSuccess
	tradeRecord.Clientip = p.GetConn().GetIPAddr()
	tradeRecord.Agent = p.GetAgent()
	tradeRecord.Atype = p.GetAtype()
	tradeRecord.Userid = p.GetUserid()
	//tradeRecord.Paytype = "ios"
	if p.GetMoney() == 0 {
		tradeRecord.First = 1
	}
	tradeRecord.Ctime = utils.BsonNow()
	if !tradeRecord.Save() { //
		glog.Errorf("apple pay save err %#v", tradeRecord)
		return false
	}
	//发货
	images.DeliverGoods(p, tradeRecord.Diamond, tradeRecord.Money, tradeRecord.First)
	return true
}
*/
