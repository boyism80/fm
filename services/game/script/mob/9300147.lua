-- Mob name (String.wz/Mob.img.xml): 호문쿨루

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001132, min = 1, max = 1, chance = 350000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
