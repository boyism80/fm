-- Skill name (String.wz/Skill.img.xml): 엘리멘탈 리셋

function on_activated_12101005(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.ElementReset, "x")
end

function on_unbuff_12101005(me, skill)
end
