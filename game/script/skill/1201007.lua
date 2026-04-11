-- Skill name (String.wz/Skill.img.xml): 파워 가드
function on_activated_1201007(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Powerguard, "x")
end

function on_unbuff_1201007(me, skill)
end

