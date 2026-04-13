-- Skill name (String.wz/Skill.img.xml): 홀리 차지 : 검

function on_activated_1221003(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

