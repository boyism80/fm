local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local ELLINIA = 101000000
local ORBIS = 200000000
local AMORIA = 680000000
local STONE = 5040000
local VIP_STONE = 5041000
local CONSUME_STONE = 2320000
local EMPTY = 999999999

local RESULT_LIST = 3
local RESULT_CANNOT_GO = 5
local RESULT_NOT_FOUND = 6
local RESULT_CURRENT_MAP = 9
local RESULT_CANNOT_REGISTER = 10

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function contains(maps, map_id)
	for _, registered in ipairs(maps or {}) do
		if registered == map_id then
			return true
		end
	end
	return false
end

local function count(maps, map_id)
	local n = 0
	for _, registered in ipairs(maps or {}) do
		if registered == map_id then
			n = n + 1
		end
	end
	return n
end

local function setup(ctx, bot, items)
	if pq.command(bot, "/봇초기화 30 0 0 " .. items .. " -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if pq.command(bot, "/순간이동초기화", "순간이동 초기화 완료") == false then
		return ctx:fail(bot:name() .. " 순간이동 초기화 실패")
	end
	return true
end

local function move(ctx, bot, map_id)
	if bot:map_move(map_id) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

local function register(ctx, bot, vip, what)
	local p = bot:request(resp.teleport_stone_result, req.teleport_stone_list { action = 1, vip = vip }, nil, 3000)
	if p == false then
		return ctx:fail("등록 응답 없음: " .. what)
	end
	return p
end

local function use(ctx, bot, item_id, target, what)
	local use_req
	if item_id == CONSUME_STONE then
		use_req = req.use_teleport_stone { slot = bot:slot(item_id) or 0, item_id = item_id, target = target }
	else
		use_req = req.use_cash_item { slot = bot:slot(item_id) or 0, item_id = item_id, target = target }
	end
	local p = bot:request({ resp.warp, resp.teleport_stone_result }, use_req, nil, 3000)
	if p == false then
		return ctx:fail("사용 응답 없음: " .. what)
	end
	ctx:sleep(300)
	return p
end

local function expect_moved(ctx, bot, p, item_id, before, map_id, what)
	if check(ctx, p.result == nil, what .. ": 이동 실패 " .. tostring(p.result)) == false then
		return false
	end
	if check(ctx, bot:map() == map_id, what .. ": 도착 맵이 다름 " .. bot:map()) == false then
		return false
	end
	return check(ctx, (bot:items()[item_id] or 0) == before - 1, what .. ": 아이템이 1개 줄지 않음")
end

local function expect_failed(ctx, bot, p, result, item_id, before, what)
	if check(ctx, p.result == result, string.format("%s: 결과 코드가 다름 %s/%d", what, tostring(p.result), result)) == false then
		return false
	end
	return check(ctx, (bot:items()[item_id] or 0) == before, what .. ": 실패했는데 아이템이 줄어듦")
end

local function register_and_remove(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, STONE .. ":5," .. VIP_STONE .. ":2," .. CONSUME_STONE .. ":1") == false then
		return false
	end
	if move(ctx, bot, HENESYS) == false then
		return false
	end

	local p = register(ctx, bot, false, "헤네시스")
	if p == false then
		return false
	end
	if check(ctx, p.result == RESULT_LIST and #p.maps == 5 and p.maps[1] == HENESYS, "헤네시스 등록 목록이 다름") == false then
		return false
	end
	if check(ctx, contains(p.maps, EMPTY), "빈 칸이 999999999로 오지 않음") == false then
		return false
	end

	p = register(ctx, bot, false, "중복")
	if p == false then
		return false
	end
	if check(ctx, p.result == RESULT_LIST and count(p.maps, HENESYS) == 1, "같은 맵이 두 번 등록됨") == false then
		return false
	end

	if move(ctx, bot, AMORIA) == false then
		return false
	end
	p = register(ctx, bot, false, "아모리아")
	if p == false then
		return false
	end
	if check(ctx, p.result == RESULT_CANNOT_REGISTER, "순간이동 제한 맵이 등록됨: " .. p.result) == false then
		return false
	end

	if move(ctx, bot, ELLINIA) == false then
		return false
	end
	if register(ctx, bot, false, "엘리니아") == false then
		return false
	end
	p = bot:request(resp.teleport_stone_result, req.teleport_stone_list { action = 0, vip = false, map_id = ELLINIA }, nil, 3000)
	if p == false then
		return ctx:fail("삭제 응답 없음")
	end
	return check(ctx, p.result == RESULT_LIST and contains(p.maps, ELLINIA) == false and contains(p.maps, HENESYS), "삭제 후 목록이 다름")
end

local function move_by_map(ctx)
	local bot = ctx:bot(0)
	if move(ctx, bot, ELLINIA) == false then
		return false
	end

	local before = bot:items()[STONE] or 0
	local p = use(ctx, bot, STONE, { map_id = HENESYS }, "헤네시스로")
	if p == false or expect_moved(ctx, bot, p, STONE, before, HENESYS, "헤네시스로") == false then
		return false
	end

	before = bot:items()[STONE] or 0
	p = use(ctx, bot, STONE, { map_id = HENESYS }, "현재 맵")
	if p == false or expect_failed(ctx, bot, p, RESULT_CURRENT_MAP, STONE, before, "현재 맵") == false then
		return false
	end

	p = use(ctx, bot, STONE, { map_id = ELLINIA }, "등록 안 한 맵")
	if p == false or expect_failed(ctx, bot, p, RESULT_CANNOT_GO, STONE, before, "등록 안 한 맵") == false then
		return false
	end

	if move(ctx, bot, ORBIS) == false then
		return false
	end
	if register(ctx, bot, false, "오르비스 일반") == false or register(ctx, bot, true, "오르비스 VIP") == false then
		return false
	end
	if move(ctx, bot, HENESYS) == false then
		return false
	end
	p = use(ctx, bot, STONE, { map_id = ORBIS }, "다른 대륙 일반")
	if p == false or expect_failed(ctx, bot, p, RESULT_CANNOT_GO, STONE, before, "다른 대륙 일반") == false then
		return false
	end

	before = bot:items()[VIP_STONE] or 0
	p = use(ctx, bot, VIP_STONE, { map_id = ORBIS }, "다른 대륙 VIP")
	if p == false or expect_moved(ctx, bot, p, VIP_STONE, before, ORBIS, "다른 대륙 VIP") == false then
		return false
	end

	if move(ctx, bot, ELLINIA) == false then
		return false
	end
	before = bot:items()[CONSUME_STONE] or 0
	p = use(ctx, bot, CONSUME_STONE, { map_id = HENESYS }, "소비 순간이동의 돌")
	if p == false then
		return false
	end
	return expect_moved(ctx, bot, p, CONSUME_STONE, before, HENESYS, "소비 순간이동의 돌")
end

local function move_by_name(ctx)
	local bot = ctx:bot(0)
	local target = ctx:bot(1)
	if setup(ctx, target, "-") == false or move(ctx, target, ELLINIA) == false then
		return false
	end
	if move(ctx, bot, HENESYS) == false then
		return false
	end

	local before = bot:items()[STONE] or 0
	local p = use(ctx, bot, STONE, { name = target:name() }, "캐릭터에게")
	if p == false or expect_moved(ctx, bot, p, STONE, before, ELLINIA, "캐릭터에게") == false then
		return false
	end

	before = bot:items()[STONE] or 0
	p = use(ctx, bot, STONE, { name = target:name() }, "같은 맵 캐릭터")
	if p == false or expect_failed(ctx, bot, p, RESULT_CURRENT_MAP, STONE, before, "같은 맵 캐릭터") == false then
		return false
	end

	p = use(ctx, bot, STONE, { name = "없는캐릭터" }, "없는 캐릭터")
	if p == false then
		return false
	end
	return expect_failed(ctx, bot, p, RESULT_NOT_FOUND, STONE, before, "없는 캐릭터")
end

local function persistence(ctx)
	local bot = ctx:bot(0)
	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	local stones = bot:teleport_stones(false)
	local vip = bot:teleport_stones(true)
	if check(ctx, #stones == 5 and contains(stones, HENESYS) and contains(stones, ORBIS), "채널 이동 후 일반 목록이 다름") == false then
		return false
	end
	return check(ctx, #vip == 10 and contains(vip, ORBIS), "채널 이동 후 VIP 목록이 다름")
end

test_suite {
	name = "Teleport stone: 등록·삭제·이동·저장",
	bot_count = 2,

	scenarios = {
		register_and_remove,
		move_by_map,
		move_by_name,
		persistence,
	},
}
