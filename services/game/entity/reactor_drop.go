package entity

import (
	"log"
	"math/rand"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type reactorDropSpawn struct {
	isMeso bool
	count  int32
	item   Item
}

func (r *Reactor) DropItems() {
	if r == nil || r.Wz == nil || r.Map == nil || r.GameWorld == nil {
		return
	}

	resources := r.GameWorld.GetResources()
	if resources == nil {
		return
	}

	reactorDrops, ok := resources.ReactorDrops[r.Wz.ID]
	if !ok || len(reactorDrops) == 0 {
		return
	}

	trigger := r.GetTrigger()
	mesoRate := float32(r.GameWorld.GetMesoRate())
	gmDrop := trigger != nil && trigger.HasRoleAtLeast(constant.RoleAdmin)
	spawns := make([]reactorDropSpawn, 0, len(reactorDrops))

	for _, entry := range reactorDrops {
		if entry.QuestID > 0 && !questInProgress(trigger, entry.QuestID) {
			continue
		}
		if !gmDrop && entry.Prob > 0 && rand.Float32() > entry.Prob {
			continue
		}

		if entry.Item == 0 {
			min := float64(entry.Money) * 0.75
			max := float64(entry.Money)
			if entry.Max > 0 || entry.Min > 0 {
				min = float64(entry.Min)
				max = float64(entry.Max)
			}
			if max < min {
				max = min
			}
			var count int32
			if max == min {
				count = int32(max)
			} else {
				count = int32(min + rand.Float64()*(max-min))
			}
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoRate)
			if count == 0 {
				continue
			}
			spawns = append(spawns, reactorDropSpawn{isMeso: true, count: count})
		} else {
			count := uint16(1)
			if entry.Max > 0 && entry.Min > 0 {
				count = uint16(rand.Intn(int(entry.Max-entry.Min)+1) + int(entry.Min))
			} else if entry.Max > 0 {
				count = entry.Max
			} else if entry.Min > 0 {
				count = entry.Min
			}

			item, err := NewItem(entry.Item, count, r.GameWorld)
			if err != nil {
				continue
			}
			if eq, ok := item.(Equipment); ok {
				if em, ok := eq.GetModel().(wz.Equipment); ok {
					eq.GetEquipmentCore().RandomizeStats(em)
				}
			}
			spawns = append(spawns, reactorDropSpawn{isMeso: false, count: int32(count), item: item})
		}
	}

	if len(spawns) == 0 {
		return
	}

	spawnPoint := r.Position
	dropX := spawnPoint.X - int16(12*len(spawns))
	ownerID, dropType := reactorDropOwner(trigger)

	for _, spawn := range spawns {
		destPoint := types.Point[int16]{X: dropX, Y: spawnPoint.Y}
		dropX += 25

		if spawn.isMeso {
			if _, err := r.Map.SpawnMeso(spawn.count, destPoint, ownerID, dropType, false); err != nil {
				log.Printf("Failed to spawn reactor meso drop: %v", err)
			}
		} else {
			fp := &FieldPlacement{
				ObjectCore: &ObjectCore{
					Position:  destPoint,
					GameWorld: r.GameWorld,
				},
				Owner:        ownerID,
				SpawnedPoint: spawnPoint,
				DropType:     dropType,
			}
			fp.ObjectCore.self = fp
			spawn.item.BindFieldPlacement(fp)

			if err := r.Map.SpawnItem(spawn.item, ownerID, dropType); err != nil {
				log.Printf("Failed to spawn reactor item drop: %v", err)
			}
		}
	}
}

func questInProgress(ch *Character, questID uint32) bool {
	if questID == 0 {
		return true
	}
	if ch == nil {
		return false
	}
	return ch.Quests != nil && ch.Quests.Get(questID).IsStarted()
}

func reactorDropOwner(trigger *Character) (uint32, constant.DropType) {
	if trigger == nil {
		return 0, constant.DropTypeFFA
	}
	if trigger.GetPartyID() != nil {
		return trigger.GetID(), constant.DropTypeParty
	}
	return trigger.GetID(), constant.DropTypeOwnerOnly
}
