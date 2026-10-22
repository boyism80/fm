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
	passedAt time.Time
	huntedAt time.Time
}

func (ch *Character) UseLieDetectorItem(targetName string, slot int16, itemID uint32) error {
	item := ch.inventory.GetItem(constant.InventoryTypeConsume, slot)
	if itemID != lieDetectorItemID || item == nil || item.GetModel().GetID() != itemID {
		return ErrLieDetectorInvalid
	}

	if err := ch.testLieDetector(targetName, false); err != nil {
		return err
	}
	ch.inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1)
	return nil
}

func (ch *Character) UseLieDetectorSkill(targetName string) error {
	if !ch.ActsAsGM() && ch.class != constant.ClassMapleleafBrigadier {
		return ErrLieDetectorInvalid
	}
	return ch.testLieDetector(targetName, true)
}

func (ch *Character) testLieDetector(targetName string, skill bool) error {
	m := ch.GetMap()
	if m == nil {
		return ErrLieDetectorInvalid
	}
	if m.wz.Limits(constant.FieldLimitLieDetector) {
		ch.listener.OnMessage(ch, constant.MsgPinkText, "이곳에서는 거짓말 탐지기를 사용할 수 없습니다.")
		return ErrLieDetectorInvalid
	}
	target := m.FindPlayerByName(targetName)
	if target == nil || target == ch || target.role > ch.role {
		ch.listener.OnLieDetectorResult(ch, pconst.LieDetectorResultNotFound, false)
		return ErrLieDetectorInvalid
	}
	// An item cannot pull a member out of a party quest or event; the GM skill still can.
	if !skill && target.stateMachine != nil {
		ch.listener.OnMessage(ch, constant.MsgPinkText, "이곳에서는 거짓말 탐지기를 사용할 수 없습니다.")
		return ErrLieDetectorInvalid
	}

	err := target.lieDetector.Start(ch, skill)
	switch {
	case errors.Is(err, ErrLieDetectorBusy):
		ch.listener.OnLieDetectorResult(ch, pconst.LieDetectorResultBusy, false)
	case errors.Is(err, ErrLieDetectorTested):
		ch.listener.OnLieDetectorResult(ch, pconst.LieDetectorResultTested, false)
	case errors.Is(err, ErrLieDetectorNotHunting):
		ch.listener.OnLieDetectorResult(ch, pconst.LieDetectorResultNotHunting, false)
	}
	return err
}

func (d *LieDetector) Start(tester *Character, skill bool) error {
	now := clock.Now()
	switch {
	case d.testerID != 0:
		return ErrLieDetectorBusy
	case !skill && now.Sub(d.passedAt) < lieDetectorImmunity:
		return ErrLieDetectorTested
	case !skill && now.Sub(d.huntedAt) > lieDetectorHuntWindow:
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
	tester.listener.OnLieDetectorStarted(tester, d.owner.GetName())
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
	d.owner.listener.OnLieDetectorQuestion(d.owner, captcha.Image, d.skill)

	wait := lieDetectorAnswerTime
	if limit := time.Duration(d.owner.gm.timerLimit) * time.Second; limit > 0 && limit < wait {
		wait = limit
	}

	// Timer callbacks run without an actor context; time out through the character's dispatch.
	d.owner.AddTimer(lieDetectorTimer, wait, false, func() {
		d.owner.gameWorld.Systems().Get[*DispatchSystem]().Visit(d.owner.GetID(), func(actx actor.Context, current *Character) {
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
	d.passedAt = clock.Now()

	d.owner.listener.OnLieDetectorResult(d.owner, pconst.LieDetectorResultPassed, d.skill)
	d.owner.inventory.GainMeso(lieDetectorRewardMeso)
	if tester != nil {
		tester.listener.OnMessage(tester, constant.MsgPinkText, "해당 유저가 거짓말 탐지기를 통과했습니다.")
	}
}

func (d *LieDetector) fail(actx actor.Context, failure lieDetectorFailure) {
	tester := d.tester()
	log.Printf("lie detector caught: target=%d tester=%d failure=%s", d.owner.GetID(), d.testerID, failure)
	d.testerID = 0
	d.owner.RemoveTimer(lieDetectorTimer)

	// A tester who already left the map is not paid and the target keeps the meso; missing that penalty is acceptable.
	if tester != nil {
		taken := min(d.owner.inventory.meso, lieDetectorPenaltyMeso)
		d.owner.inventory.RemoveMeso(taken)
		tester.inventory.GainMeso(taken)
		tester.listener.OnLieDetectorResult(tester, pconst.LieDetectorResultReported, false)
	}
	if failure == lieDetectorEscape {
		return
	}

	d.owner.listener.OnLieDetectorResult(d.owner, pconst.LieDetectorResultCaught, d.skill)
	jail := d.owner.gameWorld.Systems().Get[*MapSystem]().Get(lieDetectorJailMapID)
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
