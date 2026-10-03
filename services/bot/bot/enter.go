package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/common/config"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/services/game/wz"
)

var ErrConnect = errors.New("connect")

type looks struct {
	face, hair, top, bottom, shoes, weapon uint32
}

type session struct {
	cfg    *config.Bot
	conn   *conn.Conn
	seat   int
	id     string
	name   string
	gender uint8
	charID uint32
	town   uint32
	wz     *wz.Resources
}

func Enter(cfg *config.Bot) error {
	runID := strconv.FormatInt(time.Now().Unix(), 36)
	name := "b" + runID + "0"
	if len(name) > 12 {
		return fmt.Errorf("character name %q is longer than 12 bytes", name)
	}
	s := &session{
		cfg:  cfg,
		seat: 1,
		id:   "bot_" + runID + "_0",
		name: name,
	}
	defer func() {
		if s.conn != nil {
			_ = s.conn.Close()
		}
	}()
	log.Printf("loading wz %s", cfg.WzPath)
	s.wz = wz.NewResources(cfg.WzPath)
	if s.wz != nil {
		if id, ok := s.wz.NameToMap("헤네시스"); ok {
			s.town = id
		}
	}
	if s.town == 0 {
		s.town = 100000000
	}

	if err := s.login(); err != nil {
		return err
	}
	if err := s.openGame(); err != nil {
		return err
	}
	return s.moveMaps()
}

func (s *session) login() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Login.Host, s.cfg.Login.Port)
	c, err := s.dial(addr)
	if err != nil {
		return fmt.Errorf("login %s: %w", addr, err)
	}
	s.conn = c

	var auth *response.Authenticate
	for attempt := 0; attempt < 5; attempt++ {
		if err := s.conn.Send(&request.Login{ID: s.id, Pw: s.cfg.Password, Mac: "00-00-00-00-00-01"}); err != nil {
			return err
		}
		pkt, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
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
	s.gender = auth.Gender
	log.Printf("authenticated account %s id=%d", s.id, auth.AccountId)

	if _, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
		_, ok := pkt.(*response.EndOfServerList)
		return ok
	}); err != nil {
		return err
	}

	if s.cfg.WorldID > 255 || s.cfg.Channel > 255 {
		return fmt.Errorf("world_id and channel must fit in one byte")
	}
	if err := s.conn.Send(&request.CharacterList{Server: uint8(s.cfg.WorldID), Channel: uint8(s.cfg.Channel)}); err != nil {
		return err
	}
	pkt, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
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
		return fmt.Errorf("character list failed reason %d", failed.Reason)
	}
	list := pkt.(*response.CharacterList)
	if len(list.Characters) == 0 {
		return s.createCharacter()
	}
	s.charID = list.Characters[0].ID
	s.name = list.Characters[0].Name
	log.Printf("using character %s id=%d", s.name, s.charID)
	return nil
}

func (s *session) createCharacter() error {
	for i := 0; i < 3; i++ {
		name := s.name
		if i > 0 {
			name = fmt.Sprintf("%s%d", s.name, i)
		}
		if len(name) > 12 {
			return fmt.Errorf("character name %q is longer than 12 bytes", name)
		}
		if err := s.conn.Send(&request.CheckName{Name: name}); err != nil {
			return err
		}
		pkt, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
			_, ok := pkt.(*response.CheckName)
			return ok
		})
		if err != nil {
			return err
		}
		checked := pkt.(*response.CheckName)
		if checked.Exists {
			continue
		}
		s.name = name
		look := maleLooks
		if s.gender == 1 {
			look = femaleLooks
		}
		if err := s.conn.Send(&request.CreateCharacter{
			Name: s.name, Face: look.face, Hair: look.hair,
			Top: look.top, Bottom: look.bottom, Shoes: look.shoes, Weapon: look.weapon,
		}); err != nil {
			return err
		}
		createdPkt, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
			_, ok := pkt.(*response.CreateCharacter)
			return ok
		})
		if err != nil {
			return err
		}
		created := createdPkt.(*response.CreateCharacter)
		if created.Success == false || created.Character == nil || created.Character.ID == 0 {
			return fmt.Errorf("create character %s failed", s.name)
		}
		s.charID = created.Character.ID
		log.Printf("created character %s id=%d", s.name, s.charID)
		return nil
	}
	return fmt.Errorf("character name %s is taken", s.name)
}

func (s *session) openGame() error {
	if err := s.conn.Send(&request.SelectCharacter{CharacterId: s.charID}); err != nil {
		return err
	}
	pkt, err := s.readUntil(s.wait(), loginOpcodes, func(pkt any) bool {
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
	_ = s.conn.Close()
	s.conn = nil

	host := transfer.IP
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", host, transfer.Port)
	c, err := s.dial(addr)
	if err != nil {
		return fmt.Errorf("game %s: %w", addr, err)
	}
	s.conn = c
	if err := s.conn.Send(&request.LoginGame{PlayerId: transfer.CharacterId}); err != nil {
		return err
	}

	var sawLogin, sawKeys bool
	deadline := time.Now().Add(s.wait())
	for sawLogin == false || sawKeys == false {
		pkt, err := s.readUntil(time.Until(deadline), gameOpcodes, func(pkt any) bool {
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
			s.charID = p.Character.ID
			s.name = p.Character.Name
			sawLogin = true
			log.Printf("in game map=%d spawn=%d", p.Character.Map, p.Character.SpawnPoint)
		case *response.KeyMap:
			sawKeys = true
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("game login timed out")
		}
	}
	return nil
}

func (s *session) moveMaps() error {
	if err := s.command(fmt.Sprintf("/맵이동 %d", s.town), s.town); err != nil {
		return err
	}
	if err := s.command(fmt.Sprintf("/인스턴스이동 %d %d", s.town, s.seat), s.town); err != nil {
		return err
	}
	log.Printf("moved %s to map %d slot %d", s.name, s.town, s.seat)
	return nil
}

func (s *session) command(text string, mapID uint32) error {
	if err := s.conn.Send(&request.NormalChat{Message: text}); err != nil {
		return err
	}
	pkt, err := s.readUntil(s.wait(), gameOpcodes, func(pkt any) bool {
		warp, ok := pkt.(*response.Warp)
		return ok && warp.Character != nil && warp.Character.Map == mapID
	})
	if err != nil {
		return fmt.Errorf("%s: %w", text, err)
	}
	log.Printf("%s -> map %d", text, pkt.(*response.Warp).Character.Map)
	return nil
}

func (s *session) wait() time.Duration {
	return time.Duration(s.cfg.TimeoutMs) * time.Millisecond
}

func (s *session) readUntil(timeout time.Duration, allow []uint16, match func(any) bool) (any, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timed out")
	}
	if err := s.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	for {
		opcode, body, err := s.conn.Read()
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
		if notice, ok := pkt.(*response.Notice); ok && strings.Contains(notice.Message, "권한이 부족합니다") {
			return nil, fmt.Errorf("권한이 부족합니다")
		}
		if match(pkt) {
			return pkt, nil
		}
	}
}

func (s *session) dial(addr string) (*conn.Conn, error) {
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
		(&response.Notice{}).Opcode(),
	}
)
