-- Skill name (String.wz/Skill.img.xml): 인피니티

function on_activated_2321004(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.Infinity, 1)
end
