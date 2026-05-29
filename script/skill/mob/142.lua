function on_mob_skill_choose_142(mob, skill)
    if mob:has_buff({MobBuff.DamageImmunity, MobBuff.MagicImmunity, MobBuff.WeaponImmunity}) then
        return false
    end
    return true
end

function on_mob_skill_142(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.DamageImmunity, effect.x, effect.duration_ms, skill, controller)
    return true
end
