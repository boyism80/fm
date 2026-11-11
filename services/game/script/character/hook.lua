local combat = require("script/lib/combat")
local skill_lib = require("script/lib/skill")
local medal = require("script/lib/medal")
local class_lib = require("script/lib/class")

local function unbuff_stationary_summons_for_map(me, map_id)
    if me == nil then
        return
    end
    for _, s in ipairs(me:summons()) do
        if s == nil then
            goto continue
        end
        local sm = s:map()
        if sm == nil then
            goto continue
        end
        local swz = sm:wz()
        if swz == nil or swz:id() ~= map_id then
            goto continue
        end
        local mt = s:movement_type()
        if mt ~= SummonMovementType.Stationary
            and mt ~= SummonMovementType.WalkStationary
            and mt ~= SummonMovementType.CircleStationary then
            goto continue
        end
        me:unbuff(s:skill_id())
        ::continue::
    end
end

local MAX_HP_MP = 30000

-- Max HP and MP ranges and the tabs that grow by four slots, on the explorer first and second advancements the
-- original server rewarded.
local CLASS_CHANGE_GAINS = {
    [100] = { hp = { 200, 250 }, tabs = { InventoryType.Equipment, InventoryType.Use, InventoryType.Installation, InventoryType.Etc } },
    [110] = { hp = { 300, 350 }, tabs = { InventoryType.Use, InventoryType.Etc } },
    [120] = { hp = { 100, 150 }, tabs = { InventoryType.Use, InventoryType.Etc } },
    [130] = { hp = { 100, 150 }, tabs = { InventoryType.Use, InventoryType.Etc } },
    [200] = { mp = { 100, 150 } },
    [210] = { mp = { 450, 500 }, tabs = { InventoryType.Use } },
    [220] = { mp = { 450, 500 }, tabs = { InventoryType.Use } },
    [230] = { mp = { 450, 500 }, tabs = { InventoryType.Use } },
    [300] = { hp = { 100, 150 }, mp = { 30, 50 }, tabs = { InventoryType.Equipment, InventoryType.Use } },
    [310] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Etc } },
    [320] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Etc } },
    [400] = { hp = { 100, 150 }, mp = { 30, 50 }, tabs = { InventoryType.Equipment, InventoryType.Etc } },
    [500] = { hp = { 100, 150 }, mp = { 30, 50 }, tabs = { InventoryType.Equipment, InventoryType.Etc } },
    [410] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Use } },
    [420] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Use } },
    [510] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Use } },
    [520] = { hp = { 300, 350 }, mp = { 150, 200 }, tabs = { InventoryType.Use } },
}

-- HP and MP a level up adds, by class line, as the original server rolled them.
local function level_up_ranges(me)
    local r = class_lib.line(me).level_up
    return r[1], r[2], r[3], r[4]
end

local function skill_x(me, skill_id)
    local s = skill_id and me:skill(skill_id)
    local effect = s and s:effect()
    return effect and effect.x or 0
end

return {
    on_mob_kill = function(attacker, mobs)
        if attacker == nil or mobs == nil then
            return
        end
        local quest = attacker:quest(29400)
        if quest == nil or not quest:started() then
            return
        end
        local char_level = attacker:level()
        local count = 0
        for _, mob in ipairs(mobs) do
            if mob == nil then
                goto continue
            end
            local wz = mob:wz()
            if wz ~= nil and wz:level() ~= nil and wz:level() >= char_level then
                count = count + 1
            end
            ::continue::
        end
        if count <= 0 then
            return
        end
        local mon = tonumber(quest:record_ex("mon")) or 0
        quest:record_ex("mon", tostring(mon + count))
    end,

    on_map_enter = function(me, map)
        medal.check(me)
    end,

    on_medal_lost = function(me)
        medal.check(me)
    end,

    on_medal_rechecked = function(me)
        medal.take_back_lost(me)
    end,

    on_map_leave = function(me, map)
        if me == nil or map == nil then
            return
        end
        local mob = me:homing()
        if mob ~= nil then
            mob:homing(me, nil)
        end
        local mwz = map:wz()
        if mwz == nil then
            return
        end
        unbuff_stationary_summons_for_map(me, mwz:id())
    end,

    on_logout = function(me)
        if me == nil then
            return
        end
        me:unbuff(Skill.MysticDoor)
        me:unbuff(BuffFlag.Puppet)
    end,

    on_script = function(me)
        local x, y = me:position()

        me:class(232)
        local wz_skills = class_learnable_skill_wzs(me:class())
        for _, wz in pairs(wz_skills) do
            local skill = me:skill(wz:id())
            if skill == nil then
                skill = me:add_skill(wz:id())
            end
            if skill ~= nil then
                skill:level(wz:max_level(), wz:max_level())
            end
        end
        me:max_hp(20000)
        me:hp(me:max_hp())
        me:max_mp(20000)
        me:mp(me:max_mp())
        local weapon = me:mkitem('우드완드')
        if weapon ~= nil then
            me:equip(weapon)
        end
        me:mkitem('수비표창', 200)
        me:mkitem('불릿', 2000)
        me:mkitem('마법의돌', 200)
        me:base_str(80)
        me:base_dex(300)
        me:base_int(900)
        me:base_luk(900)
        me:level(200)
        me:map('커닝시티')
    end,

    on_damaged = function(me, attacker, skill, damage, params)
        if me:buff_value(BuffFlag.DivineBody) ~= nil then
            return 0
        end
        -- The reflection works on the whole hit, as the original's did, and the guards split what is left of it.
        local d = combat.reflect_incoming_damage(me, attacker, skill, damage, params)
        if (me:buff_value(BuffFlag.MagicGuard) or 0) > 0 then
            return skill_lib.absorb_magic_guard(me, d, params)
        end
        if (me:buff_value(BuffFlag.MesoGuard) or 0) > 0 then
            return skill_lib.absorb_meso_guard(me, d, params)
        end
        return skill_lib.take_hit(me, d, params)
    end,

    on_poison = function(mist, targets)
        if mist == nil or targets == nil then
            return
        end
        local effect = mist:effect()
        if effect == nil then
            return
        end
        local prop = effect.prop or 0
        if prop <= 0 then
            prop = 100
        end
        local time = effect.time or 0
        if time <= 0 then
            return
        end

        if mist:from_mob() then
            local skill_id = effect.skill_id
            local skill_level = effect.level or mist:level()
            local x = 30
            for _, ch in ipairs(targets) do
                if ch == nil then
                    goto continue_char
                end
                if math.random(1, 100) > prop then
                    goto continue_char
                end
                ch:debuff(DebuffFlag.Poison, time, x, skill_id, skill_level)
                ::continue_char::
            end
        else
            local wz = mist:wz()
            local level = mist:level()
            local multiplier = mist:poison_tick_multiplier()
            for _, mob in ipairs(targets) do
                if mob == nil then
                    goto continue_mob
                end
                if math.random(1, 100) > prop then
                    goto continue_mob
                end
                local value = combat.compute_poison_tick_damage_wz_level(wz, level, mob, multiplier)
                if value <= 0 then
                    goto continue_mob
                end
                mob:buff(MobBuff.Poison, value, time, mist, mist:causer())
                ::continue_mob::
            end
        end
    end,

    on_blocked = function(me, attacker)
        if me == nil or attacker == nil or me:equipped(EquipmentPart.Shield) == nil then
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
        local effect = skill:effect()
        local prop = effect.prop or 0
        local time = effect.time or 0
        if prop <= 0 or time <= 0 then
            return
        end
        if math.random(1, 100) <= math.min(100, prop) then
            attacker:buff(MobBuff.Stun, 1, time, skill, me)
        end
    end,

    on_equipment_changed = function(me, part, before, after)
        if part == EquipmentPart.Weapon then
            me:unbuff(BuffFlag.Booster)
            me:unbuff(BuffFlag.SpiritClaw)
            me:unbuff(BuffFlag.SoulArrow)
            me:unbuff(BuffFlag.WkCharge)
        end
    end,

    -- Every advancement fills HP and MP, after the gains of the advancements that have them.
    on_class_change = function(me, old_class, new_class)
        if me == nil then
            return
        end
        local gains = CLASS_CHANGE_GAINS[new_class] or {}
        local base_hp, base_mp = me:base_hp(), me:base_mp()
        if gains.hp ~= nil then
            base_hp = base_hp + math.random(gains.hp[1], gains.hp[2])
        end
        if gains.mp ~= nil then
            base_mp = base_mp + math.random(gains.mp[1], gains.mp[2])
        end
        me:base_hp(math.min(base_hp, MAX_HP_MP), false)
        me:base_mp(math.min(base_mp, MAX_HP_MP), false)
        me:hp(me:max_hp(), false)
        me:mp(me:max_mp(), false)
        me:update_stats({ STAT.MaxHp, STAT.MaxMp, STAT.Hp, STAT.Mp })
        for _, tab in ipairs(gains.tabs or {}) do
            me:expand_inventory(tab, 4)
        end
    end,

    on_level_up = function(me, old_level, new_level)
        if me == nil or new_level <= old_level then
            return
        end
        local diff = new_level - old_level
        me:ability_point(me:ability_point() + diff * 5, false)
        if not class_lib.is_beginner(me) then
            me:skill_point(me:skill_point() + diff * 3, false)
        end
        local hp_low, hp_high, mp_low, mp_high = level_up_ranges(me)
        local int_mp = math.floor(me:total_int() / 10)
        local total_hp = 0
        local total_mp = 0
        for _ = old_level + 1, new_level do
            total_hp = total_hp + math.random(hp_low, hp_high)
            total_mp = total_mp + math.random(mp_low, mp_high) + int_mp
        end
        local bonus_hp = diff * skill_x(me, class_lib.hp_skill(me))
        local bonus_mp = diff * skill_x(me, class_lib.mp_skill(me)) * 2
        me:base_hp(math.min(me:base_hp() + total_hp + bonus_hp, MAX_HP_MP), false)
        me:base_mp(math.min(me:base_mp() + total_mp + bonus_mp, MAX_HP_MP), false)
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
}
