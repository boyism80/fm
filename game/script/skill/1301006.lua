-- Skill name (String.wz/Skill.img.xml): 아이언 월

function on_activated_1301006(me, skill, params)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local pdd = effect.pdd or 0
	local mdd = effect.mdd or 0
	if pdd <= 0 and mdd <= 0 then
		return
	end
	me:buff(skill, {
		[BuffFlag.WeaponDef] = pdd,
		[BuffFlag.MagicDef] = mdd,
	})
	-- TODO: Apply same buff to party members in range (lt/rb) when party system exists
end
