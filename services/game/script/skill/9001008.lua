-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		me:buff(skill, {
			[BuffFlag.MaxHp] = effect.x,
			[BuffFlag.MaxMp] = effect.x,
		})
	end
}
