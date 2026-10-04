-- Mob name (String.wz/Mob.img.xml): 모래 두더지

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 35, max = 55, chance = 500000 },
	{ item = 2000001, min = 1, max = 1, chance = 10000 },
	{ item = 2000003, min = 1, max = 1, chance = 10000 },
	{ item = 2060000, min = 20, max = 30, chance = 8000 },
	{ item = 2061000, min = 20, max = 30, chance = 8000 },
	{ item = 2002004, min = 1, max = 1, chance = 10000 },
	{ item = 4020007, min = 1, max = 1, chance = 150 },
	{ item = 4020000, min = 1, max = 1, chance = 150 },
	{ item = 2044501, min = 1, max = 1, chance = 70 },
	{ item = 1402002, min = 1, max = 1, chance = 40 },
	{ item = 1032004, min = 1, max = 1, chance = 40 },
	{ item = 1050001, min = 1, max = 1, chance = 40 },
	{ item = 1002055, min = 1, max = 1, chance = 40 },
	{ item = 1040049, min = 1, max = 1, chance = 40 },
	{ item = 1060037, min = 1, max = 1, chance = 40 },
	{ item = 1082068, min = 1, max = 1, chance = 40 },
	{ item = 1432001, min = 1, max = 1, chance = 40 },
	{ item = 4010007, min = 1, max = 1, chance = 1000 },
	{ item = 4000350, min = 1, max = 1, chance = 400000 },
	{ item = 1002619, min = 1, max = 1, chance = 40 },
	{ item = 2043214, min = 1, max = 1, chance = 70 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
