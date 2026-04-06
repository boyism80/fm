-- Skill name (String.wz/Skill.img.xml): 포커스

function on_activated_3001003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local acc = tonumber(effect.acc) or 0
	local eva = tonumber(effect.eva) or 0
	if acc <= 0 and eva <= 0 then
		return
	end
	me:buff(skill, {
		[BuffFlag.Acc] = acc,
		[BuffFlag.Avoid] = eva,
	})
end
