-- Skill name (String.wz/Skill.img.xml): 패닉 : 도끼

function on_activated(me, skill)
	-- Panic (Axe) shares the same combo orb consumption behavior as Panic (Sword).
	consume_combo_orbs(me)
end
