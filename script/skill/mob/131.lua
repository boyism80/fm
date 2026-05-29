function on_mob_skill_choose_131(mob, skill)
    return true
end

function on_mob_skill_131(mob, controller, skill)
    local effect = skill:effect()
    mob:create_mist(skill, effect.x * 10, MistType.Poison, effect.bounds)
    return true
end
