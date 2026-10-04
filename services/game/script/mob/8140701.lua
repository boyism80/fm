-- Mob name (String.wz/Mob.img.xml): 레드 드래곤터틀

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031872, min = 1, max = 1, chance = 200000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
