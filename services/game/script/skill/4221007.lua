-- Skill name (String.wz/Skill.img.xml): 부메랑 스텝

local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
		combat.apply_venom(me, damages, Skill.Venom4220005)
	end
}
