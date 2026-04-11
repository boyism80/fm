-- Skill name (String.wz/Skill.img.xml): 타임 리프

function on_activated_5121010(me, skill, params)
	for_each_character_in_skill_area(me, skill, function(ch)
		local skills = ch:skills()
		if skills == nil then
			return
		end
		for id, skill in pairs(skills) do
			if id ~= Skill.TimeLeap and skill ~= nil then
				skill:cooldown(0)
			end
		end
	end)
end
