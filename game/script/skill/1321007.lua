-- Skill name (String.wz/Skill.img.xml): 비홀더

function on_activated(me, skill, params)
	if me == nil or skill == nil then
		return
	end
	local effect = get_skill_effect(skill)
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
	local skill_id = wz.id
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	me:create_summon(skill_id, level, duration_ms, SummonMovementType.Follow, SummonType.Buff)
end
