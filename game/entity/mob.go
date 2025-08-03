package entity

import (
	"log"
	"math/rand"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
)

type Mob struct {
	Life
	Spec     *data.MobSpec
	Foothold int16
	MapID    uint32 // Map ID where this mob is located
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
