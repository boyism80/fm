-- Skill name (String.wz/Skill.img.xml): 샤이닝 레이

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	end
}
