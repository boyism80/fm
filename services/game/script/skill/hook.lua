local combat = require("script/lib/combat")

local function damages_has_positive_damage(hits)
    if hits == nil then
        return false
    end
    for _, amount in ipairs(hits) do
        if (amount or 0) > 0 then
            return true
        end
    end
    return false
end

local function roll_percent(prob)
    if prob == nil or prob <= 0 then
        return false
    end
    if prob >= 100 then
        return true
    end
    return math.random(1, 100) <= prob
end

local function get_mp_eater_skill(me)
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

local function on_mp_eater(me, damages)
    local mp_eater_skill = get_mp_eater_skill(me)
    if mp_eater_skill == nil then
        return
    end

    local effect = mp_eater_skill:effect()
    if effect == nil then
        return
    end

    if effect.prop <= 0 or effect.x <= 0 then
        return
    end

    local total_absorb_mp = 0

    for mob, hits in pairs(damages) do
        if not damages_has_positive_damage(hits) then
            goto continue_mob
        end

        local wz = mob:wz()
        local is_boss = wz and wz:boss()
        if is_boss then
            goto continue_mob
        end

        if not roll_percent(effect.prop) then
            goto continue_mob
        end

        local mob_mp = mob:mp()
        if mob_mp == nil or mob_mp <= 0 then
            goto continue_mob
        end

        local absorb_mp = math.floor(mob:max_mp() * (effect.x / 100.0))
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
        me:show_skill_effect(mp_eater_skill)
    end
end

local function is_magic_attack_skill(skill)
    if skill == nil then
        return false
    end
    local effect = skill:effect()
    if effect == nil then
        return false
    end
    return effect.mad > 0
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
    local amp_effect = amp:effect()
    if amp_effect == nil then
        return mp_con
    end
    local mult = amp_effect.x
    if mult <= 0 then
        return mp_con
    end
    return math.floor(mp_con * mult / 100)
end

local function apply_concentrate_mp_cost(me, base_mp_con, mp_con_after_amp)
    if base_mp_con <= 0 or mp_con_after_amp <= 0 then
        return mp_con_after_amp
    end
    local s = me:buff_value(BuffFlag.Concentrate)
    if s == nil then
        return mp_con_after_amp
    end
    if s <= 0 then
        return mp_con_after_amp
    end
    local reduce = math.floor(base_mp_con * s / 100)
    local out = mp_con_after_amp - reduce
    if out < 0 then
        out = 0
    end
    return out
end

-- on_activating checks this cost before anything is spent, and on_activated pays it.
local function skill_cost(me, skill)
    local effect = skill:effect()
    if effect == nil then
        return 0, 0
    end
    local mp_con = apply_element_amplification_mp_cost(me, skill, effect.mp_con)
    if me:buff_value(BuffFlag.Infinity) ~= nil then
        mp_con = 0
    else
        mp_con = apply_concentrate_mp_cost(me, effect.mp_con, mp_con)
    end
    return mp_con, effect.hp_con
end

local function can_pay_skill_cost(me, skill)
    local mp_con, hp_con = skill_cost(me, skill)
    return (mp_con <= 0 or me:mp() >= mp_con) and (hp_con <= 0 or me:hp() > hp_con)
end

local function pay_skill_cost(me, skill)
    local mp_con, hp_con = skill_cost(me, skill)
    if mp_con > 0 then
        me:mp(math.max(me:mp() - mp_con, 0), false)
    end
    if hp_con > 0 then
        me:hp(math.max(me:hp() - hp_con, 1), false)
    end
    me:update_stats({ STAT.Hp, STAT.Mp })
end

local function on_energy_charge(me, damages)
    local ec
    if me:class_of(Class.Brawler) then
        ec = me:skill(Skill.EnergyCharge)
    elseif me:class_of(Class.ThunderBreaker2) then
        ec = me:skill(Skill.EnergyChargeCygnus)
    else
        return
    end

    if ec == nil then
        return
    end
    local effect = ec:effect()
    if effect == nil then
        return
    end
    local value = me:buff_value(BuffFlag.EnergyCharge) or 0
    if value >= 10000 then
        return
    end
    local targets = 0
    for _, hits in pairs(damages or {}) do
        targets = targets + #(hits or {})
    end
    if targets <= 0 then
        return
    end
    local inc = effect.x * targets
    value = math.min(value + inc, 10000)
    if value < 10000 then
        me:buff(ec, BuffFlag.EnergyCharge, value, { time = 0 })
    else
        me:buff(ec, BuffFlag.EnergyCharge, value)
    end
end

local function on_blind_acc_debuff(me, damages)
    local blind_buff = me:buff(BuffFlag.Blind)
    if blind_buff == nil then
        return
    end
    local effect = blind_buff:effect()
    if effect == nil then
        return
    end
    if effect.prop <= 0 then
        return
    end
    local acc = me:buff_value(BuffFlag.Blind) or 0
    if acc == 0 then
        return
    end
    local duration_ms = effect.y * 1000
    if duration_ms <= 0 then
        return
    end
    for mob, hits in pairs(damages) do
        if not damages_has_positive_damage(hits) then
            goto continue_blind
        end
        if roll_percent(effect.prop) then
            mob:buff(MobBuff.Acc, acc, duration_ms, blind_buff, me)
        end
        ::continue_blind::
    end
end

local function on_hamstring_slow(me, damages)
    local ham_buff = me:buff(BuffFlag.Hamstring)
    if ham_buff == nil then
        return
    end

    local effect = ham_buff:effect()
    if effect == nil then
        return
    end

    if effect.prop <= 0 then
        return
    end

    local slow = me:buff_value(BuffFlag.Hamstring) or 0
    if slow == 0 then
        return
    end

    local duration_ms = effect.y * 1000
    if duration_ms <= 0 then
        return
    end

    for mob, hits in pairs(damages) do
        if not damages_has_positive_damage(hits) then
            goto continue_mob
        end
        if roll_percent(effect.prop) then
            mob:buff(MobBuff.Speed, slow, duration_ms, ham_buff, me)
        end
        ::continue_mob::
    end
end

local function on_dark_sight(me, damages)
    if not (me:class_of(Class.NightWalker1) or me:class_of(Class.WindArcher2)) then
        return
    end

    if me:buff(BuffFlag.Darksight) == nil then
        return
    end

    local is_hit = false
    for _, hits in pairs(damages) do
        for _, hit in pairs(hits) do
            if hit > 0 then
                is_hit = true
                break
            end
        end
    end
    if not is_hit then
        return
    end

    me:unbuff(BuffFlag.Darksight)
end

return {
    on_passive = function(me, skill)
    end,

    on_unpassive = function(me, skill)
    end,

    on_activating = function(me, skill, params)
        if me:forbidden("skill") then
            me:message("이곳에서 스킬을 사용할 수 없습니다.", Msg.PinkText)
            return false
        end
        return can_pay_skill_cost(me, skill)
    end,

    on_activated = function(me, skill, params)
        pay_skill_cost(me, skill)
        return true
    end,

    on_attack = function(me, skill, damages, attack_info)
        local targets = combat.damages_to_targets(damages)
        combat.increment_combo_orbs(me, targets, skill)
        combat.spawn_pickpocket_meso(me, skill, damages)
        combat.freeze_from_ice_charge(me, damages)
        on_mp_eater(me, damages)
        on_hamstring_slow(me, damages)
        on_blind_acc_debuff(me, damages)
        combat.roll_mortal_blow(me, damages)
        on_energy_charge(me, damages)
        on_dark_sight(me, damages)
        combat.reflect_mob_attack_damage(me, damages, attack_info)
    end
}
