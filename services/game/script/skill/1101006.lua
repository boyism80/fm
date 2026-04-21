-- Skill name (String.wz/Skill.img.xml): 분노

function on_activated_1101006(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local vals = { [BuffFlag.WeaponAtk] = effect.pad }
	if effect.pdd ~= nil and effect.pdd > 0 then
		vals[BuffFlag.WeaponDef] = -effect.pdd
	end

	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, vals)
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end
