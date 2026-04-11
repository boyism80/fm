-- Skill name (String.wz/Skill.img.xml): 아머 크래쉬

function on_activated_1111007(me, skill, params)
	for_each_mob_in_skill_area(me, skill, function(mob)
		mob:clear_buffs(MobBuff.WeaponDefenseUp)
	end)
end
