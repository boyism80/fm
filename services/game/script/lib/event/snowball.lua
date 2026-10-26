-- Channel event game (old/src/server/events/MapleSnowball.java): 눈덩이 굴리기

local M = {}

local MAP = 109060000
-- The legacy guide gives the game a time limit; the snowball that rolled further wins when it runs out.
local ROUND_MS = 600000
local FINISH = 899
local STAGES = { 255, 511, 767 }
local SNOWMAN_HP = 7500
-- The snowmen stand between these x on both tracks; the bottom track is below this y.
local SNOWMAN_X = { -560, -360 }
local TRACK_Y = -80
local DAMAGE = { ball = 10, snowman = 15, snowman_critical = 45 }
local STOP_MS = 10000
local STUN_MS = 4000
local STATE = { Normal = 0, BottomGone = 2, TopGone = 3, Move = 4 }
local MESSAGE = { Stopped = 4, Rolling = 5 }
local TEAM_NAMES = { [0] = "story", [1] = "maple" }

M.MAPS = { MAP }
M.ENTRY = "sp"

local function team_of(player)
    local _, y = player:position()
    if y > TRACK_Y then
        return 0
    end
    return 1
end

local function number(sm, team, key)
    return tonumber(sm:get_property("snowball" .. team .. "_" .. key)) or 0
end

local function set(sm, team, key, value)
    sm:set_property("snowball" .. team .. "_" .. key, tostring(value))
end

local function show(sm, state)
    sm:map(MAP):snowball_state(state,
        { math.floor(number(sm, 0, "hp") / 75), math.floor(number(sm, 1, "hp") / 75) },
        { number(sm, 0, "position"), number(sm, 1, "position") },
        { number(sm, 0, "stage"), number(sm, 1, "stage") },
        { DAMAGE.ball, DAMAGE.snowman, DAMAGE.snowman_critical })
end

local function players(sm)
    local list = {}
    for _, player in pairs(sm:map(MAP):characters()) do
        list[#list + 1] = player
    end
    return list
end

local function knock_back(player)
    player:debuff(DebuffFlag.Stun, STUN_MS)
    player:snowball_knockback()
end

local function finish(sm, winner)
    sm:cancel("snowball")
    if winner >= 0 then
        for _, player in ipairs(players(sm)) do
            if team_of(player) == winner then
                sm:set_property("winner:" .. player:id(), "true")
            end
        end
        sm:map(MAP):snowball_state(STATE.BottomGone, { 0, 0 }, { 0, 0 }, { 0, 0 }, { 0, 0, 0 })
        sm:map(MAP):snowball_state(STATE.TopGone, { 0, 0 }, { 0, 0 }, { 0, 0 }, { 0, 0, 0 })
        for _, player in ipairs(players(sm)) do
            player:message("축하합니다! " .. TEAM_NAMES[winner] .. " 팀이 눈덩이굴리기 이벤트에서 우승했습니다!", Msg.LightBlueText)
        end
    end
    sm:call_hook("on_game_over")
end

function M.start(sm)
    for team = 0, 1 do
        set(sm, team, "hp", SNOWMAN_HP)
        set(sm, team, "position", 0)
        set(sm, team, "stage", 0)
        set(sm, team, "stopped", "false")
    end
    show(sm, STATE.Normal)
    for team = 0, 1 do
        sm:map(MAP):snowball_message(team, MESSAGE.Rolling)
    end
    for _, player in ipairs(players(sm)) do
        player:clock(ROUND_MS / 1000)
    end
    sm:set_property("clock_until", tostring(now() + ROUND_MS / 1000))
    sm:after("snowball", ROUND_MS, "on_snowball_end")
end

local function resume(sm, team)
    set(sm, team, "stopped", "false")
    sm:map(MAP):snowball_message(team, MESSAGE.Rolling)
end

M.hooks = {
    -- The server judges a hit by where the attacker stands: on its track near the snowball, or at the snowmen.
    on_snowball_hit = function(sm, player)
        if sm:get_property("phase") ~= "running" then
            return
        end
        local team = team_of(player)
        local x = player:position()
        local map = sm:map(MAP)
        if x > SNOWMAN_X[1] and x < SNOWMAN_X[2] then
            local damage = DAMAGE.snowman
            if math.random() < 0.3 then
                damage = 0
            end
            if math.random() < 0.05 then
                damage = DAMAGE.snowman_critical
            end
            map:snowball_hit(team + 2, damage, 0)
            if damage == 0 then
                return
            end
            local hp = number(sm, team, "hp") - damage
            if hp > 0 then
                set(sm, team, "hp", hp)
                show(sm, STATE.Normal)
                return
            end
            -- Bringing a snowman down stops the other team's snowball for a while.
            set(sm, team, "hp", SNOWMAN_HP)
            local other = 1 - team
            set(sm, other, "stopped", "true")
            map:snowball_message(other, MESSAGE.Stopped)
            sm:cancel("snowball_stop" .. other)
            sm:after("snowball_stop" .. other, STOP_MS, "on_snowball_resume" .. other)
            show(sm, STATE.Normal)
            return
        end

        local position = number(sm, team, "position")
        local left = position * 3 + 175
        local rolling = sm:get_property("snowball" .. team .. "_stopped") ~= "true"
        local damage = 0
        if rolling and (math.random() < 0.01 or (x > left and x < left + 275)) then
            damage = DAMAGE.ball
        end
        map:snowball_hit(team, damage, 256)
        if damage == 0 then
            if math.random() < 0.2 then
                knock_back(player)
            end
            return
        end
        position = position + 1
        set(sm, team, "position", position)
        if position >= FINISH then
            finish(sm, team)
            return
        end
        for stage, mark in ipairs(STAGES) do
            if position == mark then
                set(sm, team, "stage", stage)
                map:snowball_message(team, stage)
            end
        end
        show(sm, STATE.Move)
    end,

    on_snowball_touch = function(sm, player)
        if sm:get_property("phase") ~= "running" then
            return
        end
        knock_back(player)
    end,

    on_snowball_resume0 = function(sm)
        resume(sm, 0)
    end,

    on_snowball_resume1 = function(sm)
        resume(sm, 1)
    end,

    on_snowball_end = function(sm)
        local bottom = number(sm, 0, "position")
        local top = number(sm, 1, "position")
        local winner = -1
        if bottom > top then
            winner = 0
        elseif top > bottom then
            winner = 1
        end
        finish(sm, winner)
    end,
}

return M
