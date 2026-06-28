package entity

import (
	"fmt"
	"sync"
	"time"

	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Character struct {
	LifeCore
	Sendable

	dialog            *lua.LState
	id                uint32
	name              string
	gender            uint8
	skinColor         uint8
	face              uint32
	hair              uint32
	level             uint8
	rank              uint32
	rankDiff          int32
	classRank         uint32
	classRankDiff     int32
	exp               uint32
	famePoint         uint16
	mega              bool
	random1           stream.RandomStream
	random2           stream.RandomStream
	random3           stream.RandomStream
	questStatuses     map[int]*QuestStatus
	marriageId        uint32
	regRocks          []uint32
	rocks             []uint32
	monsterBookCover  uint32
	monsterBook       *MonsterBook
	quests            map[uint16]string
	luaDialog         *lua.LState
	dialogMutex       sync.Mutex
	hidden            bool
	Listener          CharacterListener
	Class             uint16
	Role              constant.CharacterRole
	AccountID         uint32
	AbilityPoint      uint16
	SkillPoint        uint16
	HpApUsed          uint16
	Meso              int32
	Inventory         map[constant.InventoryType]*Inventory
	Equipments        map[constant.EquipmentPartsType]Equipment
	Rings             RingContainer
	Skills            *SkillContainer
	keyLayout         *KeyLayout
	CurrentShopID     uint32
	Chair             uint32
	LastHealHPTime    time.Time
	LastHealMPTime    time.Time
	BaseStats         BaseStats
	BonusStats        BonusStats
	Buffs             *BuffContainer
	debuffs           map[constant.DebuffFlag]*Debuff
	summons           map[constant.SkillID]*Summon
	doors             map[constant.SkillID]*Door
	HomingTargetOID   *uint32
	partyID           *uint32
	guildID           *uint32
	GuildInvites      map[uint32]time.Time
	partySearchConfig *PartySearchConfig
	buddyList         *BuddyList
	InstantKill       bool
}

type Debuff struct {
	Flag      constant.DebuffFlag
	StartTime time.Time
	Duration  time.Duration
}

func (ch *Character) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeCharacter
}

func (ch *Character) Is(typ constant.ObjectType) bool {
	return ch.GetObjectType().Has(typ)
}

func (ch *Character) SendSpawnSyncToViewer(viewer *Character) {
	if ch == nil || viewer == nil {
		return
	}
	if ch.GetID() == viewer.GetID() {
		return
	}
	if ch.IsHidden() && !viewer.HasRoleAtLeast(ch.Role) {
		return
	}
	spawnBuffData := ch.GetSpawnPlayerBuffData()
	spawnPacket := &response.SpawnPlayer{
		Character:         ch.ToDTO(),
		BuffStates:        spawnBuffData.BuffStates,
		Diseases:          ch.GetDiseaseMask(),
		SpeedBuff:         spawnBuffData.SpeedBuff,
		ComboCount:        spawnBuffData.ComboCount,
		WKChargeSkillId:   spawnBuffData.WKChargeSkillID,
		MorphId:           spawnBuffData.MorphID,
		SpiritClawSkillId: spawnBuffData.SpiritClawSkillID,
		MountLevel:        spawnBuffData.MountLevel,
		MountExp:          spawnBuffData.MountExp,
		MountFatigue:      spawnBuffData.MountFatigue,
		CrushRings:        RingsToDTO(ch.Rings.Left),
		FriendshipRings:   RingsToDTO(ch.Rings.Mid),
		MarriageRings:     RingsToDTO(ch.Rings.Right),
	}
	if guildID, ok := ch.GetGuildID(); ok && ch.GameWorld != nil {
		if guild := ch.GameWorld.GetGuildSystem().Get(guildID); guild != nil {
			spawnPacket.GuildName = guild.Name
			if guild.Logo != nil {
				spawnPacket.GuildLogoBG = uint16(guild.Logo.LogoBG)
				spawnPacket.GuildLogoBGColor = uint8(guild.Logo.LogoBGColor)
				spawnPacket.GuildLogo = uint16(guild.Logo.Logo)
				spawnPacket.GuildLogoColor = uint8(guild.Logo.LogoColor)
			}
		}
	}
	viewer.Send(spawnPacket, types.SEND_POLICY_ENCRYPT)

	if mountID, active := ch.GetRiddingInfo(); active {
		viewer.Send(&response.UpdateRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagMonsterRiding, Value: mountID}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if buff, _, active := ch.Buffs.GetBuffValue(constant.BuffFlagEnergyCharge); active {
		viewer.Send(&response.UpdateBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buff.GetBuffID(),
			Duration:    50 * time.Second,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagEnergyCharge, Value: 10000}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if dashBuff, dashValue, active := ch.Buffs.GetBuffValue(constant.BuffFlagDashSpeed); active {
		buffs := []dto.BuffEntry{{Buff: constant.BuffFlagDashSpeed, Value: dashValue}}
		if _, dashJumpValue, hasDashJump := ch.Buffs.GetBuffValue(constant.BuffFlagDashJump); hasDashJump {
			buffs = append(buffs, dto.BuffEntry{Buff: constant.BuffFlagDashJump, Value: dashJumpValue})
		}
		viewer.Send(&response.UpdateBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      dashBuff.GetBuffID(),
			Duration:    dashBuff.RemainingDuration(time.Now()),
			Buffs:       buffs,
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func summonTimerKey(skillID constant.SkillID) string {
	return fmt.Sprintf("summon:%d", skillID)
}

func (ch *Character) SpawnSummon(skillID constant.SkillID, skillLevel uint8, movementType constant.SummonMovementType, summonType constant.SummonType, position types.Point[int16], duration time.Duration) *Summon {
	if ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	if ch.summons != nil {
		if current, ok := ch.summons[skillID]; ok && current != nil {
			ch.RemoveSummon(current, true)
		}
	}

	s := &Summon{
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				Position:  position,
				GameWorld: m.GameWorld,
				Map:       nil,
			},
			hp:     1,
			BaseHp: 1,
			BaseMp: 1,
		},
		Owner:        ch,
		OwnerID:      ch.GetID(),
		SkillID:      skillID,
		SkillLevel:   skillLevel,
		MovementType: movementType,
		SummonType:   summonType,
	}
	s.LifeCore.ObjectCore.self = s
	if ch.summons == nil {
		ch.summons = make(map[constant.SkillID]*Summon)
	}
	ch.summons[skillID] = s
	m.AddSummon(s)
	if duration > 0 {
		_ = ch.AddTimer(summonTimerKey(s.SkillID), duration, false, func() {
			ch.handleSummonExpireBySkill(s.SkillID)
		})
	}
	return s
}

func (ch *Character) SpawnMist(skill *SkillEntry, position types.Point[int16], mistType constant.MistType, bounds types.Rect[int32], duration time.Duration, initialDelay time.Duration, poisonTickMultiplier float64) *Mist {
	if ch == nil || ch.GameWorld == nil {
		return nil
	}
	if skill == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	wzSkill := ch.GameWorld.GetResources().GetSkill(skill.Wz.ID)
	if wzSkill == nil {
		return nil
	}
	skillLevel := uint8(skill.Level())
	ld := wzSkill.GetLevelData(int(skillLevel))
	b := bounds
	if b.Left == 0 && b.Right == 0 && b.Top == 0 && b.Bottom == 0 {
		b = MistWorldBounds(position, ld)
	}
	if poisonTickMultiplier <= 0 {
		poisonTickMultiplier = 1.0
	}
	mist := &Mist{
		ObjectCore: ObjectCore{
			Position:  position,
			GameWorld: m.GameWorld,
			Map:       nil,
		},
		Causer:               ch.GetID(),
		SkillWz:              wzSkill,
		SkillLevel:           skillLevel,
		MistType:             mistType,
		MobMist:              false,
		MobSkill:             false,
		SkillDelay:           8,
		Bounds:               b,
		ExpiresAt:            time.Time{},
		PoisonTickMultiplier: poisonTickMultiplier,
	}
	mist.ObjectCore.self = mist
	if initialDelay > 0 {
		mist.NextPoisonTickAt = time.Now().Add(initialDelay)
	}
	if duration > 0 {
		mist.ExpiresAt = time.Now().Add(duration)
	}
	m.AddMist(mist)
	return mist
}

func (ch *Character) SpawnDoor(skillID constant.SkillID) {
	if ch == nil {
		return
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return
	}
	if ch.doors != nil {
		if current, ok := ch.doors[skillID]; ok && current != nil {
			ch.RemoveDoor(current, true)
		}
	}

	if skillID != constant.SkillMysticDoor {
		return
	}
	destMapID := uint32(m.Wz.ReturnMapId)
	if destMapID == 0 || destMapID == uint32(m.Wz.ID) || ch.GameWorld == nil {
		return
	}
	if gw := ch.GameWorld; gw != nil {
		gw.GetMapSystem().CreateReturnDoor(ch, skillID)
	}
}

func (ch *Character) SpawnFieldMapDoor(skillID constant.SkillID, returnPortalID uint8, townPortalPos types.Vector2[int16], fieldPortalID uint8) *Door {
	if ch == nil {
		return nil
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return nil
	}
	destMapID := uint32(m.Wz.ReturnMapId)
	fieldPos := ch.Position
	door := &Door{
		ObjectCore: ObjectCore{
			Position:  ch.Position,
			GameWorld: m.GameWorld,
			Map:       nil,
		},
		OwnerID:            ch.GetID(),
		SkillID:            skillID,
		ReturnMapID:        destMapID,
		FieldMapID:         uint32(m.Wz.ID),
		ReturnPortalID:     returnPortalID,
		FieldPortalID:      fieldPortalID,
		PartyID:            ch.partyID,
		FieldPosition:      fieldPos,
		TownPortalPosition: townPortalPos,
	}
	door.ObjectCore.self = door
	if ch.doors == nil {
		ch.doors = make(map[constant.SkillID]*Door)
	}
	ch.doors[skillID] = door
	m.AddDoor(door)
	ch.Listener.OnPartyMemberFieldsChanged(ch)
	return door
}

func (ch *Character) forgetDoorRegistrationIfSame(door *Door) {
	if ch == nil || door == nil {
		return
	}
	if ch.doors != nil {
		key := constant.SkillID(door.SkillID)
		if d := ch.doors[key]; d == door {
			delete(ch.doors, key)
		}
	}
}

func (ch *Character) RemoveMist(mist *Mist) {
	if mist == nil {
		return
	}
	if m := mist.GetMap(); m != nil && mist.OID != 0 {
		m.RemoveMist(mist.OID)
	}
}

func (ch *Character) RemoveSummon(target *Summon, animated bool) {
	if target == nil {
		return
	}
	_ = ch.RemoveTimer(summonTimerKey(target.SkillID))
	if m := target.GetMap(); m != nil && target.OID != 0 {
		m.RemoveSummon(target.OID, animated)
	}
	if ch.summons != nil {
		delete(ch.summons, constant.SkillID(target.SkillID))
	}
}

func (ch *Character) RemoveDoor(target *Door, animated bool) {
	if target == nil {
		return
	}
	if m := target.GetMap(); m != nil && target.OID != 0 {
		m.RemoveDoor(target.OID, animated)
	}
	if ch.doors != nil {
		delete(ch.doors, constant.SkillID(target.SkillID))
	}
	ch.Listener.OnPartyMemberFieldsChanged(ch)
}

func (ch *Character) RemoveDoorBySkill(skillID constant.SkillID, animated bool) {
	if ch.doors == nil {
		return
	}
	door := ch.doors[skillID]
	if door == nil {
		return
	}
	ch.RemoveDoor(door, animated)
}

func (ch *Character) GetSummons() []*Summon {
	if len(ch.summons) == 0 {
		return nil
	}
	out := make([]*Summon, 0, len(ch.summons))
	for _, s := range ch.summons {
		if s != nil {
			out = append(out, s)
		}
	}
	return out
}

func (ch *Character) GetSummon(skillId constant.SkillID) *Summon {
	return ch.summons[skillId]
}

func (ch *Character) GetSummonsSize() int {
	return len(ch.summons)
}

func (ch *Character) ClearSummons() {
	if len(ch.summons) == 0 {
		return
	}
	for _, s := range ch.GetSummons() {
		ch.RemoveSummon(s, true)
	}
}

func (ch *Character) GetDoors() []*Door {
	if len(ch.doors) == 0 {
		return nil
	}
	out := make([]*Door, 0, len(ch.doors))
	for _, d := range ch.doors {
		if d != nil {
			out = append(out, d)
		}
	}
	return out
}

func (ch *Character) GetDoor(index int) *Door {
	if ch.doors == nil {
		return nil
	}
	doors := ch.GetDoors()
	if doors == nil || index < 0 || index >= len(doors) {
		return nil
	}
	return doors[index]
}

func (ch *Character) ClearDoors() {
	if len(ch.doors) == 0 {
		return
	}
	for _, d := range ch.GetDoors() {
		ch.RemoveDoor(d, true)
	}
}

func (ch *Character) handleSummonExpireBySkill(skillID constant.SkillID) {
	if ch.summons == nil {
		return
	}
	s := ch.summons[skillID]
	if s == nil {
		return
	}
	ch.RemoveSummon(s, true)
}

func (ch *Character) GetBonusHp() int32   { return ch.BonusHp }
func (ch *Character) GetBonusMp() int32   { return ch.BonusMp }
func (ch *Character) GetInvincible() bool { return ch.Invincible }
func (ch *Character) GetGender() uint8    { return ch.gender }
func (ch *Character) GetSkinColor() uint8 { return ch.skinColor }
func (ch *Character) GetFace() uint32     { return ch.face }
func (ch *Character) GetHair() uint32     { return ch.hair }
func (ch *Character) GetLevel() uint8     { return ch.level }
func (ch *Character) GetExp() uint32      { return ch.exp }
func (ch *Character) GetSpawnPoint() uint8 {
	mapInstance := ch.GetMap()
	if mapInstance != nil && mapInstance.Wz != nil {
		return mapInstance.Wz.FindClosestPortalSpawnID(ch.Position)
	}
	return 0
}

func (ch *Character) SetHp(v uint32, notify bool) {
	ch.LifeCore.setHp(v)
	if !notify {
		return
	}
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatHP: int32(ch.GetHp())}, false)
}

func (ch *Character) SetMp(v uint32, notify bool) {
	ch.LifeCore.setMp(v)
	if !notify {
		return
	}
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatMP: int32(ch.GetMp())}, false)
}

func (ch *Character) SetBonusHp(v int32, notify bool) {
	ch.LifeCore.SetBonusHp(v, false)
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetBonusMp(v int32, notify bool) {
	ch.LifeCore.SetBonusMp(v, false)
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetMaxHpPercent(p int16, notify bool) {
	ch.BonusStats.MaxHpPercent = p
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetMaxMpPercent(p int16, notify bool) {
	ch.BonusStats.MaxMpPercent = p
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetInvincible(b bool) { ch.Invincible = b }

func (ch *Character) GetInstantKill() bool { return ch.InstantKill }

func (ch *Character) SetInstantKill(b bool) { ch.InstantKill = b }

func (ch *Character) AddHp(amount int) {
	n := int(ch.GetHp()) + amount
	if n < 0 {
		n = 0
	}
	maxHp := int(ch.GetMaxHp())
	if n > maxHp {
		n = maxHp
	}
	ch.SetHp(uint32(n), false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatHP: int32(ch.GetHp())}, false)
}

func (ch *Character) AddMp(amount int) {
	n := int(ch.GetMp()) + amount
	if n < 0 {
		n = 0
	}
	maxMp := int(ch.GetMaxMp())
	if n > maxMp {
		n = maxMp
	}
	ch.SetMp(uint32(n), false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatMP: int32(ch.GetMp())}, false)
}

func (ch *Character) AddHpMp(hpDelta, mpDelta int) {
	nh := int(ch.GetHp()) + hpDelta
	if nh < 0 {
		nh = 0
	}
	maxHp := int(ch.GetMaxHp())
	if nh > maxHp {
		nh = maxHp
	}
	nm := int(ch.GetMp()) + mpDelta
	if nm < 0 {
		nm = 0
	}
	maxMp := int(ch.GetMaxMp())
	if nm > maxMp {
		nm = maxMp
	}
	ch.SetHp(uint32(nh), false)
	ch.SetMp(uint32(nm), false)
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
		constant.StatMP: int32(ch.GetMp()),
	}, false)
}

func (ch *Character) SetBaseHp(v uint32, notify bool) {
	if v > constant.StatMaxHPMP {
		v = constant.StatMaxHPMP
	}
	ch.BaseHp = v
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetBaseMp(v uint32, notify bool) {
	if v > constant.StatMaxHPMP {
		v = constant.StatMaxHPMP
	}
	ch.BaseMp = v
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) AddBaseHp(amount uint32, notify bool) {
	ch.LifeCore.AddBaseHp(amount)
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) AddBaseMp(amount uint32, notify bool) {
	ch.LifeCore.AddBaseMp(amount)
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetAbilityPoint(v uint16, notify bool) {
	ch.AbilityPoint = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatAvailableAP: int32(ch.AbilityPoint),
		}, false)
	}
}

func (ch *Character) SetSkillPoint(v uint16, notify bool) {
	ch.SkillPoint = v
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatAvailableSP: int32(ch.SkillPoint),
		}, false)
	}
}

func (ch *Character) Warp(targetMap *Map, spawnPoint uint8) error {
	if ch.GameWorld == nil {
		return fmt.Errorf("no game world")
	}
	return ch.GameWorld.GetMapSystem().Warp(ch, targetMap, spawnPoint)
}

func (ch *Character) Relocate(spawnPoint uint8) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}
	if m.Wz == nil {
		return fmt.Errorf("map model not found")
	}
	pos, ok := m.Wz.GetSpawnPosition(spawnPoint)
	if !ok {
		return fmt.Errorf("invalid spawn point %d", spawnPoint)
	}
	beforePosition := ch.Position
	ch.Position = pos
	ch.Stance = constant.StanceDefaultValue
	for _, summon := range ch.GetSummons() {
		if summon == nil || summon.Owner != ch {
			continue
		}
		summon.Position = ch.Position
	}
	if ch.Listener != nil {
		ch.Listener.OnFieldRelocate(ch, spawnPoint)
		ch.Listener.OnPlayerMove(ch, beforePosition, []dto.MoveFragment{
			&dto.TeleportMovement{
				BasicMovement: &dto.BasicMovement{
					Command: 3,
					Stance:  ch.Stance,
				},
				Position: ch.Position,
				Velocity: types.Vector2[int16]{},
			},
		})
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) error {
	if ch.Sendable == nil {
		return nil
	}
	return ch.Sendable.Send(p, policy)
}

func (ch *Character) IsHidden() bool {
	return ch.hidden
}

func (ch *Character) SetHidden(hidden bool) {
	if ch.hidden == hidden {
		return
	}
	ch.hidden = hidden
	ch.Listener.OnHiddenChanged(ch, hidden)
}

func (ch *Character) GetID() uint32 {
	return ch.id
}

func (ch *Character) GetPK() uint32 {
	return ch.id
}

func (ch *Character) GetRole() constant.CharacterRole {
	return ch.Role
}

func (ch *Character) GetName() string {
	return ch.name
}

func (ch *Character) ToProtoPartyMember(worldID uint32, channelID int32, role internal.PartyMemberRole) *internal.PartyMember {
	if m := PartyMemberFromCharacter(ch, worldID, channelID, role); m != nil {
		return m.ToProto()
	}
	return nil
}

func (ch *Character) ToProtoGuildMember(worldID uint32, channelID int32, rank internal.GuildMemberRank) *internal.GuildMember {
	if m := GuildMemberFromCharacter(ch, worldID, channelID, rank); m != nil {
		return m.ToProto()
	}
	return nil
}

func (ch *Character) GetPartyID() *uint32 {
	if ch == nil || ch.partyID == nil {
		return nil
	}
	p := new(uint32)
	*p = *ch.partyID
	return p
}

func (ch *Character) SetPartyID(partyID *uint32) {
	ch.partyID = partyID
	if ch.doors == nil {
		return
	}
	for _, d := range ch.doors {
		if d != nil {
			d.PartyID = partyID
		}
	}
}

func (ch *Character) BuddyList() *BuddyList {
	if ch == nil {
		return nil
	}
	if ch.buddyList == nil {
		ch.buddyList = NewBuddyList()
	}
	return ch.buddyList
}

func (ch *Character) LoadBuddyList(entries []*internal.BuddyEntry, capacity uint32) {
	ch.BuddyList().LoadFromProto(entries, capacity)
}

func (ch *Character) SendBuddyLoginSync() {
	if ch == nil || ch.Listener == nil {
		return
	}
	entries := ch.BuddyList().SnapshotForClient()
	ch.Listener.OnBuddyListUpdate(ch, pconst.BuddyListSyncLogin, entries)
}

func (ch *Character) GetGuildID() (uint32, bool) {
	if ch == nil || ch.guildID == nil {
		return 0, false
	}
	return *ch.guildID, true
}

func (ch *Character) SetGuildID(guildID *uint32) {
	ch.guildID = guildID
}

func (ch *Character) Message(message string) {
	ch.Listener.OnMessage(ch, constant.MsgLightBlueText, message)
}

func (ch *Character) HasRoleAtLeast(role constant.CharacterRole) bool {
	return ch.Role >= role
}

func (ch *Character) SetMeso(meso int32) {
	if meso < 0 {
		ch.Meso = 0
		return
	}
	ch.Meso = meso
}

func (ch *Character) IsRanked() bool {
	if ch.HasRoleAtLeast(constant.RoleAdmin) {
		return false
	}

	if ch.level < 30 {
		return false
	}

	return true
}

func (ch *Character) remainingExpToMaxLevel() uint32 {
	if ch == nil || ch.level >= 200 {
		return 0
	}
	if ch.GameWorld == nil {
		return ^uint32(0)
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return ^uint32(0)
	}

	var cap uint32
	if needed := resources.GetExpNeededForLevel(ch.level); needed > ch.exp {
		cap = needed - ch.exp
	}
	for lvl := ch.level + 1; lvl < 200; lvl++ {
		need := resources.GetExpNeededForLevel(lvl)
		if cap > ^uint32(0)-need {
			return ^uint32(0)
		}
		cap += need
	}
	return cap
}

func (ch *Character) AddExp(exp uint32) {
	exp = min(exp, ^uint32(0)-ch.exp)
	exp = min(exp, ch.remainingExpToMaxLevel())
	if exp == 0 {
		return
	}
	ch.exp += exp
	ch.Listener.OnExpGain(ch, exp)

	if !ch.tryLevelUp() {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatEXP: int32(ch.exp),
		}, false)
	}
}

func (ch *Character) ComputeMobKillExp(raw uint32) uint32 {
	if ch == nil || raw == 0 {
		return 0
	}
	mulClamp := func(a, b uint32) uint32 {
		if a == 0 || b == 0 {
			return 0
		}
		if a > ^uint32(0)/b {
			return ^uint32(0)
		}
		return a * b
	}

	exp := raw
	if ch.BonusStats.ExpRate > 0 {
		exp = mulClamp(exp, uint32(ch.BonusStats.ExpRate)) / 100
	}
	exp = mulClamp(exp, uint32(ch.GetHolySymbolExpRate())) / 100
	if ch.HasDebuff(constant.DebuffFlagCurse) {
		exp /= 2
	}
	if gw := ch.GameWorld; gw != nil {
		if r := gw.GetExpRate(); r > 0 {
			exp = mulClamp(exp, uint32(r))
		}
	}
	return exp
}

func (ch *Character) GetHolySymbolExpRate() int32 {
	if ch == nil || ch.Buffs == nil {
		return 100
	}
	ent := ch.Buffs.GetEntity(constant.BuffFlagHolySymbol)
	sb, typeOk := ent.(*SkillBuff)
	if !typeOk || sb == nil || sb.Wz == nil {
		return 100
	}
	x, valOk := sb.Values[constant.BuffFlagHolySymbol]
	if !valOk || x <= 0 {
		return 100
	}
	full := sb.Wz.ID == uint32(constant.SkillGmHolySymbol)
	if !full {
		pid := ch.GetPartyID()
		m := ch.GetMap()
		full = pid != nil && m != nil && len(m.GetPartyMembers(*pid)) >= 2
	}
	if full {
		return 100 + x
	}
	return int32(100 * int64(150+x) / 150)
}

type CharacterInitData struct {
	ID           uint32
	AccountID    uint32
	Name         string
	Gender       uint8
	SkinColor    uint8
	Face         uint32
	Hair         uint32
	Level        uint8
	Class        uint16
	Role         uint8
	Str          uint16
	Dex          uint16
	Int          uint16
	Luk          uint16
	Hp           uint32
	MaxHp        uint32
	Mp           uint32
	MaxMp        uint32
	AbilityPoint uint16
	SkillPoint   uint16
	Exp          uint32
	Meso         int32
	PositionX    int16
	PositionY    int16
	Stance       uint8
	Hidden       bool
	PartyID      *uint32
	GuildID      *uint32
}

func NewCharacter(sender Sendable, listener CharacterListener, data *CharacterInitData, gw GameWorld) *Character {
	if listener == nil {
		panic("NewCharacter: listener must not be nil")
	}
	ch := &Character{
		Sendable: sender,
		Listener: listener,
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				GameWorld: gw,
				Position: types.Vector2[int16]{
					X: data.PositionX,
					Y: data.PositionY,
				},
			},
			BaseHp: data.MaxHp,
			BaseMp: data.MaxMp,
			Stance: data.Stance,
		},
		id:           data.ID,
		AccountID:    data.AccountID,
		name:         data.Name,
		gender:       data.Gender,
		skinColor:    data.SkinColor,
		face:         data.Face,
		hair:         data.Hair,
		level:        data.Level,
		Class:        data.Class,
		Role:         constant.CharacterRole(data.Role),
		hidden:       data.Hidden,
		BaseStats:    BaseStats{Str: data.Str, Dex: data.Dex, Int: data.Int, Luk: data.Luk},
		AbilityPoint: data.AbilityPoint,
		SkillPoint:   data.SkillPoint,
		exp:          data.Exp,
		Meso:         data.Meso,
		partyID:      data.PartyID,
		guildID:      data.GuildID,
		GuildInvites: make(map[uint32]time.Time),
		buddyList:    NewBuddyList(),

		random1: stream.NewRandomStream(),
		random2: stream.NewRandomStream(),
		random3: stream.NewRandomStream(),

		Inventory: map[constant.InventoryType]*Inventory{
			constant.InventoryTypeEquipment:    NewInventory(constant.InventoryTypeEquipment),
			constant.InventoryTypeConsume:      NewInventory(constant.InventoryTypeConsume),
			constant.InventoryTypeInstallation: NewInventory(constant.InventoryTypeInstallation),
			constant.InventoryTypeETC:          NewInventory(constant.InventoryTypeETC),
			constant.InventoryTypeCash:         NewInventory(constant.InventoryTypeCash),
		},
		Rings: RingContainer{
			Left:  []*Ring{},
			Right: []*Ring{},
			Mid:   []*Ring{},
		},
		Equipments: map[constant.EquipmentPartsType]Equipment{},

		regRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}
	ch.Buffs = NewBuffContainer(ch)
	ch.Skills = NewSkillContainer(ch)
	ch.keyLayout = NewKeyLayout()
	ch.LifeCore.ObjectCore.self = ch
	ch.LifeCore.ObjectCore.initTimers()
	ch.LifeCore.setHp(data.Hp)
	ch.LifeCore.setMp(data.Mp)
	return ch
}

func (ch *Character) KeyLayout() *KeyLayout {
	if ch == nil {
		return nil
	}
	return ch.keyLayout
}

func (ch *Character) GetDialog() *lua.LState {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	return ch.luaDialog
}

func (ch *Character) SetDialog(lua *lua.LState) {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.luaDialog = lua
}

func (ch *Character) ClearCurrentDialog() {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.luaDialog = nil
}

func (ch *Character) tryLevelUp() bool {
	if ch.GameWorld == nil {
		return false
	}

	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return false
	}

	if ch.level >= 200 {
		return false
	}

	oldLevel := ch.level
	remainingExp := ch.exp
	targetLevel := ch.level

	for targetLevel < 200 {
		expNeeded := resources.GetExpNeededForLevel(targetLevel)
		if expNeeded == 0 {
			break
		}

		if remainingExp < expNeeded {
			break
		}

		remainingExp -= expNeeded
		targetLevel++
	}

	if targetLevel == oldLevel {
		return false
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}

	root := mapInstance.GetLuaRoot()
	if root == nil {
		return false
	}

	ch.exp = remainingExp
	ch.SetLevel(targetLevel)
	levelDiff := int(targetLevel) - int(oldLevel)
	if levelDiff >= 1 {
		if thread, err := luax.NewThread(root, constant.CharacterHookScriptPath); err == nil {
			_, _ = luax.Call(thread, "on_level_up", ch, int32(oldLevel), int32(targetLevel))
		}

		stats := map[constant.Stat]int32{
			constant.StatLevel:       int32(ch.level),
			constant.StatEXP:         int32(ch.exp),
			constant.StatMaxHP:       int32(ch.GetMaxHp()),
			constant.StatMaxMP:       int32(ch.GetMaxMp()),
			constant.StatHP:          int32(ch.GetHp()),
			constant.StatMP:          int32(ch.GetMp()),
			constant.StatAvailableAP: int32(ch.AbilityPoint),
			constant.StatAvailableSP: int32(ch.SkillPoint),
		}
		ch.Listener.OnUpdateStats(ch, stats, false)
		for i := 0; i < levelDiff; i++ {
			ch.broadcastLevelUpEffect()
		}
	}
	return true
}

func (ch *Character) SetLevel(newLevel uint8) {
	if newLevel < 1 {
		newLevel = 1
	}
	if newLevel > 200 {
		newLevel = 200
	}

	if newLevel == ch.level {
		return
	}

	oldLevel := ch.level
	ch.level = newLevel
	ch.Listener.OnPartyMemberFieldsChanged(ch)
	if newLevel < oldLevel {
		if ch.GameWorld != nil {
			resources := ch.GameWorld.GetResources()
			if resources != nil {
				if newLevel > 1 {
					ch.exp = resources.GetExpNeededForLevel(newLevel - 1)
				} else {
					ch.exp = 0
				}
			}
		}
	}

	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatLevel: int32(ch.level),
		constant.StatEXP:   int32(ch.exp),
	}, false)
}

func (ch *Character) broadcastLevelUpEffect() {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	ch.Broadcast(&response.ShowEffect{
		CharacterID: ch.id,
		Type:        response.EffectTypeLevelUp,
	}, nil)
}

func debuffTimerKey(flag constant.DebuffFlag) string {
	return fmt.Sprintf("debuff_%d_%d", flag.Position, flag.Mask)
}

func (ch *Character) HasDebuff(flag constant.DebuffFlag) bool {
	if ch.debuffs == nil {
		return false
	}
	_, ok := ch.debuffs[flag]
	return ok
}

func (ch *Character) AddDebuff(holder *Debuff) {
	if holder == nil {
		return
	}
	if ch.debuffs == nil {
		ch.debuffs = make(map[constant.DebuffFlag]*Debuff)
	}
	ch.RemoveTimer(debuffTimerKey(holder.Flag))
	ch.debuffs[holder.Flag] = holder
	if holder.Duration > 0 {
		flag := holder.Flag
		ch.AddTimer(debuffTimerKey(flag), holder.Duration, false, func() {
			ch.RemoveTimer(debuffTimerKey(flag))
			if _, ok := ch.debuffs[flag]; ok {
				delete(ch.debuffs, flag)
				ch.Listener.OnDebuffRemoved(ch, []constant.DebuffFlag{flag})
			}
		})
	}
}

func (ch *Character) GiveDebuff(flag constant.DebuffFlag, duration time.Duration, x int16, skillID uint16, skillLevel uint16) {
	if skillID == 0 {
		skillID = flag.DiseaseSkillID
	}
	if skillLevel == 0 {
		skillLevel = 1
	}
	holder := &Debuff{
		Flag:      flag,
		StartTime: time.Now(),
		Duration:  duration,
	}
	ch.AddDebuff(holder)
	ch.Listener.OnDebuffAdded(ch, flag, x, skillID, skillLevel, int32(duration.Milliseconds()))
}

func (ch *Character) RemoveDebuff(flags ...constant.DebuffFlag) {
	var removed []constant.DebuffFlag
	if ch.debuffs != nil {
		for _, flag := range flags {
			ch.RemoveTimer(debuffTimerKey(flag))
			if _, ok := ch.debuffs[flag]; ok {
				delete(ch.debuffs, flag)
				removed = append(removed, flag)
			}
		}
	}
	if len(removed) > 0 {
		ch.Listener.OnDebuffRemoved(ch, removed)
	}
}

func (ch *Character) GetDiseaseMask() [4]uint32 {
	var mask [4]uint32
	if ch.debuffs == nil {
		return mask
	}
	for flag := range ch.debuffs {
		idx := flag.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= flag.Mask
		}
	}
	return mask
}

type SpawnPlayerBuffData struct {
	BuffStates        [4]uint32
	SpeedBuff         uint8
	ComboCount        uint8
	WKChargeSkillID   uint32
	MorphID           uint16
	SpiritClawSkillID uint32
	MountLevel        uint32
	MountExp          uint32
	MountFatigue      uint32
}

func (ch *Character) GetRiddingInfo() (mountID int32, active bool) {
	if ch == nil || ch.Buffs == nil {
		return 0, false
	}
	_, mountID, active = ch.Buffs.GetBuffValue(constant.BuffFlagMonsterRiding)
	return
}

func (ch *Character) GetSpawnPlayerBuffData() SpawnPlayerBuffData {
	data := SpawnPlayerBuffData{
		SpeedBuff:  1,
		ComboCount: 1,
		MountLevel: 1,
	}
	if ch == nil || ch.Buffs == nil {
		return data
	}

	for _, buff := range ch.Buffs.Entities() {
		if buff == nil {
			continue
		}
		values := buff.GetValues()
		if values == nil {
			continue
		}
		for flag, value := range values {
			idx := flag.Position - 1
			if idx < 0 || idx >= constant.MaxBuffFlag {
				continue
			}
			if constant.IsTargetStatFlag(flag) {
				data.BuffStates[idx] |= flag.Mask
			}
			switch flag {
			case constant.BuffFlagSpeed:
				data.SpeedBuff = uint8(value)
			case constant.BuffFlagCombo:
				data.ComboCount = uint8(value)
			case constant.BuffFlagMorph:
				data.MorphID = uint16(value)
			}
		}
		switch typed := buff.(type) {
		case *SkillBuff:
			for _, flag := range typed.GetFlags() {
				switch flag {
				case constant.BuffFlagWkCharge:
					data.WKChargeSkillID = typed.Wz.ID
				case constant.BuffFlagSpiritClaw:
					data.SpiritClawSkillID = typed.Wz.ID
				}
			}
		}
	}

	return data
}
