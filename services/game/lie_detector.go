package game

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
)

const (
	lieDetectorItemID      = 2190000
	lieDetectorAttempts    = 2
	lieDetectorAnswerTime  = 60 * time.Second
	lieDetectorImmunity    = 15 * time.Minute
	lieDetectorHuntWindow  = 60 * time.Second
	lieDetectorPenaltyMeso = 7000
	lieDetectorRewardMeso  = 5000
	lieDetectorJailMapID   = 910000000
	lieDetectorTimer       = "lie_detector"
)

var (
	ErrLieDetectorInvalid    = errors.New("invalid lie detector request")
	ErrLieDetectorBusy       = errors.New("lie detector already in progress")
	ErrLieDetectorTested     = errors.New("lie detector already passed")
	ErrLieDetectorNotHunting = errors.New("lie detector target is not hunting")
)

type lieDetectorFailure uint8

const (
	lieDetectorWrongAnswer lieDetectorFailure = iota
	lieDetectorTimeout
	lieDetectorEscape
)

func (f lieDetectorFailure) String() string {
	switch f {
	case lieDetectorWrongAnswer:
		return "wrong_answer"
	case lieDetectorTimeout:
		return "timeout"
	case lieDetectorEscape:
		return "escape"
	}
	return fmt.Sprintf("failure(%d)", uint8(f))
}

type LieDetector struct {
	owner    *Character
	testerID uint32
	answer   string
	attempts int
	skill    bool
	PassedAt time.Time
	HuntedAt time.Time
}

func (ch *Character) UseLieDetectorItem(targetName string, slot int16, itemID uint32) error {
	item := ch.Inventory.GetItem(constant.InventoryTypeConsume, slot)
	if itemID != lieDetectorItemID || item == nil || item.GetModel().GetID() != itemID {
		return ErrLieDetectorInvalid
	}

	if err := ch.testLieDetector(targetName, false); err != nil {
		return err
	}
	ch.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	return nil
}

func (ch *Character) UseLieDetectorSkill(targetName string) error {
	if !ch.ActsAsGM() && ch.Class != constant.ClassMapleleafBrigadier {
		return ErrLieDetectorInvalid
	}
	return ch.testLieDetector(targetName, true)
}

func (ch *Character) testLieDetector(targetName string, skill bool) error {
	m := ch.GetMap()
	if m == nil {
		return ErrLieDetectorInvalid
	}
	if m.Wz.Limits(constant.FieldLimitLieDetector) {
		ch.Listener.OnMessage(ch, constant.MsgPinkText, "이곳에서는 거짓말 탐지기를 사용할 수 없습니다.")
		return ErrLieDetectorInvalid
	}
	target := m.FindPlayerByName(targetName)
	if target == nil || target == ch || target.Role > ch.Role {
		ch.Listener.OnLieDetectorResult(ch, pconst.LieDetectorResultNotFound, false)
		return ErrLieDetectorInvalid
	}

	err := target.LieDetector.Start(ch, skill)
	switch {
	case errors.Is(err, ErrLieDetectorBusy):
		ch.Listener.OnLieDetectorResult(ch, pconst.LieDetectorResultBusy, false)
	case errors.Is(err, ErrLieDetectorTested):
		ch.Listener.OnLieDetectorResult(ch, pconst.LieDetectorResultTested, false)
	case errors.Is(err, ErrLieDetectorNotHunting):
		ch.Listener.OnLieDetectorResult(ch, pconst.LieDetectorResultNotHunting, false)
	}
	return err
}

func (d *LieDetector) Start(tester *Character, skill bool) error {
	now := clock.Now()
	switch {
	case d.testerID != 0:
		return ErrLieDetectorBusy
	case !skill && now.Sub(d.PassedAt) < lieDetectorImmunity:
		return ErrLieDetectorTested
	case !skill && now.Sub(d.HuntedAt) > lieDetectorHuntWindow:
		return ErrLieDetectorNotHunting
	}

	d.testerID = tester.GetID()
	d.attempts = 0
	d.skill = skill
	if err := d.ask(); err != nil {
		d.testerID = 0
		return err
	}
	log.Printf("lie detector started: target=%d tester=%d skill=%t", d.owner.GetID(), d.testerID, skill)
	tester.Listener.OnLieDetectorStarted(tester, d.owner.GetName())
	return nil
}

func (d *LieDetector) Answer(actx actor.Context, answer string) error {
	if d.testerID == 0 || answer == "" {
		return ErrLieDetectorInvalid
	}

	d.owner.RemoveTimer(lieDetectorTimer)
	if strings.EqualFold(answer, d.answer) {
		d.pass()
		return nil
	}
	d.miss(actx, lieDetectorWrongAnswer)
	return nil
}

func (d *LieDetector) Escape() {
	if d.testerID == 0 {
		return
	}
	d.fail(nil, lieDetectorEscape)
}

func (d *LieDetector) ask() error {
	captcha, err := NewCaptcha()
	if err != nil {
		return err
	}

	d.answer = captcha.Answer
	d.attempts++
	attempt := d.attempts
	d.owner.Listener.OnLieDetectorQuestion(d.owner, captcha.Image, d.skill)

	wait := lieDetectorAnswerTime
	if limit := time.Duration(d.owner.GM.TimerLimit) * time.Second; limit > 0 && limit < wait {
		wait = limit
	}

	// Timer callbacks run without an actor context; time out through the character's dispatch.
	d.owner.AddTimer(lieDetectorTimer, wait, false, func() {
		d.owner.GameWorld.Systems().Get[*DispatchSystem]().CallCharacter(d.owner.GetID(), func(actx actor.Context, current *Character) {
			if current != d.owner || d.testerID == 0 || d.attempts != attempt {
				return
			}
			d.miss(actx, lieDetectorTimeout)
		})
	})
	return nil
}

func (d *LieDetector) miss(actx actor.Context, failure lieDetectorFailure) {
	if d.attempts < lieDetectorAttempts && d.ask() == nil {
		return
	}
	d.fail(actx, failure)
}

func (d *LieDetector) pass() {
	tester := d.tester()
	log.Printf("lie detector passed: target=%d tester=%d attempt=%d", d.owner.GetID(), d.testerID, d.attempts)
	d.testerID = 0
	d.PassedAt = clock.Now()

	d.owner.Listener.OnLieDetectorResult(d.owner, pconst.LieDetectorResultPassed, d.skill)
	d.owner.Inventory.GainMeso(lieDetectorRewardMeso)
	if tester != nil {
		tester.Listener.OnMessage(tester, constant.MsgPinkText, "해당 유저가 거짓말 탐지기를 통과했습니다.")
	}
}

func (d *LieDetector) fail(actx actor.Context, failure lieDetectorFailure) {
	tester := d.tester()
	log.Printf("lie detector caught: target=%d tester=%d failure=%s", d.owner.GetID(), d.testerID, failure)
	d.testerID = 0
	d.owner.RemoveTimer(lieDetectorTimer)

	// A tester who already left the map is not paid and the target keeps the meso; missing that penalty is acceptable.
	if tester != nil {
		taken := min(d.owner.Inventory.Meso, lieDetectorPenaltyMeso)
		d.owner.Inventory.RemoveMeso(taken)
		tester.Inventory.GainMeso(taken)
		tester.Listener.OnLieDetectorResult(tester, pconst.LieDetectorResultReported, false)
	}
	if failure == lieDetectorEscape {
		return
	}

	d.owner.Listener.OnLieDetectorResult(d.owner, pconst.LieDetectorResultCaught, d.skill)
	jail := d.owner.GameWorld.Systems().Get[*MapSystem]().Get(lieDetectorJailMapID)
	if jail == nil {
		return
	}
	if err := d.owner.Warp(actx, jail, 0); err != nil {
		log.Printf("lie detector warp: target=%d: %v", d.owner.GetID(), err)
	}
}

func (d *LieDetector) tester() *Character {
	m := d.owner.GetMap()
	if m == nil {
		return nil
	}
	return m.GetPlayer(d.testerID)
}
