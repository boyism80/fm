-- Skill name (String.wz/Skill.img.xml): 블라인드

function on_activated_3221006(me, skill, params)
	if me == nil or skill == nil then
		return
	end
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local acc = tonumber(effect.x) or 0
	if acc == 0 then
		return
	end
	me:buff(skill, {[BuffFlag.Blind] = acc})
end

function on_unbuff_3221006(me, skill)
end
