-- Skill name (String.wz/Skill.img.xml): 코마 : 도끼

function on_activated(me, skill, params)
	-- Coma (Axe) shares the same combo orb consumption behavior as Coma (Sword).
	consume_combo_orbs(me)
end
