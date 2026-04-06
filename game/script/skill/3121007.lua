-- Skill name (String.wz/Skill.img.xml): 햄스트링

function on_activated_3121007(me, skill, params)
	if me == nil or skill == nil then
		return
	end

	local effect = skill:effect()
	if effect == nil then
		return
	end

	local slow = tonumber(effect.x) or 0
	if slow == 0 then
		return
	end

	me:buff(skill, {[BuffFlag.Hamstring] = slow})
end

function on_unbuff_3121007(me, skill)
end
