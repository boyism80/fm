-- Mob name (String.wz/Mob.img.xml): 바나드 그레이

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031926, min = 1, max = 1, chance = 100000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
