package game

import (
	"slices"
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

func NewCharacterFromInternalProto(sender Sendable, listener CharacterListener, reply *internal.EnterGameReply, gw GameWorld) *Character {
	if listener == nil {
		panic("NewCharacterFromInternalProto: listener must not be nil")
	}
	p := reply.GetCharacter()
	var partyID, guildID *uint32
	if reply.PartyId != nil {
		v := *reply.PartyId
		partyID = &v
	}
	if reply.GuildId != nil {
		v := *reply.GuildId
		guildID = &v
	}
	ch := &Character{
		Sendable: sender,
		Listener: listener,
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				GameWorld: gw,
				Position: types.Vector2[int16]{
					X: int16(p.GetPositionX()),
					Y: int16(p.GetPositionY()),
				},
			},
			BaseHp: p.GetMaxHp(),
			BaseMp: p.GetMaxMp(),
			Stance: uint8(p.GetStance()),
		},
		id:        p.GetCharacterId(),
		AccountID: p.GetAccountId(),
		name:      p.GetName(),
		look: Look{
			Gender:    uint8(p.GetGender()),
			SkinColor: uint8(p.GetSkinColor()),
			Face:      p.GetFace(),
			Hair:      p.GetHair(),
		},
		level: uint8(p.GetLevel()),
		Class: constant.ClassType(p.GetClassId()),
		Role:  constant.CharacterRole(p.GetRole()),
		GM:    GMMode{Hidden: p.GetHidden()},
		exp:   p.GetExp(),

		random: [3]stream.RandomStream{stream.NewRandomStream(), stream.NewRandomStream(), stream.NewRandomStream()},
	}
	ch.Stats = &Stats{owner: ch, Base: BaseStats{Str: uint16(p.GetStr()), Dex: uint16(p.GetDex()), Int: uint16(p.GetIntStat()), Luk: uint16(p.GetLuk())}, Population: int16(p.GetPopulation()), PopulationHistory: NewPopulationHistoryFromInternalProto(p.GetPopulationHistory())}
	ch.Points = &Points{owner: ch, AP: uint16(p.GetAbilityPoint()), SP: uint16(p.GetSkillPoint()), HPMPUsed: uint16(p.GetHpApUsed())}
	ch.Party = &PartyMembership{owner: ch, id: partyID}
	ch.Guild = &GuildMembership{owner: ch, id: guildID, invites: make(map[uint32]time.Time)}
	ch.Buddies = NewBuddyList(ch)
	ch.Family = &Family{owner: ch, Tree: FamilyTree{family: reply.GetFamily()}, Progress: p.GetFamilyProgress(), Reputation: make(map[uint32]uint64)}
	for id, amount := range p.GetFamilyReputation() {
		ch.Family.Reputation[id], _ = strconv.ParseUint(amount, 10, 64)
	}
	ch.Messenger = &Messenger{owner: ch}
	ch.Dialog = NewDialog(ch)
	ch.Buffs = NewBuffs(ch)
	ch.Debuffs = &Debuffs{owner: ch, entries: make(map[constant.DebuffFlag]*Debuff)}
	ch.Skills = NewSkills(ch)
	ch.Quests = NewQuests(ch)
	ch.Summons = NewSummons(ch)
	ch.Doors = NewDoors(ch)
	ch.TeleportStones = &TeleportStones{
		owner:   ch,
		regular: constant.NewTeleportStoneSlots(p.GetTeleportStones(), constant.TeleportStoneCount),
		vip:     constant.NewTeleportStoneSlots(p.GetVipTeleportStones(), constant.VipTeleportStoneCount),
	}
	ch.Wedding = &Wedding{owner: ch, Marriage: NewMarriageFromInternalProto(reply.GetMarriage())}
	ch.Pets = &Pets{owner: ch, summoned: p.GetSummonedPet(), HPItem: p.GetPetHpItem(), MPItem: p.GetPetMpItem()}
	ch.Mount = &Mount{owner: ch, Level: max(reply.GetMount().GetLevel(), 1), Exp: reply.GetMount().GetExp(), Fatigue: reply.GetMount().GetFatigue()}
	ch.Inventory = NewInventory(ch)
	ch.Inventory.Meso = p.GetMeso()
	ch.Escrow = &Escrow{owner: ch}
	for i, limit := range p.GetSlotLimits() {
		if tab := ch.Inventory.Tabs[constant.InventoryType(i+1)]; tab != nil {
			tab.SlotLimit = uint8(limit)
		}
	}
	ch.Duey = &Duey{owner: ch}
	ch.Notes = &Notes{owner: ch}
	ch.Dojo = &Dojo{owner: ch}
	ch.LieDetector = &LieDetector{owner: ch}
	ch.Carnival = &CarnivalMembership{owner: ch}
	ch.Unsettled = NewUnsettledRequestsFromInternalProto(ch, reply.GetUnsettledRequests())
	ch.StoreBank = &StoreBank{owner: ch}
	ch.RemoteShops = &RemoteShops{owner: ch}
	ch.SavedLocations = &SavedLocations{owner: ch, entries: make(map[string]uint32)}
	ch.Records = NewRecordsFromInternalProto(reply.GetRecords())
	ch.Account = NewAccountFromInternalProto(reply.GetAccount(), ch)
	ch.KeyLayout = NewKeyLayout()
	ch.LifeCore.ObjectCore.self = ch
	ch.LifeCore.ObjectCore.initTimers()
	ch.LifeCore.setHp(p.GetHp())
	ch.LifeCore.setMp(p.GetMp())

	ch.KeyLayout.LoadKeyLayoutProto(reply.GetKeyLayout())
	ch.CashWishlist = reply.GetCashWishlist()
	ch.MonsterBook = NewMonsterBookFromInternalProto(ch, p.GetMonsterBookCover(), reply.GetMonsterBook())
	ch.Inventory.Load(reply.GetInventory())
	ch.Inventory.LoadRings(reply.GetRings())
	ch.Skills.Load(reply.GetSkills())
	ch.Skills.Macros.Load(NewSkillMacrosFromInternalProto(reply.GetSkillMacros()))
	ch.Buffs.Load(reply.GetBuffs())
	ch.Debuffs.Load(reply.GetDebuffs())
	ch.Quests.Load(reply.GetQuests())
	ch.SavedLocations.Load(reply.GetSavedLocations())
	ch.Buddies.LoadFromProto(reply.GetBuddies(), reply.GetBuddyCapacity())
	return ch
}

func equipmentLooksForPersist(ch *Character) (baseLooks, overlays map[int32]uint32) {
	baseLooks = make(map[int32]uint32)
	overlays = make(map[int32]uint32)
	if ch == nil {
		return
	}
	for parts, equipment := range ch.Inventory.Equipped {
		if equipment == nil || parts < -127 {
			continue
		}
		model, ok := equipment.GetModel().(wz.Equipment)
		if !ok {
			continue
		}
		absoluteParts := int32(parts * -1)
		if absoluteParts < 100 {
			if _, exists := baseLooks[absoluteParts]; !exists {
				baseLooks[absoluteParts] = model.GetID()
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := absoluteParts - 100
			if existing, exists := baseLooks[adjustedParts]; exists {
				overlays[adjustedParts] = existing
			}
			baseLooks[adjustedParts] = model.GetID()
		} else if _, exists := baseLooks[absoluteParts]; exists {
			overlays[absoluteParts] = model.GetID()
		}
	}
	return
}

func (ch *Character) PersistMapID() uint32 {
	if ch == nil {
		return 0
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return 0
	}
	if m.Wz.ForcedReturn != 0 {
		return m.Wz.ForcedReturn
	}
	if ch.GetHp() < 1 && m.Wz.ReturnMapId > 0 {
		return uint32(m.Wz.ReturnMapId)
	}
	return m.TemplateID()
}

func (ch *Character) BuildRPC(worldID uint32) *internal.CharacterSaveEntry {
	mapID := ch.PersistMapID()
	if mapID == 0 {
		return nil
	}
	baseLooks, overlays := equipmentLooksForPersist(ch)
	slotLimits := make([]uint32, 0, constant.InventoryTypeCash)
	for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
		slotLimits = append(slotLimits, uint32(ch.Inventory.Tabs[typ].SlotLimit))
	}
	reputation := make(map[uint32]string, len(ch.Family.Reputation))
	for id, amount := range ch.Family.Reputation {
		reputation[id] = strconv.FormatUint(amount, 10)
	}

	pb := &internal.Character{
		CharacterId:       ch.GetID(),
		AccountId:         ch.AccountID,
		WorldId:           worldID,
		Name:              ch.GetName(),
		Gender:            uint32(ch.GetGender()),
		SkinColor:         uint32(ch.GetSkinColor()),
		Face:              ch.GetFace(),
		Hair:              ch.GetHair(),
		Level:             uint32(ch.GetLevel()),
		ClassId:           uint32(ch.Class),
		Role:              uint32(ch.Role),
		Hidden:            ch.IsHidden(),
		Str:               uint32(ch.Stats.Base.Str),
		Dex:               uint32(ch.Stats.Base.Dex),
		IntStat:           uint32(ch.Stats.Base.Int),
		Luk:               uint32(ch.Stats.Base.Luk),
		Hp:                ch.GetHp(),
		MaxHp:             ch.BaseHp,
		Mp:                ch.GetMp(),
		MaxMp:             ch.BaseMp,
		AbilityPoint:      uint32(ch.Points.AP),
		Exp:               ch.GetExp(),
		MapId:             mapID,
		SpawnPoint:        uint32(ch.GetSpawnPoint()),
		PositionX:         int32(ch.Position.X),
		PositionY:         int32(ch.Position.Y),
		Stance:            uint32(ch.Stance),
		Meso:              ch.Escrow.mergeMeso(ch.Inventory.Meso),
		SkillPoint:        uint32(ch.Points.SP),
		Population:        int32(ch.Stats.Population),
		FamilyProgress:    ch.Family.Progress,
		FamilyReputation:  reputation,
		PopulationHistory: ch.Stats.PopulationHistory.BuildRPC(),
		HpApUsed:          uint32(ch.Points.HPMPUsed),
		PetHpItem:         ch.Pets.HPItem,
		PetMpItem:         ch.Pets.MPItem,
		SummonedPet:       ch.Pets.summoned,
		SlotLimits:        slotLimits,
		MonsterBookCover:  ch.MonsterBook.Cover,
		TeleportStones:    ch.TeleportStones.regular.Registered(),
		VipTeleportStones: ch.TeleportStones.vip.Registered(),
		SavedAtUnixMs:     clock.Now().UnixMilli(),
	}
	return &internal.CharacterSaveEntry{
		Character:         pb,
		BaseLooks:         baseLooks,
		Overlays:          overlays,
		Inventory:         ch.Inventory.BuildRPC(),
		Skills:            ch.Skills.BuildRPC(),
		SkillMacros:       ch.Skills.Macros.BuildRPC(),
		Buffs:             ch.Buffs.BuildRPC(),
		KeyLayout:         ch.KeyLayout.BuildRPC(),
		Quests:            ch.Quests.BuildRPC(),
		SavedLocations:    ch.SavedLocations.BuildRPC(),
		MonsterBook:       ch.MonsterBook.BuildRPC(),
		Records:           ch.Records.BuildRPC(),
		Account:           ch.Account.BuildRPC(worldID),
		ClaimedShops:      slices.Clone(ch.StoreBank.claimed),
		Mount:             &internal.Mount{Level: ch.Mount.Level, Exp: ch.Mount.Exp, Fatigue: ch.Mount.Fatigue},
		UnsettledRequests: ch.Unsettled.BuildRPC(),
	}
}
