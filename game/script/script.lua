-- Main script entry: loads sub modules via run_script, defines all on_* hooks called from Go.
run_script("script/script_common.lua")
run_script("script/script_heal_recovery.lua")
run_script("script/script_level_stat.lua")
run_script("script/script_combat.lua")
run_script("script/script_damage.lua")
run_script("script/script_equip.lua")

function on_start(me)
	local npc = 9001000
	-- me:dialog_list("안녕하세요", {"hello1", "hello2", "hello3"})
	-- me:dialog_accept('안녕하세요', true)
	-- me:dialog_yes_no('안녕하세요', true, true)
	local text = me:dialog_input(npc,'안녕하세요')
	me:dialog(npc, string.format("반갑습니다. %s님", me:name()))
	for i=1, 10 do
		sleep(100)
	end
	if text ~= nil then
	else
	end
end

function on_script(me)
    local x, y = me:position()

    me:class(212)
    local wz_skills = class_learnable_skill_wzs(me:class())
    for _, wz in pairs(wz_skills) do
        local skill = me:skill(wz.id)
        if skill == nil then
            skill = me:add_skill(wz.id)
        end
        if skill ~= nil then
            skill:level(wz.max_level, wz.master_level)
        end
    end
    me:max_hp(20000)
    me:hp(me:max_hp() * 0.5)
    me:max_mp(20000)
    me:mp(me:max_mp())
    local weapon = me:mkitem('폴암')
    if weapon ~= nil then
        me:equip(weapon)
    end
    me:base_str(80)
    me:base_dex(4)
    me:base_int(999)
    me:base_luk(999)
    me:level(200)
    me:map('오르비스탑입구')
end

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

function handle_mp_eater(me, damages)
    if damages == nil then
        return
    end

    local mp_eater_skill_id = Skill.MpEater
    if mp_eater_skill_id == nil then
        return
    end

    local mp_eater_skill = me:skill(mp_eater_skill_id)
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

function on_attack(me, skill, damages, attack_info)
    local targets = damages_to_targets(damages)
    handle_combo_attack(me, targets, skill)
    handle_pickpocket(me, skill, damages)
    handle_ice_charge_freeze(me, damages)
    handle_attack_consume_item(me, skill, attack_info)
    handle_mp_eater(me, damages)
end

function on_damaged(me, attacker, skill, damage)
    local d = handle_magic_guard(me, attacker, skill, damage)
    d = handle_meso_guard(me, attacker, skill, d)
    d = handle_power_guard(me, attacker, skill, d)
    return d
end

function on_blocked(me, attacker)
    if me == nil or attacker == nil then
        return
    end
    local class = me:class()
    local block_skill_id
    if class == Class.Hero then
        block_skill_id = Skill.Guardian
    elseif class == Class.Paladin then
        block_skill_id = Skill.Guardian1220006
    else
        return
    end
    local skill = me:skill(block_skill_id)
    if skill == nil then
        return
    end
    local effect = get_skill_effect(skill)
    if effect == nil then
        return
    end
    local prop = effect.prop or 0
    local duration_ms = effect.time or 0
    if prop <= 0 or duration_ms <= 0 then
        return
    end
    if math.random(1, 100) <= math.min(100, prop) then
        attacker:set_status(MobStatus.Stun, 1, duration_ms, skill, me)
    end
end

function on_equipment_changed(me, part, before, after)
    if part == EquipmentPart.Weapon then
        me:unbuff(BuffFlag.WkCharge)
    end
end

function on_level_up(me, old_level, new_level)
    if me == nil or new_level <= old_level then
        return
    end
    local diff = new_level - old_level
    me:ability_point(me:ability_point() + diff * 5, false)
    if not me:class_of(Class.Beginner) then
        me:skill_point(me:skill_point() + diff * 3, false)
    end
    local total_hp = 0
    local total_mp = 0
    for level = old_level + 1, new_level do
        total_hp = total_hp + (20 + level * 2)
        total_mp = total_mp + (10 + level)
    end
    local bonus_hp = 0
    local bonus_mp = 0
    if me:class_of(Class.Warrior) then
        local effect = get_skill_effect(me:skill(Skill.ImprovingMaxhpIncrease))
        if effect ~= nil and effect.x ~= nil and effect.x > 0 then
            bonus_hp = bonus_hp + diff * effect.x
        end
    end
    if me:class_of(Class.Magician) then
        local effect = get_skill_effect(me:skill(Skill.ImprovingMaxMpIncrease))
        if effect ~= nil and effect.x ~= nil and effect.x > 0 then
            bonus_mp = bonus_mp + diff * effect.x
        end
    end
    me:base_hp(me:base_hp() + total_hp + bonus_hp, false)
    me:base_mp(me:base_mp() + total_mp + bonus_mp, false)
    me:hp(me:max_hp(), false)
    me:mp(me:max_mp(), false)
    me:update_stats({
        STAT.Level,
        STAT.Exp,
        STAT.MaxHp,
        STAT.MaxMp,
        STAT.Hp,
        STAT.Mp,
        STAT.AvailableAP,
        STAT.AvailableSP,
    })
end

function on_ap_to_hp(me)
    local base = ap_to_hp_base(me)
    local bonus = 0
    if me:class_of(Class.Warrior) then
        local effect = get_skill_effect(me:skill(Skill.ImprovingMaxhpIncrease))
        if effect ~= nil and effect.y ~= nil and effect.y > 0 then
            bonus = effect.y
        end
    end
    return base + bonus
end

function on_ap_to_mp(me)
    local base = ap_to_mp_base(me)

    local bonus = 0
    if me:class_of(Class.Magician) then
        local effect = get_skill_effect(me:skill(Skill.ImprovingMaxMpIncrease))
        if effect ~= nil and effect.y ~= nil and effect.y > 0 then
            bonus = effect.y
        end
    end
    return base + bonus
end
