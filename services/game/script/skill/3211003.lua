-- Skill name (String.wz/Skill.img.xml): 아이스 샷

local combat = require("script/lib/combat")

function on_attack_3211003(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Freeze)
end

function on_activated_3211003(me, skill, params)
end
