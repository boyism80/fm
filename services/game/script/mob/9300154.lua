-- Mob name (String.wz/Mob.img.xml): 실험용 네오 휴로이드

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031780, min = 1, max = 1, chance = 200000 },
	{ item = 4031781, min = 1, max = 1, chance = 200000 },
	{ item = 4031782, min = 1, max = 1, chance = 200000 },
	{ item = 4031783, min = 1, max = 1, chance = 200000 },
	{ item = 4031784, min = 1, max = 1, chance = 200000 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
