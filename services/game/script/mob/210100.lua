-- Mob name (String.wz/Mob.img.xml): 슬라임

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4030000, min = 1, max = 1, chance = 200 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
