-- Skill name (String.wz/Skill.img.xml): 패닉

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.consume_combo_orbs(me)
	end
}
