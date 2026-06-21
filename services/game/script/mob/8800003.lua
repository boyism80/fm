-- Mob name (String.wz/Mob.img.xml): 자쿰팔1

local zakum = require("script/lib/zakum")

function on_mob_die_8800003(mob, attacker, map)
	zakum.on_arm_die(mob, attacker, map)
end
