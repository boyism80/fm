function on_mob_skill_choose_144(mob, skill)
    if mob:has_buff({MobBuff.DamageImmunity, MobBuff.MagicImmunity, MobBuff.WeaponImmunity}) then
        return false
    end
    return true
end

function on_mob_skill_144(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.MagicDamageReflect, effect.x, effect.duration_ms, skill, controller)
    mob:buff(MobBuff.MagicImmunity, effect.x, effect.duration_ms, skill, controller)
    return true
end
