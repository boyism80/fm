-- Mob name (String.wz/Mob.img.xml): 갈색 모래토끼

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000325, min = 1, max = 1, chance = 400000 },
	{ item = 2002002, min = 1, max = 1, chance = 10000 },
	{ item = 4010006, min = 1, max = 1, chance = 1000 },
	{ item = 4020001, min = 1, max = 1, chance = 150 },
	{ item = 2044602, min = 1, max = 1, chance = 70 },
	{ item = 1092021, min = 1, max = 1, chance = 40 },
	{ item = 1472006, min = 1, max = 1, chance = 40 },
	{ item = 1060017, min = 1, max = 1, chance = 40 },
	{ item = 1061028, min = 1, max = 1, chance = 40 },
	{ item = 1002129, min = 1, max = 1, chance = 40 },
	{ item = 1041027, min = 1, max = 1, chance = 40 },
	{ item = 1061025, min = 1, max = 1, chance = 40 },
	{ item = 4003004, min = 1, max = 1, chance = 600 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1082183, min = 1, max = 1, chance = 40 },
	{ item = 2044210, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
