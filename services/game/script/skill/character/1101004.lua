-- Skill name (String.wz/Skill.img.xml): 소드 부스터

function on_activated_1101004(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_1101004(me, skill)
end

