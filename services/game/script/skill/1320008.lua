-- Skill name (String.wz/Skill.img.xml): 비홀더스 힐링

return {
	on_summon_skill = function(me, summon, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end

		if effect.hp <= 0 then
			return
		end
		if skill:is_cooling() then
			return
		end

		skill:cooldown(effect.x * 1000)
		me:add_hp(effect.hp)
		summon:use_skill(5)
	end
}
