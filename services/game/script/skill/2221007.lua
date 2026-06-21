-- Skill name (String.wz/Skill.img.xml): 블리자드

local combat = require("script/lib/combat")

function on_attack_2221007(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Freeze)
end

function on_activated_2221007(me, skill, params)
end
