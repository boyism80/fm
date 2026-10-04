-- Mob name (String.wz/Mob.img.xml): 네오 휴로이드

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001133, min = 1, max = 1, chance = 100000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
