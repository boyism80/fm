-- Skill name (String.wz/Skill.img.xml): 새크리파이스

function on_attack_1311005(me, skill, damages)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local total = 0
	local target = nil
	for mob, hits in pairs(damages) do
		target = mob
		for _, amount in ipairs(hits) do
			if amount and amount > 0 then
				total = total + amount
			end
		end
		break
	end
	if target == nil or total == 0 then
		return
	end
	local wz = target:wz()
	if wz and wz.boss then
		return
	end
	local hp_loss = math.floor(total * effect.x / 100)
	if hp_loss <= 0 then
		return
	end
	local cur = me:hp()
	local new_hp = math.max(cur - hp_loss, 1)
	me:hp(new_hp, true)
end

function on_activated_1311005(me, skill, params)
end
