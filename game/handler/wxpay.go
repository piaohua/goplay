package handler

import (
	"encoding/xml"
	"math"
	"strconv"

	"api/wxpay"
	"goplay/data"
	"goplay/game/config"
	"goplay/glog"
	"goplay/pb"
	"utils"

	jsoniter "github.com/json-iterator/go"
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
	ip string) (stoc *pb.SWxpayOrder, trade string) {
	stoc = new(pb.SWxpayOrder)
	var waresid uint32 = ctos.GetId()
	var body string = ctos.GetBody()
	var userid string = p.GetUserid()
	var agent string = p.GetAgent()
	glog.Info("waresid ", waresid)
	t := wxpayOrder(waresid, userid, agent, ip, body)
	if t == nil {
		glog.Info("wx order fail:", waresid, userid)
		stoc.Error = pb.PayOrderFail
		return
	}
	payRequest := config.Apppay.NewPaymentRequest(t.Transid)
	payReqJson, err := wxpay.ToJson(&payRequest)
	//retMap, err := wxpay.ToMap(&payRequest)
	if err != nil {
		glog.Error("wx order err:", waresid, userid, err)
		stoc.Error = pb.PayOrderFail
		return
	}
	//payReqStr := wxpay.ToXmlString(retMap)
	//stoc.Payreq = payReqStr
	trade1, err := jsoniter.Marshal(t)
	if err != nil {
		glog.Errorf("tradeRecord Marshal err %v", err)
		stoc.Error = pb.PayOrderFail
		return
	}
	trade = string(trade1)
	//响应
	stoc.Orderid = t.Id
	stoc.Payreq = string(payReqJson)
	//glog.Info("orderid ", t.Id, " transid ", t.Transid)
	stoc.Id = waresid
	return
}

// 下单
func wxpayOrder(waresid uint32, userid, agent,
	ip, body string) *data.TradeRecord {
	d := config.GetShop(utils.String(waresid))
	if uint32(d.Payway) != data.RMB {
		return nil
	}
	var diamond uint32 = d.Number
	var price uint32 = uint32(d.Price * 100) //转换为分
	var itemid string = utils.String(d.Propid)
	//var orderid string = data.GenCporderid(userid)
	var orderid string = data.GenOrderid()
	transid, err := config.Apppay.Submit(orderid, float64(price), body, ip)
	glog.Debugf("orderid %s, transid %s, err %v", orderid, transid, err)
	if err != nil {
		glog.Error("wx order err:", waresid, userid, err)
		return nil
	}
	//transid,下单记录
	return &data.TradeRecord{
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
}

//主动查询发货
func ActWxpayQuery(orderid string) string {
	utils.Sleep(120) //2分钟后主动查询到账情况
	queryResult, err := config.Apppay.Query(orderid)
	if err != nil {
		glog.Errorf("queryResult  %#v, err %v, orderid %s", queryResult, err, orderid)
		return ""
	}
	return actWxpayQuery2(queryResult)
}

//主动查询发货
func actWxpayQuery2(q wxpay.QueryOrderResult) string {
	if q.ReturnCode != "SUCCESS" ||
		q.ResultCode != "SUCCESS" ||
		q.TradeState != "SUCCESS" {
		glog.Errorf("actWxpayQuery2 failed %#v", q)
		return ""
	}
	t := new(wxpay.TradeResult)
	t.ReturnCode = q.ReturnCode
	t.ReturnMsg = q.ReturnMsg
	t.AppId = q.AppId
	t.MchId = q.MchId
	t.DeviceInfo = q.DeviceInfo
	t.NonceStr = q.NonceStr
	t.Sign = q.Sign
	t.ResultCode = q.ResultCode
	t.ErrCode = q.ErrCode
	t.ErrCodeDesc = q.ErrCodeDesc
	t.OpenId = q.OpenId
	t.IsSubscribe = q.IsSubscribe
	t.TradeType = q.TradeType
	t.BankType = q.BankType
	t.TotalFee = q.TotalFee
	t.FeeType = q.FeeType
	t.CashFee = q.CashFee
	t.CashFeeType = q.CashFeeType
	t.CouponFee = q.CouponFee
	t.CouponCount = q.CouponCount
	t.TransactionId = q.TransactionId
	t.OrderId = q.OrderId
	t.Attach = q.Attach
	t.TimeEnd = q.TimeEnd
	b, err := xml.Marshal(t)
	if err != nil {
		glog.Errorf("actWxpayQuery2 err %v", err)
		return ""
	}
	return string(b)
}

//回调验证,gate调用
func WxpayVerify(arg *pb.WxpayCallback) bool {
	result, err := wxpay.ParseTradeResult([]byte(arg.Result))
	if err != nil {
		glog.Errorf("WxpayVerify err %v, arg %#v", err, arg)
		return false
	}
	err = config.Apppay.RecvVerify(&result)
	if err != nil {
		glog.Errorf("recv verify %v, err:, %v", result, err)
		return false
	}
	return true
}

//回调验证或主动查询发货,dbms调用
func WxpayCallback(arg *pb.WxpayCallback) *wxpay.TradeResult {
	result, err := wxpay.ParseTradeResult([]byte(arg.Result))
	if err != nil {
		glog.Errorf("WxpayCallback err %v, arg %#v", err, arg)
		return nil
	}
	return &result
}

//发货验证
func WxpayTradeVerify(t *wxpay.TradeResult) *data.TradeRecord {
	//sign
	tradeRecord := &data.TradeRecord{
		Id: t.OrderId,
		//Transid: t.TransactionId,
	}
	//订单获取
	tradeRecord.Get()
	//glog.Infof("tradeRecord  %#v", tradeRecord)
	//glog.Infof("TradeResult  %#v", t)
	if tradeRecord.Transid == "" {
		//订单不存在或其它
		glog.Errorf("not exist orderid %v", t)
		return nil
	}
	if tradeRecord.Result == 0 {
		//重复发货
		glog.Errorf("repeat resp %v", t)
		return nil
	}
	//更新记录
	tradeRecord.Transtime = t.TimeEnd
	tradeRecord.Currency = t.FeeType
	tradeRecord.Paytype = 403 //t.TradeType == "APP"
	money, err := strconv.Atoi(t.TotalFee)
	if err != nil {
		glog.Errorf("wxpay: %v, err: %v", t, err)
	}
	tradeRecord.Money = uint32(money)      //转换为分
	tradeRecord.Result = data.TradeSuccess //交易成功
	//glog.Infof("tradeRecord  %#v", tradeRecord)
	return tradeRecord
}

//发货
func WxpaySendGoods(online bool, trade *data.TradeRecord, user *data.User) {
	if user == nil || user.Userid == "" {
		glog.Errorf("user empty %v, trade %#v", online, trade)
		trade.Result = data.TradeGoods //发货失败
	} else {
		if user.GetMoney() == 0 {
			trade.First = 1
		}
		//交易成功
		trade.Agent = user.GetAgent()
		trade.Atype = user.GetAtype()
	}
	if !online {
		//离线状态
		//TODO 日志记录
		var diamond, coin int32
		//首充
		if trade.First == 1 {
			diamond += config.GetEnv(data.ENV4)
			coin += config.GetEnv(data.ENV5)
		}
		//充值数量
		diamond += int32(trade.Diamond)
		//vip赠送
		diamond += getVipGive(user.GetVipLevel(), diamond)
		//货币变更
		user.AddDiamond(diamond)
		user.AddCoin(coin)
		//vip变更
		lev2 := config.GetVipLevel(user.GetVip() + trade.Money)
		user.SetVip(lev2, trade.Money)
		user.AddMoney(trade.Money)
		//存储
		user.Save()
	}
	//update record
	if !trade.Upsert() {
		glog.Errorf("trade save failed: %#v", trade)
	}
}

//vip赠送
func getVipGive(level int, num int32) int32 {
	if level <= 0 {
		return 0
	}
	pay := config.GetVipPay(level)
	return int32(math.Ceil(float64(num) * (float64(pay) / 100)))
}
