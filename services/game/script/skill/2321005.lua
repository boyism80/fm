-- Skill name (String.wz/Skill.img.xml): 홀리 실드

function on_activated_2321005(me, skill, params)
	for_each_near_party_member(me, skill, function(ch)
		apply_buff_from_effect(ch, skill, BuffFlag.HolyShield, "x")
	end)
end
