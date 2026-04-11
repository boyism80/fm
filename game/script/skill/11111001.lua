-- Skill name (String.wz/Skill.img.xml): 콤보 어택

function on_activated_11111001(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.Combo, 1)
end
