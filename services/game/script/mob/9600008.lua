-- Mob name (String.wz/Mob.img.xml): 검은 양

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4000194, min = 1, max = 1, chance = 500000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
