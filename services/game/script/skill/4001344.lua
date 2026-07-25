-- Skill name (String.wz/Skill.img.xml): 럭키 세븐

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_venom(me, damages, Skill.Venom)
	end
}
