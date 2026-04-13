-- Skill name (String.wz/Skill.img.xml): 소울 차지

function on_activated_11111007(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end
