function on_mob_skill_choose_100(mob, skill)
    if mob:has_buff(MobBuff.WeaponAttackUp) then
        return false
    end
    return true
end

function on_mob_skill_100(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.WeaponAttackUp, effect.x, effect.duration_ms, skill, controller)
    return true
end
