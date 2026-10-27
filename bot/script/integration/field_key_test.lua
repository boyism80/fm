local pq = require("script/integration/lib/party_quest")

local FIRST = 100000000
local SECOND = 104000000
local FAR_X = 2500
local FAR_Y = -500

local function server_position(ctx, bot)
    local p = pq.command(bot, "/좌표", "Position:")
    if p == false then
        return nil
    end
    local x, y = p.message:match("Position: (%-?%d+), (%-?%d+)")
    return tonumber(x), tonumber(y)
end

-- A move the client made on the field it left keeps that field's key, and the server drops it instead of moving the
-- character on the new field; a move with the new key still counts.
local function stale_move(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 100 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:force_instance_move(FIRST, "필드 키 시험 첫 맵") == false then
        return ctx:fail("첫 맵 이동 실패")
    end
    local old = bot:field_key()
    bot:move(FAR_X, FAR_Y)
    local x, y = server_position(ctx, bot)
    if x ~= FAR_X or y ~= FAR_Y then
        return ctx:fail(string.format("현재 필드 키의 이동이 반영되지 않음: %s, %s", tostring(x), tostring(y)))
    end

    if bot:force_instance_move(SECOND, "서버가 다른 맵으로 워프") == false then
        return ctx:fail("두 번째 맵 이동 실패")
    end
    if bot:field_key() == old then
        return ctx:fail("맵에 다시 들어왔는데 필드 키가 그대로임: " .. old)
    end
    bot:move(FAR_X, FAR_Y, 0, old)
    x, y = server_position(ctx, bot)
    if x == FAR_X and y == FAR_Y then
        return ctx:fail("이전 맵의 필드 키를 가진 이동이 새 맵 위치를 덮어씀")
    end

    bot:move(FAR_X, FAR_Y)
    x, y = server_position(ctx, bot)
    if x ~= FAR_X or y ~= FAR_Y then
        return ctx:fail(string.format("새 필드 키의 이동이 반영되지 않음: %s, %s", tostring(x), tostring(y)))
    end
    return true
end

test_suite {
    name = "Field key: 맵 전환 전 이동 무시",
    bot_count = 1,

    scenarios = {
        stale_move,
    },
}
