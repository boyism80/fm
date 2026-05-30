function on_mob_skill_choose_110(mob, skill)
    if mob:has_buff(MobBuff.WeaponAttackUp) then
        return false
    end
    return true
end

function on_mob_skill_110(mob, controller, skill)
    local effect = skill:effect()
    local targets = {}
    local bounds = effect.bounds
    if bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0 then
        for _, target in ipairs(mob:objects_in(bounds, ObjectType.Mob)) do
            table.insert(targets, target)
        end
    end
    table.insert(targets, mob)
    for _, target in ipairs(targets) do
        target:buff(MobBuff.WeaponAttackUp, effect.x, effect.time, skill, controller)
    end
    return true
end
