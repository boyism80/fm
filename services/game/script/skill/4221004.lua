-- Skill name (String.wz/Skill.img.xml): 닌자 앰부쉬

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.apply_ninja_ambush_skill(me, skill)
	end
}
