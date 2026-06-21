-- Skill name (String.wz/Skill.img.xml): 샤우트

local combat = require("script/lib/combat")

function on_attack_1111008(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_1111008(me, skill, params)
end
