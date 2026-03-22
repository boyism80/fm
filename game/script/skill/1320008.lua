-- Skill name (String.wz/Skill.img.xml): 비홀더스 힐링

function on_activated_1320008(me, skill, params)
end

function on_summon_skill_1320008(me, summon, skill, params)
	if me == nil or summon == nil or skill == nil then
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

	local hp = tonumber(effect.hp) or 0
	if hp <= 0 then
		return
	end

	me:add_hp(hp)
	summon:use_skill(5)
end
