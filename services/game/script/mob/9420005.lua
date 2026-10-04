-- Mob name (String.wz/Mob.img.xml): 흰 닭

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4000252, min = 1, max = 1, chance = 100000 },
	{ item = 4000253, min = 1, max = 1, chance = 500000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
