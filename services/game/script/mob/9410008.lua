-- Mob name (String.wz/Mob.img.xml): 예티 인형자판기

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031352, min = 1, max = 1, chance = 80000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
