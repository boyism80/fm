-- Skill name (String.wz/Skill.img.xml): 퍼펫

function on_activated_13111004(me, skill, params)
	apply_archer_puppet_activated(me, skill, params)
end

function on_buff_13111004(me, skill)
	apply_archer_puppet_buff(me, skill)
end

function on_unbuff_13111004(me, skill)
	apply_archer_puppet_unbuff(me, skill)
end
