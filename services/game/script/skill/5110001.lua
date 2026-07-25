-- Skill name (String.wz/Skill.img.xml): 에너지 차지

return {
	on_passive = function(me, skill)
		local value = me:buff_value(BuffFlag.EnergyCharge) or 0
		me:buff(skill, BuffFlag.EnergyCharge, value, { time = 0 })
	end
}
