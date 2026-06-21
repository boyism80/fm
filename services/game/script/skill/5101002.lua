-- Skill name (String.wz/Skill.img.xml): 백스핀 블로우

local combat = require("script/lib/combat")

function on_attack_5101002(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_5101002(me, skill, params)
end
