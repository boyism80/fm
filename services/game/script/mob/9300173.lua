-- Mob name (String.wz/Mob.img.xml): 중독된 스톤버그

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4001161, min = 1, max = 1, chance = 999999 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
