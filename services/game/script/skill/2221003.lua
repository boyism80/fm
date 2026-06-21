-- Skill name (String.wz/Skill.img.xml): 아이스 데몬

local combat = require("script/lib/combat")

function on_attack_2221003(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Poison)
end

function on_activated_2221003(me, skill, params)
end
