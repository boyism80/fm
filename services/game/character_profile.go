package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
)

func (ch *Character) Inspect(targetID uint32) {
	ch.listener.OnUnlockAction(ch)

	m := ch.GetMap()
	if m == nil {
		return
	}
	target := m.GetPlayer(targetID)
	if target == nil {
		return
	}
	if target.IsHidden() && !ch.HasRoleAtLeast(target.role) {
		return
	}

	profile := &response.CharacterProfile{
		CharacterID: target.GetID(),
		Level:       target.level,
		Job:         uint16(target.class),
		Population:  target.stats.population,
		Married:     target.wedding.marriage != nil && target.wedding.marriage.status == MarriageStatusMarried,
		GuildName:   "-",
		Self:        target.GetID() == ch.GetID(),
		Wishlist:    target.cashWishlist,
	}
	if guildID, ok := target.guild.ID(); ok {
		if guild := target.gameWorld.Systems().Get[*GuildSystem]().Get(guildID); guild != nil {
			profile.GuildName = guild.name
			if alliance := guild.Alliance(); alliance != nil {
				profile.AllianceName = alliance.name
			}
		}
	}
	if target.pets.active != nil {
		profile.Pet = &response.CharacterProfilePet{
			ItemID:    target.pets.active.item.GetModel().GetID(),
			Name:      target.pets.active.item.name,
			Level:     target.pets.active.item.level,
			Closeness: target.pets.active.item.closeness,
			Fullness:  target.pets.active.item.fullness,
			Skills:    uint16(target.pets.active.item.skills),
		}
		equip := target.inventory.equipped[constant.EquipmentPartsPetEquip.Cash()]
		if equip == nil {
			equip = target.inventory.equipped[constant.EquipmentPartsPetEquip]
		}
		if equip != nil {
			profile.Pet.EquipItemID = equip.GetModel().GetID()
		}
	}
	if target.inventory.equipped[constant.EquipmentPartsTamingMob] != nil && target.inventory.equipped[constant.EquipmentPartsSaddle] != nil {
		profile.Mount = &response.CharacterProfileMount{
			Level:   target.mount.level,
			Exp:     target.mount.exp,
			Fatigue: target.mount.fatigue,
		}
	}
	profile.BookLevel = target.monsterBook.Level()
	profile.BookNormalCards, profile.BookSpecialCards = target.monsterBook.Count()
	if card, ok := target.itemModel(target.monsterBook.cover).(*wz.Consume); ok {
		profile.BookCoverMobID = card.MobID
	}
	ch.listener.OnInspect(ch, profile)
}
