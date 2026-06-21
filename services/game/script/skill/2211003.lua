-- Skill name (String.wz/Skill.img.xml): 썬더 스피어

local combat = require("script/lib/combat")

function on_attack_2211003(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_2211003(me, skill, params)
end
