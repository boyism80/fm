local pq = require("script/integration/lib/party_quest")

local FREE_MARKET = 910000000
local TOWNS = {
	{ map = 100000100, portal = "market00", x = 838, y = 155 },
	{ map = 102000000, portal = "market00", x = 1157, y = 581 },
	{ map = 120000200, portal = "market00", x = 1867, y = 148 },
	{ map = 200000000, portal = "market00", x = -1181, y = -458 },
	{ map = 211000100, portal = "market00", x = -390, y = 94 },
	{ map = 220000000, portal = "market00", x = -1064, y = -613 },
	{ map = 221000000, portal = "market00", x = 3220, y = -144 },
	{ map = 222000000, portal = "market00", x = -1399, y = 149 },
	{ map = 230000000, portal = "market01", x = 1511, y = 38 },
	{ map = 240000000, portal = "market00", x = -1936, y = -32 },
	{ map = 250000000, portal = "market00", x = -197, y = -547 },
	{ map = 251000000, portal = "market00", x = 1020, y = -66 },
	{ map = 260000000, portal = "market00", x = 406, y = 278 },
	{ map = 261000000, portal = "market0", x = 336, y = -227 },
	{ map = 500000000, portal = "market00", x = 2015, y = 151 },
	{ map = 540000000, portal = "market00", x = 4582, y = 43 },
	{ map = 541000000, portal = "market00", x = 34, y = -109 },
	{ map = 550000000, portal = "market00", x = 2685, y = 649 },
	{ map = 551000000, portal = "market00", x = -233, y = 132 },
	{ map = 801000300, portal = "market00", x = 258, y = 151 },
}

local function setup(ctx, bot, level)
	if pq.command(bot, "/봇초기화 " .. level .. " 0 0 - -", "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패")
	end
	return true
end

local function walk(ctx, bot, x, y)
	if pq.move(bot, x, y) == false then
		return ctx:fail(string.format("포탈 앞으로 이동 실패 (%d, %d)", x, y))
	end
	return true
end

local scenarios = {
	function(ctx)
		local bot = ctx:bot(0)
		local town = TOWNS[2]
		if setup(ctx, bot, 7) == false then
			return false
		end
		if bot:map_move(town.map) == false then
			return ctx:fail("마을로 이동 실패: " .. town.map)
		end
		if walk(ctx, bot, town.x, town.y) == false then
			return false
		end
		local p = bot:request(resp.notice, req.warp { target = 4294967295, portal_name = town.portal }, function(p)
			return p.message:find("레벨 8 이상", 1, true) ~= nil
		end, 5000)
		if p == false then
			return ctx:fail("레벨 7 자유시장 입장이 거절되지 않음 (현재 " .. bot:map() .. ")")
		end
		return true
	end,
}
for _, town in ipairs(TOWNS) do
	scenarios[#scenarios + 1] = function(ctx)
		local bot = ctx:bot(0)
		if setup(ctx, bot, 30) == false then
			return false
		end
		if bot:map_move(town.map) == false then
			return ctx:fail("마을로 이동 실패: " .. town.map)
		end
		if walk(ctx, bot, town.x, town.y) == false then
			return false
		end
		if pq.portal(ctx, bot, town.portal, FREE_MARKET) == false then
			return false
		end
		if walk(ctx, bot, -120, 39) == false then
			return false
		end
		if pq.portal(ctx, bot, "out00", town.map) == false then
			return false
		end
		local x, y = bot:position()
		if math.abs(x - town.x) > 50 then
			return ctx:fail(string.format("자유시장 포탈 앞으로 돌아오지 않음: %d (%d, %d)", town.map, x, y))
		end
		return true
	end
end

test_suite {
	name = "Free Market: 자유시장 입장·퇴장",
	bot_count = 1,

	on_initialize = function(ctx)
		if pq.command(ctx:bot(0), "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return true
	end,

	scenarios = scenarios,
}
