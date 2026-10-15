package constant

type WeaponType uint32

const (
	WeaponTypeNone       WeaponType = 0
	WeaponTypeSword1H    WeaponType = 30
	WeaponTypeAxe1H      WeaponType = 31
	WeaponTypeBlunt1H    WeaponType = 32
	WeaponTypeDagger     WeaponType = 33
	WeaponTypeKatara     WeaponType = 34
	WeaponTypeMagicArrow WeaponType = 35
	WeaponTypeWand       WeaponType = 37
	WeaponTypeStaff      WeaponType = 38
	WeaponTypeSword2H    WeaponType = 40
	WeaponTypeAxe2H      WeaponType = 41
	WeaponTypeBlunt2H    WeaponType = 42
	WeaponTypeSpear      WeaponType = 43
	WeaponTypePoleArm    WeaponType = 44
	WeaponTypeBow        WeaponType = 45
	WeaponTypeCrossbow   WeaponType = 46
	WeaponTypeClaw       WeaponType = 47
	WeaponTypeKnuckle    WeaponType = 48
	WeaponTypeGun        WeaponType = 49
	WeaponTypeDualBow    WeaponType = 52
	WeaponTypeCannon     WeaponType = 53
)

func AllWeaponTypes() map[string]WeaponType {
	return map[string]WeaponType{
		"None": WeaponTypeNone, "Sword1H": WeaponTypeSword1H, "Axe1H": WeaponTypeAxe1H,
		"Blunt1H": WeaponTypeBlunt1H, "Dagger": WeaponTypeDagger, "Katara": WeaponTypeKatara,
		"MagicArrow": WeaponTypeMagicArrow, "Wand": WeaponTypeWand, "Staff": WeaponTypeStaff,
		"Sword2H": WeaponTypeSword2H, "Axe2H": WeaponTypeAxe2H, "Blunt2H": WeaponTypeBlunt2H,
		"Spear": WeaponTypeSpear, "PoleArm": WeaponTypePoleArm, "Bow": WeaponTypeBow,
		"Crossbow": WeaponTypeCrossbow, "Claw": WeaponTypeClaw, "Knuckle": WeaponTypeKnuckle,
		"Gun": WeaponTypeGun, "DualBow": WeaponTypeDualBow, "Cannon": WeaponTypeCannon,
	}
}
