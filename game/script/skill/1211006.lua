-- Skill name (String.wz/Skill.img.xml): 아이스 차지 : 둔기

function on_activated_1211006(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

