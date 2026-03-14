-- Attack handlers: combo, pickpocket, ice charge, ammo consume

local PICKPOCKET_SKILL_IDS = { [0] = true }
do
    local function add(id)
        if id ~= nil then
            PICKPOCKET_SKILL_IDS[id] = true
        end
    end
    add(Skill.DoubleStab)
    add(Skill.SavageBlow)
    add(Skill.Assaulter)
    add(Skill.BandOfThieves)
    add(Skill.Showdown4221003)
    add(Skill.BoomerangStep)
end

function handle_combo_attack(me, targets, skill)
    if targets == nil or #targets == 0 then
        return
    end

    local current = me:buff_value(BuffFlag.Combo)
    if current == nil then
        return
    end

    local shout_hero = Skill.Shout
    local shout_dw = Skill.DawnWarriorShout
    if skill ~= nil then
        local wz = skill:wz()
        if wz ~= nil then
            local sid = wz.id
            if sid == shout_hero or sid == shout_dw then
                return
            end
        end
    end

    local combo_skill = nil
    local adv_skill = nil
    if me:class_of(Class.Crusader) then
        combo_skill = Skill.ComboAttack
        adv_skill = Skill.AdvancedCombo
    elseif me:class_of(Class.DawnWarrior3) then
        combo_skill = Skill.DawnWarriorComboAttack
        adv_skill = Skill.DawnWarriorAdvancedCombo
    else
        return
    end

    local combo = me:skill(combo_skill)
    local adv = me:skill(adv_skill)
    if combo == nil then
        return
    end

    local ceffect_skill = combo
    if adv ~= nil and adv:level() > 0 then
        ceffect_skill = adv
    end

    local ceffect = nil
    local wz = ceffect_skill:wz()
    if wz ~= nil and wz.effects ~= nil then
        ceffect = wz.effects[ceffect_skill:level()]
    end
    if ceffect == nil then
        return
    end

    local max_orbs = (ceffect.x or 0) + 1
    if max_orbs <= 0 then
        return
    end
    if current >= max_orbs then
        return
    end

    local new_orbs = current + 1
    if adv ~= nil and adv:level() > 0 then
        local prop = ceffect.prop or 0
        if prop > 0 and math.random(1, 100) <= prop and new_orbs < max_orbs then
            new_orbs = new_orbs + 1
        end
    end

    if new_orbs > max_orbs then
        new_orbs = max_orbs
    end
    me:buff_value(BuffFlag.Combo, new_orbs)
end

function handle_pickpocket(me, skill, damages)
    if damages == nil then
        return
    end
    local maxmeso = me:buff_value(BuffFlag.Pickpocket)
    if maxmeso == nil or maxmeso < 1 then
        return
    end
    local wz = skill and skill:wz()
    local skill_id = (wz and wz.id) or 0
    if not PICKPOCKET_SKILL_IDS[skill_id] then
        return
    end
    local map = me:map()
    if map == nil then
        return
    end
    for mob, hits in pairs(damages) do
        if not mob or not hits then
            goto continue_mob
        end
        for _, amount in ipairs(hits) do
            if not amount or amount <= 0 then
                goto continue_hit
            end
            local meso = math.floor((amount / 12300) * maxmeso)
            if meso < 1 then
                meso = 1
            end
            if meso > maxmeso then
                meso = maxmeso
            end
            if math.random(100) >= 100 then
                goto continue_hit
            end
            local x, y = mob:position()
            local offset = math.random(-20, 20)
            map:spawn_meso(meso, { x + offset, y }, me)
            ::continue_hit::
        end
        ::continue_mob::
    end
end

function damages_to_targets(damages)
    if damages == nil then
        return {}
    end
    local targets = {}
    for mob, _ in pairs(damages) do
        if mob ~= nil then
            targets[#targets + 1] = mob
        end
    end
    return targets
end

local function total_damage_to_mob(hits)
    if hits == nil then
        return 0
    end
    local total = 0
    for _, amount in ipairs(hits) do
        if amount and amount > 0 then
            total = total + amount
        end
    end
    return total
end

function handle_ice_charge_freeze(me, damages)
    if damages == nil then
        return
    end
    local job = me:class()
    if job ~= Class.WhiteKnight and job ~= Class.Paladin then
        return
    end
    local buff = me:buff(BuffFlag.WkCharge)
    if buff == nil then
        return
    end
    local wz = buff:wz()
    if wz == nil or wz.effects == nil then
        return
    end
    local sid = wz.id
    if sid ~= Skill.IceChargeSword and sid ~= Skill.BlizzardChargeBw then
        return
    end
    local effect = wz.effects[buff:level()]
    if effect == nil then
        return
    end
    local y = effect.y or 0
    local duration_ms = y * 2000
    if duration_ms <= 0 then
        return
    end
    for mob, hits in pairs(damages) do
        if mob and hits and total_damage_to_mob(hits) > 0 then
            mob:set_status(MobStatus.Freeze, 1, duration_ms, buff)
        end
    end
end

