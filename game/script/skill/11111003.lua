-- Skill name (String.wz/Skill.img.xml): 코마

function on_activated(me, skill)
	-- Dawn Warrior Coma consumes combo orbs the same way as Hero Coma.
	consume_combo_orbs(me)
end
