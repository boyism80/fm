-- Skill name (String.wz/Skill.img.xml): 코마 : 도끼

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	end,

	on_activated = function(me, skill, params)
		-- Coma (Axe) shares the same combo orb consumption behavior as Coma (Sword).
		util.consume_combo_orbs(me)
	end
}
