-- Skill name (String.wz/Skill.img.xml): 자벨린 부스터

function on_activated_14101002(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_14101002(me, skill)
end

