-- Damage handlers: magic guard, meso guard (on_damaged)

function handle_magic_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end

    if me:buff_value(BuffFlag.MagicGuard) == nil then
        return damage
    end

    local effect = get_skill_effect(me:skill(Skill.MagicGuard))
    if effect == nil then
        effect = get_skill_effect(me:skill(Skill.MagicGuardCygnus))
    end
    local guard_percent = effect and effect.x
    if guard_percent == nil or guard_percent <= 0 then
        return damage
    end

    local current_mp = me:mp()
    if current_mp == nil or current_mp <= 0 then
        return damage
    end

    local mp_loss = math.floor(damage * (guard_percent / 100.0))
    if mp_loss < 0 then
        mp_loss = 0
    end
    if mp_loss > current_mp then
        mp_loss = current_mp
    end

    local hp_loss = damage - mp_loss
    if hp_loss < 0 then
        hp_loss = 0
    end

    if mp_loss > 0 then
        me:add_mp(-mp_loss)
    end

    return hp_loss
end

function handle_meso_guard(me, attacker, skill, damage)
    if damage == nil or damage <= 0 then
        return damage
    end
    if me:buff_value(BuffFlag.MesoGuard) == nil then
        return damage
    end

    local effect = get_skill_effect(me:skill(Skill.MesoGuard))
    local guard_percent = effect and effect.x
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
        me:meso(-meso_loss)
    end
    return hp_loss
end
