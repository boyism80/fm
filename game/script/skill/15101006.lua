-- Skill name (String.wz/Skill.img.xml): 라이트닝 차지

function on_activated_15101006(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

