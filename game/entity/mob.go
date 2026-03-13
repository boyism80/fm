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

type mobMobStatusEntry struct {
	value             int32
	Wz                *wz.Skill // skill that applied this debuff (nil if non-skill); packet uses Wz.ID when non-nil
	Level             uint8
	CauserCharacterID uint32 // character ID who applied the debuff (for EXP/drops when mob dies from DoT); 0 if none
	expiresAt         time.Time
	cancelTimer       *time.Timer
}

type Mob struct {
	Life
	Wz       *wz.Mob
	Foothold int16
	Spawn    *MobSpawn
	debuffs  map[constant.MobStatus]*mobMobStatusEntry
	debuffMu sync.RWMutex
}

func (m *Mob) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMob
}

func (m *Mob) Is(typ constant.ObjectType) bool {
	return m.GetObjectType().Has(typ)
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
		"set_status": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			statusTbl := L.CheckTable(2)
			maskLV := statusTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobStatus table with numeric mask expected")
				return 0
			}
			status := constant.MobStatus(uint32(lua.LVAsNumber(maskLV)))
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
			var causer *Character
			if L.GetTop() >= 6 && L.Get(6) != lua.LNil {
				if cud, ok := L.Get(6).(*lua.LUserData); ok {
					if ch, ok := cud.Value.(*Character); ok {
						causer = ch
					} else {
						L.ArgError(6, "causer must be Character or nil")
						return 0
					}
				} else {
					L.ArgError(6, "causer must be Character or nil")
					return 0
				}
			}
			mob.ApplyMobStatus(status, value, durationMs, skillWz, skillLevel, causer)
			return 0
		},
		"clear_status": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			statusTbl := L.CheckTable(2)
			maskLV := statusTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobStatus table with numeric mask expected")
				return 0
			}
			status := constant.MobStatus(uint32(lua.LVAsNumber(maskLV)))
			mob.CancelMobStatus(status)
			return 0
		},
		"has_status": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			statusTbl := L.CheckTable(2)
			maskLV := statusTbl.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobStatus table with numeric mask expected")
				return 0
			}
			status := constant.MobStatus(uint32(lua.LVAsNumber(maskLV)))
			L.Push(lua.LBool(mob.HasMobStatus(status)))
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
		"damage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc != 3 {
				L.ArgError(2, "damage(attacker, amount) requires attacker (Character or nil) and amount")
				return 0
			}

			var attacker *Character
			if L.Get(2) != lua.LNil {
				if aud, ok := L.Get(2).(*lua.LUserData); ok {
					if ch, ok := aud.Value.(*Character); ok {
						attacker = ch
					} else {
						L.ArgError(2, "attacker must be Character or nil")
						return 0
					}
				} else {
					L.ArgError(2, "attacker must be Character or nil")
					return 0
				}
			}

			amount := L.CheckInt(3)
			if amount <= 0 {
				L.Push(lua.LBool(false))
				return 1
			}

			killed := mob.ApplyDamage(attacker, uint32(amount))
			L.Push(lua.LBool(killed))
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

// ApplyMobStatus applies a debuff (e.g. FREEZE) to the mob with optional duration. If durationMs > 0, a timer clears the debuff when it expires.
// causer is the character who applied the debuff; when the mob dies from this debuff (e.g. DoT), EXP/drops can be attributed to that character.
// Internally only the character ID is stored on the debuff entry so reconnects do not break attribution.
// Re-applying the same debuff type: existing entry is removed (timer stopped, no cancel packet), then a new entry is created (no apply packet). Same observable behavior as refresh.
// skillWz and skillLevel identify the skill that applied the debuff; packet uses skillWz.ID when non-nil.
// Boss check (FREEZE/STUN/SEAL not applied to boss) is the caller's responsibility.
func (m *Mob) ApplyMobStatus(debuff constant.MobStatus, value int32, durationMs int64, skillWz *wz.Skill, skillLevel uint8, causer *Character) {
	m.debuffMu.Lock()
	if m.debuffs == nil {
		m.debuffs = make(map[constant.MobStatus]*mobMobStatusEntry)
	}
	existing := m.debuffs[debuff]
	wasRefresh := existing != nil
	if existing != nil {
		if existing.cancelTimer != nil {
			existing.cancelTimer.Stop()
		}
		delete(m.debuffs, debuff)
	}

	var causerID uint32
	if causer != nil {
		causerID = causer.GetID()
	}

	entry := &mobMobStatusEntry{
		value:             value,
		Wz:                skillWz,
		Level:             skillLevel,
		CauserCharacterID: causerID,
	}
	if durationMs > 0 {
		entry.expiresAt = time.Now().Add(time.Duration(durationMs) * time.Millisecond)
		entry.cancelTimer = time.AfterFunc(time.Duration(durationMs)*time.Millisecond, func() {
			m.CancelMobStatus(debuff)
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
		mapInstance.listener.OnMobMobStatusApplied(mapInstance, m, debuff, value, packetSkillID, durationMs)
	}
}

// CancelMobStatus removes a debuff and notifies the map listener (broadcast cancel packet).
func (m *Mob) CancelMobStatus(debuff constant.MobStatus) {
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
		mapInstance.listener.OnMobMobStatusCancelled(mapInstance, m, debuff)
	}
}

// HasMobStatus returns whether the mob currently has the given debuff.
func (m *Mob) HasMobStatus(debuff constant.MobStatus) bool {
	m.debuffMu.RLock()
	defer m.debuffMu.RUnlock()
	_, ok := m.debuffs[debuff]
	return ok
}

// GetCauserCharacterID returns the character ID who applied the given debuff, or 0 if the debuff is not present or has no causer.
// Used when the mob dies from that debuff (e.g. DoT) to attribute EXP/drops; resolve via map.GetPlayer(id).
func (m *Mob) GetCauserCharacterID(debuff constant.MobStatus) uint32 {
	m.debuffMu.RLock()
	defer m.debuffMu.RUnlock()
	entry, ok := m.debuffs[debuff]
	if !ok || entry == nil {
		return 0
	}
	return entry.CauserCharacterID
}

// ClearAllMobStatusTimers stops all debuff timers without notifying. Used when mob is removed from map.
func (m *Mob) ClearAllMobStatusTimers() {
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

func (m *Mob) ApplyDamage(attacker *Character, amount uint32) bool {
	if amount == 0 || m.Hp == 0 {
		return false
	}

	// Clamp to current HP so we do not underflow.
	damage := amount
	if damage > uint32(m.Hp) {
		damage = uint32(m.Hp)
	}

	m.Hp -= uint16(damage)
	isDead := m.Hp == 0

	// If mob is still alive, only HP bar update is needed.
	if !isDead {
		if attacker != nil && attacker.Listener != nil {
			maxHp := m.GetMaxHp()
			if maxHp == 0 {
				return false
			}
			percent := min(uint32(m.Hp)*100/uint32(maxHp), 100)
			attacker.Listener.OnShowMobHp(m, uint8(percent))
		}
		return false
	}

	// Mob dies - handle all death-related logic.
	if attacker != nil {
		log.Printf("Mob %d (ID: %d) killed by character %d", m.OID, m.Wz.ID, attacker.GetID())

		exp := uint32(m.Wz.EXP)
		if attacker.BonusStats.ExpRate > 0 {
			exp = exp * uint32(attacker.BonusStats.ExpRate) / 100
		}
		attacker.AddExp(exp)

		// Generate mob drops
		m.dropItems(attacker)
	} else {
		log.Printf("Mob %d (ID: %d) killed with no attacker", m.OID, m.Wz.ID)
	}

	// Remove mob from map
	mapInstance := m.GetObject().GetMap()
	if mapInstance != nil {
		mapInstance.RemoveMob(m.OID, constant.MOB_DIE_ANIMATION_TYPE_FADE_OUT)
	}

	// Send mob HP update to attacker via listener (0% HP)
	if attacker != nil && attacker.Listener != nil {
		attacker.Listener.OnShowMobHp(m, 0)
	}

	return true
}
