-- Skill name (String.wz/Skill.img.xml): 비홀더스 버프

function on_activated_1320009(me, skill, params)
end

function on_summon_skill_1320009(me, summon, skill, params)
	if me == nil or summon == nil or skill == nil or params == nil or params.buff_effect_index == nil then
		return
	end

	local buff_effect_index = tonumber(params.buff_effect_index) or 0
	if buff_effect_index < 0 or buff_effect_index > 4 then
		return
	end

	local skill_effect = skill:effect()
	if skill_effect == nil then
		return
	end

	local duration = tonumber(skill_effect.time) or 0
	if duration <= 0 then
		return
	end

	local buff_id = 2022125 + buff_effect_index
	local buff_flag = nil
	local buff_value = 0

	if buff_effect_index == 0 then
		buff_flag = BuffFlag.WeaponDef
		buff_value = tonumber(skill_effect.pdd) or 0
	elseif buff_effect_index == 1 then
		buff_flag = BuffFlag.MagicDef
		buff_value = tonumber(skill_effect.mdd) or 0
	elseif buff_effect_index == 2 then
		buff_flag = BuffFlag.Acc
		buff_value = tonumber(skill_effect.acc) or 0
	elseif buff_effect_index == 3 then
		buff_flag = BuffFlag.Avoid
		buff_value = tonumber(skill_effect.eva) or 0
	elseif buff_effect_index == 4 then
		buff_flag = BuffFlag.WeaponAtk
		buff_value = tonumber(skill_effect.pad) or 0
	else
		return
	end

	if buff_value <= 0 or buff_flag == nil then
		return
	end

	local consume_wz = item_wz(buff_id)
	if consume_wz == nil then
		return
	end

	me:buff(consume_wz, duration, {[buff_flag] = buff_value})
	summon:use_skill(buff_effect_index + 6)
end
