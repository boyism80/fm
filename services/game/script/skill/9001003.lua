-- Skill name (String.wz/Skill.img.xml): 블레스

function on_activated_9001003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	me:buff(skill, {
		[BuffFlag.WeaponDef] = effect.pdd,
		[BuffFlag.MagicDef] = effect.mdd,
		[BuffFlag.Acc] = effect.acc,
		[BuffFlag.Avoid] = effect.eva,
	})
end
