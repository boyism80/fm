local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local BUDDY = { add = 1, accept = 2 }
local TRANSITION_TIMEOUT_MS = 30000

-- A lost presence notice must not turn a channel change into a failed login.
local function enter_without_mq(ctx)
    local a, b = ctx:bot(0), ctx:bot(1)
    for _, bot in ipairs({ a, b }) do
        if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
            return ctx:fail(bot:name() .. " 봇 초기화 실패")
        end
        if bot:map_move(HENESYS) == false then
            return ctx:fail(bot:name() .. " 맵 이동 실패")
        end
    end
    if pq.form_party(ctx) == false then
        return false
    end

    if pq.command(a, "/MQ장애 on", "MQ 장애: ON") == false then
        return ctx:fail("MQ 장애 설정 실패")
    end
    local moved = a:transfer(1)
    if pq.command(b, "/MQ장애 off", "MQ 장애: OFF") == false then
        return ctx:fail("MQ 장애 해제 실패")
    end
    if moved == false then
        return ctx:fail("MQ 발행이 끊긴 동안 파티원이 채널을 옮기지 못함")
    end

    if a:transfer(0) == false then
        return ctx:fail("MQ 복구 뒤 원래 채널로 돌아오지 못함")
    end
    for _, bot in ipairs({ b, a }) do
        bot:request(resp.party_update_leave, req.party_operation { operation = PARTY.Leave }, nil, 3000)
    end
    return true
end

-- A character that leaves a channel and never reaches the next one is logged out once the move times out, so its
-- buddies stop seeing it on the old channel.
local function abandoned_transfer(ctx)
    local a, b = ctx:bot(0), ctx:bot(1)
    for _, bot in ipairs({ a, b }) do
        if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
            return ctx:fail(bot:name() .. " 봇 초기화 실패")
        end
    end
    b:send(req.buddy { mode = BUDDY.add, name = a:name(), group = "봇친구" })
    ctx:sleep(1500)
    a:send(req.buddy { mode = BUDDY.accept, character_id = b:id() })
    ctx:sleep(1500)

    local channel = nil
    ctx:hook("abandoned_transfer", function(_, bot, p, name)
        if bot:id() == b:id() and name == "buddy_channel_update" and p.character_id == a:id() then
            channel = p.channel
        end
    end)
    if a:abandon_transfer(1) == false then
        ctx:unhook("abandoned_transfer")
        return ctx:fail("채널 이동 요청 실패")
    end
    for _ = 1, (TRANSITION_TIMEOUT_MS + 15000) / 500 do
        if channel == -1 then
            break
        end
        ctx:sleep(500)
    end
    ctx:unhook("abandoned_transfer")
    if channel ~= -1 then
        return ctx:fail("새 채널에 오지 않은 친구가 접속 종료로 보이지 않음: channel=" .. tostring(channel))
    end
    if a:relog({ passwords = { "" } }) == false then
        return ctx:fail("이동 시간 초과 뒤 다시 접속하지 못함")
    end
    return true
end

test_suite {
    name = "Session presence: MQ 발행 장애 중 채널 이동·도착하지 않은 채널 이동",
    bot_count = 2,

    scenarios = {
        enter_without_mq,
        abandoned_transfer,
    },
}
