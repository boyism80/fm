-- Skill name (String.wz/Skill.img.xml): 콜드 빔

local combat = require("script/lib/combat")

function on_attack_2201004(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Freeze)
end

function on_activated_2201004(me, skill, params)
end
