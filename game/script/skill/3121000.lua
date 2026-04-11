-- Skill name (String.wz/Skill.img.xml): 메이플 용사

function on_activated_3121000(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.MapleWarrior, "x")
end
