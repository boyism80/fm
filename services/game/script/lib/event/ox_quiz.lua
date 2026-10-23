-- Channel event game (old/src/server/events/MapleOxQuiz.java): 스피드 OX 퀴즈

local M = {}

local OX_MAP = 109020001
local LOSER_MAP = 109050001
local ROUNDS = 6
local READY_MS = 20000
local ANSWER_MS = 30000
local EXP = 3000
-- O stands left of x = -234 and X right of it, both on the floor below y = -26.
local BORDER_X = -234
local FLOOR_Y = -26

M.MAPS = { OX_MAP }

local function contestants(sm)
    local players = {}
    for _, player in pairs(sm:map(OX_MAP):characters()) do
        if player:role() < ROLE.Admin then
            players[#players + 1] = player
        end
    end
    return players
end

local function notice(sm, text)
    for _, player in pairs(sm:map(OX_MAP):characters()) do
        player:message(text, Msg.Notice)
    end
end

local function round(sm)
    return tonumber(sm:get_property("ox_round")) or 0
end

local function pick(sm)
    local questions = ox_questions()
    local unused = {}
    for _, q in ipairs(questions) do
        if sm:get_property(string.format("ox_used:%d:%d", q.set, q.id)) ~= "true" then
            unused[#unused + 1] = q
        end
    end
    if #unused == 0 then
        return nil
    end
    local q = unused[math.random(#unused)]
    sm:set_property(string.format("ox_used:%d:%d", q.set, q.id), "true")
    return q
end

-- Standing on the side of the answer, on the floor; anywhere else is wrong.
local function answered(player, o)
    local x, y = player:position()
    if y <= FLOOR_Y then
        return false
    end
    if o then
        return x < BORDER_X
    end
    return x > BORDER_X
end

local function prepare(sm)
    if round(sm) == ROUNDS then
        notice(sm, "축하합니다~ 준비된 퀴즈를 모두 맞추셨습니다! 잠시 후 이벤트 상품을 지급받을 수 있습니다.")
    else
        notice(sm, "20초 후 문제가 출제됩니다. 모두 준비해 주세요!")
    end
    for _, player in pairs(sm:map(OX_MAP):characters()) do
        player:clock(READY_MS / 1000)
    end
    sm:after("ox", READY_MS, "on_ox_ask")
end

local function over(sm)
    for _, player in ipairs(contestants(sm)) do
        sm:set_property("winner:" .. player:id(), "true")
    end
    sm:call_hook("on_game_over")
end

function M.start(sm)
    sm:set_property("ox_round", "0")
    prepare(sm)
end

M.hooks = {
    on_ox_ask = function(sm)
        if #contestants(sm) == 0 or round(sm) == ROUNDS then
            over(sm)
            return
        end
        local q = pick(sm)
        if q == nil then
            over(sm)
            return
        end
        sm:set_property("ox_set", tostring(q.set))
        sm:set_property("ox_id", tostring(q.id))
        sm:set_property("ox_o", tostring(q.o))
        local m = sm:map(OX_MAP)
        for _, player in pairs(m:characters()) do
            player:stop_clock()
        end
        m:ox_quiz(true, q.set, q.id)
        sm:after("ox", ANSWER_MS, "on_ox_answer")
    end,

    on_ox_answer = function(sm)
        local set = tonumber(sm:get_property("ox_set"))
        local id = tonumber(sm:get_property("ox_id"))
        local o = sm:get_property("ox_o") == "true"
        sm:map(OX_MAP):ox_quiz(false, set, id)
        sm:set_property("ox_round", tostring(round(sm) + 1))
        for _, player in ipairs(contestants(sm)) do
            if answered(player, o) then
                player:message("[OX퀴즈] 정답을 맞추었습니다!", Msg.LightBlueText)
                player:exchange(nil, { exp = EXP })
            else
                player:message("[OX퀴즈] 틀렸습니다!", Msg.LightBlueText)
                player:map(LOSER_MAP)
            end
        end
        prepare(sm)
    end,
}

return M
