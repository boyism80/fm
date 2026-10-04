-- Mob name (String.wz/Mob.img.xml): 돼지

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4030011, min = 1, max = 1, chance = 600 },
	{ item = 4031846, min = 1, max = 1, chance = 100000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
