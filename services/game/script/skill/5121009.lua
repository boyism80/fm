-- Skill name (String.wz/Skill.img.xml): 윈드 부스터

function on_activated_5121009(me, skill, params)
	for_each_near_party_member(me, skill, function(ch)
		apply_buff_from_effect(ch, skill, BuffFlag.WindBooster, "x")
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end
