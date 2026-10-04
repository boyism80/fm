-- Mob name (String.wz/Mob.img.xml): 옥토퍼스

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4030010, min = 1, max = 1, chance = 600 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
