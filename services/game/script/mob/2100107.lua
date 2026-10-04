-- Mob name (String.wz/Mob.img.xml): 목도리 프릴드

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 4000327, min = 1, max = 1, chance = 400000 },
	{ item = 4020003, min = 1, max = 1, chance = 150 },
	{ item = 4020002, min = 1, max = 1, chance = 150 },
	{ item = 2048001, min = 1, max = 1, chance = 70 },
	{ item = 1332004, min = 1, max = 1, chance = 40 },
	{ item = 1382017, min = 1, max = 1, chance = 40 },
	{ item = 1041054, min = 1, max = 1, chance = 40 },
	{ item = 1061050, min = 1, max = 1, chance = 40 },
	{ item = 1002141, min = 1, max = 1, chance = 40 },
	{ item = 1051011, min = 1, max = 1, chance = 40 },
	{ item = 1040059, min = 1, max = 1, chance = 40 },
	{ item = 1060045, min = 1, max = 1, chance = 40 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 1492003, min = 1, max = 1, chance = 40 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
