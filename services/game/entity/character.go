package entity

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

type Look struct {
	Gender    uint8
	SkinColor uint8
	Face      uint32
	Hair      uint32
}

type Character struct {
	LifeCore
	Sendable

	id              uint32
	name            string
	look            Look
	level           uint8
	Rank            Ranking
	exp             uint32
	Stats           *Stats
	Points          *Points
	mega            bool
	random          [3]stream.RandomStream
	Quests          *Quests
	Wedding         *Wedding
	CashWishlist    []uint32
	TeleportStones  *TeleportStones
	MonsterBook     *MonsterBook
	Dialog          *Dialog
	Listener        CharacterListener
	Class           uint16
	Role            constant.CharacterRole
	AccountID       uint32
	Inventory       *Inventory
	Storage         *Storage
	Duey            *Duey
	StoreBank       *StoreBank
	MiniRoom        *EntrustedShop
	miniRoomPending bool
	Skills          *Skills
	KeyLayout       *KeyLayout
	Chair           uint32
	lastHeal        lastHeal
	Buffs           *Buffs
	Debuffs         *Debuffs
	Summons         *Summons
	Pets            *Pets
	Doors           *Doors
	HomingTargetOID *uint32
	Party           *PartyMembership
	Guild           *GuildMembership
	Buddies         *BuddyList
	GM              GMMode
	stateMachine    *StateMachine
	carnivalTeam    *CarnivalTeam
	dojoEnergy      int
	SavedLocations  *SavedLocations
	Records         *Records
	AccountRecords  *Records
	session         session
}

type session struct {
	loggedOut     atomic.Bool
	logoutEntry   chan *internal.CharacterSaveEntry
	destination   atomic.Pointer[Map]
	cashItemInUse atomic.Bool
}

func (ch *Character) Destination() *Map {
	return ch.session.destination.Load()
}

func (ch *Character) BeginMove(target *Map) bool {
	return ch.session.destination.CompareAndSwap(nil, target)
}

func (ch *Character) FinishMove(target *Map) {
	ch.session.destination.CompareAndSwap(target, nil)
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
	carnivalTeam := constant.CarnivalTeamNone
	if team := ch.CarnivalTeam(); team != nil {
		carnivalTeam = team.TeamID
	}
	spawnPacket := &response.SpawnPlayer{
		Character:         ch.ToDTO(),
		BuffStates:        spawnBuffData.BuffStates,
		Diseases:          ch.Debuffs.DiseaseMask(),
		SpeedBuff:         spawnBuffData.SpeedBuff,
		ComboCount:        spawnBuffData.ComboCount,
		WKChargeSkillId:   spawnBuffData.WKChargeSkillID,
		MorphId:           spawnBuffData.MorphID,
		SpiritClawSkillId: spawnBuffData.SpiritClawSkillID,
		MountLevel:        spawnBuffData.MountLevel,
		MountExp:          spawnBuffData.MountExp,
		MountFatigue:      spawnBuffData.MountFatigue,
		CrushRing:         ch.Inventory.WornRing(ch.Inventory.Rings.Left),
		FriendshipRing:    ch.Inventory.WornRing(ch.Inventory.Rings.Mid),
		MarriageRing:      ch.Wedding.RingToDTO(),
		HasTeam:           ch.GetMap() != nil && ch.GetMap().Wz.HasTeam(),
		Team:              carnivalTeam,
	}
	if ch.Pets.Active != nil {
		spawnPacket.Pet = ch.Pets.Active.ToDTO()
	}
	if guildID, ok := ch.Guild.ID(); ok && ch.GameWorld != nil {
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
			Duration:    dashBuff.RemainingDuration(clock.Now()),
			Buffs:       buffs,
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (ch *Character) SendDestroySyncToViewer(viewer *Character) {
	if ch == nil || viewer == nil {
		return
	}
	if ch.GetID() == viewer.GetID() {
		return
	}
	if ch.IsHidden() && !viewer.HasRoleAtLeast(ch.Role) {
		return
	}
	viewer.Send(&response.LeavePlayer{ID: ch.GetID()}, types.SEND_POLICY_ENCRYPT)
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
		mist.NextPoisonTickAt = clock.Now().Add(initialDelay)
	}
	if duration > 0 {
		mist.ExpiresAt = clock.Now().Add(duration)
	}
	m.AddMist(mist)
	return mist
}

func (ch *Character) RemoveMist(mist *Mist) {
	if mist == nil {
		return
	}
	if m := mist.GetMap(); m != nil && mist.OID != 0 {
		m.RemoveMist(mist.OID)
	}
}

func (ch *Character) GetBonusHp() int32   { return ch.BonusHp }
func (ch *Character) GetBonusMp() int32   { return ch.BonusMp }
func (ch *Character) GetInvincible() bool { return ch.Invincible }
func (ch *Character) GetGender() uint8    { return ch.look.Gender }
func (ch *Character) GetSkinColor() uint8 { return ch.look.SkinColor }
func (ch *Character) GetFace() uint32     { return ch.look.Face }
func (ch *Character) GetHair() uint32     { return ch.look.Hair }

func (ch *Character) SetGender(gender uint8) {
	if ch.look.Gender == gender {
		return
	}
	ch.look.Gender = gender
	ch.Listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetHair(id uint32) {
	if ch.look.Hair == id {
		return
	}
	ch.look.Hair = id
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHair: int32(id),
	}, false)
	ch.Listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetFace(id uint32) {
	if ch.look.Face == id {
		return
	}
	ch.look.Face = id
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatFace: int32(id),
	}, false)
	ch.Listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetSkin(color uint8) {
	if ch.look.SkinColor == color {
		return
	}
	ch.look.SkinColor = color
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatSkin: int32(color),
	}, false)
	ch.Listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) MapMessage(messageType constant.ServerMessageType, message string) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	for _, player := range mapInstance.GetAllPlayers() {
		other, ok := player.(*Character)
		if !ok {
			continue
		}
		other.Listener.OnMessage(other, messageType, message)
	}
}

func (ch *Character) GetLevel() uint8 { return ch.level }
func (ch *Character) GetExp() uint32  { return ch.exp }
func (ch *Character) GetSpawnPoint() uint8 {
	mapInstance := ch.GetMap()
	if mapInstance != nil && mapInstance.Wz != nil {
		return mapInstance.Wz.FindClosestPortalSpawnID(ch.Position)
	}
	return 0
}

func (ch *Character) SetHp(v uint32, notify bool) {
	wasAlive := ch.IsAlive()
	ch.LifeCore.setHp(v)
	if notify {
		ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatHP: int32(ch.GetHp())}, false)
	}

	if wasAlive && ch.IsAlive() == false {
		ch.die()
	}
}

func (ch *Character) die() {
	ch.Summons.Clear()
	if sm := ch.StateMachine(); sm != nil {
		sm.CallHook("on_player_dead", ch)
	}
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
	ch.Stats.Bonus.MaxHpPercent = p
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
	ch.Stats.Bonus.MaxMpPercent = p
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

func (ch *Character) updateMaxHpMp() {
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP:    int32(ch.GetHp()),
		constant.StatMaxHP: int32(ch.GetMaxHp()),
		constant.StatMP:    int32(ch.GetMp()),
		constant.StatMaxMP: int32(ch.GetMaxMp()),
	}, false)
}

func (ch *Character) SetInvincible(b bool) { ch.Invincible = b }

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

const clockTimer = "clock"

func (ch *Character) Warp(ctx actor.Context, targetMap *Map, spawnPoint uint8) error {
	if ch.GameWorld == nil {
		return fmt.Errorf("no game world")
	}
	return ch.GameWorld.GetMapSystem().Warp(ctx, ch, targetMap, spawnPoint, nil)
}

func (ch *Character) EnterPortal(ctx actor.Context, portal *Portal) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}

	if scriptName := portal.Script(); scriptName != "" {
		root := m.GetLuaRoot()
		if root == nil {
			ch.Listener.OnUnlockAction(ch)
			return fmt.Errorf("root lua state not found")
		}
		scriptPath := fmt.Sprintf("script/portal/%s.lua", scriptName)
		thread, err := luax.NewThread(root, scriptPath)
		if err != nil {
			ch.Listener.OnScriptError(ch, scriptPath, err)
			ch.Listener.OnUnlockAction(ch)
			return fmt.Errorf("portal script thread: %w", err)
		}
		luax.SetConfiguration(thread, luax.Configuration{
			ActorContext: ctx,
			ActorPID:     m.LogicActorPID(),
		})
		luax.CallAsync(ctx, root, thread, "on_enter", ch, portal).Do(func([]lua.LValue) error {
			if ch.Dialog.Thread() == nil {
				ch.Listener.OnUnlockAction(ch)
			}
			return nil
		}).OnError(func(err error) {
			log.Printf("portal script %s failed: %v", scriptPath, err)
			ch.Listener.OnScriptError(ch, scriptPath, err)
			if ch.Dialog.Thread() == nil {
				ch.Listener.OnUnlockAction(ch)
			}
		})
		return nil
	}

	targetMap := ch.GameWorld.GetMapSystem().Find(m.StateMachine(), uint32(portal.Wz.TargetMapId))
	if targetMap == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	targetPortal := targetMap.FindPortalByName(portal.Wz.Target)
	if targetPortal == nil || targetPortal.Wz == nil {
		return ch.Warp(ctx, targetMap, 0)
	}
	return ch.Warp(ctx, targetMap, targetPortal.Wz.ID)
}

func (ch *Character) Revive(ctx actor.Context) error {
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return fmt.Errorf("not on a map")
	}
	if ch.GetHp() > 0 {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	targetMap := ch.GameWorld.GetMapSystem().Find(m.StateMachine(), uint32(m.Wz.ReturnMapId))
	if targetMap == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return fmt.Errorf("return map %d not found", m.Wz.ReturnMapId)
	}

	ch.SetHp(50, false)
	ch.Stance = constant.StanceDefaultValue
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, true)
	if sm := ch.StateMachine(); sm != nil {
		sm.CallHook("on_player_revive", ch)
	}
	return ch.Warp(ctx, targetMap, 0)
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
	ch.Stance = constant.StanceDefaultValue
	beforePosition := ch.moveTo(m, pos)
	if ch.Listener != nil {
		ch.Listener.OnFieldRelocate(ch, spawnPoint)
		ch.Listener.OnPlayerMove(ch, beforePosition, []dto.MoveFragment{
			&dto.TeleportMovement{
				BasicMovement: &dto.BasicMovement{
					Command: 3,
					Stance:  ch.Stance,
				},
				Position: ch.Position,
			},
		})
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (ch *Character) EnterInnerPortal(portalName string, to types.Vector2[int16]) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}
	portal := m.FindPortalByName(portalName)
	if portal == nil || portal.Wz == nil {
		return fmt.Errorf("portal %q not found", portalName)
	}
	if portal.Wz.Position.DistanceSq(ch.Position) > 22500 && ch.HasRoleAtLeast(constant.RoleAdmin) == false {
		return fmt.Errorf("portal %q is too far", portalName)
	}
	ch.moveTo(m, to)
	return nil
}

func (ch *Character) moveTo(m *Map, pos types.Vector2[int16]) types.Vector2[int16] {
	before := ch.Position
	ch.Position = pos
	for _, summon := range ch.Summons.All() {
		if summon == nil || summon.Owner != ch {
			continue
		}
		summonBefore := summon.Position
		summon.Position = pos
		m.OnMoved(summon, summonBefore)
	}
	m.OnMoved(ch, before)
	return before
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) error {
	if ch.Sendable == nil {
		return nil
	}
	return ch.Sendable.Send(p, policy)
}

func (ch *Character) IsHidden() bool {
	return ch.GM.Hidden
}

func (ch *Character) SetHidden(hidden bool) {
	if ch.GM.Hidden == hidden {
		return
	}
	ch.GM.Hidden = hidden
	ch.Listener.OnHiddenChanged(ch, hidden)
}

func (ch *Character) GetID() uint32 {
	return ch.id
}

// MarkLoggedOut runs on the disconnect goroutine; whichever actor holds the character next logs it out.
func (ch *Character) MarkLoggedOut() <-chan *internal.CharacterSaveEntry {
	ch.session.logoutEntry = make(chan *internal.CharacterSaveEntry, 1)
	ch.session.loggedOut.Store(true)
	return ch.session.logoutEntry
}

func (ch *Character) LoggedOut() bool {
	return ch.session.loggedOut.Load()
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

func (ch *Character) Message(message string) {
	ch.Listener.OnMessage(ch, constant.MsgLightBlueText, message)
}

func (ch *Character) HasRoleAtLeast(role constant.CharacterRole) bool {
	return ch.Role >= role
}

func (inv *Inventory) SetMeso(meso int32) {
	if meso < 0 {
		inv.Meso = 0
		return
	}
	inv.Meso = meso
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
	if exp == 0 {
		return
	}
	if ch.validateExpExchange(0, exp) != ExchangeOK {
		return
	}
	ch.addExpUnchecked(exp)
}

func (ch *Character) addExpUnchecked(exp uint32) {
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
	if ch.Stats.Bonus.ExpRate > 0 {
		exp = mulClamp(exp, uint32(ch.Stats.Bonus.ExpRate)) / 100
	}
	exp = mulClamp(exp, uint32(ch.GetHolySymbolExpRate())) / 100
	if ch.Debuffs.Has(constant.DebuffFlagCurse) {
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
	if ch == nil {
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
		pid := ch.Party.ID()
		m := ch.GetMap()
		full = pid != nil && m != nil && len(m.GetPartyMembers(*pid)) >= 2
	}
	if full {
		return 100 + x
	}
	return int32(100 * int64(150+x) / 150)
}

func (ch *Character) BindKey(slot int, typ byte, action int32) {
	ch.KeyLayout.SetKey(slot, typ, action)
	ch.Listener.OnKeyMap(ch)
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
			luax.CallAsync(nil, root, thread, "on_level_up", ch, int32(oldLevel), int32(targetLevel))
		}

		stats := map[constant.Stat]int32{
			constant.StatLevel:       int32(ch.level),
			constant.StatEXP:         int32(ch.exp),
			constant.StatMaxHP:       int32(ch.GetMaxHp()),
			constant.StatMaxMP:       int32(ch.GetMaxMp()),
			constant.StatHP:          int32(ch.GetHp()),
			constant.StatMP:          int32(ch.GetMp()),
			constant.StatAvailableAP: int32(ch.Points.AP),
			constant.StatAvailableSP: int32(ch.Points.SP),
		}
		ch.Listener.OnUpdateStats(ch, stats, false)
		for i := 0; i < levelDiff; i++ {
			ch.broadcastLevelUpEffect()
		}
		ch.Quests.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
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
	if ch == nil {
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
	if ch == nil {
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
