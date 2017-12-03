package handler

import (
	"goplay/data"
	"goplay/pb"
	"utils"
)

//test
func NewMail(userid, id string) {
	from := "系统消息"
	content := `北京时间10月19日凌晨，时隔20余月，
	AlphaGo再次登上科学杂志《nature》。
	谷歌人工智能团队DeepMind 发布了他们的最新论文
	Mastering the game of Go without human knowledge，
	向人们介绍了阿尔法狗家族的新成员——阿尔法元。
	与之前几个版本的阿尔法狗不同，阿尔法元除了解围棋规则外，
	完全不依靠棋谱和人类数据，从零开始“自学成才”，
	成为全世界最厉害的（人工智能）围棋手。`
	m := newMail2(id, from, userid, content)
	m.Save()
}

//测试
func newMail2(id, From, To, Content string) *data.Mail {
	m := new(data.Mail)
	m.Id = id
	m.From = From
	m.To = To
	m.Content = Content
	m.Etime = utils.Stamp2Time(utils.TimestampTomorrow())
	item := data.Items{
		Rtype:  int(data.DIAMOND),
		Number: 100,
	}
	m.Attachment = []data.Items{item}
	return m
}

// 新邮件列表
func GetMailList(ctos *pb.CMailList, userid string) (stoc *pb.SMailList) {
	stoc = new(pb.SMailList)
	if userid == "" {
		return
	}
	maxid := ctos.GetMaxid()
	list := data.GetNewMailList(maxid, userid)
	for _, v := range list {
		l := &pb.MailList{
			Id:      v.Id,
			From:    v.From,
			Content: v.Content,
			Status:  uint32(v.Status),
			Etime:   uint32(utils.Time2Stamp(v.Etime.Local())),
			Ctime:   uint32(utils.Time2Stamp(v.Ctime.Local())),
		}
		for _, val := range v.Attachment {
			i := &pb.Items{
				Rtype:  uint32(val.Rtype),
				Number: val.Number,
			}
			l.Attachment = append(l.Attachment, i)
		}
		stoc.List = append(stoc.List, l)
	}
	return
}

// 删除邮件
func DelMail(ctos *pb.CDeleteMail, userid string) (stoc *pb.SDeleteMail) {
	stoc = new(pb.SDeleteMail)
	if userid == "" {
		return
	}
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.To = userid //只能删除自己的
	d.Del()
	stoc.Id = ctos.GetId()
	return
}

// 收取附件
func GetMailItem(ctos *pb.CGetMailItem,
	userid string) (stoc *pb.SGetMailItem, list []data.Items) {
	stoc = new(pb.SGetMailItem)
	if userid == "" {
		return
	}
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.Get()
	//没有附件
	if len(d.Attachment) == 0 {
		stoc.Error = pb.MailNotAttachment
		return
	}
	//已经领取过
	if d.Status == 2 {
		stoc.Error = pb.MailAlreadyGet
		return
	}
	//已经过期
	if d.Etime.Before(utils.LocalTime()) { //是否当前时间之前
		stoc.Error = pb.MailAlreadyExpire
		return
	}
	if d.To != userid { //只能是自己的
		stoc.Error = pb.MailNotExist
		return
	}
	d.Status = 2
	d.Update()
	stoc.Id = ctos.GetId()
	list = d.Attachment
	return
}
