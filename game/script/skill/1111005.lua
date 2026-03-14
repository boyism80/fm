-- Skill name (String.wz/Skill.img.xml): 코마 : 검

function on_activated(me, skill, params)
	-- Coma (Sword) consumes all combo orbs (leaving at least 1) after the finisher attack.
	consume_combo_orbs(me)
end
