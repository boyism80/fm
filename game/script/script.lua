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

    me:class(312)
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
    me:mkitem('활전용화살', 200)
    me:mkitem('석궁전용화살', 200)
    me:mkitem('석궁')
    local weapon = me:mkitem('워보우')
    if weapon ~= nil then
        me:equip(weapon)
    end
    me:base_str(80)
    me:base_dex(300)
    me:base_int(4)
    me:base_luk(4)
    me:level(200)
    me:map('오르비스탑입구')
end

function on_damaged(me, attacker, skill, damage, params)
    local d = handle_magic_guard(me, attacker, skill, damage)
    d = handle_meso_guard(me, attacker, skill, d)
    d = handle_reflect_damage(me, attacker, skill, d, params)
    return d
end

function on_poison(mist, mobs)
    if mist == nil or mobs == nil then
        return
    end
    local skill = mist:skill()
    if skill == nil then
        return
    end
    local multiplier = mist:poison_tick_multiplier() or 1.0
    if multiplier <= 0 then
        multiplier = 1.0
    end
    local effect = get_skill_effect(skill)
    if effect == nil then
        return
    end
    local prop = effect.prop or 0
    if prop <= 0 then
        prop = 100
    end
    local duration_ms = effect.time or 0
    if duration_ms <= 0 then
        return
    end
    for _, mob in ipairs(mobs) do
        if mob ~= nil and not mob:has_status(MobStatus.Poison) then
            if math.random(1, 100) <= prop then
                local value = compute_poison_tick_damage(skill, mob, multiplier)
                if value > 0 then
                    mob:set_status(MobStatus.Poison, value, duration_ms, skill, mist:owner_id())
                end
            end
        end
    end
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
