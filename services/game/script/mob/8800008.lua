run_script("script/mob/zakum_arm.lua")

function on_mob_die_8800008(mob, attacker, map)
	on_zakum_arm_die(mob, attacker, map)
end
