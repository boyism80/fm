package entity

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
	lua "github.com/yuin/gopher-lua"
)

type Mob struct {
	Life
	Spec     *data.MobSpec
	Foothold int16
	MapID    uint32 // Map ID where this mob is located
}

// Luable interface implementation
func (m *Mob) LuaTypeName() string {
	return "LuaMob"
}

func (m *Mob) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return id
				L.Push(lua.LNumber(mob.Spec.ID))
				return 1
			} else {
				L.ArgError(2, "id() is read-only")
				return 0
			}
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return name (using ID as name for now)
				L.Push(lua.LString(fmt.Sprintf("Mob_%d", mob.Spec.ID)))
				return 1
			} else {
				L.ArgError(2, "name() is read-only")
				return 0
			}
		},
		"exp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return exp
				L.Push(lua.LNumber(mob.Spec.EXP))
				return 1
			} else {
				L.ArgError(2, "exp() is read-only")
				return 0
			}
		},
		"foothold": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return foothold
				L.Push(lua.LNumber(mob.Foothold))
				return 1
			} else if argc == 2 {
				// Setter: foothold(value)
				foothold := L.CheckInt(2)
				mob.Foothold = int16(foothold)
				return 0
			} else {
				L.ArgError(2, "foothold() requires 0 or 1 arguments")
				return 0
			}
		},
		"map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				// Getter: return map_id
				L.Push(lua.LNumber(mob.MapID))
				return 1
			} else if argc == 2 {
				// Setter: map_id(value)
				mapID := L.CheckInt(2)
				if mapID < 0 {
					mapID = 0
				}
				mob.MapID = uint32(mapID)
				return 0
			} else {
				L.ArgError(2, "map_id() requires 0 or 1 arguments")
				return 0
			}
		},
	}
}

func (m *Mob) String() string {
	return m.LuaTypeName()
}

func (m *Mob) Type() lua.LValueType {
	return lua.LTUserData
}

func (m *Mob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(0)
	return nil
}

func (m *Mob) dropItems(attacker *Character) {
	// Get map instance
	mapInstance := m.Context.GetMap(m.MapID)
	if mapInstance == nil {
		return
	}

	resources := m.Context.GetResources()
	mobDrops, ok := resources.Drops[m.Spec.ID]
	if !ok {
		return
	}

	// Collect all drops first to calculate positions
	var drops []struct {
		isMeso bool
		count  int32
		item   Item
	}

	for _, drop := range mobDrops {
		// Check drop probability
		if rand.Float32() > drop.Prob {
			continue
		}

		if drop.Item == 0 {
			// Generate meso drop
			min := float64(drop.Money) * 0.75
			max := float64(drop.Money)
			count := int32(min + rand.Float64()*(max-min))
			if count == 0 {
				continue
			}
			drops = append(drops, struct {
				isMeso bool
				count  int32
				item   Item
			}{isMeso: true, count: count, item: nil})
		} else {
			// Generate item drop
			count := uint16(1)
			if drop.Max != 0 && drop.Min != 0 {
				count = uint16(rand.Intn(int(drop.Max-drop.Min)+1) + int(drop.Min))
			}

			// Create item using NewItem function
			item, err := NewItem(drop.Item, count, m.Context)
			if err != nil {
				continue // Skip if item creation fails
			}
			drops = append(drops, struct {
				isMeso bool
				count  int32
				item   Item
			}{isMeso: false, count: int32(count), item: item})
		}
	}

	// Spawn drops with position spreading (following old server pattern)
	spawnPoint := m.Position
	spacing := int16(15)

	for i, drop := range drops {
		destPoint := spawnPoint
		if len(drops) > 1 {
			offset := spacing * int16(i/2+1)
			if i%2 == 0 {
				destPoint.X += offset
			} else {
				destPoint.X -= offset
			}
		}

		if drop.isMeso {
			// Spawn meso drop
			if err := mapInstance.SpawnMeso(drop.count, destPoint, attacker.GetID(), constant.DROP_TYPE_OWNED); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {
			// Bind drop information to item
			drop.item.BindDrop(&Drop{
				Object: &Object{
					OID:      0,         // Will be set by Map.SpawnItem
					Position: destPoint, // Use calculated position
					Context:  m.Context,
				},
				Owner:        attacker.GetID(),
				SpawnedPoint: spawnPoint,
				DropType:     constant.DROP_TYPE_OWNED,
			})

			// Spawn item drop
			if err := mapInstance.SpawnItem(drop.item, attacker.GetID(), constant.DROP_TYPE_OWNED); err != nil {
				log.Printf("Failed to spawn item drop: %v", err)
			}
		}
	}
}

// Damage applies damage to the mob and handles all death-related logic
func (m *Mob) Damage(damage uint16, attacker *Character) bool {
	if damage > m.Hp {
		damage = m.Hp
	}
	m.Hp -= damage
	isDead := m.Hp == 0
	if !isDead {
		return false
	}

	// Mob dies - handle all death-related logic
	log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Spec.ID, attacker.GetID())

	// Add experience to character
	attacker.AddExp(uint32(m.Spec.EXP))

	// Generate mob drops
	m.dropItems(attacker)

	// Remove mob from map
	mapInstance := m.Context.GetMap(m.MapID)
	if mapInstance != nil {
		mapInstance.RemoveMob(m.OID, constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT)
	}

	return true
}
