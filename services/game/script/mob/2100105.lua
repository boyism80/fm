-- Mob name (String.wz/Mob.img.xml): 벨라모아

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000328, min = 1, max = 1, chance = 400000 },
	{ item = 4010004, min = 1, max = 1, chance = 1000 },
	{ item = 4006001, min = 1, max = 1, chance = 600 },
	{ item = 2040902, min = 1, max = 1, chance = 70 },
	{ item = 2044002, min = 1, max = 1, chance = 70 },
	{ item = 1452003, min = 1, max = 1, chance = 40 },
	{ item = 1302006, min = 1, max = 1, chance = 40 },
	{ item = 1372004, min = 1, max = 1, chance = 40 },
	{ item = 1002013, min = 1, max = 1, chance = 40 },
	{ item = 1072086, min = 1, max = 1, chance = 40 },
	{ item = 1002119, min = 1, max = 1, chance = 40 },
	{ item = 1072007, min = 1, max = 1, chance = 40 },
	{ item = 1082020, min = 1, max = 1, chance = 40 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1052101, min = 1, max = 1, chance = 40 },
	{ item = 1482003, min = 1, max = 1, chance = 40 },
	{ item = 2040316, min = 1, max = 1, chance = 70 },
	{ item = 2040319, min = 1, max = 1, chance = 70 },
	{ item = 2044412, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
