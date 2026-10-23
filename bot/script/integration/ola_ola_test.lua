local pq = require("script/integration/lib/party_quest")

local LITH_HARBOR = 104000000
local GUIDE = 9000001
local DIALOG_LIST = 4
local JOIN = 2
local FIRST_STAGE = 109030001
local WINNER_MAP = 109050000
-- Stage 1 has portals ch00-ch04, stage 2 ch00-ch07 and stage 3 ch00-ch15.
local PORTALS = { 5, 8, 16 }
local ALL_PLAYERS = 4294967295

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function step(bot, portal)
    bot:send(req.warp { target = ALL_PLAYERS, portal_name = portal })
end

-- Every portal of a stage is tried; the wrong ones send the bot back to the start of the map.
local function climb(ctx, bot, stage)
    local from = bot:map()
    for i = 0, PORTALS[stage] - 1 do
        step(bot, string.format("ch%02d", i))
        ctx:sleep(800)
        if bot:map() ~= from then
            return i + 1
        end
    end
    return ctx:fail("올라올라 " .. stage .. "단계의 올바른 포탈을 찾지 못함")
end

local function ola_ola(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 0 - -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(LITH_HARBOR) == false then
        return ctx:fail("리스항구 이동 실패")
    end

    local announced = bot:request(resp.notice, req.normal_chat { message = "/일정이벤트 OLA" }, function(p)
        return p.message:find("올라올라", 1, true) ~= nil and p.message:find("잠시 후", 1, true) ~= nil
    end, 15000)
    if check(ctx, announced ~= false, "올라올라 모집 공지가 오지 않음") == false then
        return false
    end

    local oid = pq.npc(ctx, bot, GUIDE)
    if oid == false then
        return false
    end
    if bot:npc_click(oid) == false then
        return ctx:fail("도우미 대화가 열리지 않음")
    end
    if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = JOIN }, nil, 5000) == false then
        return ctx:fail("이벤트 참가로 이동하지 않음")
    end
    if check(ctx, bot:map() == FIRST_STAGE, "올라올라 1단계 맵으로 이동하지 않음: " .. bot:map()) == false then
        return false
    end

    local closed = bot:request(resp.notice, req.warp { target = ALL_PLAYERS, portal_name = "join00" }, function(p)
        return p.message:find("시작되지 않았", 1, true) ~= nil
    end, 5000)
    if check(ctx, closed ~= false, "시작 전에 출발 포탈이 열려 있음") == false then
        return false
    end

    local started = bot:request(resp.notice, req.normal_chat { message = "/이벤트시작" }, function(p)
        return p.message:find("포탈이 열렸습니다", 1, true) ~= nil
    end, 60000)
    if check(ctx, started ~= false, "카운트다운 뒤 이벤트가 시작되지 않음") == false then
        return false
    end

    step(bot, "join00")
    ctx:sleep(800)
    local tries = 0
    for stage = 1, #PORTALS do
        local used = climb(ctx, bot, stage)
        if used == false then
            return false
        end
        tries = tries + used
    end
    -- Every stage solved on the first try has a chance of one in 640: a portal that always passed would show here.
    if check(ctx, tries > #PORTALS, "모든 포탈이 통과시킴: 시도 " .. tries) == false then
        return false
    end
    local finished = check(ctx, bot:map() == WINNER_MAP, "완주했는데 우승자 맵이 아님: " .. bot:map())
    pq.command(bot, "/이벤트취소", "이벤트를 취소했습니다")
    return finished
end

test_suite {
    name = "Ola Ola: 올라올라 채널 이벤트 (모집·시작 전 포탈·단계별 올바른 포탈·완주)",
    bot_count = 1,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/이벤트취소", "이벤트")
    end,

    scenarios = {
        ola_ola,
    },
}
