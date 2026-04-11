-- Skill name (String.wz/Skill.img.xml): 메디테이션

function on_activated_2201001(me, skill, params)
	apply_buff_from_effect(me, skill, BuffFlag.MagicAtk, "mad")
end
