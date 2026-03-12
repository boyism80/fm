-- Skill name (String.wz/Skill.img.xml): 홀리 심볼

function on_activated(me, skill)
	Skill.apply_holy_symbol(me, skill)
end

function on_buff(me, skill)
	Skill.add_holy_symbol_bonus(me, skill)
end

function on_unbuff(me, skill)
	Skill.remove_holy_symbol_bonus(me, skill)
end
