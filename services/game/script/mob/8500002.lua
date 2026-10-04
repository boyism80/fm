-- Mob name (String.wz/Mob.img.xml): 파풀라투스

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031869, min = 1, max = 1, chance = 999999, quest = 6360 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
