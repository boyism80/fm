-- Skill name (String.wz/Skill.img.xml): 블레스

function on_activated_2301004(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local vals = {
		[BuffFlag.WeaponDef] = effect.pdd,
		[BuffFlag.MagicDef] = effect.mdd,
		[BuffFlag.Acc] = effect.acc,
		[BuffFlag.Avoid] = effect.eva,
	}
	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, vals)
	end)
end
