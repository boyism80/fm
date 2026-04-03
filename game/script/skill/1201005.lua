-- Skill name (String.wz/Skill.img.xml): 메이스 부스터

function on_activated_1201005(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_1201005(me, skill)
end

