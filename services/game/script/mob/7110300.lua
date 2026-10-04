-- Mob name (String.wz/Mob.img.xml): D.로이

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031741, min = 1, max = 1, chance = 1000000, quest = 3345 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
