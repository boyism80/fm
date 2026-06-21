-- Skill name (String.wz/Skill.img.xml): 매직 컴포지션

local combat = require("script/lib/combat")

function on_attack_2111006(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Poison)
end

function on_activated_2111006(me, skill, params)
end
