-- Skill name (String.wz/Skill.img.xml): 마나 리플렉션

function on_activated_2121002(me, skill, params)
	apply_mana_reflection(me, skill)
end

function on_unbuff_2121002(me, skill)
end
