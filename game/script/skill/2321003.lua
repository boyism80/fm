-- Skill name (String.wz/Skill.img.xml): 바하뮤트

function on_activated_2321003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	me:buff(skill, BuffFlag.Summon, 1)
end

function on_buff_2321003(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	local wz = skill:wz()
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	me:create_summon(wz.id, level, effect.time, SummonMovementType.Follow, SummonType.Normal)
end

function on_unbuff_2321003(me, skill)
	local wz = skill:wz()
	me:remove_summon(wz.id)
end
