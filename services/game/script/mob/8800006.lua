-- Mob name (String.wz/Mob.img.xml): 자쿰팔4

local zakum = require("script/lib/zakum")

return {
	on_mob_die = function(mob, attacker, map)
		zakum.on_arm_die(mob, attacker, map)
	end
}
