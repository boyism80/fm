-- Skill name (String.wz/Skill.img.xml): 홀리 심볼

function on_activated_9001002(me, skill, params)
	apply_holy_symbol(me, skill)
end

function on_buff_9001002(me, skill)
	add_holy_symbol_bonus(me, skill)
end

function on_unbuff_9001002(me, skill)
	remove_holy_symbol_bonus(me, skill)
end
