-- Mob name (String.wz/Mob.img.xml): 주니어 카투스

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000329, min = 1, max = 1, chance = 400000 },
	{ item = 4010002, min = 1, max = 1, chance = 1000 },
	{ item = 4020004, min = 1, max = 1, chance = 150 },
	{ item = 2040705, min = 1, max = 1, chance = 70 },
	{ item = 2044102, min = 1, max = 1, chance = 70 },
	{ item = 1051004, min = 1, max = 1, chance = 40 },
	{ item = 1050024, min = 1, max = 1, chance = 40 },
	{ item = 1032006, min = 1, max = 1, chance = 40 },
	{ item = 1002048, min = 1, max = 1, chance = 40 },
	{ item = 1072027, min = 1, max = 1, chance = 40 },
	{ item = 1040044, min = 1, max = 1, chance = 40 },
	{ item = 1060033, min = 1, max = 1, chance = 40 },
	{ item = 1462003, min = 1, max = 1, chance = 40 },
	{ item = 2022155, min = 1, max = 1, chance = 6000 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1082183, min = 1, max = 1, chance = 40 },
	{ item = 1082186, min = 1, max = 1, chance = 40 },
	{ item = 2044901, min = 1, max = 1, chance = 70 },
	{ item = 2043212, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
