function on_mob_skill_choose_154(mob, skill)
    return true
end

function on_mob_skill_154(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.Acc, effect.x, effect.duration_ms, skill, controller)
    return true
end
