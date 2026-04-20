-- Skill name (String.wz/Skill.img.xml): 매직 크래쉬

function on_activated_1211009(me, skill, params)
	for_each_mob_in_skill_area(me, skill, function(mob)
		mob:clear_buffs(MobBuff.MagicDefenseUp)
	end)
end
