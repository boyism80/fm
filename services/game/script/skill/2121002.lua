-- Skill name (String.wz/Skill.img.xml): 마나 리플렉션

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		util.apply_mana_reflection(me, skill)
	end
}
