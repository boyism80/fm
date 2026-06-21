-- Skill name (String.wz/Skill.img.xml): 매직 아머

function on_activated_2001003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	if effect.pdd <= 0 then
		return
	end

	me:buff(skill, BuffFlag.WeaponDef, effect.pdd)
end
