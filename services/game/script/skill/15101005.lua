-- Skill name (String.wz/Skill.img.xml): 에너지 버스터

return {
	on_activating = function(me, skill, params)
		local energy = me:buff_value(BuffFlag.EnergyCharge) or 0
		return energy >= 10000
	end
}
