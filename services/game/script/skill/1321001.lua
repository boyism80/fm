-- Skill name (String.wz/Skill.img.xml): 몬스터 마그넷

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

return {
	on_attack = function(me, skill, damages)
		combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
	end,

	on_activating = function(me, skill, params)
		return true
	end,

	on_activated = function(me, skill, params)
		util.apply_monster_magnet(me, skill, params)
	end
}
