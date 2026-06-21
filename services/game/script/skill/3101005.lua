-- Skill name (String.wz/Skill.img.xml): 에로우 봄 : 활

local combat = require("script/lib/combat")

function on_attack_3101005(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_3101005(me, skill, params)
end
