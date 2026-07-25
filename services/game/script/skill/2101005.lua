-- Skill name (String.wz/Skill.img.xml): 포이즌 브레스

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Poison)
	end
}
