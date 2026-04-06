-- Heal-over-time and Endure (ladder/rope HP recovery) logic

local function get_hp_recover_check(me)
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

local function get_mp_recover_check(me)
    if me:class_of(Class.Magician) then
        local s = me:skill(Skill.ImprovingMpRecovery2000000)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
        return 0
    end
    if me:class_of(Class.Assassin) then
        local s = me:skill(Skill.Endure4100002)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    end
    if me:class_of(Class.Bandit) then
        local s = me:skill(Skill.Endure4200001)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    end
    if me:class_of(Class.Crusader) then
        local s = me:skill(Skill.ImprovingMpRecovery)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    end
    if me:class_of(Class.WhiteKnight) then
        local s = me:skill(Skill.ImprovingMpRecovery1210000)
        if s ~= nil then
            local effect = s:effect()
            if effect.mp > 0 then
                return effect.mp
            end
        end
    end
    return 0
end

function get_endure_hp_interval(me)
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

function get_heal_over_time_cap(me)
    local recovery_rate = 1.0
    local map = me:map()
    if map ~= nil then
        recovery_rate = map:recovery_rate()
    end
    local stance_mult = me:stance_of(Stance.Sit) and 1.5 or 1.0
    local chair = me:chair()

    local hp_check = get_hp_recover_check(me) + 10 * recovery_rate
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

    local mp_check = get_mp_recover_check(me) + 3 * recovery_rate
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
    }
end
