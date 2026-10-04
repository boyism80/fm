-- Mob name (String.wz/Mob.img.xml): 미요캐츠

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031568, min = 1, max = 1, chance = 11000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
