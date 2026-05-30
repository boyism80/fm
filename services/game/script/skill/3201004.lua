-- Skill name (String.wz/Skill.img.xml): 소울 에로우 : 석궁

function on_activated_3201004(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.SoulArrow, 1)
end
