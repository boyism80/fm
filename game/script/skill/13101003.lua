-- Skill name (String.wz/Skill.img.xml): 소울 에로우
function on_activated_13101003(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.SoulArrow, 1)
end
