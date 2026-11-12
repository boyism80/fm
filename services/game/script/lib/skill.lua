local combat = require("script/lib/combat")
local M = {}

M.apply_buff_from_effect = function(me, skill, flag, value_key)
    value_key = value_key or "x"
    local effect = skill:effect()
    if effect == nil then
        return
    end
    local val = effect[value_key]
    if val == nil then
        return
    end
    me:buff(skill, flag, val)
end

M.apply_buff_fixed = function(me, skill, flag, value)
    value = value or 1
    me:buff(skill, flag, value)
end

-- Homing marks the first mob the attack hurt; a new mark replaces the old one.
M.mark_homing = function(me, skill, damages)
    if skill == nil or me == nil or damages == nil then
        return
    end
    for mob, hits in pairs(damages) do
        if combat.damaged(hits) then
            mob:homing(me, skill)
            return
        end
    end
end

-- Dark sight slows its user by the skill's WZ speed, which the buff carries.
M.apply_dark_sight = function(me, skill)
    me:buff(skill, { [BuffFlag.Darksight] = 1 })
end

M.apply_iron_body = function(me, skill)
    me:buff(skill)
end

-- A transformation carries the skill's WZ stats along with the look.
M.apply_morph = function(me, skill, morph)
    me:buff(skill, { [BuffFlag.Morph] = morph })
end

M.apply_monster_magnet = function(me, skill, params)
    if params == nil or params.magnet == nil then
        return
    end
    for _, entry in ipairs(params.magnet.mobs) do
        local mob = entry.mob
        local success = entry.success
        if mob then
            me:show_magnet(mob, success)
            mob:controller(me)
        end
    end
end

M.consume_combo_orbs = function(me, howmany)
    local current = me:buff_value(BuffFlag.Combo)
    if current == nil or current <= 1 then
        return
    end

    if howmany == nil then
        howmany = current - 1
    end
    if howmany <= 0 then
        return
    end

    local new_orbs = current - howmany
    if new_orbs < 1 then
        new_orbs = 1
    end

    me:buff_value(BuffFlag.Combo, new_orbs)
end

M.apply_hero_will = function(me)
    me:remove_debuff(DebuffFlag.Seduce)
end

M.apply_dispel = function(me, skill)
    local effect = skill:effect()
    if effect == nil then
        return
    end
    if effect.prop <= 0 then
        return
    end
    if effect.prop < 100 and math.random(1, 100) > effect.prop then
        return
    end
    me:remove_debuff(
        DebuffFlag.Curse,
        DebuffFlag.Darkness,
        DebuffFlag.Poison,
        DebuffFlag.Seal,
        DebuffFlag.Slow,
        DebuffFlag.Weaken
    )
end

M.get_heal_recovery_amount = function(me, skill)
    local effect = skill:effect()
    if effect == nil then
        return 0
    end
    local hp_rate = effect.hp / 100.0
    if hp_rate <= 0 then
        return 0
    end

    local total_magic = me:total_magic()
    if total_magic <= 0 then
        total_magic = 1
    end

    local min_heal = math.floor(total_magic * 3.0 * hp_rate)
    local max_heal = math.floor(total_magic * 5.0 * hp_rate)
    if min_heal < 1 then
        min_heal = 1
    end
    if max_heal < min_heal then
        max_heal = min_heal
    end

    return math.random(min_heal, max_heal)
end

M.apply_mana_reflection = function(me, skill)
    local effect = skill:effect()
    if effect == nil then
        return
    end

    local val = effect.prop - 30
    if val < 0 then
        val = 0
    end

    me:buff(skill, BuffFlag.ManaReflection, val)
end

-- A summon's hit puts the status on each mob it damaged for 4 seconds, at the skill's prop (100 when the WZ has none).
M.summon_status_attack = function(me, skill, damages, flag)
    local effect = skill:effect()
    local prop = 100
    if effect ~= nil and effect.prop > 0 then
        prop = effect.prop
    end
    for mob, hits in pairs(damages) do
        if combat.damaged(hits) and math.random(1, 100) <= prop then
            mob:buff(flag, 1, 4000, skill, me)
        end
    end
end

M.apply_resurrection = function(me, skill)
    combat.for_each_near_party_member(me, skill, function(ch)
        if ch:is_alive() then
            return
        end
        ch:stance(0)
        ch:hp(ch:max_hp())
        ch:mp(ch:max_mp())
        ch:show_skill_effect(skill, SkillEffectType.Affected)
    end, { allow_dead = true, exclude_caster = true })
end

M.mark_archer_puppet = function(me, skill, params)
    local effect = skill:effect()
    if effect == nil then
        return
    end
    if effect.time <= 0 then
        return
    end
    me:buff(skill, BuffFlag.Puppet, 1)
end

M.apply_archer_puppet_buff = function(me, skill)
    local wz = skill:wz()
    local effect = skill:effect()
    if effect == nil then
        return
    end
    if effect.time <= 0 then
        return
    end
    local level = skill:level()
    if level == nil or level <= 0 then
        return
    end
    if effect.x <= 0 then
        return
    end
    local s = me:create_summon(wz:id(), level, effect.time, SummonMovementType.Stationary, SummonType.Puppet)
    if s ~= nil then
        s:max_hp(effect.x, false)
        s:hp(effect.x, false)
    end
end

M.apply_archer_puppet_unbuff = function(me, skill)
    local wz = skill:wz()
    me:remove_summon(wz:id())
end

M.absorb_magic_guard = function(me, damage, params)
    if damage == nil or damage <= 0 then
        return damage
    end

    local infinity = me:buff_value(BuffFlag.Infinity) ~= nil
    if params.deadly then
        if not infinity and me:mp() > 1 then
            me:add_mp(1 - me:mp())
        end
        return math.max(me:hp() - 1, 0)
    end

    local guard_percent = me:buff_value(BuffFlag.MagicGuard)
    local mp_burn = params.mp_burn or 0
    local mp_loss = math.max(math.floor(damage * (guard_percent / 100.0)), 0) + mp_burn
    local hp_loss = damage - mp_loss
    if infinity then
        mp_loss = 0
    elseif mp_loss > me:mp() then
        mp_loss = me:mp()
        hp_loss = damage - mp_loss + mp_burn
    end

    if mp_loss > 0 then
        me:add_mp(-mp_loss)
    end

    return math.max(hp_loss, 0)
end

M.absorb_meso_guard = function(me, damage, params)
    if damage == nil or damage <= 0 then
        return damage
    end

    local mp_loss = params.mp_burn or 0
    if params.deadly then
        mp_loss = math.max(me:mp() - 1, 0)
    end
    if mp_loss > 0 then
        me:add_mp(-mp_loss)
    end

    local guard_percent = me:buff_value(BuffFlag.MesoGuard)
    if guard_percent == nil or guard_percent <= 0 then
        return damage
    end
    local current_meso = me:meso()
    if current_meso == nil or current_meso <= 0 then
        return damage
    end
    local meso_loss = math.floor(damage * (guard_percent / 100.0))
    if meso_loss < 0 then
        meso_loss = 0
    end
    if meso_loss > current_meso then
        meso_loss = current_meso
    end
    local hp_loss = damage - meso_loss
    if hp_loss < 0 then
        hp_loss = 0
    end
    if meso_loss > 0 then
        me:meso(current_meso - meso_loss)
    end
    return hp_loss
end

M.take_hit = function(me, damage, params)
    if damage == nil or damage <= 0 then
        return damage
    end
    if params.deadly then
        if me:mp() > 1 then
            me:add_mp(1 - me:mp())
        end
        return math.max(me:hp() - 1, 0)
    end
    local mp_burn = params.mp_burn or 0
    if mp_burn > 0 then
        me:add_mp(-mp_burn)
    end
    return damage
end

return M
