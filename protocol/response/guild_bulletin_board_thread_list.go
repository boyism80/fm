package response

import (
	"time"

	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildBulletinBoardThreadEntry struct {
	LocalThreadID     int32
	PosterCharacterID uint32
	Title             string
	Timestamp         time.Time
	Icon              int32
	ReplyCount        int32
}

type GuildBulletinBoardThreadList struct {
	Start      int
	TotalCount int
	Notice     *GuildBulletinBoardThreadEntry
	Threads    []GuildBulletinBoardThreadEntry
}

func (p *GuildBulletinBoardThreadList) Opcode() uint16 {
	return 0x48
}

func (p *GuildBulletinBoardThreadList) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildBulletinBoardS2CThreadList))

	if p.Notice == nil {
		w.WriteU8(0)
	} else {
		w.WriteU8(1)
		writeGuildBulletinBoardThreadSummary(w, *p.Notice)
	}

	threadCount := p.TotalCount
	start := p.Start
	if threadCount < start {
		start = 0
	}
	w.Write32(int32(threadCount))

	pages := constant.GuildBulletinBoardThreadsPerPage
	if threadCount-start < pages {
		pages = threadCount - start
	}
	if pages < 0 {
		pages = 0
	}
	if len(p.Threads) < pages {
		pages = len(p.Threads)
	}
	w.Write32(int32(pages))

	for i := 0; i < pages; i++ {
		writeGuildBulletinBoardThreadSummary(w, p.Threads[i])
	}
	return nil
}

func (p *GuildBulletinBoardThreadList) Deserialize(_ *stream.StreamReader) {}

func writeGuildBulletinBoardThreadSummary(w *stream.StreamWriter, entry GuildBulletinBoardThreadEntry) {
	w.Write32(entry.LocalThreadID)
	w.Write32(int32(entry.PosterCharacterID))
	w.WriteStr16(entry.Title)
	w.WriteDateTime(entry.Timestamp)
	w.Write32(entry.Icon)
	w.Write32(entry.ReplyCount)
}
