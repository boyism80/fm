-- Skill name (String.wz/Skill.img.xml): 닌자 스톰

local combat = require("script/lib/combat")

function on_attack_4121008(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_4121008(me, skill, params)
end
