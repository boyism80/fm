package game

import (
	"log"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type Skills struct {
	owner           *Character
	entries         map[uint32]*SkillEntry
	macros          *SkillMacros
	homingTargetOID *uint32
}

func NewSkills(owner *Character) *Skills {
	return &Skills{
		owner:   owner,
		entries: make(map[uint32]*SkillEntry),
		macros:  NewSkillMacros(owner),
	}
}

func (sc *Skills) Bind(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 && entry.masterLevel < 1 {
		return
	}
	entry.owner = sc.owner
	sc.entries[skillID] = entry
}

func (sc *Skills) Register(skillID uint32, entry *SkillEntry) {
	if entry == nil || entry.Level() < 1 && entry.masterLevel < 1 {
		return
	}
	sc.Bind(skillID, entry)
	if entry.Level() > 0 {
		sc.CallHook(nil, entry, "on_passive")
	}
	sc.owner.listener.OnUpdateSkill(sc.owner, skillID, int32(entry.Level()), int32(entry.masterLevel))
}

func (sc *Skills) RestorePassives() {
	for _, entry := range sc.entries {
		if entry.Level() > 0 {
			sc.CallHook(nil, entry, "on_passive")
		}
	}
}

func (sc *Skills) Remove(skillID uint32) bool {
	entry := sc.entries[skillID]
	if entry == nil {
		return false
	}
	if entry.Level() > 0 {
		sc.CallHook(nil, entry, "on_unpassive")
	}
	delete(sc.entries, skillID)
	return true
}

func (sc *Skills) Get(skillID uint32) *SkillEntry {
	if entry := sc.entries[skillID]; entry != nil && entry.Level() > 0 && !entry.wz.Dojo {
		return entry
	}
	return sc.dojoSkill(skillID)
}

func (sc *Skills) Learned(skillID uint32) *SkillEntry {
	return sc.entries[skillID]
}

func (sc *Skills) dojoSkill(skillID uint32) *SkillEntry {
	if sc.owner.gameWorld == nil {
		return nil
	}
	w := sc.owner.gameWorld.GetResources().GetSkill(skillID)
	if w == nil || !w.Dojo {
		return nil
	}
	if sc.owner.dojo.energy < DojoEnergyFull || !sc.owner.dojo.onField() {
		return nil
	}
	return NewSkillEntry(sc.owner, w, 1, 1)
}

func (sc *Skills) ForEach(fn func(skillID uint32, entry *SkillEntry)) {
	for id, e := range sc.entries {
		fn(id, e)
	}
}

func (sc *Skills) Level(skillID uint32) int {
	entry := sc.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (sc *Skills) UseAttack(skillID uint32, activate func() bool) bool {
	wzSkill := sc.owner.gameWorld.GetResources().GetSkill(skillID)
	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}
	skillLevel := sc.Level(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}
	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}
	skillEntry := sc.Get(skillID)
	if levelData.Cooldown > 0 && (skillEntry == nil || skillEntry.IsCooling()) {
		return false
	}
	if !activate() {
		return false
	}
	if levelData.Cooldown > 0 {
		skillEntry.StartCooldown(levelData.Cooldown)
	}
	return true
}

func (sc *Skills) LevelData(id constant.SkillID) *wz.SkillLevelData {
	entry := sc.Get(uint32(id))
	if entry == nil || entry.Level() <= 0 {
		return nil
	}
	return entry.wz.GetLevelData(entry.Level())
}

func (sc *Skills) masteries(weapon constant.WeaponType) []constant.SkillID {
	switch weapon {
	case constant.WeaponTypeSword1H, constant.WeaponTypeSword2H:
		return []constant.SkillID{constant.SkillSwordMastery, constant.SkillSwordMastery1200000, constant.SkillSwordMasteryCygnus}
	case constant.WeaponTypeAxe1H, constant.WeaponTypeAxe2H:
		return []constant.SkillID{constant.SkillAxeMastery}
	case constant.WeaponTypeBlunt1H, constant.WeaponTypeBlunt2H:
		return []constant.SkillID{constant.SkillBwMastery}
	case constant.WeaponTypeDagger:
		return []constant.SkillID{constant.SkillDaggerMastery}
	case constant.WeaponTypeSpear:
		return []constant.SkillID{constant.SkillSpearMastery}
	case constant.WeaponTypePoleArm:
		return []constant.SkillID{constant.SkillPoleArmMastery}
	case constant.WeaponTypeBow:
		if sc.owner.class >= constant.ClassNoblesse {
			return []constant.SkillID{constant.SkillBowMasteryCygnus, constant.SkillBowExpertCygnus}
		}
		return []constant.SkillID{constant.SkillBowMastery, constant.SkillBowExpert}
	case constant.WeaponTypeCrossbow:
		return []constant.SkillID{constant.SkillCrossbowMastery, constant.SkillCrossbowExpert}
	case constant.WeaponTypeClaw:
		if sc.owner.class >= constant.ClassNoblesse {
			return []constant.SkillID{constant.SkillClawMasteryCygnus}
		}
		return []constant.SkillID{constant.SkillClawMastery}
	case constant.WeaponTypeKnuckle:
		if sc.owner.class >= constant.ClassNoblesse {
			return []constant.SkillID{constant.SkillKnuckleMasteryCygnus}
		}
		return []constant.SkillID{constant.SkillKnuckleMastery}
	case constant.WeaponTypeGun:
		return []constant.SkillID{constant.SkillGunMastery}
	}
	return nil
}

func (sc *Skills) criticalLevelData(weapon constant.WeaponType) *wz.SkillLevelData {
	switch weapon {
	case constant.WeaponTypeBow, constant.WeaponTypeCrossbow:
		if sc.owner.class >= constant.ClassNoblesse {
			return sc.LevelData(constant.SkillCriticalShotCygnus)
		}
		return sc.LevelData(constant.SkillCriticalShot)
	case constant.WeaponTypeClaw:
		if sc.owner.class >= constant.ClassNoblesse {
			return sc.LevelData(constant.SkillCriticalThrowCygnus)
		}
		return sc.LevelData(constant.SkillCriticalThrow)
	case constant.WeaponTypeKnuckle:
		return sc.LevelData(constant.SkillCriticalPunch)
	}
	return nil
}

func (sc *Skills) makerLevel() int {
	skillID := uint32(constant.SkillMaker)
	if sc.owner.class >= constant.ClassNoblesse && sc.owner.class < constant.ClassLegend {
		skillID = uint32(constant.SkillMakerCygnus)
	}
	entry := sc.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (sc *Skills) grant(skill ExchangeSkill) {
	if skill.SkillID == 0 {
		return
	}
	if sc.owner.gameWorld == nil {
		return
	}
	resources := sc.owner.gameWorld.GetResources()
	if resources == nil {
		return
	}
	wzSkill := resources.GetSkill(skill.SkillID)
	if wzSkill == nil {
		return
	}
	targetLevel := exchangeSkillTargetLevel(0, skill.Level)
	targetMaster := skill.MasterLevel
	if targetMaster <= 0 {
		targetMaster = wzSkill.DefaultMasterLevel()
	}
	entry := sc.Learned(skill.SkillID)
	if entry == nil {
		entry = NewSkillEntry(sc.owner, wzSkill, targetLevel, targetMaster)
		entry.expiration = time.Time{}
		sc.Register(skill.SkillID, entry)
	} else {
		targetLevel = exchangeSkillTargetLevel(entry.Level(), skill.Level)
		targetMaster = max(entry.masterLevel, targetMaster)
		if targetLevel != entry.Level() || targetMaster != entry.masterLevel {
			entry.SetLevels(targetLevel, targetMaster)
		}
	}
}
