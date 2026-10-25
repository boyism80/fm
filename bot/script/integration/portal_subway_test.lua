local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")

local STATION = 103000310
local PORTAL = "out00"
local SUBWAY = 103000301
local ARRIVE = 103000100
local HENESYS = 100000000
local FAR_PORTAL = "east00"
local FAR_PORTAL_TARGET = 100010000

test_suite {
    name = "Portal: 커닝 지하철·먼 포탈 거절",
    bot_count = 1,

    scenarios = {
        function(ctx)
            local bot = ctx:bot(0)
            if route.travel(ctx, bot, STATION) == false then
                return ctx:fail("커닝 스퀘어 역 이동 실패")
            end

            local warp = bot:warp(PORTAL)
            if warp == false then
                return ctx:fail("포탈 이동 응답 없음")
            end
            if warp.character.map ~= SUBWAY then
                return ctx:fail("지하철 객차가 아닌 맵으로 이동: " .. warp.character.map)
            end

            local arrive = bot:request(resp.warp, nil, function(p)
                return p.character.map == ARRIVE
            end, 25000)
            if arrive == false then
                return ctx:fail("20초 뒤 도착 맵으로 이동하지 않음")
            end
            return true
        end,

        -- A portal is taken only within 200 pixels of it, so one named from across the map is refused.
        function(ctx)
            local bot = ctx:bot(0)
            if bot:instance_move(HENESYS) == false then
                return ctx:fail("헤네시스 이동 실패")
            end
            if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
                return ctx:fail("플레이어 모드 설정 실패")
            end
            if bot:request(resp.warp, req.direct_warp { portal_name = FAR_PORTAL }, nil, 1500) ~= false then
                return ctx:fail("맵 반대편의 포탈로 이동함")
            end
            local warp = bot:warp(FAR_PORTAL)
            if warp == false or warp.character.map ~= FAR_PORTAL_TARGET then
                return ctx:fail("포탈 앞에서 탄 포탈로 이동하지 않음")
            end
            return true
        end,
    },
}
