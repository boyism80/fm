-- Mob name (String.wz/Mob.img.xml): 자쿰팔6

local zakum = require("script/lib/zakum")

function on_mob_die_8800008(mob, attacker, map)
	zakum.on_arm_die(mob, attacker, map)
end
