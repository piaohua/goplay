package handler

import (
	"api/wxpay"
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
	"utils"
)

//微信支付查询
func WxQuery(ctos *pb.CWxpayQuery) (stoc *pb.SWxpayQuery) {
	stoc = new(pb.SWxpayQuery)
	var transid string = ctos.GetTransid()
	if transid == "" {
		stoc.Error = pb.PayOrderError
		return
	}
	queryResult, err := config.Apppay.Query(transid)
	//glog.Infof("queryResult  %#v, err %v, transid %s", queryResult, err, transid)
	if err != nil {
		stoc.Error = pb.PayOrderError
		return
	}
	if queryResult.ReturnCode == "SUCCESS" &&
		queryResult.ResultCode == "SUCCESS" &&
		queryResult.TradeState == "SUCCESS" {
		//glog.Infof("queryResult  %#v, err %v, transid %s", queryResult, err, transid)
		stoc.Result = 0
		stoc.Orderid = queryResult.OrderId
	} else {
		stoc.Error = pb.PayOrderError
	}
	return
}

//微信支付下单
func WxOrder(ctos *pb.CWxpayOrder, p *data.User,
	ip string) (stoc *pb.SWxpayOrder, t *data.TradeRecord) {
	stoc = new(pb.SWxpayOrder)
	var waresid uint32 = ctos.GetId()
	var body string = ctos.GetBody()
	var userid string = p.GetUserid()
	var agent string = p.GetAgent()
	glog.Info("waresid ", waresid)
	var transid, orderid string
	transid, orderid, t = wxpayOrder(waresid, userid, agent, ip, body)
	if transid == "" || orderid == "" {
		glog.Info("wx order fail:", waresid, userid)
		stoc.Error = pb.PayOrderFail
		return
	}
	payRequest := config.Apppay.NewPaymentRequest(transid)
	payReqJson, err := wxpay.ToJson(&payRequest)
	//retMap, err := wxpay.ToMap(&payRequest)
	if err != nil {
		glog.Error("wx order err:", waresid, userid, err)
		stoc.Error = pb.PayOrderFail
		return
	}
	//payReqStr := wxpay.ToXmlString(retMap)
	//stoc.Payreq = payReqStr
	stoc.Orderid = orderid
	stoc.Payreq = string(payReqJson)
	//glog.Info("orderid ", orderid, " transid ", transid)
	//go wxpayQuery(orderid) //查询
	stoc.Id = waresid
	return
}

// 下单
func wxpayOrder(waresid uint32, userid, agent, ip,
	body string) (transid, orderid string, t *data.TradeRecord) {
	d := config.GetShop(utils.String(waresid))
	if uint32(d.Payway) != data.RMB {
		return
	}
	var diamond uint32 = d.Number
	var price uint32 = uint32(d.Price * 100) //转换为分
	var itemid string = utils.String(d.Propid)
	//var orderid string = data.GenCporderid(userid)
	orderid = data.GenOrderid()
	transid, err := config.Apppay.Submit(orderid, float64(price), body, ip)
	glog.Debugf("orderid %s, transid %s, err %v", orderid, transid, err)
	if err != nil {
		glog.Error("wx order err:", waresid, userid, err)
		return
	}
	//transid,下单记录
	t = &data.TradeRecord{
		Id:       orderid,
		Transid:  transid,
		Userid:   userid,
		Itemid:   itemid,
		Amount:   "1",
		Diamond:  diamond,
		Money:    price,
		Ctime:    utils.BsonNow(),
		DayStamp: utils.TimestampTodayTime(),
		Result:   data.Tradeing, //下单状态
		Clientip: ip,
		Agent:    agent,
	}
	if !t.Save() { //TODO:优化,未支付订单
		glog.Error("trade record save err:", orderid, userid)
		return
	}
	return
}

/*
//主动查询发货
func wxpayQuery(orderid string) {
	utils.Sleep(120)
	queryResult, err := config.Apppay.Query(orderid)
	if err != nil {
		glog.Errorf("wxpayQuery  %#v, err %v, orderid %s", queryResult, err, orderid)
		return
	}
	//主动查询发货
	config.WxpayQuery(queryResult)
}
*/
