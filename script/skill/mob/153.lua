function on_mob_skill_choose_153(mob, skill)
    if mob:has_buff(MobBuff.MagicDefenseUp) then
        return false
    end
    return true
end

function on_mob_skill_153(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.MagicDefenseUp, effect.x, effect.duration_ms, skill, controller)
    return true
end
