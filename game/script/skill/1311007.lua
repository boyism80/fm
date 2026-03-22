-- Skill name (String.wz/Skill.img.xml): 파워 크래쉬

function on_activated_1311007(me, skill, params)
	for_each_mob_in_skill_area(me, skill, function(mob)
		mob:clear_status(MobStatus.WeaponAttackUp)
	end)
end
