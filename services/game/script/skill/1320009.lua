-- Skill name (String.wz/Skill.img.xml): 비홀더스 버프

return {
	on_summon_skill = function(me, summon, skill, params)
		if params == nil or params.buff_effect_index == nil then
			return
		end

		local buff_effect_index = params.buff_effect_index
		if buff_effect_index < 0 or buff_effect_index > 4 then
			return
		end

		local skill_effect = skill:effect()
		if skill_effect == nil then
			return
		end

		if skill_effect.time <= 0 then
			return
		end

		local buff_id = 2022125 + buff_effect_index
		local buff_flag = nil
		local buff_value = 0

		if buff_effect_index == 0 then
			buff_flag = BuffFlag.WeaponDef
			buff_value = skill_effect.pdd
		elseif buff_effect_index == 1 then
			buff_flag = BuffFlag.MagicDef
			buff_value = skill_effect.mdd
		elseif buff_effect_index == 2 then
			buff_flag = BuffFlag.Acc
			buff_value = skill_effect.acc
		elseif buff_effect_index == 3 then
			buff_flag = BuffFlag.Avoid
			buff_value = skill_effect.eva
		elseif buff_effect_index == 4 then
			buff_flag = BuffFlag.WeaponAtk
			buff_value = skill_effect.pad
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

		me:buff(consume_wz, skill_effect.time, {[buff_flag] = buff_value})
		summon:use_skill(buff_effect_index + 6)
	end
}
