-- Mob name (String.wz/Mob.img.xml): 스콜피언

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 2002005, min = 1, max = 1, chance = 10000 },
	{ item = 4010002, min = 1, max = 1, chance = 1000 },
	{ item = 4010001, min = 1, max = 1, chance = 1000 },
	{ item = 4004002, min = 1, max = 1, chance = 80 },
	{ item = 2043002, min = 1, max = 1, chance = 70 },
	{ item = 2043701, min = 1, max = 1, chance = 70 },
	{ item = 1332009, min = 1, max = 1, chance = 40 },
	{ item = 1412004, min = 1, max = 1, chance = 40 },
	{ item = 1040062, min = 1, max = 1, chance = 40 },
	{ item = 1060051, min = 1, max = 1, chance = 40 },
	{ item = 1050029, min = 1, max = 1, chance = 40 },
	{ item = 1002164, min = 1, max = 1, chance = 40 },
	{ item = 1082005, min = 1, max = 1, chance = 40 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 4000351, min = 1, max = 1, chance = 400000 },
	{ item = 1492003, min = 1, max = 1, chance = 40 },
	{ item = 2043114, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
