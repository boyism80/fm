local class_lib = require("script/lib/class")

local SAUNA_ROBES = {
    [1050018] = true,
    [1051017] = true,
}

local function hp_recover_check(me)
    if me:class_of(Class.Warrior) then
        local s = me:skill(Skill.ImprovingHpRecovery)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    if me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.hp > 0 then
                return effect.hp
            end
        end
    end
    return 0
end

local function mp_recover_check(me)
    if me:class_of(Class.Magician) then
        -- WZ gives this skill no mp; the amount is skill level x character level / 10, rounded.
        local s = me:skill(Skill.ImprovingMpRecovery2000000)
        if s == nil then
            return 0
        end
        return math.floor(s:level() * me:level() * 0.1 + 0.5)
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.Crusader) then
        local s = me:skill(Skill.ImprovingMpRecovery)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.WhiteKnight) then
        local s = me:skill(Skill.ImprovingMpRecovery1210000)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    elseif me:class_of(Class.DawnWarrior3) then
        local s = me:skill(Skill.ImprovingMpRecoveryCygnus)
        if s ~= nil then
            local effect = s:effect()
            if effect ~= nil and effect.mp > 0 then
                return effect.mp
            end
        end
    end
    return 0
end

local function ap_to_hp_base(me)
    if class_lib.is_beginner(me) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) or me:class_of(Class.DawnWarrior1) then
        return math.random(12, 20)
    end
    if me:class_of(Class.Magician) or me:class_of(Class.BlazeWizard1) then
        return math.random(6, 11)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) or me:class_of(Class.WindArcher1) or me:class_of(Class.NightWalker1) then
        return math.random(14, 18)
    end
    if me:class_of(Class.Pirate) or me:class_of(Class.ThunderBreaker1) then
        return math.random(16, 20)
    end
    if me:class_of(Class.Legend) then
        return math.random(34, 38)
    end
    return math.random(50, 100)
end

local function ap_to_mp_base(me)
    if class_lib.is_beginner(me) then
        return math.random(6, 8)
    end
    if me:class_of(Class.Magician) or me:class_of(Class.BlazeWizard1) then
        return math.random(10, 20)
    end
    if me:class_of(Class.Bowman) or me:class_of(Class.Thief) or me:class_of(Class.WindArcher1) or me:class_of(Class.NightWalker1) then
        return math.random(8, 12)
    end
    if me:class_of(Class.Warrior) or me:class_of(Class.DawnWarrior1) then
        return math.random(4, 7)
    end
    if me:class_of(Class.Pirate) or me:class_of(Class.ThunderBreaker1) then
        return math.random(10, 12)
    end
    if me:class_of(Class.Legend) then
        return math.random(4, 6)
    end
    return math.random(50, 100)
end

local function skill_y(me, skill_id)
    local s = me:skill(skill_id)
    if s == nil then
        return 0
    end
    local effect = s:effect()
    if effect == nil then
        return 0
    end
    return effect.y
end

local function hp_skill_bonus(me)
    if me:class_of(Class.Warrior) then
        return skill_y(me, Skill.ImprovingMaxhpIncrease)
    end
    if me:class_of(Class.Pirate) then
        return skill_y(me, Skill.HpIncrease)
    end
    if me:class_of(Class.ThunderBreaker1) then
        return skill_y(me, Skill.HpIncreaseCygnus)
    end
    if me:class_of(Class.DawnWarrior1) then
        return skill_y(me, Skill.ImprovingMaxhpIncreaseCygnus)
    end
    return 0
end

local function mp_skill_bonus(me)
    if me:class_of(Class.Magician) then
        return skill_y(me, Skill.ImprovingMaxMpIncrease)
    end
    if me:class_of(Class.BlazeWizard1) then
        return skill_y(me, Skill.ImprovingMaxMpIncreaseCygnus)
    end
    return 0
end

local function endure_hp_interval(me)
    if me:class_of(Class.Warrior) then
        local s = me:skill(Skill.Endure)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    if me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.time > 0 then
                return effect.time / 1000
            end
        end
    end
    return 0
end

local function return_scroll_continent_allowed(from_map_id, to_map_id)
    local from = math.floor(from_map_id / 10000000)
    local to = math.floor(to_map_id / 10000000)
    if to == 60 or from == 61 then
        return true
    end
    if to == 21 or from == 20 then
        return true
    end
    if to == 12 or from == 10 then
        return true
    end
    if to == 10 or from == 12 then
        return true
    end
    return from == to
end

local function teleport_stone_blocked(map_id)
    if map_id < 100000000 then
        return true
    end
    return math.floor(map_id / 1000000) % 100 == 9
end

local function teleport_stone_region_allowed(item_id, from_map_id, to_map_id, by_name)
    if item_id ~= 5040000 or by_name == false then
        return true
    end
    local region = math.floor(to_map_id / 10000000)
    if region >= 20 and region <= 23 then
        return math.floor(from_map_id / 100000000) == math.floor(to_map_id / 100000000)
    end
    if region >= 24 and region <= 26 then
        return math.floor(from_map_id / 10000000) == region
    end
    return true
end

return {
    can_use_return_scroll = function(me, item_id, move_to)
        local map = me:map()
        if map == nil then
            return false
        end
        local wz = map:wz()
        if wz == nil or wz:blocks_potion_use() then
            return false
        end
        local current = map:template_id()
        if current == 749040100 then
            return false
        end
        if move_to == current and item_id ~= 2031010 and item_id ~= 2030021 then
            return false
        end
        if move_to == 999999999 then
            return true
        end
        return return_scroll_continent_allowed(current, move_to)
    end,

    can_register_teleport_stone = function(me, map_id)
        return teleport_stone_blocked(map_id) == false
    end,

    can_teleport_stone = function(me, item_id, from_map_id, to_map_id, by_name)
        if teleport_stone_blocked(from_map_id) or teleport_stone_blocked(to_map_id) then
            return false
        end
        return teleport_stone_region_allowed(item_id, from_map_id, to_map_id, by_name)
    end,


    get_heal_over_time_cap = function(me)
        local recovery_rate = 1.0
        local map = me:map()
        if map ~= nil then
            recovery_rate = map:recovery_rate()
        end
        local overall = me:equipped(EquipmentPart.Top)
        if recovery_rate ~= 1.0 and overall ~= nil and SAUNA_ROBES[overall:wz():id()] then
            recovery_rate = recovery_rate * 1.5
        end
        local stance_mult = me:stance_of(Stance.Sit) and 1.5 or 1.0
        local chair = me:chair()

        local hp_check = hp_recover_check(me) + 10 * recovery_rate
        hp_check = hp_check * stance_mult
        if chair >= 3000000 then
            if chair == 3010000 then
                hp_check = hp_check + 50
            elseif chair == 3010001 then
                hp_check = hp_check + 35
            elseif chair == 3010007 then
                hp_check = hp_check + 60
            elseif chair == 3010009 then
                hp_check = hp_check + 20
            end
        end

        local mp_check = mp_recover_check(me) + 3 * recovery_rate
        mp_check = mp_check * stance_mult
        if chair >= 3000000 then
            if chair == 3010008 then
                mp_check = mp_check + 60
            elseif chair == 3010009 then
                mp_check = mp_check + 20
            end
        end

        return {
            max_hp = math.floor(hp_check + 0.5),
            max_mp = math.floor(mp_check + 0.5),
            endure_hp_interval = endure_hp_interval(me),
        }
    end,

    get_ap_to_hp = function(me)
        return ap_to_hp_base(me) + hp_skill_bonus(me)
    end,

    get_ap_to_mp = function(me)
        return ap_to_mp_base(me) + mp_skill_bonus(me)
    end,

    -- An AP reset takes back a fixed amount of HP or MP, as the original server did.
    get_ap_reset_hp = function(me)
        local bonus = hp_skill_bonus(me)
        if class_lib.is_beginner(me) then
            return 12
        end
        if me:class_of(Class.Warrior) or me:class_of(Class.DawnWarrior1) then
            return 24 + bonus
        end
        if me:class_of(Class.Magician) or me:class_of(Class.BlazeWizard1) then
            return 10
        end
        if me:class_of(Class.Bowman) or me:class_of(Class.Thief) or me:class_of(Class.WindArcher1) or me:class_of(Class.NightWalker1) then
            return 15
        end
        if me:class_of(Class.Pirate) or me:class_of(Class.ThunderBreaker1) then
            return 20 + bonus
        end
        return 50
    end,

    get_ap_reset_mp = function(me)
        local bonus = mp_skill_bonus(me)
        if class_lib.is_beginner(me) then
            return 8
        end
        if me:class_of(Class.Warrior) or me:class_of(Class.DawnWarrior1) then
            return 4
        end
        if me:class_of(Class.Magician) or me:class_of(Class.BlazeWizard1) then
            return 20 + bonus
        end
        if me:class_of(Class.Bowman) or me:class_of(Class.Thief) or me:class_of(Class.Pirate) or me:class_of(Class.WindArcher1) or me:class_of(Class.NightWalker1) or me:class_of(Class.ThunderBreaker1) then
            return 10
        end
        return 50
    end,

    get_recharge_capacity = function(me, capacity)
        local s
        if me:class_of(Class.Assassin) then
            s = me:skill(Skill.ClawMastery)
        elseif me:class_of(Class.NightWalker2) then
            s = me:skill(Skill.ClawMasteryCygnus)
        elseif me:class_of(Class.Gunslinger) then
            s = me:skill(Skill.GunMastery)
        end
        local effect = s ~= nil and s:effect() or nil
        if effect == nil then
            return capacity
        end
        return capacity + (effect.y or 0)
    end
}
