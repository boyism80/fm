-- Mob name (String.wz/Mob.img.xml): 주황버섯

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4030001, min = 1, max = 1, chance = 300 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
