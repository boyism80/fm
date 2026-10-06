-- Skill name (String.wz/Skill.img.xml): 쇼다운

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_showdown(me, skill, damages)
	end
}
