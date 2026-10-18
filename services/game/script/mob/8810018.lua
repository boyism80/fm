-- Mob name (String.wz/Mob.img.xml): 혼테일

local horntail = require("script/lib/horntail")
local medal = require("script/lib/medal")

return {
	on_mob_die = function(mob, attacker, map)
		medal.count_boss_kill(map, 29501)
		horntail.clear_map_after(map, 3000)
	end
}
