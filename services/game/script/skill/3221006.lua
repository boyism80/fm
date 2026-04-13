-- Skill name (String.wz/Skill.img.xml): 블라인드

function on_activated_3221006(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.x == 0 then
		return
	end
	me:buff(skill, {[BuffFlag.Blind] = effect.x})
end

function on_unbuff_3221006(me, skill)
end
