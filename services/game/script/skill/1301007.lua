-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

function on_activated_1301007(me, skill, params)
	apply_hyper_body(me, skill)
end

function on_buff_1301007(me, skill)
	add_hyper_body_bonus(me, skill)
end

function on_unbuff_1301007(me, skill)
	remove_hyper_body_bonus(me, skill)
end
