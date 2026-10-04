-- Mob name (String.wz/Mob.img.xml): 마뇽

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031860, min = 1, max = 1, chance = 399999, quest = 6944 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
