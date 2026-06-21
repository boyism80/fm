-- Skill name (String.wz/Skill.img.xml): 에너지 오브

function on_activating_5121002(me, skill, params)
	local energy = me:buff_value(BuffFlag.EnergyCharge) or 0
	return energy >= 10000
end
