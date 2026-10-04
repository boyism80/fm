-- Mob name (String.wz/Mob.img.xml): 차원의 블록퍼스

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001022, min = 1, max = 1, chance = 1000000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
