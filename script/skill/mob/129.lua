function on_mob_skill_choose_129(mob, skill)
    return true
end

local function effect_has_bounds(effect)
    local bounds = effect.bounds
    return bounds.left ~= 0 or bounds.top ~= 0 or bounds.right ~= 0 or bounds.bottom ~= 0
end

local function world_bounds_from_effect(mob, effect)
    local mob_x, mob_y = mob:position()
    local bounds = effect.bounds
    return {
        left = mob_x + bounds.left,
        top = mob_y + bounds.top,
        right = mob_x + bounds.right,
        bottom = mob_y + bounds.bottom,
    }
end

local function position_in_bounds(x, y, bounds)
    return x >= bounds.left and x <= bounds.right and y >= bounds.top and y <= bounds.bottom
end

function on_mob_skill_129(mob, controller, skill)
    local effect = skill:effect()
    local targets = {}

    if effect_has_bounds(effect) then
        if controller ~= nil then
            local map = mob:map()
            local world_bounds = world_bounds_from_effect(mob, effect)

            for _, target in pairs(map:characters()) do
                local target_x, target_y = target:position()
                if position_in_bounds(target_x, target_y, world_bounds) then
                    table.insert(targets, target)
                end
            end
        end
    else
        if controller ~= nil then
            table.insert(targets, controller)
        end
    end

    for _, target in ipairs(targets) do
        mob:banish(target)
    end

    return true
end
