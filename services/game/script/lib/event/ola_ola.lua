-- Channel event game (old/src/server/events/MapleOla.java): 올라올라

local M = {}

local MAPS_BY_STAGE = { 109030001, 109030002, 109030003 }
local WINNER_MAP = 109050000
local DURATION_MS = 360000
-- Stage 1 has portals ch00-ch04, stage 2 ch00-ch07 and stage 3 ch00-ch15.
local PORTALS = { 5, 8, 16 }

M.MAPS = MAPS_BY_STAGE

function M.prepare(sm)
    local closed = sm:map(MAPS_BY_STAGE[1]):portal("join00")
    if closed ~= nil then
        closed:script("event_closed")
    end
end

local function contestants(sm)
    local players = {}
    for _, map_id in ipairs(MAPS_BY_STAGE) do
        for _, player in pairs(sm:map(map_id):characters()) do
            players[#players + 1] = player
        end
    end
    return players
end

local function stage_portal(sm, stage)
    return tonumber(sm:get_property("ola_stage" .. stage))
end

function M.start(sm)
    for stage, count in ipairs(PORTALS) do
        local correct = math.random(0, count - 1)
        -- The second portal of stage 1 cannot be reached.
        if stage == 1 and correct == 2 then
            correct = 3
        end
        sm:set_property("ola_stage" .. stage, tostring(correct))
    end
    sm:map(MAPS_BY_STAGE[1]):portal("join00"):script("")
    for _, player in ipairs(contestants(sm)) do
        player:message("포탈이 열렸습니다. 포탈 위에서 ↑를 눌러주세요.", Msg.Notice)
        player:clock(DURATION_MS / 1000)
    end
    sm:set_property("clock_until", tostring(now() + DURATION_MS / 1000))
    sm:after("ola", DURATION_MS, "on_ola_end")
end

-- Called by the rand_ola portals: the right portal of a stage leads up, any other one back to the start of the map.
function M.pass(me, portal)
    local sm = me:state_machine()
    local map_id = me:map():wz():id()
    local stage = map_id % 10
    if sm == nil or sm:get_property("phase") ~= "running" then
        me:map(me:map(), 0, { relocate = true })
        return
    end
    if portal:name() ~= string.format("ch%02d", stage_portal(sm, stage)) then
        me:map(me:map(), 0, { relocate = true })
        return
    end
    if stage < #MAPS_BY_STAGE then
        me:map(sm:map(MAPS_BY_STAGE[stage + 1]), 0)
        return
    end

    sm:set_property("winner:" .. me:id(), "true")
    me:stop_clock()
    me:message("축하합니다! [올라올라] 이벤트를 완주하셨습니다! NPC [피에트로] 를 더블클릭하면 보상을 받을 수 있습니다.", Msg.Notice)
    me:map(sm:map(WINNER_MAP))
    if #contestants(sm) == 0 then
        sm:cancel("ola")
        sm:call_hook("on_game_over")
    end
end

M.hooks = {
    on_ola_end = function(sm)
        sm:call_hook("on_game_over")
    end,
}

return M
