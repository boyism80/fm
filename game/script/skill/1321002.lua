-- Skill name (String.wz/Skill.img.xml): 스탠스
function on_activated_1321002(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.Stance, "prop")
end
