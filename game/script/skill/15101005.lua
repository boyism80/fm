-- Skill name (String.wz/Skill.img.xml): 에너지 버스터

function on_activating_15101005(me, skill, params)
	local energy = me:buff_value(BuffFlag.EnergyCharge) or 0
	return energy >= 10000
end

function on_activated_15101005(me, skill, params)
end
