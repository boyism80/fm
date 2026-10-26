package game

import (
	"slices"
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/maple/constant"
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
		fieldTeam: constant.CarnivalTeamNone,
		sendable:  sender,
		listener:  listener,
		gameWorld: gw,
		position: types.Vector2[int16]{
			X: int16(p.GetPositionX()),
			Y: int16(p.GetPositionY()),
		},
		baseHp:    p.GetMaxHp(),
		baseMp:    p.GetMaxMp(),
		stance:    uint8(p.GetStance()),
		id:        p.GetCharacterId(),
		accountID: p.GetAccountId(),
		name:      p.GetName(),
		look: Look{
			gender:    uint8(p.GetGender()),
			skinColor: uint8(p.GetSkinColor()),
			face:      p.GetFace(),
			hair:      p.GetHair(),
		},
		level:        uint8(p.GetLevel()),
		class:        constant.ClassType(p.GetClassId()),
		role:         constant.CharacterRole(p.GetRole()),
		gm:           GMMode{hidden: p.GetHidden()},
		exp:          p.GetExp(),
		expReachedAt: time.UnixMilli(p.GetExpReachedUnixMs()),

		random: [3]stream.RandomStream{stream.NewRandomStream(), stream.NewRandomStream(), stream.NewRandomStream()},
	}
	ch.stats = &Stats{owner: ch, base: BaseStats{str: uint16(p.GetStr()), dex: uint16(p.GetDex()), int: uint16(p.GetIntStat()), luk: uint16(p.GetLuk())}, population: int16(p.GetPopulation()), populationHistory: NewPopulationHistoryFromInternalProto(p.GetPopulationHistory())}
	ch.points = &Points{owner: ch, ap: uint16(p.GetAbilityPoint()), sp: uint16(p.GetSkillPoint()), hpmpUsed: uint16(p.GetHpApUsed())}
	ch.party = &PartyMembership{owner: ch}
	if partyID != nil {
		ch.party.id = *partyID
		ch.party.joined = true
	}
	ch.guild = &GuildMembership{owner: ch, id: guildID, invites: make(map[uint32]time.Time)}
	ch.buddies = NewBuddyList(ch)
	ch.family = &Family{owner: ch, tree: FamilyTreeFromProto(reply.GetFamily()), progress: p.GetFamilyProgress(), reputation: make(map[uint32]uint64)}
	for id, amount := range p.GetFamilyReputation() {
		ch.family.reputation[id], _ = strconv.ParseUint(amount, 10, 64)
	}
	ch.messenger = &Messenger{owner: ch, roomID: reply.GetMessengerRoomId()}
	ch.dialog = NewDialog(ch)
	ch.buffs = NewBuffs(ch)
	ch.cheat = &CheatTracker{owner: ch}
	ch.debuffs = &Debuffs{owner: ch, entries: make(map[constant.DebuffFlag]*Debuff)}
	ch.skills = NewSkills(ch)
	ch.quests = NewQuests(ch)
	ch.summons = NewSummons(ch)
	ch.doors = NewDoors(ch)
	ch.teleportStones = &TeleportStones{
		owner:   ch,
		regular: constant.NewTeleportStoneSlots(p.GetTeleportStones(), constant.TeleportStoneCount),
		vip:     constant.NewTeleportStoneSlots(p.GetVipTeleportStones(), constant.VipTeleportStoneCount),
	}
	ch.wedding = &Wedding{owner: ch, marriage: NewMarriageFromInternalProto(reply.GetMarriage())}
	ch.pets = &Pets{owner: ch, summoned: p.GetSummonedPet(), hpItem: p.GetPetHpItem(), mpItem: p.GetPetMpItem()}
	ch.mount = &Mount{owner: ch, level: max(reply.GetMount().GetLevel(), 1), exp: reply.GetMount().GetExp(), fatigue: reply.GetMount().GetFatigue()}
	ch.inventory = NewInventory(ch)
	ch.inventory.meso = p.GetMeso()
	ch.escrow = &Escrow{owner: ch}
	for i, limit := range p.GetSlotLimits() {
		if tab := ch.inventory.tabs[constant.InventoryType(i+1)]; tab != nil {
			tab.slotLimit = uint8(limit)
		}
	}
	ch.duey = &Duey{owner: ch}
	ch.notes = &Notes{owner: ch}
	ch.dojo = &Dojo{owner: ch}
	ch.lieDetector = &LieDetector{owner: ch}
	ch.carnival = &CarnivalMembership{owner: ch}
	ch.unsettled = NewUnsettledRequestsFromInternalProto(ch, reply.GetUnsettledRequests())
	ch.storeBank = &StoreBank{owner: ch}
	ch.remoteShops = &RemoteShops{owner: ch}
	ch.savedLocations = &SavedLocations{owner: ch, entries: make(map[string]uint32)}
	ch.records = NewRecordsFromInternalProto(reply.GetRecords())
	ch.account = NewAccountFromInternalProto(reply.GetAccount(), ch)
	ch.keyLayout = NewKeyLayout()
	ch.lifeCore.objectCore.self = ch
	ch.lifeCore.objectCore.initTimers()

	ch.keyLayout.LoadKeyLayoutProto(reply.GetKeyLayout())
	ch.cashWishlist = reply.GetCashWishlist()
	ch.monsterBook = NewMonsterBookFromInternalProto(ch, p.GetMonsterBookCover(), reply.GetMonsterBook())
	ch.inventory.Load(reply.GetInventory())
	ch.inventory.LoadRings(reply.GetRings())
	ch.skills.Load(reply.GetSkills())
	ch.skills.macros.Load(NewSkillMacrosFromInternalProto(reply.GetSkillMacros()))
	ch.buffs.Load(reply.GetBuffs())
	ch.debuffs.Load(reply.GetDebuffs())
	ch.SetHp(p.GetHp(), false)
	ch.SetMp(p.GetMp(), false)
	ch.quests.Load(reply.GetQuests())
	ch.savedLocations.Load(reply.GetSavedLocations())
	ch.buddies.LoadFromProto(reply.GetBuddies(), reply.GetBuddyCapacity())
	return ch
}

func equipmentLooksForPersist(ch *Character) (baseLooks, overlays map[int32]uint32, cashWeapon uint32) {
	baseLooks = make(map[int32]uint32)
	overlays = make(map[int32]uint32)
	if ch == nil {
		return
	}
	base, covered, cashWeapon := ch.inventory.look().split()
	for parts, itemID := range base {
		baseLooks[int32(parts)] = itemID
	}
	for parts, itemID := range covered {
		overlays[int32(parts)] = itemID
	}
	return
}

func (ch *Character) PersistMapID() uint32 {
	if ch == nil {
		return 0
	}
	m := ch.GetMap()
	if m == nil || m.wz == nil {
		return 0
	}
	if m.wz.ForcedReturn != 0 {
		return m.wz.ForcedReturn
	}
	if ch.GetHp() < 1 && m.wz.ReturnMapId > 0 {
		return uint32(m.wz.ReturnMapId)
	}
	return m.TemplateID()
}

func (ch *Character) BuildRPC(worldID uint32) *internal.CharacterSaveEntry {
	mapID := ch.PersistMapID()
	if mapID == 0 {
		return nil
	}
	baseLooks, overlays, cashWeapon := equipmentLooksForPersist(ch)
	inventory := ch.inventory.BuildRPC()
	slotLimits := make([]uint32, 0, constant.InventoryTypeCash)
	for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
		slotLimits = append(slotLimits, uint32(ch.inventory.tabs[typ].slotLimit))
	}
	reputation := make(map[uint32]string, len(ch.family.reputation))
	for id, amount := range ch.family.reputation {
		reputation[id] = strconv.FormatUint(amount, 10)
	}

	pb := &internal.Character{
		CharacterId:       ch.GetID(),
		AccountId:         ch.accountID,
		WorldId:           worldID,
		Name:              ch.GetName(),
		Gender:            uint32(ch.GetGender()),
		SkinColor:         uint32(ch.GetSkinColor()),
		Face:              ch.GetFace(),
		Hair:              ch.GetHair(),
		Level:             uint32(ch.GetLevel()),
		ClassId:           uint32(ch.class),
		Role:              uint32(ch.role),
		Hidden:            ch.IsHidden(),
		Str:               uint32(ch.stats.base.str),
		Dex:               uint32(ch.stats.base.dex),
		IntStat:           uint32(ch.stats.base.int),
		Luk:               uint32(ch.stats.base.luk),
		Hp:                ch.GetHp(),
		MaxHp:             ch.baseHp,
		Mp:                ch.GetMp(),
		MaxMp:             ch.baseMp,
		AbilityPoint:      uint32(ch.points.ap),
		Exp:               ch.GetExp(),
		MapId:             mapID,
		SpawnPoint:        uint32(ch.GetSpawnPoint()),
		PositionX:         int32(ch.position.X),
		PositionY:         int32(ch.position.Y),
		Stance:            uint32(ch.stance),
		Meso:              ch.escrow.mergeMeso(ch.inventory.meso),
		SkillPoint:        uint32(ch.points.sp),
		Population:        int32(ch.stats.population),
		FamilyProgress:    ch.family.progress,
		FamilyReputation:  reputation,
		PopulationHistory: ch.stats.populationHistory.BuildRPC(),
		HpApUsed:          uint32(ch.points.hpmpUsed),
		PetHpItem:         ch.pets.hpItem,
		PetMpItem:         ch.pets.mpItem,
		SummonedPet:       ch.pets.summoned,
		SlotLimits:        slotLimits,
		MonsterBookCover:  ch.monsterBook.cover,
		TeleportStones:    ch.teleportStones.regular.Registered(),
		VipTeleportStones: ch.teleportStones.vip.Registered(),
		SavedAtUnixMs:     clock.Now().UnixMilli(),
		ExpReachedUnixMs:  ch.expReachedAt.UnixMilli(),
	}
	return &internal.CharacterSaveEntry{
		Character:         pb,
		BaseLooks:         baseLooks,
		Overlays:          overlays,
		Inventory:         inventory,
		Skills:            ch.skills.BuildRPC(),
		SkillMacros:       ch.skills.macros.BuildRPC(),
		Buffs:             ch.buffs.BuildRPC(),
		KeyLayout:         ch.keyLayout.BuildRPC(),
		Quests:            ch.quests.BuildRPC(),
		SavedLocations:    ch.savedLocations.BuildRPC(),
		MonsterBook:       ch.monsterBook.BuildRPC(),
		Records:           ch.records.BuildRPC(),
		Account:           ch.account.BuildRPC(worldID),
		ClaimedShops:      slices.Clone(ch.storeBank.claimed),
		Mount:             &internal.Mount{Level: ch.mount.level, Exp: ch.mount.exp, Fatigue: ch.mount.fatigue},
		UnsettledRequests: ch.unsettled.BuildRPC(),
		CashWeapon:        cashWeapon,
	}
}
