-- Skill name (String.wz/Skill.img.xml): 퍼펫

function on_activated_3211002(me, skill, params)
	apply_archer_puppet_activated(me, skill, params)
end

function on_buff_3211002(me, skill)
	apply_archer_puppet_buff(me, skill)
end

function on_unbuff_3211002(me, skill)
	apply_archer_puppet_unbuff(me, skill)
end
