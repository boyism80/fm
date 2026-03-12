-- Skill name (String.wz/Skill.img.xml): 소울 차지

function on_activated(me, skill)
	Skill.apply_wk_charge(me, skill)
end

