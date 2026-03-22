-- Skill name (String.wz/Skill.img.xml): 콤보 어택

function on_activated_1111002(me, skill, params)
	apply_combo(me, skill)
end

function on_unbuff_1111002(me, skill)
end
