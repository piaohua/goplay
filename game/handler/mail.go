package handler

import (
	"goplay/data"
	"goplay/pb"
	"niu/mail"
	"utils"
)

//test
func newMail(p *data.User) {
	from := "系统消息"
	content := `北京时间10月19日凌晨，时隔20余月，
	AlphaGo再次登上科学杂志《nature》。
	谷歌人工智能团队DeepMind 发布了他们的最新论文
	Mastering the game of Go without human knowledge，
	向人们介绍了阿尔法狗家族的新成员——阿尔法元。
	与之前几个版本的阿尔法狗不同，阿尔法元除了解围棋规则外，
	完全不依靠棋谱和人类数据，从零开始“自学成才”，
	成为全世界最厉害的（人工智能）围棋手。`
	m := mail.New2(from, p.GetUserid(), content)
	m.Save()
	//msg := new(pb.SMailNotice)
	//p.Send(msg)
}

// 新邮件列表
func getMailList(ctos *pb.CMailList, p *data.User) (stoc *pb.SMailList) {
	stoc = new(pb.SMailList)
	maxid := ctos.GetMaxid()
	list := data.GetNewMailList(maxid, p.GetUserid())
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
func delMail(ctos *pb.CDeleteMail, p *data.User) (stoc *pb.SDeleteMail) {
	stoc = new(pb.SDeleteMail)
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.Del()
	stoc.Id = ctos.GetId()
	return
}

// 收取附件
func getMailItem(ctos *pb.CGetMailItem, p *data.User) (stoc *pb.SGetMailItem) {
	stoc = new(pb.SGetMailItem)
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.Get()
	//没有附件
	if len(d.Attachment) == 0 {
		stoc.Error = pb.NotInRoom
		return
	}
	//已经领取过
	if d.Status == 2 {
		stoc.Error = pb.NotInRoom
		return
	}
	//已经过期
	if d.Etime.Before(utils.LocalTime()) { //是否当前时间之前
		stoc.Error = pb.NotInRoom
		return
	}
	for _, v := range d.Attachment {
		addPrize(p, int(v.Rtype), data.LogType22, int32(v.Number))
	}
	d.Status = 2
	d.Update()
	stoc.Id = ctos.GetId()
	return
}
