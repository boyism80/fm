-- Skill name (String.wz/Skill.img.xml): 힐

function on_activated_2301002(me, skill, params)
	local amount = get_heal_recovery_amount(me, skill)
	if amount == 0 then
		return
	end
	me:add_hp(amount)
end
