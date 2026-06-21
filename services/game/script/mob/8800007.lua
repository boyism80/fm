-- Mob name (String.wz/Mob.img.xml): 자쿰팔5

local zakum = require("script/mob/zakum")

function on_mob_die_8800007(mob, attacker, map)
	zakum.on_arm_die(mob, attacker, map)
end
