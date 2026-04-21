-- Skill name (String.wz/Skill.img.xml): 힐 + 디스펠

function on_activated_9001000(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.prop <= 0 then
		return
	end
	if effect.prop < 100 and math.random(1, 100) > effect.prop then
		return
	end
	for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		apply_dispel(ch, skill)
		ch:hp(ch:max_hp())
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end
