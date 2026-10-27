package game

import (
	"math"
	"math/rand"

	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/util"
)

func (c *EquipmentCore) rolledStatTotal(base uint16, maxRange int) int16 {
	if c == nil || base == 0 || maxRange < 0 {
		return 0
	}
	lMax := int(math.Min(math.Ceil(float64(base)*0.1), float64(maxRange)))
	offset := rand.Intn(2*lMax + 1)
	v := int(base) - lMax + offset
	return int16(util.ClampInt(v, -32768, 32767))
}

func (c *EquipmentCore) rollBonusDelta(base uint16, maxRange int) int16 {
	if c == nil {
		return 0
	}
	total := int32(c.rolledStatTotal(base, maxRange))
	return int16(util.ClampInt32(total-int32(base), -32768, 32767))
}

func (c *EquipmentCore) RandomizeStats(model wz.Equipment) {
	if c == nil || model == nil {
		return
	}
	c.bonusStats = &EquipmentBonusStats{}
	ab := model.GetAbility()
	const rMain = 5
	const rDefHPMP = 10
	if ab.Str != 0 {
		c.bonusStats.Str = c.rollBonusDelta(ab.Str, rMain)
	}
	if ab.Dex != 0 {
		c.bonusStats.Dex = c.rollBonusDelta(ab.Dex, rMain)
	}
	if ab.Int != 0 {
		c.bonusStats.Int = c.rollBonusDelta(ab.Int, rMain)
	}
	if ab.Luk != 0 {
		c.bonusStats.Luk = c.rollBonusDelta(ab.Luk, rMain)
	}
	if ab.PAD != 0 {
		c.bonusStats.PAD = c.rollBonusDelta(ab.PAD, rMain)
	}
	if ab.MAD != 0 {
		c.bonusStats.MAD = c.rollBonusDelta(ab.MAD, rMain)
	}
	if ab.ACC != 0 {
		c.bonusStats.ACC = c.rollBonusDelta(ab.ACC, rMain)
	}
	if ab.Avoid != 0 {
		c.bonusStats.Avoid = c.rollBonusDelta(ab.Avoid, rMain)
	}
	if ab.Jump != 0 {
		c.bonusStats.Jump = c.rollBonusDelta(ab.Jump, rMain)
	}
	if ab.Hands != 0 {
		c.bonusStats.Hands = c.rollBonusDelta(ab.Hands, rMain)
	}
	if ab.Speed != 0 {
		c.bonusStats.Speed = c.rollBonusDelta(ab.Speed, rMain)
	}
	if ab.PDD != 0 {
		c.bonusStats.PDD = c.rollBonusDelta(ab.PDD, rDefHPMP)
	}
	if ab.MDD != 0 {
		c.bonusStats.MDD = c.rollBonusDelta(ab.MDD, rDefHPMP)
	}
	if ab.MaxHP != 0 {
		c.bonusStats.MaxHP = c.rollBonusDelta(ab.MaxHP, rDefHPMP)
	}
	if ab.MaxMP != 0 {
		c.bonusStats.MaxMP = c.rollBonusDelta(ab.MaxMP, rDefHPMP)
	}
}

// A maker catalyst raises each stat the result already has, gems included, by up to a tenth; it never lowers one.
func (c *EquipmentCore) RandomizeStatsAbove(model wz.Equipment) {
	if c.bonusStats == nil {
		c.bonusStats = &EquipmentBonusStats{}
	}
	const rMain = 5
	const rDefHPMP = 10
	ab := model.GetAbility()
	for _, s := range []struct {
		base     uint16
		bonus    *int16
		maxRange int
	}{
		{ab.Str, &c.bonusStats.Str, rMain},
		{ab.Dex, &c.bonusStats.Dex, rMain},
		{ab.Int, &c.bonusStats.Int, rMain},
		{ab.Luk, &c.bonusStats.Luk, rMain},
		{ab.MAD, &c.bonusStats.MAD, rMain},
		{ab.PAD, &c.bonusStats.PAD, rMain},
		{ab.ACC, &c.bonusStats.ACC, rMain},
		{ab.Avoid, &c.bonusStats.Avoid, rMain},
		{ab.Jump, &c.bonusStats.Jump, rMain},
		{ab.Hands, &c.bonusStats.Hands, rMain},
		{ab.Speed, &c.bonusStats.Speed, rMain},
		{ab.PDD, &c.bonusStats.PDD, rDefHPMP},
		{ab.MDD, &c.bonusStats.MDD, rDefHPMP},
		{ab.MaxHP, &c.bonusStats.MaxHP, rDefHPMP},
		{ab.MaxMP, &c.bonusStats.MaxMP, rDefHPMP},
	} {
		total := int(s.base) + int(*s.bonus)
		if total <= 0 {
			continue
		}
		lMax := int(math.Min(math.Ceil(float64(total)*0.1), float64(s.maxRange)))
		*s.bonus = int16(util.ClampInt(int(*s.bonus)+rand.Intn(lMax+1), -32768, 32767))
	}
}

// A black crystal (option) or dark crystal (stat) moves each attack or main stat the result already has up or down.
func (c *EquipmentCore) RandomizeMakerOption(model wz.Equipment, option int16, stat int16) {
	if option == 0 && stat == 0 {
		return
	}
	if c.bonusStats == nil {
		c.bonusStats = &EquipmentBonusStats{}
	}
	ab := model.GetAbility()
	for _, s := range []struct {
		base  uint16
		bonus *int16
		shift int16
	}{
		{ab.PAD, &c.bonusStats.PAD, option},
		{ab.MAD, &c.bonusStats.MAD, option},
		{ab.Str, &c.bonusStats.Str, stat},
		{ab.Dex, &c.bonusStats.Dex, stat},
		{ab.Int, &c.bonusStats.Int, stat},
		{ab.Luk, &c.bonusStats.Luk, stat},
	} {
		if s.shift == 0 || int(s.base)+int(*s.bonus) <= 0 {
			continue
		}
		if rand.Intn(2) == 0 {
			*s.bonus += s.shift
		} else {
			*s.bonus -= s.shift
		}
	}
}
