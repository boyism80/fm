-- Skill name (String.wz/Skill.img.xml): 비홀더

function on_activated(me, skill, params)
	if me == nil or skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end
	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	me:buff(skill, BuffFlag.Summon, 1)
end

function on_buff(me, skill)
	if me == nil or skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end
	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.id == nil then
		return
	end
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	me:create_summon(wz.id, level, duration_ms, SummonMovementType.Follow, SummonType.Buff)
end

function on_unbuff(me, skill)
	if me == nil or skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.id == nil then
		return
	end
	me:remove_summon(wz.id)
end
