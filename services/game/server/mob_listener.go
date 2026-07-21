package server

import (
	"slices"
	"time"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type MobListenerImpl struct{}

func (l *MobListenerImpl) OnMobBuffApplied(mob *entity.Mob, ent *entity.MobBuff, addedReflections []int32, remaining time.Duration) {
	if mob == nil || ent == nil || len(ent.Values) == 0 {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}
	buffTime := int16(32767)
	if remaining > 0 {
		sec := remaining / time.Second
		if sec < 32767 {
			buffTime = int16(sec)
		}
	}
	skillID := uint16(0)
	if ent.Wz != nil {
		skillID = uint16(ent.Wz.ID)
	}
	flags := make([]constant.MobBuffFlag, 0, len(ent.Values))
	for f := range ent.Values {
		flags = append(flags, f)
	}
	slices.Sort(flags)
	if len(addedReflections) > 0 {
		var mask int32
		entries := make([]response.ApplyMobBuffReflectEntry, 0, len(flags))
		for _, flag := range flags {
			mask |= int32(flag)
			entries = append(entries, response.ApplyMobBuffReflectEntry{
				X:          int16(ent.Values[flag]),
				SkillID:    skillID,
				SkillLevel: uint16(ent.SkillLevel),
				BuffTime:   buffTime,
			})
		}
		pkt := &response.ApplyMobBuffReflect{
			OID:         mob.OID,
			StatusMask:  mask,
			Entries:     entries,
			Reflections: addedReflections,
			Delay:       0,
			StatusSize:  byte(len(flags)),
		}
		mapInstance.Broadcast(pkt, nil)
		return
	}
	for _, flag := range flags {
		pkt := &response.ApplyMobBuff{
			OID:        mob.OID,
			Status:     int32(flag),
			X:          int16(ent.Values[flag]),
			SkillID:    skillID,
			SkillLevel: uint16(ent.SkillLevel),
			BuffTime:   buffTime,
			Delay:      0,
			StatusSize: 1,
		}
		mapInstance.Broadcast(pkt, nil)
	}
}

func (l *MobListenerImpl) OnMobBuffCancelled(mob *entity.Mob, buff constant.MobBuffFlag) {
	if mob == nil {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}
	pkt := &response.CancelMobBuff{
		OID:    mob.OID,
		Status: int32(buff),
		Size:   1,
	}
	mapInstance.Broadcast(pkt, nil)
}

func (l *MobListenerImpl) OnMobDamaged(mob *entity.Mob, amount int32) {
	if mob == nil {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}
	pkt := &response.DamageMob{
		OID:     mob.OID,
		Display: constant.MobDamageDisplayNormal,
		Damage:  amount,
	}
	mapInstance.Broadcast(pkt, nil)
}

func (l *MobListenerImpl) OnMobAllyDamaged(mob *entity.Mob, amount int32) {
	if mob == nil {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}
	pkt := &response.DamageMob{
		OID:     mob.OID,
		Display: constant.MobDamageDisplayAllyShowHp,
		Damage:  amount,
		HP:      int32(mob.GetHp()),
		MaxHP:   int32(mob.GetMaxHp()),
	}
	mapInstance.Broadcast(pkt, nil)
}

func (l *MobListenerImpl) OnShowBossHp(mob *entity.Mob, clear bool) {
	if mob == nil {
		return
	}
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}

	mobID := uint32(0)
	tagColor := uint8(0)
	tagBgColor := uint8(0)
	if mob.Wz != nil {
		mobID = mob.Wz.ID
		tagColor = mob.Wz.HpTagColor
		tagBgColor = mob.Wz.HpTagBgColor
	}

	maxHpVal := mob.GetMaxHp()
	maxHp := int32(maxHpVal)
	if maxHpVal > 0x7fffffff {
		maxHp = 0x7fffffff
	}

	currentHp := int32(-1)
	if !clear {
		hp := mob.GetHp()
		if hp > 0x7fffffff {
			ratio := float64(hp) / float64(maxHpVal)
			currentHp = int32(ratio * float64(0x7fffffff))
		} else {
			currentHp = int32(hp)
		}
	}

	mapInstance.Broadcast(&response.ShowBossHp{
		MobID:      mobID,
		CurrentHP:  currentHp,
		MaxHP:      maxHp,
		TagColor:   tagColor,
		TagBgColor: tagBgColor,
	}, nil)
}
