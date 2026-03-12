-- Skill name (String.wz/Skill.img.xml): 스탠스

function on_activated(me, skill)
	Skill.apply_stance(me, skill)
end
