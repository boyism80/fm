function on_mob_skill_choose_152(mob, skill)
    if mob:has_buff(MobBuff.WeaponDefenseUp) then
        return false
    end
    return true
end

function on_mob_skill_152(mob, controller, skill)
    local effect = skill:effect()
    mob:buff(MobBuff.WeaponDefenseUp, effect.x, effect.duration_ms, skill, controller)
    return true
end
