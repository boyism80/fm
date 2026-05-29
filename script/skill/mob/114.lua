function on_mob_skill_choose_114(mob, skill)
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

local function random_range_heal_hp(effect)
    local hp_scale = math.floor(effect.x / 1000)
    local hp_roll = math.floor(950 + 1050 * math.random())
    return hp_scale * hp_roll
end

function on_mob_skill_114(mob, controller, skill)
    local effect = skill:effect()
    local is_range = effect_has_bounds(effect)
    local targets = {}

    if is_range then
        local map = mob:map()
        local world_bounds = world_bounds_from_effect(mob, effect)

        for _, target in pairs(map:mobs()) do
            local target_x, target_y = target:position()
            if position_in_bounds(target_x, target_y, world_bounds) then
                table.insert(targets, target)
            end
        end
    else
        table.insert(targets, mob)
    end

    local hp_amount = effect.x
    if is_range then
        hp_amount = random_range_heal_hp(effect)
    end

    for _, target in ipairs(targets) do
        if hp_amount ~= 0 then
            target:add_hp(hp_amount)
        end
        if effect.y ~= 0 then
            target:add_mp(effect.y)
        end
    end

    return true
end
