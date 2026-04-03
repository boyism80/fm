-- Skill name (String.wz/Skill.img.xml): 건 부스터

function on_activated_5201003(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_5201003(me, skill)
end

