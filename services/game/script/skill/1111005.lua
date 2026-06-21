-- Skill name (String.wz/Skill.img.xml): 코마 : 검

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_attack_1111005(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_1111005(me, skill, params)
	-- Coma (Sword) consumes all combo orbs (leaving at least 1) after the finisher attack.
	util.consume_combo_orbs(me)
end
