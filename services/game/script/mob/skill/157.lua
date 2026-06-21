-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 157

function on_mob_skill_choose_157(mob, skill)
    return true
end

function on_mob_skill_157(mob, controller, skill)
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
        target:buff(MobBuff.Seal, effect.x, effect.time, skill, controller)
    end
    return true
end
