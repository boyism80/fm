-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 132

function on_mob_skill_choose_132(mob, skill)
    return true
end

function on_mob_skill_132(mob, controller, skill)
    local effect = skill:effect()
    local time = effect.time
    if time <= 0 then
        return true
    end

    local targets = {}
    local bounds = effect.bounds
    if bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0 then
        for _, ch in ipairs(mob:objects_in(bounds, ObjectType.Character)) do
            table.insert(targets, ch)
        end
    elseif controller ~= nil then
        table.insert(targets, controller)
    end

    local x = effect.x
    if x <= 0 then
        x = 1
    end
    for _, ch in ipairs(targets) do
        ch:debuff(DebuffFlag.ReverseDirection, time, x, skill:id(), skill:level())
    end
    return true
end
