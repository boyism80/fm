-- Skill name (String.wz/Skill.img.xml): 보우 부스터

function on_activated_13101001(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_13101001(me, skill)
end

