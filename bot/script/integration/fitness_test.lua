local pq = require("script/integration/lib/party_quest")

local LITH_HARBOR = 104000000
local GUIDE = 9000001
local DIALOG_LIST = 4
local JOIN = 2
local LOBBY = 109040000
-- The lobby and stages 1-4, each left through in00; the last in00 leads to the winner map.
local MAPS = { 109040000, 109040001, 109040002, 109040003, 109040004 }
local WINNER_MAP = 109050000
local ALL_PLAYERS = 4294967295

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function fitness(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 0 - -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(LITH_HARBOR) == false then
        return ctx:fail("리스항구 이동 실패")
    end

    local announced = bot:request(resp.notice, req.normal_chat { message = "/일정이벤트 FITNESS" }, function(p)
        return p.message:find("고지를 향해서", 1, true) ~= nil and p.message:find("잠시 후", 1, true) ~= nil
    end, 15000)
    if check(ctx, announced ~= false, "고지를 향해서 모집 공지가 오지 않음") == false then
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
    if check(ctx, bot:map() == LOBBY, "고지를 향해서 대기 맵으로 이동하지 않음: " .. bot:map()) == false then
        return false
    end

    local closed = bot:request(resp.notice, req.warp { target = ALL_PLAYERS, portal_name = "in00" }, function(p)
        return p.message:find("시작되지 않았", 1, true) ~= nil
    end, 5000)
    if check(ctx, closed ~= false, "시작 전에 1단계 포탈이 열려 있음") == false then
        return false
    end

    local started = bot:request(resp.notice, req.normal_chat { message = "/이벤트시작" }, function(p)
        return p.message:find("포탈이 열렸습니다", 1, true) ~= nil
    end, 60000)
    if check(ctx, started ~= false, "카운트다운 뒤 이벤트가 시작되지 않음") == false then
        return false
    end

    for i, map_id in ipairs(MAPS) do
        local next_map = MAPS[i + 1] or WINNER_MAP
        local moved = bot:request(resp.warp, req.warp { target = ALL_PLAYERS, portal_name = "in00" }, function(p)
            return p.character.map == next_map
        end, 5000)
        if check(ctx, moved ~= false, map_id .. "에서 다음 단계로 이동하지 않음") == false then
            return false
        end
    end
    local finished = check(ctx, bot:map() == WINNER_MAP, "완주했는데 우승자 맵이 아님: " .. bot:map())
    pq.command(bot, "/이벤트취소", "이벤트를 취소했습니다")
    return finished
end

test_suite {
    name = "Fitness: 고지를 향해서 채널 이벤트 (모집·시작 전 포탈·단계 이동·완주)",
    bot_count = 1,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/이벤트취소", "이벤트")
    end,

    scenarios = {
        fitness,
    },
}
