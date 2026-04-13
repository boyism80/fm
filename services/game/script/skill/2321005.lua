-- Skill name (String.wz/Skill.img.xml): 홀리 실드

function on_activated_2321005(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.HolyShield, "x")
end
