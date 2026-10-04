-- Mob name (String.wz/Mob.img.xml): 폭주족 원숭이

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4000202, min = 1, max = 1, chance = 500000 },
	{ item = 4031296, min = 1, max = 1, chance = 50000, quest = 4010 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
