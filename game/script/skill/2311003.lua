-- Skill name (String.wz/Skill.img.xml): 홀리 심볼

function on_activated_2311003(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.HolySymbol, "x")
end

function on_buff_2311003(me, skill)
	add_holy_symbol_bonus(me, skill)
end

function on_unbuff_2311003(me, skill)
	remove_holy_symbol_bonus(me, skill)
end
