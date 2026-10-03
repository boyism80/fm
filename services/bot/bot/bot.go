package bot

import (
	"errors"
	"fmt"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/stream"
)

var ErrConnect = errors.New("connect")

type outbound interface {
	Opcode() byte
	Serialize(writer *stream.StreamWriter) error
}

type Bot struct {
	Index  int
	ID     string
	Name   string
	Gender uint8
	CharID uint32
	Gen    int
	Map    uint32
	cfg    *config.Bot
	conn   *conn.Conn
}

func New(cfg *config.Bot, runID string, n int) (*Bot, error) {
	name := fmt.Sprintf("b%s%d", runID, n)
	if len(name) > 12 {
		return nil, fmt.Errorf("character name %q is longer than 12 bytes", name)
	}
	return &Bot{
		Index: n,
		ID:    fmt.Sprintf("bot_%s_%d", runID, n),
		Name:  name,
		cfg:   cfg,
	}, nil
}

func (b *Bot) Send(pkt outbound) error {
	if b.conn == nil {
		return fmt.Errorf("bot %s is not connected", b.Name)
	}
	return b.conn.Send(pkt)
}

func (b *Bot) Close() {
	if b.conn == nil {
		return
	}
	_ = b.conn.Close()
	b.conn = nil
}

func (b *Bot) Listen(onPacket func(gen int, opcode uint16, body []byte), onClose func(gen int, err error)) {
	c, gen := b.conn, b.Gen
	go func() {
		for {
			opcode, body, err := c.Read()
			if err != nil {
				onClose(gen, err)
				return
			}
			onPacket(gen, opcode, body)
		}
	}()
}

func (b *Bot) Update(pkt any) {
	switch p := pkt.(type) {
	case *response.Login:
		if p.Character != nil {
			b.Map = p.Character.Map
		}
	case *response.Warp:
		if p.Character != nil {
			b.Map = p.Character.Map
		}
	}
}

func (b *Bot) wait() time.Duration {
	return time.Duration(b.cfg.TimeoutMs) * time.Millisecond
}
