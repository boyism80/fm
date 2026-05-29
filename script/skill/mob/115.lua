function on_mob_skill_choose_115(mob, skill)
    return true
end

function on_mob_skill_115(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.Speed, effect.x, effect.duration_ms, skill, controller)
    return true
end
