package entity

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type mobDebuffEntry struct {
	value       int32
	Wz          *wz.Skill // skill that applied this debuff (nil if non-skill); packet uses Wz.ID when non-nil
	Level       uint8
	expiresAt   time.Time
	cancelTimer *time.Timer
}

type Mob struct {
	Life
	Wz       *wz.Mob
	Foothold int16
	Spawn    *MobSpawn
	debuffs  map[constant.Debuff]*mobDebuffEntry
	debuffMu sync.RWMutex
}

func (m *Mob) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMob
}

// GetObject implements ObjectProvider (Mob embeds Life which embeds Object).
func (m *Mob) GetObject() *Object {
	return &m.Life.Object
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
				L.Push(lua.LNumber(mob.Wz.ID))
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
				L.Push(lua.LString(fmt.Sprintf("Mob_%d", mob.Wz.ID)))
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
				L.Push(lua.LNumber(mob.Wz.EXP))
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
			if L.GetTop() != 1 {
				L.ArgError(2, "map_id() is read-only (map is fixed at spawn)")
				return 0
			}
			mapInstance := mob.GetObject().GetMap()
			if mapInstance == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(mapInstance.ID))
			return 1
		},
		"set_debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			debuffTbl := L.CheckTable(2)
			maskLV := debuffTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "Debuff table with numeric mask expected")
				return 0
			}
			debuff := constant.Debuff{Mask: uint32(lua.LVAsNumber(maskLV))}
			value := int32(L.CheckInt(3))
			durationMs := L.CheckInt64(4)
			if durationMs < 0 {
				durationMs = 0
			}
			var skillWz *wz.Skill
			var skillLevel uint8
			if L.GetTop() >= 5 {
				if lv, ok := L.Get(5).(*lua.LUserData); ok {
					if se, ok := lv.Value.(*SkillEntry); ok && se != nil && se.Skill != nil {
						skillWz = se.Skill
						skillLevel = uint8(se.SkillLevel)
					}
				}
			}
			mob.ApplyDebuff(debuff, value, durationMs, skillWz, skillLevel)
			return 0
		},
		"clear_debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			debuffTbl := L.CheckTable(2)
			maskLV := debuffTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "Debuff table with numeric mask expected")
				return 0
			}
			debuff := constant.Debuff{Mask: uint32(lua.LVAsNumber(maskLV))}
			mob.CancelDebuff(debuff)
			return 0
		},
		"has_debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			debuffTbl := L.CheckTable(2)
			maskLV := debuffTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "Debuff table with numeric mask expected")
				return 0
			}
			debuff := constant.Debuff{Mask: uint32(lua.LVAsNumber(maskLV))}
			L.Push(lua.LBool(mob.HasDebuff(debuff)))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			// Convert wz.Mob to Lua table
			tbl := L.NewTable()
			if mob.Wz != nil {
				tbl.RawSetString("id", lua.LNumber(mob.Wz.ID))
				tbl.RawSetString("level", lua.LNumber(mob.Wz.Level))
				tbl.RawSetString("max_hp", lua.LNumber(mob.Wz.MaxHP))
				tbl.RawSetString("max_mp", lua.LNumber(mob.Wz.MaxMP))
				tbl.RawSetString("exp", lua.LNumber(mob.Wz.EXP))
			}
			L.Push(tbl)
			return 1
		},
	}
}

func (m *Mob) String() string {
	return m.LuaTypeName()
}

func (m *Mob) Type() lua.LValueType {
	return lua.LTUserData
}

// ApplyDebuff applies a debuff (e.g. FREEZE) to the mob with optional duration. If durationMs > 0, a timer clears the debuff when it expires.
// Re-applying the same debuff type: existing entry is removed (timer stopped, no cancel packet), then a new entry is created (no apply packet). Same observable behavior as refresh.
// skillWz and skillLevel identify the skill that applied the debuff; packet uses skillWz.ID when non-nil.
// Boss check (FREEZE/STUN/SEAL not applied to boss) is the caller's responsibility.
func (m *Mob) ApplyDebuff(debuff constant.Debuff, value int32, durationMs int64, skillWz *wz.Skill, skillLevel uint8) {
	m.debuffMu.Lock()
	if m.debuffs == nil {
		m.debuffs = make(map[constant.Debuff]*mobDebuffEntry)
	}
	existing := m.debuffs[debuff]
	wasRefresh := existing != nil
	if existing != nil {
		if existing.cancelTimer != nil {
			existing.cancelTimer.Stop()
		}
		delete(m.debuffs, debuff)
	}

	entry := &mobDebuffEntry{
		value: value,
		Wz:    skillWz,
		Level: skillLevel,
	}
	if durationMs > 0 {
		entry.expiresAt = time.Now().Add(time.Duration(durationMs) * time.Millisecond)
		entry.cancelTimer = time.AfterFunc(time.Duration(durationMs)*time.Millisecond, func() {
			m.CancelDebuff(debuff)
		})
	}
	m.debuffs[debuff] = entry
	m.debuffMu.Unlock()

	if wasRefresh {
		return
	}
	packetSkillID := uint32(0)
	if skillWz != nil {
		packetSkillID = skillWz.ID
	}
	mapInstance := m.GetObject().GetMap()
	if mapInstance != nil && mapInstance.listener != nil {
		mapInstance.listener.OnMobDebuffApplied(mapInstance, m, debuff, value, packetSkillID, durationMs)
	}
}

// CancelDebuff removes a debuff and notifies the map listener (broadcast cancel packet).
func (m *Mob) CancelDebuff(debuff constant.Debuff) {
	m.debuffMu.Lock()
	entry, ok := m.debuffs[debuff]
	if !ok {
		m.debuffMu.Unlock()
		return
	}
	if entry.cancelTimer != nil {
		entry.cancelTimer.Stop()
	}
	delete(m.debuffs, debuff)
	m.debuffMu.Unlock()

	mapInstance := m.GetObject().GetMap()
	if mapInstance != nil && mapInstance.listener != nil {
		mapInstance.listener.OnMobDebuffCancelled(mapInstance, m, debuff)
	}
}

// HasDebuff returns whether the mob currently has the given debuff.
func (m *Mob) HasDebuff(debuff constant.Debuff) bool {
	m.debuffMu.RLock()
	defer m.debuffMu.RUnlock()
	_, ok := m.debuffs[debuff]
	return ok
}

// ClearAllDebuffTimers stops all debuff timers without notifying. Used when mob is removed from map.
func (m *Mob) ClearAllDebuffTimers() {
	m.debuffMu.Lock()
	defer m.debuffMu.Unlock()
	for _, entry := range m.debuffs {
		if entry.cancelTimer != nil {
			entry.cancelTimer.Stop()
		}
	}
	m.debuffs = nil
}

// Serialize method removed - use DTO instead

func (m *Mob) dropItems(attacker *Character) {
	mapInstance := m.GetObject().GetMap()
	if mapInstance == nil {
		return
	}

	resources := m.Context.GetResources()
	mobDrops, ok := resources.Drops[m.Wz.ID]
	if !ok {
		return
	}

	// Collect all drops first to calculate positions
	var drops []struct {
		isMeso bool
		count  int32
		item   Item
	}

	dropRate := float32(m.Context.GetDropRate())
	mesoRate := float32(m.Context.GetMesoRate())
	dropRateMul := int16(100)
	mesoAmountMul := int16(100)
	if attacker.BonusStats.DropRate > 0 {
		dropRateMul = attacker.BonusStats.DropRate
	}
	if attacker.BonusStats.MesoMultiplier > 0 {
		mesoAmountMul = attacker.BonusStats.MesoMultiplier
	}
	dropRateMulF := float32(dropRateMul) / 100.0
	mesoAmountMulF := float32(mesoAmountMul) / 100.0

	for _, drop := range mobDrops {
		adjustedProb := drop.Prob * dropRate * dropRateMulF
		if adjustedProb > 1.0 {
			adjustedProb = 1.0
		}

		if rand.Float32() > adjustedProb {
			continue
		}

		if drop.Item == 0 {
			min := float64(drop.Money) * 0.75
			max := float64(drop.Money)
			count := int32(min + rand.Float64()*(max-min))
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoRate)
			if count == 0 {
				continue
			}
			count = int32(float32(count) * mesoAmountMulF)
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
			if _, err := mapInstance.SpawnMeso(drop.count, destPoint, attacker.GetID(), constant.DROP_TYPE_OWNED); err != nil {
				log.Printf("Failed to spawn meso drop: %v", err)
			}
		} else {
			// Bind drop information to item
			drop.item.BindDrop(&Drop{
				Object: &Object{
					OID:      0,         // Will be set by Map.SpawnItem
					Position: destPoint, // Use calculated position (X calculated, Y is mob Y)
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
	log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())

	exp := uint32(m.Wz.EXP)
	if attacker.BonusStats.ExpRate > 0 {
		exp = exp * uint32(attacker.BonusStats.ExpRate) / 100
	}
	attacker.AddExp(exp)

	// Generate mob drops
	m.dropItems(attacker)

	// Remove mob from map
	mapInstance := m.GetObject().GetMap()
	if mapInstance != nil {
		mapInstance.RemoveMob(m.OID, constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT)
	}

	// Send mob HP update to attacker via listener
	attacker.Listener.OnShowMobHp(m, uint8(m.Hp*100/m.Life.GetMaxHp()))

	return true
}
