-- Skill name (String.wz/Skill.img.xml): 디스오더

local function hits_have_damage(hits)
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

function on_attack_14001002(me, skill, damages)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	for mob, hits in pairs(damages) do
		if hits_have_damage(hits) then
			mob:buff({ [MobBuff.Watk] = effect.x, [MobBuff.Wdef] = effect.y }, effect.time, skill, me)
		end
	end
end
