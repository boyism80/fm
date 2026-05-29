function on_mob_skill_choose_145(mob, skill)
    if mob:has_buff({MobBuff.DamageImmunity, MobBuff.MagicImmunity, MobBuff.WeaponImmunity}) then
        return false
    end
    return true
end

function on_mob_skill_145(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.WeaponDamageReflect, effect.x, effect.duration_ms, skill, controller)
    mob:buff(MobBuff.WeaponImmunity, effect.x, effect.duration_ms, skill, controller)
    mob:buff(MobBuff.MagicDamageReflect, effect.x, effect.duration_ms, skill, controller)
    mob:buff(MobBuff.MagicImmunity, effect.x, effect.duration_ms, skill, controller)
    return true
end
