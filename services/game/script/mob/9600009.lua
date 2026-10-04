-- Mob name (String.wz/Mob.img.xml): 대왕지네

local drop = require("script/lib/drop")

local DROPS = {
	{ meso = true, min = 1553, max = 3156, chance = 1000000 },
	{ item = 4031227, min = 1, max = 1, chance = 1000000, quest = 4103 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
