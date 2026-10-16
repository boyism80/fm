package game

import (
	"log"
	"math/rand"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/types"
)

type reactorDropSpawn struct {
	isMeso bool
	count  int32
	item   Item
}

func (r *Reactor) DropItems() {
	if r == nil || r.wz == nil || r.onMap == nil || r.gameWorld == nil {
		return
	}

	resources := r.gameWorld.GetResources()
	if resources == nil {
		return
	}

	reactorDrops, ok := resources.ReactorDrops[r.wz.ID]
	if !ok || len(reactorDrops) == 0 {
		return
	}

	trigger := r.GetTrigger()
	mesoRate := float32(r.gameWorld.GetMesoRate())
	gmDrop := trigger != nil && trigger.ActsAsGM()
	spawns := make([]reactorDropSpawn, 0, len(reactorDrops))

	for _, entry := range reactorDrops {
		if quest := resources.QuestItems[entry.Item]; quest > 0 && !questInProgress(trigger, quest) {
			continue
		}
		if !gmDrop && rand.Float32() > entry.Prob {
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

			item, err := NewItem(entry.Item, count, r.gameWorld)
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

	spawnPoint := r.position
	dropX := spawnPoint.X - int16(12*len(spawns))
	ownerID, dropType := reactorDropOwner(trigger)

	for _, spawn := range spawns {
		destPoint := types.Point[int16]{X: dropX, Y: spawnPoint.Y}
		dropX += 25

		if spawn.isMeso {
			if _, err := r.onMap.SpawnMeso(spawn.count, destPoint, spawnPoint, ownerID, dropType, false); err != nil {
				log.Printf("Failed to spawn reactor meso drop: %v", err)
			}
		} else {
			if err := r.onMap.SpawnItem(spawn.item, ItemSpawn{
				Position: destPoint,
				From:     spawnPoint,
				Owner:    ownerID,
				DropType: dropType,
			}); err != nil {
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
	return ch.quests != nil && ch.quests.Get(questID).IsStarted()
}

func reactorDropOwner(trigger *Character) (uint32, constant.DropType) {
	if trigger == nil {
		return 0, constant.DropTypeFFA
	}
	if trigger.party.ID() != nil {
		return trigger.GetID(), constant.DropTypeParty
	}
	return trigger.GetID(), constant.DropTypeOwnerOnly
}
