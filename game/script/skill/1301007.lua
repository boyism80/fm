-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

function on_activated(me, skill)
	Skill.apply_hyper_body(me, skill)
end

function on_buff(me, skill)
	Skill.add_hyper_body_bonus(me, skill)
end

function on_unbuff(me, skill)
	Skill.remove_hyper_body_bonus(me, skill)
end
