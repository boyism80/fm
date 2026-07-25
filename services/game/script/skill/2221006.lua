-- Skill name (String.wz/Skill.img.xml): 체인 라이트닝

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	end
}
