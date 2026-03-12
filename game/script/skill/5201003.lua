-- Skill name (String.wz/Skill.img.xml): 건 부스터

function on_activated(me, skill)
	Skill.apply_booster(me, skill)
end

function on_unbuff(me, skill)
end

