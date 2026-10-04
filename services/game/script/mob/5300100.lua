-- Mob name (String.wz/Mob.img.xml): 멜러디

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031925, min = 1, max = 1, chance = 60000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
