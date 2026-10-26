-- Channel event game (old/src/server/events/MapleCoconut.java): 코코넛 시즌 and 코크 플레이

local M = {}

-- The map's WZ coconut node: how many may fall, blow up or stick, out of every coconut on the map.
local FALLING = 401
local BOMBING = 80
local STOPPED = 20
local COCONUTS = FALLING + BOMBING + STOPPED + 5
local ROUND_MS = 300000
local EXTRA_MS = 120000
local RESULT_MS = 10000
local STATE = { Reset = 0, Hit = 1, Bomb = 2, Fall = 3 }
local TEAM_MAPLE = 0

local function number(sm, key)
    return tonumber(sm:get_property(key)) or 0
end

local function field(sm)
    return sm:map(number(sm, "coconut_map"))
end

local function players(sm)
    local list = {}
    for _, player in pairs(field(sm):characters()) do
        list[#list + 1] = player
    end
    return list
end

local function show_score(sm)
    field(sm):coconut_score(number(sm, "coconut_score0"), number(sm, "coconut_score1"))
end

local function start_round(sm, ms)
    for _, player in ipairs(players(sm)) do
        player:clock(ms / 1000)
    end
    sm:set_property("clock_until", tostring(now() + ms / 1000))
    sm:after("coconut", ms, "on_coconut_round_end")
end

-- Each map is its own game: the event names the map and the portal its contestants enter by.
function M.game(map_id)
    local game = { MAPS = { map_id }, ENTRY = "sp", hooks = M.hooks }

    function game.start(sm)
        sm:set_property("coconut_map", tostring(map_id))
        sm:set_property("coconut_falling", tostring(FALLING))
        sm:set_property("coconut_bombing", tostring(BOMBING))
        sm:set_property("coconut_stopped", tostring(STOPPED))
        local map = field(sm)
        map:coconut_hit(-1, 0, STATE.Reset)
        show_score(sm)
        for _, player in ipairs(players(sm)) do
            player:message("이벤트가 시작되었습니다!!", Msg.PinkText)
        end
        start_round(sm, ROUND_MS)
    end

    return game
end

local function take(sm, key)
    local left = number(sm, key)
    if left <= 0 then
        return false
    end
    sm:set_property(key, tostring(left - 1))
    return true
end

M.hooks = {
    -- A coconut shakes on each hit; after its third it may come down, and only a falling one scores for the hitter's team.
    on_coconut_hit = function(sm, player, id)
        if sm:get_property("phase") ~= "running" or sm:get_property("coconut_over") == "true" then
            return
        end
        if id < 0 or id >= COCONUTS or sm:get_property("coconut_gone:" .. id) == "true" then
            return
        end
        if now() < number(sm, "coconut_next:" .. id) then
            return
        end
        local map = field(sm)
        local hits = number(sm, "coconut_hits:" .. id)
        if hits <= 2 or math.random() >= 0.4 then
            sm:set_property("coconut_hits:" .. id, tostring(hits + 1))
            sm:set_property("coconut_next:" .. id, tostring(now() + 1))
            map:coconut_hit(id, 0, STATE.Hit)
            return
        end

        sm:set_property("coconut_gone:" .. id, "true")
        if math.random() < 0.01 and take(sm, "coconut_stopped") then
            map:coconut_hit(id, 0, STATE.Hit)
            return
        end
        if math.random() < 0.05 and take(sm, "coconut_bombing") then
            map:coconut_hit(id, 0, STATE.Bomb)
            return
        end
        if take(sm, "coconut_falling") then
            map:coconut_hit(id, 0, STATE.Fall)
            local team = player:field_team()
            sm:set_property("coconut_score" .. team, tostring(number(sm, "coconut_score" .. team) + 1))
            show_score(sm)
        end
    end,

    -- A tie after the first round plays an extra round; a tie after that ends with no winner.
    on_coconut_round_end = function(sm)
        local maple = number(sm, "coconut_score0")
        local story = number(sm, "coconut_score1")
        if maple == story and sm:get_property("coconut_extra") ~= "true" then
            sm:set_property("coconut_extra", "true")
            start_round(sm, EXTRA_MS)
            return
        end

        sm:set_property("coconut_over", "true")
        local winner = -1
        if maple > story then
            winner = TEAM_MAPLE
        elseif story > maple then
            winner = 1
        end
        for _, player in ipairs(players(sm)) do
            if player:field_team() == winner then
                sm:set_property("winner:" .. player:id(), "true")
                player:map_effect("event/coconut/victory")
                player:play_sound("Coconut/Victory")
            else
                player:map_effect("event/coconut/lose")
                player:play_sound("Coconut/Failed")
            end
        end
        sm:after("coconut_result", RESULT_MS, "on_coconut_result")
    end,

    on_coconut_result = function(sm)
        sm:call_hook("on_game_over")
    end,
}

return M
