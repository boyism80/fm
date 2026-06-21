-- Mob name (String.wz/Mob.img.xml): 자쿰팔2

local zakum = require("script/mob/zakum")

function on_mob_die_8800004(mob, attacker, map)
	zakum.on_arm_die(mob, attacker, map)
end
