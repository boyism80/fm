-- Skill name (String.wz/Skill.img.xml): 에너지 차지

function on_passive_5110001(me, skill)
	local value = me:buff_value(BuffFlag.EnergyCharge) or 0
	me:buff(skill, BuffFlag.EnergyCharge, value, { time = 0 })
end
