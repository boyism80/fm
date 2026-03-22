-- Skill name (String.wz/Skill.img.xml): 매직 아머
function on_activated_2001003(me, skill, params)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end

	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end

	local pdd = effect.pdd or 0
	if pdd <= 0 then
		return
	end

	me:buff(skill, BuffFlag.WeaponDef, pdd)
end

function on_unbuff_2001003(me, skill)
end
