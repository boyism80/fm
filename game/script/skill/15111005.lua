-- Skill name (String.wz/Skill.img.xml): 윈드 부스터

function on_activated_15111005(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.WindBooster, "x")
end
