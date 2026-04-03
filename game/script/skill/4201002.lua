-- Skill name (String.wz/Skill.img.xml): 대거 부스터

function on_activated_4201002(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_4201002(me, skill)
end

