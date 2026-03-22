run_script("script/script_common.lua")
run_script("script/script_combat.lua")

local function damages_has_positive_damage(hits)
	if hits == nil then
		return false
	end
	for _, amount in ipairs(hits) do
		if amount ~= nil and amount > 0 then
			return true
		end
	end
	return false
end

local function roll_percent(prob)
	local p = tonumber(prob) or 0
	if p <= 0 then
		return false
	end
	if p >= 100 then
		return true
	end
	return math.random(0, 99) < p
end

function handle_attack_consume_item(me, skill, attack_info)
	if attack_info == nil or not attack_info.ranged then
		return
	end

	if attack_info.consume_slot == nil or attack_info.consume_slot <= 0 then
		return
	end

	local weapon = me:equipped(EquipmentPart.Weapon)
	if weapon == nil then
		return
	end

	local wz_weapon = weapon:wz()
	if wz_weapon == nil then
		return
	end

	local weapon_type = wz_weapon:weapon_type()
	if weapon_type == nil then
		return
	end

	local consume_item = me:item(InventoryType.Use, attack_info.consume_slot)
	if consume_item == nil then
		return
	end

	local wz_consume = consume_item:wz()
	if wz_consume == nil then
		return
	end

	local consume_type = wz_consume:consume_type()
	if consume_type == nil then
		return
	end

	local valid_ammo = (weapon_type == WeaponType.Bow and consume_type == ConsumeType.ArrowBow)
		or (weapon_type == WeaponType.Crossbow and consume_type == ConsumeType.ArrowCrossBow)
		or (weapon_type == WeaponType.Claw and consume_type == ConsumeType.Shuriken)
		or (weapon_type == WeaponType.Gun and consume_type == ConsumeType.Bullet)

	if not valid_ammo then
		return
	end

	local skip_consume = false
	if weapon_type == WeaponType.Bow or weapon_type == WeaponType.Crossbow then
		if me:buff_value(BuffFlag.SoulArrow) ~= nil then
			skip_consume = true
		end
	elseif weapon_type == WeaponType.Claw then
		if me:buff_value(BuffFlag.SpiritClaw) ~= nil then
			skip_consume = true
		end
	end

	if skip_consume then
		return
	end

	local count = 1
	if skill ~= nil then
		local wz = skill:wz()
		if wz ~= nil and wz.effects ~= nil then
			local lv = skill:level()
			local effect = wz.effects[lv]
			if effect ~= nil then
				local bc = (effect.bullet_count and effect.bullet_count > 0) and effect.bullet_count or 1
				local ac = (effect.attack_count and effect.attack_count > 0) and effect.attack_count or 1
				count = math.max(bc, ac)
			end
		end
	end

	if me:buff_value(BuffFlag.ShadowPartner) ~= nil then
		count = count * 2
	end

	me:rmitem(InventoryType.Use, attack_info.consume_slot, count)
end

local function get_mp_eater_skill(me)
	if me == nil then
		return nil
	end
	if me:class_of(Class.FpWizard) then
		return me:skill(Skill.MpEater)
	end
	if me:class_of(Class.IlWizard) then
		return me:skill(Skill.MpEater2200000)
	end
	if me:class_of(Class.Cleric) then
		return me:skill(Skill.MpEater2300000)
	end
	if me:class_of(Class.Magician) then
		return me:skill(Skill.MpEater)
	end
	return nil
end

function handle_mp_eater(me, damages)
	if damages == nil then
		return
	end

	local mp_eater_skill = get_mp_eater_skill(me)
	if mp_eater_skill == nil then
		return
	end

	local effect = get_skill_effect(mp_eater_skill)
	if effect == nil then
		return
	end

	local chance = tonumber(effect.prop) or 0
	local absorb_percent = tonumber(effect.x) or 0

	if chance <= 0 or absorb_percent <= 0 then
		return
	end

	local total_absorb_mp = 0

	for mob, hits in pairs(damages) do
		if mob == nil then
			goto continue_mob
		end
		if hits == nil then
			goto continue_mob
		end
		if not damages_has_positive_damage(hits) then
			goto continue_mob
		end

		local wz = mob:wz()
		local is_boss = wz and wz.boss
		if is_boss then
			goto continue_mob
		end

		if not roll_percent(chance) then
			goto continue_mob
		end

		local mob_mp = mob:mp()
		if mob_mp == nil or mob_mp <= 0 then
			goto continue_mob
		end

		local absorb_mp = math.floor(mob:max_mp() * (absorb_percent / 100.0))
		if absorb_mp > mob_mp then
			absorb_mp = mob_mp
		end
		if absorb_mp <= 0 then
			goto continue_mob
		end

		mob:add_mp(-absorb_mp)
		total_absorb_mp = total_absorb_mp + absorb_mp

		::continue_mob::
	end

	if total_absorb_mp > 0 then
		me:add_mp(total_absorb_mp)
		me:show_buff_effect(mp_eater_skill)
	end
end

local function is_magic_attack_skill(skill)
	if skill == nil then
		return false
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return false
	end
	local mad = tonumber(effect.mad) or 0
	return mad > 0
end

local function apply_element_amplification_mp_cost(me, skill, mp_con)
	if mp_con <= 0 then
		return mp_con
	end
	if not is_magic_attack_skill(skill) then
		return mp_con
	end

	local amp
	if me:class_of(Class.FpMage) then
		amp = me:skill(Skill.ElementAmplification)
	elseif me:class_of(Class.IlMage) then
		amp = me:skill(Skill.ElementAmplification2210001)
	elseif me:class_of(Class.BlazeWizard3) then
		amp = me:skill(Skill.ElementAmplificationCygnus)
	end

	if amp == nil or amp:level() <= 0 then
		return mp_con
	end
	local amp_effect = get_skill_effect(amp)
	if amp_effect == nil then
		return mp_con
	end
	local mult = tonumber(amp_effect.x) or 100
	if mult <= 0 then
		return mp_con
	end
	return math.floor(mp_con * mult / 100)
end

function apply_default_skill_cost(me, skill)
	if me == nil or skill == nil then
		return true
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return true
	end
	local mp_con = tonumber(effect.mp_con) or 0
	mp_con = apply_element_amplification_mp_cost(me, skill, mp_con)
	if me:buff_value(BuffFlag.Infinity) ~= nil then
		mp_con = 0
	end
	local hp_con = tonumber(effect.hp_con) or 0
	if mp_con > 0 then
		local mp = me:mp()
		if mp < mp_con then
			return false
		end
		me:mp(mp - mp_con, false)
	end
	if hp_con > 0 then
		local hp = me:hp()
		if hp <= hp_con then
			return false
		end
		me:hp(hp - hp_con, false)
	end
	me:update_stats({ STAT.Hp, STAT.Mp })
	return true
end

function on_activating(me, skill, params)
	return apply_default_skill_cost(me, skill)
end

function on_activated(me, skill, params)
	return true
end

function on_attack(me, skill, damages, attack_info)
	local targets = damages_to_targets(damages)
	handle_combo_attack(me, targets, skill)
	handle_pickpocket(me, skill, damages)
	handle_ice_charge_freeze(me, damages)
	handle_attack_consume_item(me, skill, attack_info)
	handle_mp_eater(me, damages)
end
