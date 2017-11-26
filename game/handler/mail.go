package handle

import (
	"niu/data"
	"niu/errorcode"
	"niu/inter"
	"niu/mail"
	"niu/protocol"
	"niu/socket"
	"utils"

	"github.com/golang/protobuf/proto"
)

func init() {
	p1 := protocol.CMailList{}
	socket.Regist(p1.GetCode(), p1, getMailList)
	p2 := protocol.CDeleteMail{}
	socket.Regist(p2.GetCode(), p2, delMail)
	p3 := protocol.CGetMailItem{}
	socket.Regist(p3.GetCode(), p3, getMailItem)
}

//test
func newMail(p inter.IPlayer) {
	from := "系统消息"
	content := "北京时间10月19日凌晨，时隔20余月，AlphaGo再次登上科学杂志《nature》。谷歌人工智能团队DeepMind 发布了他们的最新论文Mastering the game of Go without human knowledge，向人们介绍了阿尔法狗家族的新成员——阿尔法元。与之前几个版本的阿尔法狗不同，阿尔法元除了解围棋规则外，完全不依靠棋谱和人类数据，从零开始“自学成才”，成为全世界最厉害的（人工智能）围棋手。"
	m := mail.New2(from, p.GetUserid(), content)
	m.Save()
	msg := new(protocol.SMailNotice)
	p.Send(msg)
}

// 新邮件列表
func getMailList(ctos *protocol.CMailList, p inter.IPlayer) {
	stoc := new(protocol.SMailList)
	maxid := ctos.GetMaxid()
	list := data.GetNewMailList(maxid, p.GetUserid())
	for _, v := range list {
		l := &protocol.MailList{
			Id:      proto.String(v.Id),
			From:    proto.String(v.From),
			Content: proto.String(v.Content),
			Status:  proto.Uint32(uint32(v.Status)),
			Etime:   proto.Uint32(uint32(utils.Time2Stamp(v.Etime.Local()))),
			Ctime:   proto.Uint32(uint32(utils.Time2Stamp(v.Ctime.Local()))),
		}
		for _, val := range v.Attachment {
			i := &protocol.Items{
				Rtype:  proto.Uint32(uint32(val.Rtype)),
				Number: proto.Uint32(val.Number),
			}
			l.Attachment = append(l.Attachment, i)
		}
		stoc.List = append(stoc.List, l)
	}
	p.Send(stoc)
}

// 删除邮件
func delMail(ctos *protocol.CDeleteMail, p inter.IPlayer) {
	stoc := new(protocol.SDeleteMail)
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.Del()
	stoc.Id = proto.String(ctos.GetId())
	p.Send(stoc)
}

// 收取附件
func getMailItem(ctos *protocol.CGetMailItem, p inter.IPlayer) {
	stoc := new(protocol.SGetMailItem)
	d := new(data.Mail)
	d.Id = ctos.GetId()
	d.Get()
	//没有附件
	if len(d.Attachment) == 0 {
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		p.Send(stoc)
		return
	}
	//已经领取过
	if d.Status == 2 {
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		p.Send(stoc)
		return
	}
	//已经过期
	if d.Etime.Before(utils.LocalTime()) { //是否当前时间之前
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		p.Send(stoc)
		return
	}
	for _, v := range d.Attachment {
		addPrize(p, int(v.Rtype), data.LogType22, int32(v.Number))
	}
	d.Status = 2
	d.Update()
	stoc.Id = proto.String(ctos.GetId())
	p.Send(stoc)
}
