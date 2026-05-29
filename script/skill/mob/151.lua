function on_mob_skill_choose_151(mob, skill)
    if mob:has_buff(MobBuff.MagicAttackUp) then
        return false
    end
    return true
end

function on_mob_skill_151(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.MagicAttackUp, effect.x, effect.duration_ms, skill, controller)
    return true
end
