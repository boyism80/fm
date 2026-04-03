-- Skill name (String.wz/Skill.img.xml): 스피어 부스터

function on_activated_1301004(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_1301004(me, skill)
end

