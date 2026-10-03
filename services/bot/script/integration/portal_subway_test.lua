local STATION = 103000310
local PORTAL = "out00"
local SUBWAY = 103000301
local ARRIVE = 103000100

test_suite {
	name = "Portal: 커닝 지하철",
	bot_count = 1,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(STATION) == false then
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
	},
}
