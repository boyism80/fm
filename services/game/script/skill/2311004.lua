-- Skill name (String.wz/Skill.img.xml): 샤이닝 레이

local combat = require("script/lib/combat")

function on_attack_2311004(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end
