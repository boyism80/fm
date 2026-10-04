-- Mob name (String.wz/Mob.img.xml): 폐쇄된 연구실의 호문

local drop = require("script/lib/drop")

local DROPS = {
	{ item = 4031698, min = 1, max = 1, chance = 199999 },
}

return {
	on_mob_die = function(mob, attacker, map)
		drop.from_mob(DROPS, mob, attacker, map)
	end
}
