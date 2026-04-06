-- Skill name (String.wz/Skill.img.xml): 집중

function on_activated_3121008(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local pad = math.floor(tonumber(effect.pad) or 0)
	local x = math.floor(tonumber(effect.x) or 0)
	local vals = {}
	if pad > 0 then
		vals[BuffFlag.WeaponAtk] = pad
	end
	if x > 0 then
		vals[BuffFlag.Concentrate] = x
	end
	if next(vals) == nil then
		return
	end
	me:buff(skill, vals)
end
