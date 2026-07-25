-- Skill name (String.wz/Skill.img.xml): 패닉 : 도끼

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		-- Panic (Axe) shares the same combo orb consumption behavior as Panic (Sword).
		util.consume_combo_orbs(me)
	end
}
