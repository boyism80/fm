-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

function on_activated(me, skill, params)
	apply_hyper_body(me, skill)
end

function on_buff(me, skill)
	add_hyper_body_bonus(me, skill)
end

function on_unbuff(me, skill)
	remove_hyper_body_bonus(me, skill)
end
