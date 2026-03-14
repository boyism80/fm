-- Skill name (String.wz/Skill.img.xml): 샤우트

function on_attack(me, skill, damages)
	if damages == nil or skill == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = effect.prop or 0
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	local mob_count = effect.mobCount or 6
	local applied = 0
	for mob, hits in pairs(damages) do
		if applied >= mob_count then
			break
		end
		if mob and hits then
			local total = 0
			for _, amount in ipairs(hits) do
				if amount and amount > 0 then
					total = total + amount
				end
			end
			if total > 0 and math.random(1, 100) <= prop then
				mob:set_status(MobStatus.Stun, 1, duration_ms, skill, me)
				applied = applied + 1
			end
		end
	end
end

function on_activated(me, skill, params)
end
