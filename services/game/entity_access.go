package game

import (
	time "time"

	constant "github.com/boyism80/fm/maple/constant"
	wz "github.com/boyism80/fm/maple/wz"
	types "github.com/boyism80/fm/types"
)

func (x *Account) Storage() *Storage {
	return x.storage
}

func (x *ActivePet) Item() *Pet {
	return x.item
}

func (x *Alliance) LeaderCharacterID() uint32 {
	return x.leaderCharacterID
}

func (x *Alliance) Name() string {
	return x.name
}

func (x *CarnivalMatch) FieldMapID() uint32 {
	return x.fieldMapID
}

func (x *CarnivalMatch) ReviveMapID() uint32 {
	return x.reviveMapID
}

func (x *Character) Account() *Account {
	return x.account
}

func (x *Character) AccountID() uint32 {
	return x.accountID
}

func (x *Character) Buddies() *BuddyList {
	return x.buddies
}

func (x *Character) Buffs() *Buffs {
	return x.buffs
}

func (x *Character) Carnival() *CarnivalMembership {
	return x.carnival
}

func (x *Character) Chalkboard() string {
	return x.chalkboard
}

func (x *Character) Debuffs() *Debuffs {
	return x.debuffs
}

func (x *Character) Dialog() *Dialog {
	return x.dialog
}

func (x *Character) Doors() *Doors {
	return x.doors
}

func (x *Character) Duey() *Duey {
	return x.duey
}

func (x *Character) Family() *Family {
	return x.family
}

func (x *Character) Guild() *GuildMembership {
	return x.guild
}

func (x *Character) Inventory() *Inventory {
	return x.inventory
}

func (x *Character) ItemEffectID() uint32 {
	return x.itemEffectID
}

func (x *Character) KeyLayout() *KeyLayout {
	return x.keyLayout
}

func (x *Character) LieDetector() *LieDetector {
	return x.lieDetector
}

func (x *Character) Listener() CharacterListener {
	return x.listener
}

func (x *Character) Messenger() *Messenger {
	return x.messenger
}

func (x *Character) MiniRoom() MiniRoom {
	return x.miniRoom
}

func (x *Character) MonsterBook() *MonsterBook {
	return x.monsterBook
}

func (x *Character) Mount() *Mount {
	return x.mount
}

func (x *Character) Notes() *Notes {
	return x.notes
}

func (x *Character) Party() *PartyMembership {
	return x.party
}

func (x *Character) Pets() *Pets {
	return x.pets
}

func (x *Character) Points() *Points {
	return x.points
}

func (x *Character) Quests() *Quests {
	return x.quests
}

func (x *Character) RemoteShops() *RemoteShops {
	return x.remoteShops
}

func (x *Character) Client() Sendable {
	return x.sendable
}

func (x *Character) Skills() *Skills {
	return x.skills
}

func (x *Character) StateMachine() *StateMachine {
	return x.stateMachine
}

func (x *Character) StoreBank() *StoreBank {
	return x.storeBank
}

func (x *Character) Summons() *Summons {
	return x.summons
}

func (x *Character) TeleportStones() *TeleportStones {
	return x.teleportStones
}

func (x *Character) Unsettled() *UnsettledRequests {
	return x.unsettled
}

func (x *Character) Wedding() *Wedding {
	return x.wedding
}

func (x *Door) OwnerID() uint32 {
	return x.ownerID
}

func (x *Door) PartyID() *uint32 {
	return x.partyID
}

func (x *EntrustedShopView) Elapsed() time.Duration {
	return x.elapsed
}

func (x *EntrustedShopView) ItemID() uint32 {
	return x.itemID
}

func (x *EntrustedShopView) Items() []*ShopItem {
	return x.items
}

func (x *EntrustedShopView) Meso() int32 {
	return x.meso
}

func (x *EntrustedShopView) OwnerName() string {
	return x.ownerName
}

func (x *EntrustedShopView) Sold() []ShopSale {
	return x.sold
}

func (x *EntrustedShopView) Title() string {
	return x.title
}

func (x *EntrustedShopView) Visitors() [3]*Character {
	return x.visitors
}

func (x *Family) Tree() FamilyTree {
	return x.tree
}

func (x *FieldPlacement) DropType() constant.DropType {
	return x.dropType
}

func (x *FieldPlacement) Owner() uint32 {
	return x.owner
}

func (x *FieldPlacement) PlayerDrop() bool {
	return x.playerDrop
}

func (x *FieldPlacement) Quest() uint32 {
	return x.quest
}

func (x *FieldPlacement) SpawnedPoint() types.Point[int16] {
	return x.spawnedPoint
}

func (x *Guild) LeaderCharacterID() uint32 {
	return x.leaderCharacterID
}

func (x *Guild) Logo() *GuildLogo {
	return x.logo
}

func (x *Guild) Name() string {
	return x.name
}

func (x *GuildLogo) Logo() uint32 {
	return x.logo
}

func (x *GuildLogo) LogoBG() uint32 {
	return x.logoBG
}

func (x *GuildLogo) LogoBGColor() uint32 {
	return x.logoBGColor
}

func (x *GuildLogo) LogoColor() uint32 {
	return x.logoColor
}

func (x *Homing) SkillLevel() uint8 {
	return x.skillLevel
}

func (x *Homing) SkillWz() *wz.Skill {
	return x.skillWz
}

func (x *Inventory) Rings() Rings {
	return x.rings
}

func (x *Inventory) Tabs() map[constant.InventoryType]*InventoryTab {
	return x.tabs
}

func (x *Map) GameWorld() GameWorld {
	return x.gameWorld
}

func (x *Map) Wz() *wz.Map {
	return x.wz
}

func (x *Meso) Count() int32 {
	return x.count
}

func (x *MessageBox) ItemID() uint32 {
	return x.itemID
}

func (x *MessageBox) Owner() string {
	return x.owner
}

func (x *MessageBox) Text() string {
	return x.text
}

func (x *Messenger) RoomID() uint32 {
	return x.roomID
}

func (x *MiniGame) Cards() []uint32 {
	return x.cards
}

func (x *MiniGame) Kind() uint8 {
	return x.kind
}

func (x *MiniGame) Piece() uint8 {
	return x.piece
}

func (x *MiniGame) Title() string {
	return x.title
}

func (x *MiniGame) Turn() uint8 {
	return x.turn
}

func (x *Mist) Bounds() types.Rect[int32] {
	return x.bounds
}

func (x *Mist) Causer() uint32 {
	return x.causer
}

func (x *Mist) MistType() constant.MistType {
	return x.mistType
}

func (x *Mist) MobMist() bool {
	return x.mobMist
}

func (x *Mist) SkillDelay() uint16 {
	return x.skillDelay
}

func (x *Mist) SkillLevel() uint8 {
	return x.skillLevel
}

func (x *Mist) SkillWz() *wz.Skill {
	return x.skillWz
}

func (x *Mob) Buffs() *MobBuffs {
	return x.buffs
}

func (x *Mob) Listener() MobListener {
	return x.listener
}

func (x *Mob) Skills() *MobSkills {
	return x.skills
}

func (x *Mob) Wz() *wz.Mob {
	return x.wz
}

func (x *MobBuff) SkillLevel() uint8 {
	return x.skillLevel
}

func (x *MobBuff) Values() map[constant.MobBuffFlag]int32 {
	return x.values
}

func (x *MobBuff) Wz() *wz.Skill {
	return x.wz
}

func (x *MobSkill) Slot() wz.MobSkillSlot {
	return x.slot
}

func (x *Notes) Entries() []Note {
	return x.entries
}

func (x *NoticeSystem) Muted() bool {
	return x.muted.Load()
}

func (x *NoticeSystem) SetMuted(muted bool) {
	x.muted.Store(muted)
}

func (x *Npc) Wz() *wz.NpcSpawn {
	return x.wz
}

func (x *PartyMembership) Search() *PartySearchConfig {
	return x.search
}

func (x *PartyMembership) SetSearch(v *PartySearchConfig) {
	x.search = v
}

func (x *PersonalShop) MaxItems() uint8 {
	return x.maxItems
}

func (x *Pets) Active() *ActivePet {
	return x.active
}

func (x *Pets) HPItem() uint32 {
	return x.hpItem
}

func (x *Pets) SetHPItem(v uint32) {
	x.hpItem = v
}

func (x *Pets) MPItem() uint32 {
	return x.mpItem
}

func (x *Pets) SetMPItem(v uint32) {
	x.mpItem = v
}

func (x *Portal) Wz() *wz.Portal {
	return x.wz
}

func (x *ReactorSpawn) Spawned() bool {
	return x.spawned
}

func (x Rings) Left() []*Ring {
	return x.left
}

func (x Rings) Mid() []*Ring {
	return x.mid
}

func (x *ShopItem) Bundles() uint16 {
	return x.bundles
}

func (x *ShopItem) Item() Item {
	return x.item
}

func (x *ShopItem) PerBundle() uint16 {
	return x.perBundle
}

func (x *ShopItem) Price() int32 {
	return x.price
}

func (x *SkillEntry) CooldownEnd() *time.Time {
	return x.cooldownEnd
}

func (x *Skills) Macros() *SkillMacros {
	return x.macros
}

func (x *StateMachine) Group() *StateMachineGroup {
	return x.group
}

func (x *StateMachineGroup) Name() string {
	return x.name
}

func (x *StateMachineGroup) ScriptPath() string {
	return x.scriptPath
}

func (x *Storage) Meso() int32 {
	return x.meso
}

func (x *Storage) Slots() uint8 {
	return x.slots
}

func (x *Storage) Tabs() map[constant.InventoryType][]Item {
	return x.tabs
}

func (x *StoreBank) Items() []*ShopItem {
	return x.items
}

func (x *StoreBank) Meso() int64 {
	return x.meso
}

func (x *Summon) MovementType() constant.SummonMovementType {
	return x.movementType
}

func (x *Summon) OwnerID() uint32 {
	return x.ownerID
}

func (x *Summon) SkillID() constant.SkillID {
	return x.skillID
}

func (x *Summon) SkillLevel() uint8 {
	return x.skillLevel
}

func (x *Summon) SummonType() constant.SummonType {
	return x.summonType
}

func (x *UnsettledRequest) InFlight() bool {
	return x.inFlight
}

func (x *UnsettledRequest) SetInFlight(v bool) {
	x.inFlight = v
}

func (x *Wedding) Marriage() *Marriage {
	return x.marriage
}
