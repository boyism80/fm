-- Skill name (String.wz/Skill.img.xml): 드레인

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.drain_hp_from_damage(me, skill, damages)
	end
}
