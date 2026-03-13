package constant

// Class code constants (job IDs). Used for Lua Class table and class_of checks.
const (
	ClassBeginner           uint16 = 0
	ClassWarrior            uint16 = 100
	ClassFighter            uint16 = 110
	ClassCrusader           uint16 = 111
	ClassHero               uint16 = 112
	ClassPage               uint16 = 120
	ClassWhiteKnight        uint16 = 121
	ClassPaladin            uint16 = 122
	ClassSpearman           uint16 = 130
	ClassDragonKnight       uint16 = 131
	ClassDarkKnight         uint16 = 132
	ClassMagician           uint16 = 200
	ClassFpWizard           uint16 = 210
	ClassFpMage             uint16 = 211
	ClassFpArchMage         uint16 = 212
	ClassIlWizard           uint16 = 220
	ClassIlMage             uint16 = 221
	ClassIlArchMage         uint16 = 222
	ClassCleric             uint16 = 230
	ClassPriest             uint16 = 231
	ClassBishop             uint16 = 232
	ClassBowman             uint16 = 300
	ClassHunter             uint16 = 310
	ClassRanger             uint16 = 311
	ClassBowmaster          uint16 = 312
	ClassCrossbowman        uint16 = 320
	ClassSniper             uint16 = 321
	ClassCrossbowmaster     uint16 = 322
	ClassThief              uint16 = 400
	ClassAssassin           uint16 = 410
	ClassHermit             uint16 = 411
	ClassNightlord          uint16 = 412
	ClassBandit             uint16 = 420
	ClassChiefBandit        uint16 = 421
	ClassShadower           uint16 = 422
	ClassPirate             uint16 = 500
	ClassBrawler            uint16 = 510
	ClassGunslinger         uint16 = 520
	ClassMarauder           uint16 = 511
	ClassOutlaw             uint16 = 521
	ClassBuccaneer          uint16 = 512
	ClassCorsair            uint16 = 522
	ClassMapleleafBrigadier uint16 = 800
	ClassGM                 uint16 = 900
	ClassSuperGM            uint16 = 910
	ClassNoblesse           uint16 = 1000
	ClassDawnWarrior1       uint16 = 1100
	ClassDawnWarrior2       uint16 = 1110
	ClassDawnWarrior3       uint16 = 1111
	ClassDawnWarrior4       uint16 = 1112
	ClassBlazeWizard1       uint16 = 1200
	ClassBlazeWizard2       uint16 = 1210
	ClassBlazeWizard3       uint16 = 1211
	ClassBlazeWizard4       uint16 = 1212
	ClassWindArcher1        uint16 = 1300
	ClassWindArcher2        uint16 = 1310
	ClassWindArcher3        uint16 = 1311
	ClassWindArcher4        uint16 = 1312
	ClassNightWalker1       uint16 = 1400
	ClassNightWalker2       uint16 = 1410
	ClassNightWalker3       uint16 = 1411
	ClassNightWalker4       uint16 = 1412
	ClassThunderBreaker1    uint16 = 1500
	ClassThunderBreaker2    uint16 = 1510
	ClassThunderBreaker3    uint16 = 1511
	ClassThunderBreaker4    uint16 = 1512
	ClassLegend             uint16 = 2000
	ClassAran2              uint16 = 2100
	ClassAran3              uint16 = 2110
	ClassAran4              uint16 = 2111
	ClassAran5              uint16 = 2112
)

// AllClassConstants returns name -> class code for Lua Class table injection.
func AllClassConstants() map[string]uint16 {
	return map[string]uint16{
		"Beginner":           ClassBeginner,
		"Warrior":            ClassWarrior,
		"Fighter":            ClassFighter,
		"Crusader":           ClassCrusader,
		"Hero":               ClassHero,
		"Page":               ClassPage,
		"WhiteKnight":        ClassWhiteKnight,
		"Paladin":            ClassPaladin,
		"Spearman":           ClassSpearman,
		"DragonKnight":       ClassDragonKnight,
		"DarkKnight":         ClassDarkKnight,
		"Magician":           ClassMagician,
		"FpWizard":           ClassFpWizard,
		"FpMage":             ClassFpMage,
		"FpArchMage":         ClassFpArchMage,
		"IlWizard":           ClassIlWizard,
		"IlMage":             ClassIlMage,
		"IlArchMage":         ClassIlArchMage,
		"Cleric":             ClassCleric,
		"Priest":             ClassPriest,
		"Bishop":             ClassBishop,
		"Bowman":             ClassBowman,
		"Hunter":             ClassHunter,
		"Ranger":             ClassRanger,
		"Bowmaster":          ClassBowmaster,
		"Crossbowman":        ClassCrossbowman,
		"Sniper":             ClassSniper,
		"Crossbowmaster":     ClassCrossbowmaster,
		"Thief":              ClassThief,
		"Assassin":           ClassAssassin,
		"Hermit":             ClassHermit,
		"Nightlord":          ClassNightlord,
		"Bandit":             ClassBandit,
		"ChiefBandit":        ClassChiefBandit,
		"Shadower":           ClassShadower,
		"Pirate":             ClassPirate,
		"Brawler":            ClassBrawler,
		"Gunslinger":         ClassGunslinger,
		"Marauder":           ClassMarauder,
		"Outlaw":             ClassOutlaw,
		"Buccaneer":          ClassBuccaneer,
		"Corsair":            ClassCorsair,
		"MapleleafBrigadier": ClassMapleleafBrigadier,
		"GM":                 ClassGM,
		"SuperGM":            ClassSuperGM,
		"Noblesse":           ClassNoblesse,
		"DawnWarrior1":       ClassDawnWarrior1,
		"DawnWarrior2":       ClassDawnWarrior2,
		"DawnWarrior3":       ClassDawnWarrior3,
		"DawnWarrior4":       ClassDawnWarrior4,
		"BlazeWizard1":       ClassBlazeWizard1,
		"BlazeWizard2":       ClassBlazeWizard2,
		"BlazeWizard3":       ClassBlazeWizard3,
		"BlazeWizard4":       ClassBlazeWizard4,
		"WindArcher1":        ClassWindArcher1,
		"WindArcher2":        ClassWindArcher2,
		"WindArcher3":        ClassWindArcher3,
		"WindArcher4":        ClassWindArcher4,
		"NightWalker1":       ClassNightWalker1,
		"NightWalker2":       ClassNightWalker2,
		"NightWalker3":       ClassNightWalker3,
		"NightWalker4":       ClassNightWalker4,
		"ThunderBreaker1":    ClassThunderBreaker1,
		"ThunderBreaker2":    ClassThunderBreaker2,
		"ThunderBreaker3":    ClassThunderBreaker3,
		"ThunderBreaker4":    ClassThunderBreaker4,
		"Legend":             ClassLegend,
		"Aran2":              ClassAran2,
		"Aran3":              ClassAran3,
		"Aran4":              ClassAran4,
		"Aran5":              ClassAran5,
	}
}
