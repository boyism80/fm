package entity

import (
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
)

func (ch *Character) Inspect(targetID uint32) {
	ch.Listener.OnUnlockAction(ch)

	m := ch.GetMap()
	if m == nil {
		return
	}
	target := m.GetPlayer(targetID)
	if target == nil {
		return
	}
	if target.IsHidden() && ch.HasRoleAtLeast(target.Role) == false {
		return
	}

	profile := &response.CharacterProfile{
		CharacterID: target.GetID(),
		Level:       target.level,
		Job:         target.Class,
		Fame:        target.population,
		Married:     target.Marriage != nil && target.Marriage.Status == MarriageStatusMarried,
		GuildName:   "-",
		Self:        target.GetID() == ch.GetID(),
		Wishlist:    target.CashWishlist,
	}
	if guildID, ok := target.GetGuildID(); ok {
		if guild := target.GameWorld.GetGuildSystem().Get(guildID); guild != nil {
			profile.GuildName = guild.Name
			if alliance := guild.Alliance(); alliance != nil {
				profile.AllianceName = alliance.Name
			}
		}
	}
	if target.Pet != nil {
		profile.Pet = &response.CharacterProfilePet{
			ItemID:    target.Pet.Item.GetModel().GetID(),
			Name:      target.Pet.Item.Name,
			Level:     target.Pet.Item.Level,
			Closeness: target.Pet.Item.Closeness,
			Fullness:  target.Pet.Item.Fullness,
			Skills:    uint16(target.Pet.Item.Skills),
		}
		equip := target.Inventory.Equipped[constant.EquipmentPartsPetEquip.Cash()]
		if equip == nil {
			equip = target.Inventory.Equipped[constant.EquipmentPartsPetEquip]
		}
		if equip != nil {
			profile.Pet.EquipItemID = equip.GetModel().GetID()
		}
	}
	if target.Inventory.Equipped[constant.EquipmentPartsTamingMob] != nil && target.Inventory.Equipped[constant.EquipmentPartsSaddle] != nil {
		buffData := target.GetSpawnPlayerBuffData()
		profile.Mount = &response.CharacterProfileMount{
			Level:   buffData.MountLevel,
			Exp:     buffData.MountExp,
			Fatigue: buffData.MountFatigue,
		}
	}
	// TODO: send monster book level, card counts and cover mob once the monster book is loaded (#254)
	ch.Listener.OnInspect(ch, profile)
}
