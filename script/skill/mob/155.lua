function on_mob_skill_choose_155(mob, skill)
    return true
end

function on_mob_skill_155(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.Avoid, effect.x, effect.duration_ms, skill, controller)
    return true
end
