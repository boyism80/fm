-- Channel event game (old/src/server/events/MapleFitness.java): 고지를 향해서

local M = {}

-- The lobby, then stages 1-4; each in00 portal leads up, and the last one to the winner map.
local MAPS = { 109040000, 109040001, 109040002, 109040003, 109040004 }
local DURATION_MS = 900000
local HINT_MS = 60000
-- Each hint is shown once a minute while the time left, in seconds, is above its mark.
local HINTS = {
    { 841, "[고지를 향해서] 는 총 4단계의 스테이지로 구성되어 있습니다. 게임 중 캐릭터가 죽게되면 탈락이니 주의해 주세요." },
    { 781, "1단계는 사용자가 많이 몰리는 관계로 렉이 조금 심합니다…1단계는 어렵지 않으니 렉 때문에 아래로 떨어지지 않도록 천천히 이동해 주세요." },
    { 701, "고지를 향해서는 성공 순위에 관계없이 시간내에 성공하신 모든 분들께 상품이 지급되니 너무 서두르지 마시고 이동해 주세요." },
    { 661, "게임 중 캐릭터가 죽게 되면 탈락이니 죽지 않게 조심해 주세요. 아직 시간은 충분하니 에너지가 없는 분들은 물약을 드시거나 HP를 회복 한 뒤에 게임에 임해주세요." },
    { 601, "2단계 에서는 원숭이들이 던지는 바나나를 잘 피해가셔야 합니다. 타이밍을 잘 맞춰서 바나나를 위해 위로 쭉쭉 이동해 주세요." },
    { 501, "원숭이들이 던지는 바나나를 피하기 위해 가장 중요한 것은 \"타이밍\" 바로 타이밍 입니다. 힘내세요." },
    { 361, "게임 중 캐릭터가 죽게 되면 탈락이니 죽지 않게 조심해 주세요. 아직 시간은 충분하니 에너지가 없는 분들은 물약을 드시거나 HP를 회복 한 뒤에 게임에 임해주세요." },
    { 301, "렉이 심한 분들은 렉 때문에 아래로 떨어지지 않도록 천천히 이동해 주세요" },
    { 201, "3단계에는 눈에는 보이지만 실제로 밟을 수 없는 함정이 있습니다. 그 함정을 조심히 이동해 주세요." },
    { 101, "4단계는 [고지를 향해서]의 마지막 단계입니다. 끝까지 포기하지 마시고 최선을 다해주세요. 위에서 상품이 기다리고 있습니다^^" },
    { 11, "자 이제 시간이 얼마 남지 않았습니다. 조금 서둘러 주세요~" },
    { 0, "제한시간이 10초 남았습니다. 성공하지 못 하신 분들은 다음 이벤트에서 꼭 성공하시길 바랍니다. 그럼 안녕히 가세요." },
}

M.MAPS = MAPS

function M.prepare(sm)
    sm:map(MAPS[1]):portal("in00"):script("event_closed")
end

local function contestants(sm)
    local players = {}
    for _, map_id in ipairs(MAPS) do
        for _, player in pairs(sm:map(map_id):characters()) do
            players[#players + 1] = player
        end
    end
    return players
end

local function notice(sm, text)
    for _, player in ipairs(contestants(sm)) do
        player:message(text, Msg.Notice)
    end
end

function M.start(sm)
    sm:map(MAPS[1]):portal("in00"):script("")
    for _, player in ipairs(contestants(sm)) do
        player:clock(DURATION_MS / 1000)
    end
    notice(sm, "포탈이 열렸습니다. 포탈 위에서 ↑를 눌러주세요.")
    sm:set_property("clock_until", tostring(now() + DURATION_MS / 1000))
    sm:after("fitness", DURATION_MS, "on_fitness_end")
    sm:after("fitness_hint", HINT_MS, "on_fitness_hint")
end

-- Called when a contestant walks through the last stage's portal into the winner map.
function M.finish(sm, me)
    if sm:get_property("phase") ~= "running" then
        return
    end
    sm:set_property("winner:" .. me:id(), "true")
    me:stop_clock()
    me:message("축하합니다! [고지를향해서] 이벤트를 완주하셨습니다! NPC [피에트로] 를 더블클릭하면 보상을 받을 수 있습니다.", Msg.Notice)
    if #contestants(sm) == 0 then
        sm:cancel("fitness")
        sm:cancel("fitness_hint")
        sm:call_hook("on_game_over")
    end
end

M.hooks = {
    on_fitness_end = function(sm)
        sm:cancel("fitness_hint")
        sm:call_hook("on_game_over")
    end,

    on_fitness_hint = function(sm)
        local left = (tonumber(sm:get_property("clock_until")) or 0) - now()
        for _, hint in ipairs(HINTS) do
            if left > hint[1] then
                notice(sm, hint[2])
                break
            end
        end
        sm:after("fitness_hint", HINT_MS, "on_fitness_hint")
    end,
}

return M
