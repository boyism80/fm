-- Mob name (String.wz/Mob.img.xml): 포장마차

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 1553, max = 3156, chance = 500000 },
	{ meso = true, min = 1553, max = 3156, chance = 1000000 },
	{ meso = true, min = 1553, max = 3156, chance = 1000000 },
	{ item = 2000005, min = 1, max = 5, chance = 1000000 },
	{ item = 2000004, min = 3, max = 6, chance = 1000000 },
	{ item = 4031354, min = 1, max = 1, chance = 1000000, quest = 4013 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
