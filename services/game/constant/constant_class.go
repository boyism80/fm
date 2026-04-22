package constant

type ClassType uint16

const (
	ClassBeginner           ClassType = 0
	ClassWarrior            ClassType = 100
	ClassFighter            ClassType = 110
	ClassCrusader           ClassType = 111
	ClassHero               ClassType = 112
	ClassPage               ClassType = 120
	ClassWhiteKnight        ClassType = 121
	ClassPaladin            ClassType = 122
	ClassSpearman           ClassType = 130
	ClassDragonKnight       ClassType = 131
	ClassDarkKnight         ClassType = 132
	ClassMagician           ClassType = 200
	ClassFpWizard           ClassType = 210
	ClassFpMage             ClassType = 211
	ClassFpArchMage         ClassType = 212
	ClassIlWizard           ClassType = 220
	ClassIlMage             ClassType = 221
	ClassIlArchMage         ClassType = 222
	ClassCleric             ClassType = 230
	ClassPriest             ClassType = 231
	ClassBishop             ClassType = 232
	ClassBowman             ClassType = 300
	ClassHunter             ClassType = 310
	ClassRanger             ClassType = 311
	ClassBowmaster          ClassType = 312
	ClassCrossbowman        ClassType = 320
	ClassSniper             ClassType = 321
	ClassCrossbowmaster     ClassType = 322
	ClassThief              ClassType = 400
	ClassAssassin           ClassType = 410
	ClassHermit             ClassType = 411
	ClassNightlord          ClassType = 412
	ClassBandit             ClassType = 420
	ClassChiefBandit        ClassType = 421
	ClassShadower           ClassType = 422
	ClassPirate             ClassType = 500
	ClassBrawler            ClassType = 510
	ClassGunslinger         ClassType = 520
	ClassMarauder           ClassType = 511
	ClassOutlaw             ClassType = 521
	ClassBuccaneer          ClassType = 512
	ClassCorsair            ClassType = 522
	ClassMapleleafBrigadier ClassType = 800
	ClassGM                 ClassType = 900
	ClassSuperGM            ClassType = 910
	ClassNoblesse           ClassType = 1000
	ClassDawnWarrior1       ClassType = 1100
	ClassDawnWarrior2       ClassType = 1110
	ClassDawnWarrior3       ClassType = 1111
	ClassDawnWarrior4       ClassType = 1112
	ClassBlazeWizard1       ClassType = 1200
	ClassBlazeWizard2       ClassType = 1210
	ClassBlazeWizard3       ClassType = 1211
	ClassBlazeWizard4       ClassType = 1212
	ClassWindArcher1        ClassType = 1300
	ClassWindArcher2        ClassType = 1310
	ClassWindArcher3        ClassType = 1311
	ClassWindArcher4        ClassType = 1312
	ClassNightWalker1       ClassType = 1400
	ClassNightWalker2       ClassType = 1410
	ClassNightWalker3       ClassType = 1411
	ClassNightWalker4       ClassType = 1412
	ClassThunderBreaker1    ClassType = 1500
	ClassThunderBreaker2    ClassType = 1510
	ClassThunderBreaker3    ClassType = 1511
	ClassThunderBreaker4    ClassType = 1512
	ClassLegend             ClassType = 2000
	ClassAran2              ClassType = 2100
	ClassAran3              ClassType = 2110
	ClassAran4              ClassType = 2111
	ClassAran5              ClassType = 2112
)

func AllClassConstants() map[string]ClassType {
	return map[string]ClassType{
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
