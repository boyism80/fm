-- Skill name (String.wz/Skill.img.xml): 에너지 드레인

local combat = require("script/lib/combat")

return {
	on_activating = function(me, skill, params)
		local energy = me:buff_value(BuffFlag.EnergyCharge) or 0
		return energy >= 10000
	end,

	on_attack = function(me, skill, damages)
		combat.drain_hp_from_damage(me, skill, damages)
	end
}
