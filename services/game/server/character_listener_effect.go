package server

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnChargeSkillEffect(ch *game.Character, skillID uint32, level, flags, speed uint8) {
	ch.Broadcast(&response.ChargeSkillEffect{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		Level:       level,
		Flags:       flags,
		Speed:       speed,
	}, nil)
}

func (l *CharacterListenerImpl) OnChargeSkillCancelled(ch *game.Character, skillID uint32) {
	ch.Broadcast(&response.CancelChargeSkillEffect{CharacterID: ch.GetID(), SkillID: skillID}, nil)
}

func (l *CharacterListenerImpl) OnPlaySound(ch *game.Character, sound string, broadcast bool) {
	if sound == "" {
		return
	}
	pkt := &response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeSound,
		Env:  sound,
	}
	if broadcast {
		if m := ch.GetMap(); m != nil {
			m.Broadcast(pkt, nil)
			return
		}
	}
	ch.Send(pkt, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPlayPortalSound(ch *game.Character) {
	ch.Send(&response.ShowSelfSkillEffect{
		Type:       pconst.SkillEffectTypePortal,
		SkillID:    0,
		SkillLevel: 1,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMapEffect(ch *game.Character, path string) {
	ch.Send(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeMapEffect,
		Env:  path,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowSelfSkillEffect(ch *game.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfSkillEffect{
		Type:       effectType,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		Additional: additional,
	}, types.SEND_POLICY_ENCRYPT)

	l.OnShowSkillEffect(ch, effectType, skillID, skillLevel, additional)
}

func (l *CharacterListenerImpl) OnShowSkillEffect(ch *game.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowSkillEffect{
		CharacterID: ch.GetID(),
		Type:        effectType,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
		Additional:  additional,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfEffect(ch *game.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfEffect{
		Type: effectType,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowEffect(ch, effectType)
}

func (l *CharacterListenerImpl) OnShowEffect(ch *game.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowEffect{
		CharacterID: ch.GetID(),
		Type:        effectType,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfDragonBloodEffect(ch *game.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfDragonBloodEffect{
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowDragonBloodEffect(ch, skillID, skillLevel)
}

func (l *CharacterListenerImpl) OnShowDragonBloodEffect(ch *game.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowDragonBloodEffect{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfHPHealedEffect(ch *game.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfHPHealedEffect{
		Amount: amount,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowHPHealedEffect(ch, amount)
}

func (l *CharacterListenerImpl) OnShowHPHealedEffect(ch *game.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowHPHealedEffect{
		CharacterID: ch.GetID(),
		Amount:      amount,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfRewardItemAnimation(ch *game.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfRewardItemAnimation{
		ItemID: itemID,
		Effect: effect,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowRewardItemAnimation(ch, itemID, effect)
}

func (l *CharacterListenerImpl) OnShowRewardItemAnimation(ch *game.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowRewardItemAnimation{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
		Effect:      effect,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfItemMakerSuccessEffect(ch *game.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfItemMakerSuccessEffect{}, types.SEND_POLICY_ENCRYPT)
	l.OnShowItemMakerSuccessEffect(ch)
}

func (l *CharacterListenerImpl) OnShowItemMakerSuccessEffect(ch *game.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowItemMakerSuccessEffect{
		CharacterID: ch.GetID(),
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfCraftingEffect(ch *game.Character, effect string, time int32, mode int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfCraftingEffect{
		Effect: effect,
		Time:   time,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowCraftingEffect(ch, effect, time, mode)
}

func (l *CharacterListenerImpl) OnShowCraftingEffect(ch *game.Character, effect string, time int32, mode int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowCraftingEffect{
		CharacterID: ch.GetID(),
		Effect:      effect,
		Time:        time,
		Mode:        mode,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfDiceEffect(ch *game.Character, effectID int32, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfDiceEffect{
		EffectID:   effectID,
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowDiceEffect(ch, effectID, skillID, skillLevel)
}

func (l *CharacterListenerImpl) OnShowDiceEffect(ch *game.Character, effectID int32, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowDiceEffect{
		CharacterID: ch.GetID(),
		EffectID:    effectID,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowScrollEffect(ch *game.Character, success bool, destroyedByCurse bool, legendarySpirit bool) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowScrollEffect{
		CharacterID:      ch.GetID(),
		Success:          success,
		DestroyedByCurse: destroyedByCurse,
		LegendarySpirit:  legendarySpirit,
	}, &game.ObjectBroadcastOption{
		WithMe: true,
	})
}
