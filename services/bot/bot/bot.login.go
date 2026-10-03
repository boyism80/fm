package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/conn"
)

type looks struct {
	face, hair, top, bottom, shoes, weapon uint32
}

func (b *Bot) Enter() error {
	if err := b.login(); err != nil {
		b.Close()
		return err
	}
	if err := b.openGame(); err != nil {
		b.Close()
		return err
	}
	return b.conn.SetDeadline(time.Time{})
}

func (b *Bot) login() error {
	addr := fmt.Sprintf("%s:%d", b.cfg.Login.Host, b.cfg.Login.Port)
	c, err := b.dial(addr)
	if err != nil {
		return fmt.Errorf("login %s: %w", addr, err)
	}
	b.conn = c

	var auth *response.Authenticate
	for attempt := 0; attempt < 5; attempt++ {
		if err := b.conn.Send(&request.Login{ID: b.ID, Pw: b.cfg.Password, Mac: "00-00-00-00-00-01"}); err != nil {
			return err
		}
		pkt, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
			switch pkt.(type) {
			case *response.Authenticate, *response.LoginFailed:
				return true
			default:
				return false
			}
		})
		if err != nil {
			return err
		}
		if failed, ok := pkt.(*response.LoginFailed); ok {
			if failed.Reason == response.LoginFailedReasonAlreadyLoggedIn && attempt < 4 {
				time.Sleep(time.Second)
				continue
			}
			return fmt.Errorf("login failed reason %d", failed.Reason)
		}
		auth = pkt.(*response.Authenticate)
		break
	}
	if auth == nil {
		return fmt.Errorf("login failed")
	}
	b.Gender = auth.Gender

	if _, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
		_, ok := pkt.(*response.EndOfServerList)
		return ok
	}); err != nil {
		return err
	}

	if b.cfg.WorldID > 255 || b.cfg.Channel > 255 {
		return fmt.Errorf("world_id and channel must fit in one byte")
	}
	var list *response.CharacterList
	for attempt := 0; attempt < 20; attempt++ {
		if err := b.conn.Send(&request.CharacterList{Server: uint8(b.cfg.WorldID), Channel: uint8(b.cfg.Channel)}); err != nil {
			return err
		}
		pkt, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
			switch pkt.(type) {
			case *response.CharacterList, *response.LoginFailed:
				return true
			default:
				return false
			}
		})
		if err != nil {
			return err
		}
		if failed, ok := pkt.(*response.LoginFailed); ok {
			if failed.Reason == response.LoginFailedReasonTooManyConnections && attempt < 19 {
				time.Sleep(time.Second)
				continue
			}
			return fmt.Errorf("character list failed reason %d", failed.Reason)
		}
		list = pkt.(*response.CharacterList)
		break
	}
	if len(list.Characters) == 0 {
		return b.createCharacter()
	}
	b.CharID = list.Characters[0].ID
	b.Name = list.Characters[0].Name
	return nil
}

func (b *Bot) createCharacter() error {
	base := b.Name
	for i := 0; i < 3; i++ {
		name := base
		if i > 0 {
			name = fmt.Sprintf("%s%d", base, i)
		}
		if len(name) > 12 {
			return fmt.Errorf("character name %q is longer than 12 bytes", name)
		}
		if err := b.conn.Send(&request.CheckName{Name: name}); err != nil {
			return err
		}
		pkt, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
			_, ok := pkt.(*response.CheckName)
			return ok
		})
		if err != nil {
			return err
		}
		if pkt.(*response.CheckName).Exists {
			continue
		}
		b.Name = name
		look := maleLooks
		if b.Gender == 1 {
			look = femaleLooks
		}
		if err := b.conn.Send(&request.CreateCharacter{
			Name: b.Name, Face: look.face, Hair: look.hair,
			Top: look.top, Bottom: look.bottom, Shoes: look.shoes, Weapon: look.weapon,
		}); err != nil {
			return err
		}
		createdPkt, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
			_, ok := pkt.(*response.CreateCharacter)
			return ok
		})
		if err != nil {
			return err
		}
		created := createdPkt.(*response.CreateCharacter)
		if created.Success == false || created.Character == nil || created.Character.ID == 0 {
			return fmt.Errorf("create character %s failed", b.Name)
		}
		b.CharID = created.Character.ID
		return nil
	}
	return fmt.Errorf("character name %s is taken", base)
}

func (b *Bot) openGame() error {
	if err := b.conn.Send(&request.SelectCharacter{CharacterId: b.CharID}); err != nil {
		return err
	}
	pkt, err := b.readUntil(b.wait(), loginOpcodes, func(pkt any) bool {
		switch pkt.(type) {
		case *response.Transfer, *response.LoginFailed:
			return true
		default:
			return false
		}
	})
	if err != nil {
		return err
	}
	if failed, ok := pkt.(*response.LoginFailed); ok {
		return fmt.Errorf("select character failed reason %d", failed.Reason)
	}
	transfer := pkt.(*response.Transfer)
	b.Close()

	host := transfer.IP
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", host, transfer.Port)
	c, err := b.dial(addr)
	if err != nil {
		return fmt.Errorf("game %s: %w", addr, err)
	}
	b.conn = c
	b.Gen++
	if err := b.conn.Send(&request.LoginGame{PlayerId: transfer.CharacterId}); err != nil {
		return err
	}

	var sawLogin, sawKeys bool
	deadline := time.Now().Add(b.wait())
	for sawLogin == false || sawKeys == false {
		pkt, err := b.readUntil(time.Until(deadline), gameOpcodes, func(pkt any) bool {
			switch pkt.(type) {
			case *response.Login, *response.KeyMap, *response.LoginFailed:
				return true
			default:
				return false
			}
		})
		if err != nil {
			return fmt.Errorf("game login: %w", err)
		}
		switch p := pkt.(type) {
		case *response.LoginFailed:
			return fmt.Errorf("game login failed reason %d", p.Reason)
		case *response.Login:
			if p.Character == nil {
				return fmt.Errorf("game login has no character")
			}
			b.CharID = p.Character.ID
			b.Name = p.Character.Name
			b.Update(p)
			sawLogin = true
		case *response.KeyMap:
			sawKeys = true
		}
	}
	log.Printf("%s entered map %d", b.Name, b.Map)
	return nil
}

func (b *Bot) readUntil(timeout time.Duration, allow []uint16, match func(any) bool) (any, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timed out")
	}
	if err := b.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	for {
		opcode, body, err := b.conn.Read()
		if err != nil {
			return nil, err
		}
		keep := false
		for _, item := range allow {
			if item == opcode {
				keep = true
				break
			}
		}
		if keep == false {
			continue
		}
		pkt, err := conn.Decode(opcode, body)
		if err != nil {
			continue
		}
		if match(pkt) {
			return pkt, nil
		}
	}
}

func (b *Bot) dial(addr string) (*conn.Conn, error) {
	delay := 200 * time.Millisecond
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		c, err := conn.Dial(addr)
		if err == nil {
			return c, nil
		}
		last = err
		time.Sleep(delay)
		delay *= 2
	}
	return nil, fmt.Errorf("%v: %w", last, ErrConnect)
}

var (
	maleLooks    = looks{face: 20100, hair: 30030, top: 1040002, bottom: 1060002, shoes: 1072001, weapon: 1302000}
	femaleLooks  = looks{face: 21700, hair: 31002, top: 1041002, bottom: 1061002, shoes: 1072001, weapon: 1302000}
	loginOpcodes = []uint16{
		(&response.LoginFailed{}).Opcode(),
		(&response.ServerList{}).Opcode(),
		(&response.CharacterList{}).Opcode(),
		(&response.Transfer{}).Opcode(),
		(&response.CheckName{}).Opcode(),
		(&response.CreateCharacter{}).Opcode(),
	}
	gameOpcodes = []uint16{
		(&response.Login{}).Opcode(),
		(&response.KeyMap{}).Opcode(),
	}
)
