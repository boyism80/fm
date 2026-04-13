-- Skill name (String.wz/Skill.img.xml): 집중

function on_activated_3121008(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local vals = {}
	if effect.pad > 0 then
		vals[BuffFlag.WeaponAtk] = effect.pad
	end
	if effect.x > 0 then
		vals[BuffFlag.Concentrate] = effect.x
	end
	if next(vals) == nil then
		return
	end
	me:buff(skill, vals)
end
