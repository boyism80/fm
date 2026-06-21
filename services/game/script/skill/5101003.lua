-- Skill name (String.wz/Skill.img.xml): 더블 어퍼

local combat = require("script/lib/combat")

function on_attack_5101003(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_5101003(me, skill, params)
end
