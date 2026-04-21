-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

function on_activated_1301007(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, {
			[BuffFlag.MaxHp] = effect.x,
			[BuffFlag.MaxMp] = effect.x,
		})
		if ch ~= me then
			ch:show_buff_effect(skill, 2)
		end
	end)
end

function on_buff_1301007(me, skill)
	add_hyper_body_bonus(me, skill)
end

function on_unbuff_1301007(me, skill)
	remove_hyper_body_bonus(me, skill)
end
