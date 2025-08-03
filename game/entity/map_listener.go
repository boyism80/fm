package entity

import "github.com/boyism80/fm/game/constant"

// RemoveItemType constants (matching resp.RemoveItemType)
const (
	REMOVE_ITEM_TYPE_EXPIRED     = 0
	REMOVE_ITEM_TYPE_NO_ANIMATED = 1
	REMOVE_ITEM_TYPE_ANIMATED    = 2
	REMOVE_ITEM_TYPE_EXPLOSION   = 3
	REMOVE_ITEM_TYPE_LOOT_BY_PET = 4
)

// MapListener defines interface for map events
type MapListener interface {
	OnPlayerAdded(mapID uint32, playerID uint32, character *Character, init bool)
	OnPlayerRemoved(mapID uint32, playerID uint32)
	OnPlayerMoved(mapID uint32, playerID uint32, character *Character)
	OnPlayerChat(mapID uint32, playerID uint32, message string)
	OnItemSpawned(mapID uint32, itemID uint32, item Item, drop *Drop)
	OnMesoSpawned(mapID uint32, itemID uint32, meso *Meso)
	OnItemRemoved(mapID uint32, itemID uint32, characterID uint32, mode uint8)
	OnMobSpawned(mapID uint32, mobID uint32, mob *Mob)
	OnMobRemoved(mapID uint32, mobID uint32, animationType constant.MobDieAnimationType)
	OnMobControllerChange(mob *Mob, before *Character, after *Character)
}

// MapListenerImpl implements MapListener interface
type MapListenerImpl struct {
	// This will be implemented to handle packet sending
	// without creating import cycles
}

// NewMapListenerImpl creates a new MapListenerImpl instance
func NewMapListenerImpl() *MapListenerImpl {
	return &MapListenerImpl{}
}

// OnPlayerAdded is called when a player is added to a map
func (l *MapListenerImpl) OnPlayerAdded(mapID uint32, playerID uint32, character *Character, init bool) {
	// TODO: Send spawn player packet to other players on the map
	// This will be implemented by the game server
}

// OnPlayerRemoved is called when a player is removed from a map
func (l *MapListenerImpl) OnPlayerRemoved(mapID uint32, playerID uint32) {
	// TODO: Send leave player packet to other players on the map
	// This will be implemented by the game server
}

// OnPlayerMoved is called when a player moves on a map
func (l *MapListenerImpl) OnPlayerMoved(mapID uint32, playerID uint32, character *Character) {
	// TODO: Send move player packet to other players on the map
	// This will be implemented by the game server
}

// OnPlayerChat is called when a player sends a chat message
func (l *MapListenerImpl) OnPlayerChat(mapID uint32, playerID uint32, message string) {
	// TODO: Send chat packet to other players on the map
	// This will be implemented by the game server
}

// OnMobControllerChange is called when a mob's controller changes
func (l *MapListenerImpl) OnMobControllerChange(mob *Mob, before *Character, after *Character) {
	// TODO: Send StartControlMob packet to new controller
	// This will be implemented by the game server
}
