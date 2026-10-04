-- Mob name (String.wz/Mob.img.xml): 트리플 루모

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031694, min = 1, max = 1, chance = 300000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
