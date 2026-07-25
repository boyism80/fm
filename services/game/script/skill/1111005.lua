-- Skill name (String.wz/Skill.img.xml): 코마 : 검

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	end,

	on_activated = function(me, skill, params)
		-- Coma (Sword) consumes all combo orbs (leaving at least 1) after the finisher attack.
		util.consume_combo_orbs(me)
	end
}
