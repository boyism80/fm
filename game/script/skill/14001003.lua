-- Skill name (String.wz/Skill.img.xml): 다크 사이트
function on_activated_14001003(me, skill, params)
	apply_buff_fixed(me, skill, BuffFlag.Darksight, 1)
end
