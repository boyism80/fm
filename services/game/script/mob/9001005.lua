-- Mob name (String.wz/Mob.img.xml): 파이렛 옥토

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031857, min = 1, max = 1, chance = 400000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
