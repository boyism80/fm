package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildBulletinBoardOperation struct {
	Action        constant.GuildBulletinBoardC2SAction
	Edit          bool
	Notice        bool
	LocalThreadID int32
	Page          int32
	ReplyID       int32
	Title         string
	Text          string
	Icon          int32
}

func (*GuildBulletinBoardOperation) Opcode() byte {
	return 0xC9
}

func (p *GuildBulletinBoardOperation) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *GuildBulletinBoardOperation) Deserialize(reader *stream.StreamReader) {
	p.Action = constant.GuildBulletinBoardC2SAction(reader.ReadU8())
	switch p.Action {
	case constant.GuildBulletinBoardC2SWriteThread:
		p.Edit = reader.ReadU8() > 0
		if p.Edit {
			p.LocalThreadID = reader.Read32()
		}
		p.Notice = reader.ReadU8() > 0
		p.Title = reader.ReadStr16()
		p.Text = reader.ReadStr16()
		p.Icon = reader.Read32()
	case constant.GuildBulletinBoardC2SDeleteThread, constant.GuildBulletinBoardC2SShowThread:
		p.LocalThreadID = reader.Read32()
	case constant.GuildBulletinBoardC2SListThreads:
		p.Page = reader.Read32()
	case constant.GuildBulletinBoardC2SWriteReply:
		p.LocalThreadID = reader.Read32()
		p.Text = reader.ReadStr16()
	case constant.GuildBulletinBoardC2SDeleteReply:
		p.LocalThreadID = reader.Read32()
		p.ReplyID = reader.Read32()
	}
}
