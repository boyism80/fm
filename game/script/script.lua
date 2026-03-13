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
		me:chat(tostring(i))
		sleep(100)
	end
	if text ~= nil then
		me:chat(text)
	else
		me:chat('cancel')
	end
end

function on_script(me)
    local x, y = me:position()
    me:chat(string.format("position: %d, %d", x, y))

    me:class(412)
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
    me:hp(me:max_hp())
    me:max_mp(20000)
    me:mp(me:max_mp())
    me:mkitem('후루츠 대거')
    me:mkitem('소환의 돌', 200)
    me:mkitem('화비표창', 2000)
    me:mkitem('메바')
    me:base_dex(128)
    me:base_luk(128)
    me:level(200)
end

function on_attack(me, skill, damages, attack_info)
    local targets = damages_to_targets(damages)
    handle_combo_attack(me, targets, skill)
    handle_pickpocket(me, skill, damages)
    handle_ice_charge_freeze(me, damages)

    if attack_info and attack_info.ranged and attack_info.consume_slot and attack_info.consume_slot > 0 then
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

        me:chat(string.format("consume_slot: %d", attack_info.consume_slot))
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

        if not skip_consume then
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

            local ok = me:rmitem(InventoryType.Use, attack_info.consume_slot, count)
            if not ok then
            end
        end
    end
end

function on_damaged(me, attacker, skill, damage)
    local d = handle_magic_guard(me, attacker, skill, damage)
    return handle_meso_guard(me, attacker, skill, d)
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
    if me:class_of(Class.Beginner) or me:class_of(Class.Noblesse) or me:class_of(Class.Legend) then
        return math.random(6, 8)
    end
    if me:class_of(Class.Magician) then
        return math.random(10, 20)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) then
        return math.random(4, 7)
    end
    return math.random(50, 100)
end
