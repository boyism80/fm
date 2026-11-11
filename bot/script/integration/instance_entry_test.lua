local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")

local LOBBY = 103040400
local CLERK = 1052125
local ZONE_FIRST = 103040410
local ZONE_LAST = 103040439
local HELPING_HYUK = 0
local DIALOG_NORMAL = 0

local function in_zone(map)
    return map >= ZONE_FIRST and map <= ZONE_LAST
end

-- Two characters asking the clerk at once never share a general zone: one entry at a time goes through, and the
-- other is sent to another zone or turned away.
local function concurrent_zones(ctx)
    local first, second = ctx:bot(0), ctx:bot(1)
    for _, bot in ipairs({ first, second }) do
        if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
            return ctx:fail(bot:name() .. " 봇 초기화 실패")
        end
        if route.travel_instance(ctx, bot, LOBBY) == false then
            return ctx:fail(bot:name() .. " 백화점 이동 실패")
        end
    end
    for _, bot in ipairs({ first, second }) do
        local oid = pq.npc(ctx, bot, CLERK)
        if oid == false then
            return false
        end
        local list = bot:npc_click(oid)
        if list == false or list == nil or list.selections == nil then
            return ctx:fail(bot:name() .. " 준이 선택지가 오지 않음")
        end
        local welcome = bot:dialog(true, HELPING_HYUK)
        if welcome == nil or welcome.text == nil or welcome.text:find("일반존", 1, true) == nil then
            return ctx:fail(bot:name() .. " 일반존 안내가 오지 않음")
        end
    end

    first:send(req.dialog { dialog_type = DIALOG_NORMAL, next = true })
    second:send(req.dialog { dialog_type = DIALOG_NORMAL, next = true })
    for _ = 1, 30 do
        if in_zone(first:map()) and in_zone(second:map()) then
            break
        end
        ctx:sleep(100)
    end
    local a, b = first:map(), second:map()
    for _, bot in ipairs({ first, second }) do
        bot:dialog(false)
    end
    if in_zone(a) == false and in_zone(b) == false then
        return ctx:fail("아무도 일반존에 들어가지 못함: " .. a .. ", " .. b)
    end
    if a == b then
        return ctx:fail("두 캐릭터가 같은 일반존에 들어감: " .. a)
    end
    return true
end

test_suite {
    name = "Instance entry: 커닝 백화점 일반존 동시 입장",
    bot_count = 2,

    scenarios = {
        concurrent_zones,
    },
}
