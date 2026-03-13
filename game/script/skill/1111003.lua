-- Skill name (String.wz/Skill.img.xml): 패닉 : 검

function on_activated(me, skill)
	-- Panic consumes all combo orbs (leaving at least 1) after the attack is fully resolved.
	consume_combo_orbs(me)
end
