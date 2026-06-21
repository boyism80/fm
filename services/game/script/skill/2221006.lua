-- Skill name (String.wz/Skill.img.xml): 체인 라이트닝

local combat = require("script/lib/combat")

function on_attack_2221006(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end
