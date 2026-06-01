package response

import (
	"sort"
	"time"

	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildBulletinBoardReplyEntry struct {
	ReplyID           int32
	PosterCharacterID uint32
	Timestamp         time.Time
	Content           string
}

type GuildBulletinBoardShowThread struct {
	LocalThreadID     int32
	PosterCharacterID uint32
	Timestamp         time.Time
	Title             string
	Body              string
	Icon              int32
	Replies           []GuildBulletinBoardReplyEntry
}

func (p *GuildBulletinBoardShowThread) Opcode() uint16 {
	return 0x48
}

func (p *GuildBulletinBoardShowThread) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildBulletinBoardS2CShowThread))
	w.Write32(p.LocalThreadID)
	w.Write32(int32(p.PosterCharacterID))
	w.WriteDateTime(p.Timestamp)
	w.WriteStr16(p.Title)
	w.WriteStr16(p.Body)
	w.Write32(p.Icon)

	replies := append([]GuildBulletinBoardReplyEntry(nil), p.Replies...)
	sort.Slice(replies, func(i, j int) bool {
		return replies[i].ReplyID < replies[j].ReplyID
	})

	w.Write32(int32(len(replies)))
	for _, reply := range replies {
		w.Write32(reply.ReplyID)
		w.Write32(int32(reply.PosterCharacterID))
		w.WriteDateTime(reply.Timestamp)
		w.WriteStr16(reply.Content)
	}
	return nil
}

func (p *GuildBulletinBoardShowThread) Deserialize(_ *stream.StreamReader) {}
