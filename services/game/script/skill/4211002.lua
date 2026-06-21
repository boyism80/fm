-- Skill name (String.wz/Skill.img.xml): 어썰터

local combat = require("script/lib/combat")

function on_attack_4211002(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	combat.apply_venom(me, damages, Skill.Venom4220005)
end
