-- Skill name (String.wz/Skill.img.xml): 홀리 심볼

function on_activated_9001002(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.HolySymbol, "x")
end
