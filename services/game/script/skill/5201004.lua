-- Skill name (String.wz/Skill.img.xml): 페이크 샷

local combat = require("script/lib/combat")

function on_attack_5201004(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end
