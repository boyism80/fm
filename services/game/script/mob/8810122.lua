-- Mob name (String.wz/Mob.img.xml): 카오스 혼테일

local horntail = require("script/lib/horntail")

return {
	on_mob_die = function(mob, attacker, map)
		horntail.clear_map_after(map, 3000)
	end,

	on_revive = function(mob, map, x, y, revives)
	end
}
