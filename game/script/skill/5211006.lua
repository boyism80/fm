-- Skill name (String.wz/Skill.img.xml): 호밍

local function has_positive_hit(hits)
	if hits == nil then
		return false
	end
	for i = 1, #hits do
		if (hits[i] or 0) > 0 then
			return true
		end
	end
	return false
end

function on_attack_5211006(me, skill, damages)
	if skill == nil or me == nil or damages == nil then
		return
	end
	for mob, hits in pairs(damages) do
		if has_positive_hit(hits) then
			mob:homing(me, skill)
			return
		end
	end
end

function on_activated_5211006(me, skill, params)
end
