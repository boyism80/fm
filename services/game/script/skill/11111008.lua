-- Skill name (String.wz/Skill.img.xml): skill 11111008

local combat = require("script/lib/combat")

function on_attack_11111008(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end
