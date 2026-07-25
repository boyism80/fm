-- Skill name (String.wz/Skill.img.xml): 쉐도우 웹

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.apply_shadow_web_skill(me, skill)
	end
}
