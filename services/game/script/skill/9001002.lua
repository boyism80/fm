-- Skill name (String.wz/Skill.img.xml): 홀리 심볼

function on_activated_9001002(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		ch:buff(skill, BuffFlag.HolySymbol, effect.x)
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end
