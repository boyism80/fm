-- Mob name (String.wz/Mob.img.xml): 크로노스

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 2022355, min = 1, max = 1, chance = 200000 },
	{ item = 4031992, min = 1, max = 1, chance = 30000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
