-- Mob name (String.wz/Mob.img.xml): [★] 다크 클라크

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031992, min = 1, max = 1, chance = 30000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
