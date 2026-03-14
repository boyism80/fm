-- Skill name (String.wz/Skill.img.xml): 소드 부스터

function on_activated(me, skill, params)
	apply_booster(me, skill)
end

function on_unbuff(me, skill)
end

