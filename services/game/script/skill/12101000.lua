-- Skill name (String.wz/Skill.img.xml): 메디테이션

function on_activated_12101000(me, skill, params)
	for_each_near_party_member(me, skill, function(ch)
		apply_buff_from_effect(ch, skill, BuffFlag.MagicAtk, "mad")
	end)
end
