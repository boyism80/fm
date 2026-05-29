function on_mob_skill_choose_157(mob, skill)
    return true
end

function on_mob_skill_157(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.Seal, effect.x, effect.duration_ms, skill, controller)
    return true
end
