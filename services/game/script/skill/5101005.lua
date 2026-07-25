-- Skill name (String.wz/Skill.img.xml): MP 리커버리

return {
	on_activating = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return false
		end

		local decrease = math.floor(me:max_hp() / 100) * 10
		if me:hp() <= decrease then
			return false
		end
		me:hp(me:hp() - decrease, false)

		local mp_gain = math.floor(decrease / 100) * effect.y
		local new_mp = math.min(me:max_mp(), me:mp() + mp_gain)
		me:mp(new_mp, false)
		me:update_stats({STAT.Hp, STAT.Mp})
		return true
	end
}
