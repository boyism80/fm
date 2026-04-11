-- Skill name (String.wz/Skill.img.xml): 대쉬

function on_activated_5001005(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	me:buff(skill, {
		[BuffFlag.DashSpeed] = effect.x,
		[BuffFlag.DashJump] = effect.y,
	})
end
