-- Skill name (String.wz/Skill.img.xml): 햄스트링

function on_activated_3121007(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	if effect.x == 0 then
		return
	end

	me:buff(skill, {[BuffFlag.Hamstring] = effect.x})
end

function on_unbuff_3121007(me, skill)
end
