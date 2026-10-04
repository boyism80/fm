-- Mob name (String.wz/Mob.img.xml): 그린 버블티

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4000256, min = 1, max = 1, chance = 500000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
