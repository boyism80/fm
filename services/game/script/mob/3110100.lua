-- Mob name (String.wz/Mob.img.xml): 리게이터

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031405, min = 1, max = 1, chance = 10000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
