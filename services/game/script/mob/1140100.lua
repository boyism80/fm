-- Mob name (String.wz/Mob.img.xml): 고스텀프

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031773, min = 1, max = 1, chance = 199999 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
