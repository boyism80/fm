function on_mob_skill_choose_200(mob, skill)
    return true
end

local function summon_position(mob, summon_id)
    local mob_x, mob_y = mob:position()
    local xpos = mob_x
    local ypos = mob_y
    local foothold = nil

    local summon_wz = id2mob(summon_id)
    if summon_wz ~= nil then
        if summon_wz.name == "로우 다크스타" then
            foothold = math.random(1, 19)
            ypos = -590
        elseif summon_wz.name == "하이 다크스타" then
            xpos = mob_x + math.random(1, 1000) - 500
            ypos = mob_y
        elseif summon_wz.name == "블러드붐" then
            if math.random(1, 5) == 1 then
                ypos = 78
                xpos = math.random(1, 5) + (math.random(1, 2) == 1 and 180 or 0)
            else
                xpos = mob_x + math.random(1, 1000) - 500
            end
        end
    end

    local map_wz = id2map(mob:map_id())
    if map_wz ~= nil then
        if map_wz.name == "시계탑의 근원" then
            if xpos < -890 then
                xpos = -890 + math.random(1, 150)
            elseif xpos > 230 then
                xpos = 230 - math.random(1, 150)
            end
        elseif map_wz.name == "피아누스의 동굴" then
            if xpos < -239 then
                xpos = -239 + math.random(1, 150)
            elseif xpos > 371 then
                xpos = 371 - math.random(1, 150)
            end
        end
    end

    return xpos, ypos, foothold
end

function on_mob_skill_200(mob, controller, skill)
    local effect = skill:effect()
    local map = mob:map()
    if map == nil then
        return true
    end
    local summons = effect.summons
    if summons == nil then
        return true
    end
    for i = 1, #summons do
        local summon_id = summons[i]
        local xpos, ypos, foothold = summon_position(mob, summon_id)
        local ground = map:point_below({x = xpos, y = ypos - 1})
        if ground ~= nil then
            xpos = ground.x
            ypos = ground.y
        end
        local spawned = map:spawn_mob(summon_id, xpos, ypos, effect.spawn_effect)
        if spawned ~= nil and foothold ~= nil then
            spawned:foothold(foothold)
        end
    end
    return true
end
