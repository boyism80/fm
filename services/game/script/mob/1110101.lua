-- Mob name (String.wz/Mob.img.xml): 다크 스텀프

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4030009, min = 1, max = 1, chance = 3000 },
	{ item = 4031773, min = 1, max = 1, chance = 199999 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
