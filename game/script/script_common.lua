-- Shared helpers for script modules (skill effect, etc.)

function get_skill_effect(skill_entry)
    if skill_entry == nil then
        return nil
    end
    local wz = skill_entry:wz()
    if wz == nil or wz.effects == nil then
        return nil
    end
    local level = skill_entry:level()
    if level == nil then
        return nil
    end
    return wz.effects[level]
end
