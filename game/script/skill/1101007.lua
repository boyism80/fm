-- Skill name (String.wz/Skill.img.xml): 파워 가드
function on_activated(me, skill, params)
	apply_powerguard(me, skill)
end

function on_unbuff(me, skill)
end

