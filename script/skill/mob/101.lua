function on_mob_skill_choose_101(mob, skill)
    if mob:has_buff(MobBuff.MagicAttackUp) then
        return false
    end
    return true
end

function on_mob_skill_101(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.MagicAttackUp, effect.x, effect.duration_ms, skill, controller)
    return true
end
