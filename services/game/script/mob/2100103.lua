-- Mob name (String.wz/Mob.img.xml): 카투스

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000330, min = 1, max = 1, chance = 400000 },
	{ item = 4010003, min = 1, max = 1, chance = 1000 },
	{ item = 4020004, min = 1, max = 1, chance = 150 },
	{ item = 2040501, min = 1, max = 1, chance = 70 },
	{ item = 1322009, min = 1, max = 1, chance = 40 },
	{ item = 1092002, min = 1, max = 1, chance = 40 },
	{ item = 1432002, min = 1, max = 1, chance = 40 },
	{ item = 1072054, min = 1, max = 1, chance = 40 },
	{ item = 1082016, min = 1, max = 1, chance = 40 },
	{ item = 1072078, min = 1, max = 1, chance = 40 },
	{ item = 1061054, min = 1, max = 1, chance = 40 },
	{ item = 1041058, min = 1, max = 1, chance = 40 },
	{ item = 1002096, min = 1, max = 1, chance = 40 },
	{ item = 1050011, min = 1, max = 1, chance = 40 },
	{ item = 2022155, min = 1, max = 1, chance = 6000 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1072291, min = 1, max = 1, chance = 40 },
	{ item = 2044314, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
