-- Skill name (String.wz/Skill.img.xml): 에너지 드레인

local combat = require("script/lib/combat")

function on_activating_15111001(me, skill, params)
	local energy = me:buff_value(BuffFlag.EnergyCharge) or 0
	return energy >= 10000
end

function on_attack_15111001(me, skill, damages)
	combat.drain_hp_from_damage(me, skill, damages)
end

function on_activated_15111001(me, skill, params)
end
