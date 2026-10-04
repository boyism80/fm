-- Mob name (String.wz/Mob.img.xml): 호문스큘러

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031737, min = 1, max = 1, chance = 1000000, quest = 3343 },
	{ item = 4031740, min = 1, max = 1, chance = 1000000, quest = 3345 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
