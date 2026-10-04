-- Mob name (String.wz/Mob.img.xml): 흰 모래토끼

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000324, min = 1, max = 1, chance = 400000 },
	{ item = 4020006, min = 1, max = 1, chance = 150 },
	{ item = 4010001, min = 1, max = 1, chance = 1000 },
	{ item = 2040601, min = 1, max = 1, chance = 70 },
	{ item = 1382002, min = 1, max = 1, chance = 40 },
	{ item = 1002119, min = 1, max = 1, chance = 40 },
	{ item = 1332010, min = 1, max = 1, chance = 40 },
	{ item = 1050025, min = 1, max = 1, chance = 40 },
	{ item = 1082002, min = 1, max = 1, chance = 40 },
	{ item = 1072007, min = 1, max = 1, chance = 40 },
	{ item = 1032009, min = 1, max = 1, chance = 40 },
	{ item = 1050005, min = 1, max = 1, chance = 40 },
	{ item = 4003004, min = 1, max = 1, chance = 600 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1072288, min = 1, max = 1, chance = 40 },
	{ item = 2040420, min = 1, max = 1, chance = 70 },
	{ item = 2049000, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
