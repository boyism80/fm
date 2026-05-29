function on_mob_skill_choose_156(mob, skill)
    return true
end

function on_mob_skill_156(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.Speed, effect.x, effect.duration_ms, skill, controller)
    return true
end
