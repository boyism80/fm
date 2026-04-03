-- Skill name (String.wz/Skill.img.xml): 블레스

function on_activated_2301004(me, skill, params)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end

	local pdd = effect.pdd or 0
	local mdd = effect.mdd or 0
	local acc = effect.acc or 0
	local eva = effect.eva or 0

	me:buff(skill, {
		[BuffFlag.WeaponDef] = pdd,
		[BuffFlag.MagicDef] = mdd,
		[BuffFlag.Acc] = acc,
		[BuffFlag.Avoid] = eva,
	})
end
