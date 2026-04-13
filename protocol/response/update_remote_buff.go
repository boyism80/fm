package response

import (
	"time"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type UpdateRemoteBuff struct {
	CharacterID int32
	BuffID      int32
	Duration    time.Duration
	Buffs       []dto.BuffEntry
}

func (p *UpdateRemoteBuff) Opcode() uint16 { return 0x90 }

func (p *UpdateRemoteBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []dto.BuffEntry{}
	}
	buffs := make([]dto.BuffEntry, len(p.Buffs))
	copy(buffs, p.Buffs)
	SortBuffEntries(buffs)
	writer.Write32(p.CharacterID)
	buffFlags := make([]constant.BuffFlag, len(buffs))
	for i := range buffs {
		buffFlags[i] = buffs[i].Buff
	}
	WriteBuffs(writer, buffFlags)

	sec := p.Duration / time.Second
	var durSec uint16
	if sec > 65535 {
		durSec = 65535
	} else if sec <= 0 && p.Duration > 0 {
		durSec = 1
	} else {
		durSec = uint16(sec)
	}

	switch constant.SkillID(uint32(p.BuffID)) {
	case constant.SkillEnergyCharge, constant.SkillEnergyChargeCygnus:
		energy := int32(0)
		for _, buff := range buffs {
			if buff.Buff == constant.BuffFlagEnergyCharge {
				energy = buff.Value
				break
			}
		}
		energy = max(int32(0), min(energy, int32(10000)))
		writer.WriteU16(0)
		writer.Write32(energy)
		writer.Write64(0)
		writer.WriteU8(0)
		if energy >= 10000 {
			writer.Write32(int32(durSec))
		} else {
			writer.Write32(0)
		}
		return nil
	case constant.SkillDash, constant.SkillDashCygnus:
		writer.WriteU16(0)
		sid := int64(p.BuffID)
		for _, buff := range buffs {
			writer.Write32(buff.Value)
			writer.Write64(sid)
			writer.Write(make([]byte, 1))
			writer.WriteU16(durSec)
		}
		writer.WriteU16(0)
		writer.WriteU16(0)
		writer.WriteU8(1)
		writer.WriteU8(1)
		return nil
	case constant.SkillWindBooster, constant.SkillWindBoosterCygnus, constant.SkillTimeLeap:
		writer.WriteU16(0)
		sid := int64(p.BuffID)
		for _, buff := range buffs {
			writer.Write32(buff.Value)
			writer.Write64(sid)
			writer.Write(make([]byte, 6))
			writer.WriteU16(durSec)
		}
		writer.WriteU16(0)
		writer.WriteU16(0)
		writer.WriteU8(1)
		writer.WriteU8(1)
		return nil
	default:
		for _, e := range buffs {
			writer.WriteU16(uint16(e.Value))
		}
		writer.WriteU16(0)
		writer.WriteU16(0)
		return nil
	}
}

func (p *UpdateRemoteBuff) Deserialize(_ *stream.StreamReader) {}
