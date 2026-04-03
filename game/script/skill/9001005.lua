-- Skill name (String.wz/Skill.img.xml): 부활

function on_activated_9001005(me, skill, params)
	apply_resurrection_in_range(me, skill)
end
