-- Skill name (String.wz/Skill.img.xml): 디스펠

function on_activated_2311001(me, skill, params)
	for_each_near_party_member(me, skill, function(ch)
		apply_dispel(ch, skill)
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end
