-- Mob skill name (Skill.wz/MobSkill.img.xml): mob skill 105

function on_mob_skill_choose_105(mob, skill)
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

local function heal_mob(target, effect)
    if effect.x ~= 0 then
        target:add_hp(effect.x)
    end
    if effect.y ~= 0 then
        target:add_mp(effect.y)
    end
end

function on_mob_skill_105(mob, controller, skill)
    local effect = skill:effect()
    local is_range = effect_has_bounds(effect)
    local targets = {}

    if is_range then
        local map = mob:map()
        local world_bounds = world_bounds_from_effect(mob, effect)
        local self_oid = mob:oid()

        for _, other in pairs(map:mobs()) do
            if other:oid() ~= self_oid then
                local other_x, other_y = other:position()
                if position_in_bounds(other_x, other_y, world_bounds) then
                    table.insert(targets, other)
                    break
                end
            end
        end
    else
        table.insert(targets, mob)
    end

    local map = mob:map()
    for _, target in ipairs(targets) do
        if is_range then
            map:remove_mob(target:oid(), MobDieAnimation.FadeOut)
            heal_mob(mob, effect)
        else
            heal_mob(target, effect)
        end
    end

    return true
end
