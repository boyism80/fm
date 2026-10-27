package game

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

const (
	spawnBalloonItemID        = 5110000
	spawnBalloonMax    uint16 = 250
)

type Look struct {
	gender    uint8
	skinColor uint8
	face      uint32
	hair      uint32
}

type sendable = Sendable

type Character struct {
	lifeCore
	sendable

	id             uint32
	name           string
	look           Look
	level          uint8
	exp            uint32
	expReachedAt   time.Time
	stats          *Stats
	points         *Points
	mega           bool
	random         [3]stream.RandomStream
	quests         *Quests
	wedding        *Wedding
	cashWishlist   []uint32
	teleportStones *TeleportStones
	monsterBook    *MonsterBook
	dialog         *Dialog
	listener       CharacterListener
	class          constant.ClassType
	role           constant.CharacterRole
	accountID      uint32
	account        *Account
	inventory      *Inventory
	escrow         *Escrow
	duey           *Duey
	unsettled      *UnsettledRequests
	storeBank      *StoreBank
	miniRoom       MiniRoom
	remoteShops    *RemoteShops
	skills         *Skills
	keyLayout      *KeyLayout
	itemEffectID   uint32
	chalkboard     string
	chair          uint32
	// fieldTeam is the coconut team the character plays for on its current map.
	fieldTeam constant.CarnivalTeam
	// The client puts the key of the last field it was sent into its moves, portals, attacks and loots.
	fieldKey uint8
	// The map's consume cooldown is shared by the use window and the pet's auto potion.
	nextConsumeAt  time.Time
	lastHeal       lastHeal
	buffs          *Buffs
	cheat          *CheatTracker
	debuffs        *Debuffs
	summons        *Summons
	pets           *Pets
	mount          *Mount
	doors          *Doors
	party          *PartyMembership
	guild          *GuildMembership
	buddies        *BuddyList
	family         *Family
	messenger      *Messenger
	notes          *Notes
	gm             GMMode
	stateMachine   *StateMachine
	carnival       *CarnivalMembership
	dojo           *Dojo
	lieDetector    *LieDetector
	savedLocations *SavedLocations
	records        *Records
	session        session
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
	if ch.IsHidden() && !viewer.HasRoleAtLeast(ch.role) {
		return
	}
	spawnBuffData := ch.buffs.SpawnPlayerData()
	diseaseReasons, poisonDamage := ch.debuffs.SpawnReasons()
	carnivalTeam := constant.CarnivalTeamNone
	if team := ch.carnival.Team(); team != nil {
		carnivalTeam = team.teamID
	} else if ch.fieldTeam != constant.CarnivalTeamNone {
		carnivalTeam = ch.fieldTeam
	}
	spawnPacket := &response.SpawnPlayer{
		Character:         ch.BuildS2C(),
		BuffStates:        spawnBuffData.BuffStates,
		Diseases:          ch.debuffs.DiseaseMask(),
		DiseaseReasons:    diseaseReasons,
		PoisonDamage:      poisonDamage,
		Balloons:          uint32(min(ch.inventory.GetCountByItemID(spawnBalloonItemID), spawnBalloonMax)),
		SpeedBuff:         spawnBuffData.SpeedBuff,
		ComboCount:        spawnBuffData.ComboCount,
		WKChargeSkillId:   spawnBuffData.WKChargeSkillID,
		MorphId:           spawnBuffData.MorphID,
		SpiritClawSkillId: spawnBuffData.SpiritClawSkillID,
		ItemEffectId:      ch.itemEffectID,
		ChairId:           ch.chair,
		Chalkboard:        ch.chalkboard,
		MountLevel:        ch.mount.level,
		MountExp:          ch.mount.exp,
		MountFatigue:      ch.mount.fatigue,
		CrushRing:         ch.inventory.WornRing(ch.inventory.rings.left),
		FriendshipRing:    ch.inventory.WornRing(ch.inventory.rings.mid),
		MarriageRing:      ch.wedding.BuildRingS2C(),
		HasTeam:           ch.GetMap() != nil && ch.GetMap().wz.HasTeam(),
		Team:              carnivalTeam,
	}
	switch room := ch.miniRoom.(type) {
	case *PersonalShop:
		if room.published && room.owner == ch {
			balloon := room.balloon()
			spawnPacket.MiniRoomBalloon = &balloon
		}
	case *MiniGame:
		if room.owner == ch {
			balloon := room.balloon()
			spawnPacket.MiniRoomBalloon = &balloon
		}
	}
	if ch.pets.active != nil {
		spawnPacket.Pet = ch.pets.active.BuildS2C()
	}
	if guildID, ok := ch.guild.ID(); ok && ch.gameWorld != nil {
		if guild := ch.gameWorld.Systems().Get[*GuildSystem]().Get(guildID); guild != nil {
			spawnPacket.GuildName = guild.name
			if guild.logo != nil {
				spawnPacket.GuildLogoBG = uint16(guild.logo.logoBG)
				spawnPacket.GuildLogoBGColor = uint8(guild.logo.logoBGColor)
				spawnPacket.GuildLogo = uint16(guild.logo.logo)
				spawnPacket.GuildLogoColor = uint8(guild.logo.logoColor)
			}
		}
	}
	viewer.Send(spawnPacket, types.SEND_POLICY_ENCRYPT)

	if buff, mountID, active := ch.buffs.GetBuffValue(constant.BuffFlagMonsterRiding); active {
		viewer.Send(&response.UpdateRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			SkillID:     buff.GetBuffID(),
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagMonsterRiding, Value: mountID}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if buff, _, active := ch.buffs.GetBuffValue(constant.BuffFlagEnergyCharge); active {
		viewer.Send(&response.UpdateBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buff.GetBuffID(),
			Duration:    50 * time.Second,
			Buffs:       []dto.BuffEntry{{Buff: constant.BuffFlagEnergyCharge, Value: 10000}},
		}, types.SEND_POLICY_ENCRYPT)
	}

	if dashBuff, dashValue, active := ch.buffs.GetBuffValue(constant.BuffFlagDashSpeed); active {
		buffs := []dto.BuffEntry{{Buff: constant.BuffFlagDashSpeed, Value: dashValue}}
		if _, dashJumpValue, hasDashJump := ch.buffs.GetBuffValue(constant.BuffFlagDashJump); hasDashJump {
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
	if ch.IsHidden() && !viewer.HasRoleAtLeast(ch.role) {
		return
	}
	viewer.Send(&response.LeavePlayer{ID: ch.GetID()}, types.SEND_POLICY_ENCRYPT)
}

func (ch *Character) GetBonusHp() int32   { return ch.bonusHp }
func (ch *Character) GetBonusMp() int32   { return ch.bonusMp }
func (ch *Character) GetInvincible() bool { return ch.invincible }
func (ch *Character) GetGender() uint8    { return ch.look.gender }
func (ch *Character) GetSkinColor() uint8 { return ch.look.skinColor }
func (ch *Character) GetFace() uint32     { return ch.look.face }
func (ch *Character) GetHair() uint32     { return ch.look.hair }

func (ch *Character) SetGender(gender uint8) {
	if ch.look.gender == gender {
		return
	}
	ch.look.gender = gender
	ch.listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetHair(id uint32) {
	if ch.look.hair == id {
		return
	}
	ch.look.hair = id
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHair: int32(id),
	}, false)
	ch.listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetFace(id uint32) {
	if ch.look.face == id {
		return
	}
	ch.look.face = id
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatFace: int32(id),
	}, false)
	ch.listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) SetSkin(color uint8) {
	if ch.look.skinColor == color {
		return
	}
	ch.look.skinColor = color
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatSkin: int32(color),
	}, false)
	ch.listener.OnUpdateCharacterLook(ch)
}

func (ch *Character) GetLevel() uint8 { return ch.level }
func (ch *Character) GetExp() uint32  { return ch.exp }
func (ch *Character) GetSpawnPoint() uint8 {
	mapInstance := ch.GetMap()
	if mapInstance != nil && mapInstance.wz != nil {
		return mapInstance.wz.FindClosestPortalSpawnID(ch.position)
	}
	return 0
}

func (ch *Character) SetHp(v uint32, notify bool) {
	wasAlive := ch.IsAlive()
	if max := ch.GetMaxHp(); v > max {
		v = max
	}
	ch.hp = v
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatHP: int32(ch.GetHp())}, false)
	}

	if wasAlive && !ch.IsAlive() {
		ch.die()
	}
}

func (ch *Character) die() {
	ch.summons.Clear()
	if ch.stateMachine != nil {
		ch.stateMachine.CallHook("on_player_dead", ch)
	}
}

func (ch *Character) SetMp(v uint32, notify bool) {
	if max := ch.GetMaxMp(); v > max {
		v = max
	}
	ch.mp = v
	if !notify {
		return
	}
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatMP: int32(ch.GetMp())}, false)
}

func (ch *Character) SetBonusHp(v int32, notify bool) {
	ch.bonusHp = v
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetBonusMp(v int32, notify bool) {
	ch.bonusMp = v
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) SetMaxHpPercent(p int16, notify bool) {
	ch.stats.bonus.maxHpPercent = p
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetMaxMpPercent(p int16, notify bool) {
	ch.stats.bonus.maxMpPercent = p
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
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
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP:    int32(ch.GetHp()),
		constant.StatMaxHP: int32(ch.GetMaxHp()),
		constant.StatMP:    int32(ch.GetMp()),
		constant.StatMaxMP: int32(ch.GetMaxMp()),
	}, false)
}

func (ch *Character) SetInvincible(b bool) { ch.invincible = b }

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
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatHP: int32(ch.GetHp())}, false)
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
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{constant.StatMP: int32(ch.GetMp())}, false)
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
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
		constant.StatMP: int32(ch.GetMp()),
	}, false)
}

func (ch *Character) SetBaseHp(v uint32, notify bool) {
	if v > constant.StatMaxHPMP {
		v = constant.StatMaxHPMP
	}
	ch.baseHp = v
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) SetBaseMp(v uint32, notify bool) {
	if v > constant.StatMaxHPMP {
		v = constant.StatMaxHPMP
	}
	ch.baseMp = v
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

func (ch *Character) AddBaseHp(amount uint32, notify bool) {
	ch.lifeCore.AddBaseHp(amount)
	if ch.GetHp() > ch.GetMaxHp() {
		ch.SetHp(ch.GetMaxHp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatHP:    int32(ch.GetHp()),
			constant.StatMaxHP: int32(ch.GetMaxHp()),
		}, false)
	}
}

func (ch *Character) AddBaseMp(amount uint32, notify bool) {
	ch.lifeCore.AddBaseMp(amount)
	if ch.GetMp() > ch.GetMaxMp() {
		ch.SetMp(ch.GetMaxMp(), false)
	}
	if notify {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
			constant.StatMP:    int32(ch.GetMp()),
			constant.StatMaxMP: int32(ch.GetMaxMp()),
		}, false)
	}
}

const clockTimer = "clock"

func (ch *Character) Warp(ctx actor.Context, targetMap *Map, spawnPoint uint8) error {
	if ch.gameWorld == nil {
		return fmt.Errorf("no game world")
	}
	return ch.gameWorld.Systems().Get[*MapSystem]().Warp(ctx, ch, targetMap, spawnPoint, nil)
}

func (ch *Character) EnterPortal(ctx actor.Context, portal *Portal) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}
	if portal.wz.Position.DistanceSq(ch.position) > 40000 && !ch.ActsAsGM() {
		ch.listener.OnUnlockAction(ch)
		return fmt.Errorf("portal %q is too far", portal.wz.Name)
	}

	if scriptName := portal.Script(); scriptName != "" {
		root := m.GetLuaRoot()
		if root == nil {
			ch.listener.OnUnlockAction(ch)
			return fmt.Errorf("root lua state not found")
		}
		scriptPath := fmt.Sprintf("script/portal/%s.lua", scriptName)
		thread, err := luax.NewThread(root, scriptPath)
		if err != nil {
			ch.listener.OnScriptError(ch, scriptPath, err)
			ch.listener.OnUnlockAction(ch)
			return fmt.Errorf("portal script thread: %w", err)
		}
		luax.SetConfiguration(thread, luax.Configuration{
			ActorContext: ctx,
			ActorPID:     m.LogicActorPID(),
		})
		luax.CallAsync(ctx, root, thread, "on_enter", ch, portal).Do(func([]lua.LValue) error {
			if ch.dialog.Thread() == nil {
				ch.listener.OnUnlockAction(ch)
			}
			return nil
		}).OnError(func(err error) {
			log.Printf("portal script %s failed: %v", scriptPath, err)
			ch.listener.OnScriptError(ch, scriptPath, err)
			if ch.dialog.Thread() == nil {
				ch.listener.OnUnlockAction(ch)
			}
		})
		return nil
	}

	targetMap := ch.gameWorld.Systems().Get[*MapSystem]().Find(m.StateMachine(), uint32(portal.wz.TargetMapId))
	if targetMap == nil {
		ch.listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if targetMap.wz != nil && ch.level < targetMap.wz.LevelLimit && !ch.ActsAsGM() {
		ch.listener.OnUnlockAction(ch)
		return fmt.Errorf("map %d needs level %d", targetMap.wz.ID, targetMap.wz.LevelLimit)
	}
	targetPortal := targetMap.FindPortalByName(portal.wz.Target)
	if targetPortal == nil || targetPortal.wz == nil {
		return ch.Warp(ctx, targetMap, 0)
	}
	return ch.Warp(ctx, targetMap, targetPortal.wz.ID)
}

func (ch *Character) Revive(ctx actor.Context) error {
	m := ch.GetMap()
	if m == nil || m.wz == nil {
		return fmt.Errorf("not on a map")
	}
	if ch.GetHp() > 0 {
		ch.listener.OnUpdateStats(ch, nil, true)
		return nil
	}

	targetMap := ch.gameWorld.Systems().Get[*MapSystem]().Find(m.StateMachine(), uint32(m.wz.ReturnMapId))
	if targetMap == nil {
		ch.listener.OnUpdateStats(ch, nil, true)
		return fmt.Errorf("return map %d not found", m.wz.ReturnMapId)
	}

	ch.SetHp(50, false)
	ch.stance = constant.StanceDefaultValue
	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatHP: int32(ch.GetHp()),
	}, true)
	if ch.stateMachine != nil {
		ch.stateMachine.CallHook("on_player_revive", ch)
	}
	return ch.Warp(ctx, targetMap, 0)
}

func (ch *Character) Relocate(spawnPoint uint8) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}
	if m.wz == nil {
		return fmt.Errorf("map model not found")
	}
	pos, ok := m.wz.GetSpawnPosition(spawnPoint)
	if !ok {
		return fmt.Errorf("invalid spawn point %d", spawnPoint)
	}
	ch.stance = constant.StanceDefaultValue
	beforePosition := ch.moveTo(m, pos)
	if ch.listener != nil {
		ch.listener.OnFieldRelocate(ch, spawnPoint)
		ch.listener.OnPlayerMove(ch, beforePosition, []dto.MoveFragment{
			&dto.TeleportMovement{
				BasicMovement: &dto.BasicMovement{
					Command: 3,
					Stance:  ch.stance,
				},
				Position: ch.position,
			},
		})
		ch.listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}

func (ch *Character) EnterInnerPortal(portalName string, to types.Vector2[int16]) error {
	m := ch.GetMap()
	if m == nil {
		return fmt.Errorf("not on a map")
	}
	portal := m.FindPortalByName(portalName)
	if portal == nil || portal.wz == nil {
		return fmt.Errorf("portal %q not found", portalName)
	}
	if portal.wz.Position.DistanceSq(ch.position) > 22500 && !ch.HasRoleAtLeast(constant.RoleAdmin) {
		return fmt.Errorf("portal %q is too far", portalName)
	}
	ch.moveTo(m, to)
	return nil
}

func (ch *Character) moveTo(m *Map, pos types.Vector2[int16]) types.Vector2[int16] {
	before := ch.position
	ch.position = pos
	for _, summon := range ch.summons.All() {
		if summon == nil || summon.owner != ch {
			continue
		}
		if summon.movementType != constant.SummonMoveFollow && summon.movementType != constant.SummonMoveCircleFollow {
			continue
		}
		summonBefore := summon.position
		summon.position = pos
		m.OnMoved(summon, summonBefore)
	}
	m.OnMoved(ch, before)
	return before
}

func (ch *Character) Send(p types.Packet, policy types.SendPolicy) error {
	if ch.sendable == nil {
		return nil
	}
	return ch.sendable.Send(p, policy)
}

func (ch *Character) IsHidden() bool {
	return ch.gm.hidden
}

func (ch *Character) SetHidden(hidden bool) {
	if ch.gm.hidden == hidden {
		return
	}
	ch.gm.hidden = hidden
	ch.listener.OnHiddenChanged(ch, hidden)
}

// Occupied is true while the character is dead, in an NPC conversation or shop, or in a shop or mini game room;
// it may not then talk to NPCs, move items, open a room, or leave the map. A trade does not count.
func (ch *Character) Occupied() bool {
	if !ch.IsAlive() || ch.dialog.Busy() {
		return true
	}
	switch ch.miniRoom.(type) {
	case nil, *Trade:
		return false
	}
	return true
}

func (ch *Character) GetID() uint32 {
	return ch.id
}

func (ch *Character) MarkLoggedOut() <-chan *internal.CharacterSaveEntry {
	ch.session.logoutEntry = make(chan *internal.CharacterSaveEntry, 1)
	ch.session.loggedOut.Store(true)
	return ch.session.logoutEntry
}

func (ch *Character) LoggedOut() bool {
	return ch.session.loggedOut.Load()
}

// An RPC reply returns to the actor that sent the request, which the character may have left by moving maps.
// Post runs fn on the actor that owns the character now. It skips fn once the character has left the channel or
// logged in again; unsettled requests are settled after that login.
func (ch *Character) Post(fn func(actx actor.Context)) {
	ch.gameWorld.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(actx actor.Context, current *Character) {
		if current != ch {
			return
		}
		fn(actx)
	})
}

func (ch *Character) GetPK() uint32 {
	return ch.id
}

func (ch *Character) GetRole() constant.CharacterRole {
	return ch.role
}

func (ch *Character) GetName() string {
	return ch.name
}

func (ch *Character) Message(message string) {
	ch.listener.OnMessage(ch, constant.MsgLightBlueText, message)
}

func (ch *Character) HasRoleAtLeast(role constant.CharacterRole) bool {
	return ch.role >= role
}

func (inv *Inventory) SetMeso(meso int32) {
	if meso < 0 {
		inv.meso = 0
		return
	}
	inv.meso = meso
}

func (ch *Character) remainingExpToMaxLevel() uint32 {
	if ch == nil || ch.level >= 200 {
		return 0
	}
	if ch.gameWorld == nil {
		return ^uint32(0)
	}
	resources := ch.gameWorld.GetResources()
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
	ch.family.gainProgress(exp)
	ch.exp += exp
	ch.expReachedAt = clock.Now()
	ch.listener.OnExpGain(ch, exp)

	if !ch.tryLevelUp() {
		ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
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
	if rate := max(ch.stats.bonus.expRate, int16(ch.family.ExpRate())); rate > 0 {
		exp = mulClamp(exp, uint32(rate)) / 100
	}
	exp = mulClamp(exp, uint32(ch.buffs.HolySymbolExpRate())) / 100
	if ch.debuffs.Has(constant.DebuffFlagCurse) {
		exp /= 2
	}
	if gw := ch.gameWorld; gw != nil {
		if r := gw.GetExpRate(); r > 0 {
			exp = mulClamp(exp, uint32(r))
		}
	}
	return exp
}

func (ch *Character) tryLevelUp() bool {
	if ch.gameWorld == nil {
		return false
	}

	resources := ch.gameWorld.GetResources()
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
			constant.StatAvailableAP: int32(ch.points.ap),
			constant.StatAvailableSP: int32(ch.points.sp),
		}
		ch.listener.OnUpdateStats(ch, stats, false)
		for i := 0; i < levelDiff; i++ {
			ch.broadcastLevelUpEffect()
		}
		ch.quests.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
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
	ch.expReachedAt = clock.Now()
	ch.NotifyMemberFieldsChanged()
	if newLevel < oldLevel {
		if ch.gameWorld != nil {
			resources := ch.gameWorld.GetResources()
			if resources != nil {
				if newLevel > 1 {
					ch.exp = resources.GetExpNeededForLevel(newLevel - 1)
				} else {
					ch.exp = 0
				}
			}
		}
	}

	ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatLevel: int32(ch.level),
		constant.StatEXP:   int32(ch.exp),
	}, false)
}

func (ch *Character) SpawnMist(skill *SkillEntry, spec mistSpec) *Mist {
	m := ch.GetMap()
	if m == nil {
		return nil
	}
	ld := skill.wz.GetLevelData(skill.Level())
	if ld == nil {
		return nil
	}
	if spec.bounds == (types.Rect[int32]{}) {
		spec.bounds = ld.Area(int32(ch.position.X), int32(ch.position.Y), ch.stance%2 != 0)
	}

	mist := &Mist{
		objectCore:   ObjectCore{position: ch.position},
		causer:       ch.GetID(),
		skillWz:      skill.wz,
		skillLevel:   uint8(skill.Level()),
		carnivalTeam: constant.CarnivalTeamNone,
		skillDelay:   8,
	}
	if team := ch.carnival.Team(); team != nil {
		mist.carnivalTeam = team.teamID
	}
	mist.spread(m, spec)
	return mist
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
