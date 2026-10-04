-- Mob name (String.wz/Mob.img.xml): [★] 게이트키퍼

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031992, min = 1, max = 1, chance = 30000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
