-- Mob name (String.wz/Mob.img.xml): 구미호

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031793, min = 1, max = 1, chance = 1000000, quest = 3647 },
	{ item = 4031793, min = 1, max = 1, chance = 1000000, quest = 3647 },
	{ item = 4031793, min = 1, max = 1, chance = 1000000, quest = 3647 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
