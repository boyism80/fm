-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 129

function on_mob_skill_choose_129(mob, skill)
    return true
end

function on_mob_skill_129(mob, controller, skill)
    local effect = skill:effect()
    local targets = {}
    local bounds = effect.bounds

    if bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0 then
        for _, target in ipairs(mob:objects_in(bounds, ObjectType.Character)) do
            table.insert(targets, target)
        end
    elseif controller ~= nil then
        table.insert(targets, controller)
    end

    for _, target in ipairs(targets) do
        mob:banish(target)
    end

    return true
end
