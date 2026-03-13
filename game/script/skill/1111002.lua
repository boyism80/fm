-- Skill name (String.wz/Skill.img.xml): 콤보 어택

function on_activated(me, skill)
	apply_combo(me, skill)
end

function on_unbuff(me, skill)
end
