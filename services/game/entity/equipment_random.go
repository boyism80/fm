package entity

import (
	"math"
	"math/rand"

	"github.com/boyism80/fm/services/game/wz"
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
	c.BonusStats = &EquipmentBonusStats{}
	ab := model.GetAbility()
	const rMain = 5
	const rDefHPMP = 10
	if ab.Str != 0 {
		c.BonusStats.Str = c.rollBonusDelta(ab.Str, rMain)
	}
	if ab.Dex != 0 {
		c.BonusStats.Dex = c.rollBonusDelta(ab.Dex, rMain)
	}
	if ab.Int != 0 {
		c.BonusStats.Int = c.rollBonusDelta(ab.Int, rMain)
	}
	if ab.Luk != 0 {
		c.BonusStats.Luk = c.rollBonusDelta(ab.Luk, rMain)
	}
	if ab.PAD != 0 {
		c.BonusStats.PAD = c.rollBonusDelta(ab.PAD, rMain)
	}
	if ab.MAD != 0 {
		c.BonusStats.MAD = c.rollBonusDelta(ab.MAD, rMain)
	}
	if ab.ACC != 0 {
		c.BonusStats.ACC = c.rollBonusDelta(ab.ACC, rMain)
	}
	if ab.Avoid != 0 {
		c.BonusStats.Avoid = c.rollBonusDelta(ab.Avoid, rMain)
	}
	if ab.Jump != 0 {
		c.BonusStats.Jump = c.rollBonusDelta(ab.Jump, rMain)
	}
	if ab.Hands != 0 {
		c.BonusStats.Hands = c.rollBonusDelta(ab.Hands, rMain)
	}
	if ab.Speed != 0 {
		c.BonusStats.Speed = c.rollBonusDelta(ab.Speed, rMain)
	}
	if ab.PDD != 0 {
		c.BonusStats.PDD = c.rollBonusDelta(ab.PDD, rDefHPMP)
	}
	if ab.MDD != 0 {
		c.BonusStats.MDD = c.rollBonusDelta(ab.MDD, rDefHPMP)
	}
	if ab.MaxHP != 0 {
		c.BonusStats.MaxHP = c.rollBonusDelta(ab.MaxHP, rDefHPMP)
	}
	if ab.MaxMP != 0 {
		c.BonusStats.MaxMP = c.rollBonusDelta(ab.MaxMP, rDefHPMP)
	}
}
