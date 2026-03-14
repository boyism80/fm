-- Skill name (String.wz/Skill.img.xml): 메소 업

function on_activated(me, skill, params)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 100
	me:buff(skill, BuffFlag.MesoUp, percent)
end

function on_buff(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 100
	local current = me:bonus_meso_multiplier()
	if current <= 0 then
		current = 100
	end
	me:bonus_meso_multiplier(current + (percent - 100))
end

function on_unbuff(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 100
	local current = me:bonus_meso_multiplier()
	me:bonus_meso_multiplier(current - (percent - 100))
end
