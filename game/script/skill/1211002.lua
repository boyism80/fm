-- Skill name (String.wz/Skill.img.xml): 차지 블로우
-- 차지된 상태에서만 사용 가능. prop% 확률로 스턴, 사용 후 차지 소모 (어드밴스드 차지 시 x% 확률로 유지).

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
	local mob_count = effect.mob_count or 6
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

	-- Consume WkCharge after use. If Advanced Charge (1220010) learned, effect.x is keep chance (%).
	local adv_charge = me:skill(Skill.AdvancedCharge)
	local keep_chance = 0
	if adv_charge ~= nil then
		local adv_effect = get_skill_effect(adv_charge)
		if adv_effect ~= nil and adv_effect.x ~= nil then
			keep_chance = adv_effect.x
		end
	end
	if math.random(1, 100) > keep_chance then
		me:unbuff(BuffFlag.WkCharge)
	end
end

function on_activated(me, skill, params)
end
