package domain

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/util"
)

type EquipmentBonusStats struct {
	Str   int16 `json:"str,omitempty"`
	Dex   int16 `json:"dex,omitempty"`
	Int   int16 `json:"int,omitempty"`
	Luk   int16 `json:"luk,omitempty"`
	MaxHP int16 `json:"maxHp,omitempty"`
	MaxMP int16 `json:"maxMp,omitempty"`
	PAD   int16 `json:"pad,omitempty"`
	MAD   int16 `json:"mad,omitempty"`
	PDD   int16 `json:"pdd,omitempty"`
	MDD   int16 `json:"mdd,omitempty"`
	ACC   int16 `json:"acc,omitempty"`
	Avoid int16 `json:"avoid,omitempty"`
	Hands int16 `json:"hands,omitempty"`
	Speed int16 `json:"speed,omitempty"`
	Jump  int16 `json:"jump,omitempty"`
}

func (s *EquipmentBonusStats) stat(key string) *int16 {
	switch key {
	case "str":
		return &s.Str
	case "dex":
		return &s.Dex
	case "int":
		return &s.Int
	case "luk":
		return &s.Luk
	case "hp":
		return &s.MaxHP
	case "mp":
		return &s.MaxMP
	case "watk":
		return &s.PAD
	case "matk":
		return &s.MAD
	case "wdef":
		return &s.PDD
	case "mdef":
		return &s.MDD
	case "acc":
		return &s.ACC
	case "avoid":
		return &s.Avoid
	case "hands":
		return &s.Hands
	case "speed":
		return &s.Speed
	case "jump":
		return &s.Jump
	default:
		return nil
	}
}

func (c *EquipmentCore) SetBonus(stats map[string]int16) {
	if c == nil || len(stats) == 0 {
		return
	}
	if c.BonusStats == nil {
		c.BonusStats = &EquipmentBonusStats{}
	}
	for key, value := range stats {
		if stat := c.BonusStats.stat(key); stat != nil {
			*stat = value
		}
	}
}

func (c *EquipmentCore) AddBonus(stats map[string]int16) {
	if c == nil || len(stats) == 0 {
		return
	}
	if c.BonusStats == nil {
		c.BonusStats = &EquipmentBonusStats{}
	}
	for key, value := range stats {
		if stat := c.BonusStats.stat(key); stat != nil {
			*stat += value
		}
	}
}

// Matches the client's maker grade: hands are excluded and HP/MP count a tenth.
func (s *EquipmentBonusStats) quality() int {
	if s == nil {
		return 0
	}
	total := int(s.Str) + int(s.Dex) + int(s.Int) + int(s.Luk) + int(s.PAD) + int(s.MAD) + int(s.PDD) + int(s.MDD) + int(s.ACC) + int(s.Avoid) + int(s.Speed) + int(s.Jump) + int(s.MaxHP)/10 + int(s.MaxMP)/10
	switch {
	case total < 0:
		return -1
	case total < 6:
		return 0
	case total < 23:
		return 1
	case total < 40:
		return 2
	default:
		return 3
	}
}

func EquipmentBonusStatsFromProto(pb *internal.EquipmentBonusStats) *EquipmentBonusStats {
	if pb == nil {
		return nil
	}
	var s EquipmentBonusStats
	s.Str = int16(util.ClampInt32(pb.GetStr(), -32768, 32767))
	s.Dex = int16(util.ClampInt32(pb.GetDex(), -32768, 32767))
	s.Int = int16(util.ClampInt32(pb.GetIntStat(), -32768, 32767))
	s.Luk = int16(util.ClampInt32(pb.GetLuk(), -32768, 32767))
	s.MaxHP = int16(util.ClampInt32(pb.GetMaxHp(), -32768, 32767))
	s.MaxMP = int16(util.ClampInt32(pb.GetMaxMp(), -32768, 32767))
	s.PAD = int16(util.ClampInt32(pb.GetPad(), -32768, 32767))
	s.MAD = int16(util.ClampInt32(pb.GetMad(), -32768, 32767))
	s.PDD = int16(util.ClampInt32(pb.GetPdd(), -32768, 32767))
	s.MDD = int16(util.ClampInt32(pb.GetMdd(), -32768, 32767))
	s.ACC = int16(util.ClampInt32(pb.GetAcc(), -32768, 32767))
	s.Avoid = int16(util.ClampInt32(pb.GetAvoid(), -32768, 32767))
	s.Hands = int16(util.ClampInt32(pb.GetHands(), -32768, 32767))
	s.Speed = int16(util.ClampInt32(pb.GetSpeed(), -32768, 32767))
	s.Jump = int16(util.ClampInt32(pb.GetJump(), -32768, 32767))
	return &s
}

func (s *EquipmentBonusStats) BuildRPC() *internal.EquipmentBonusStats {
	if s == nil {
		return nil
	}
	return &internal.EquipmentBonusStats{
		Str:     int32(s.Str),
		Dex:     int32(s.Dex),
		IntStat: int32(s.Int),
		Luk:     int32(s.Luk),
		MaxHp:   int32(s.MaxHP),
		MaxMp:   int32(s.MaxMP),
		Pad:     int32(s.PAD),
		Mad:     int32(s.MAD),
		Pdd:     int32(s.PDD),
		Mdd:     int32(s.MDD),
		Acc:     int32(s.ACC),
		Avoid:   int32(s.Avoid),
		Hands:   int32(s.Hands),
		Speed:   int32(s.Speed),
		Jump:    int32(s.Jump),
	}
}
