local pq = require("script/integration/lib/party_quest")

local FREE_MARKET = 910000000
local ENTRUSTED_ROOM = 910000009
local PERSONAL_ROOM = 910000002
local FREDRICK = 9030000

local ENTRUSTED_PERMIT = 5030000
local PERSONAL_PERMIT = 5140000
local OWL = 2310000
local CASH_OWL = 5230000
local TARGET = 2000001
local MISSING = 2022003

local MODE = {
	create = 0x00,
	visit = 0x04,
	exit = 0x0A,
	open = 0x0B,
	personal_add_item = 0x12,
	entrusted_add_item = 0x1D,
	close = 0x25,
}
local STORE_BANK = { withdraw = 0x19, confirm = 0x1A, claimed = 0x1D }
local NOTHING_MAP = 999999999

local SPOT = { x = -1000, y = 102 }

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, meso, items)
	local text = string.format("/봇초기화 30 0 %d %s -", meso, items)
	if pq.command(bot, text, "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	return true
end

local function has(ctx, bot, item_id, count)
	local have = bot:items()[item_id] or 0
	return check(ctx, have == count, string.format("%s 인벤토리 %d 개수 %d (기대 %d)", bot:name(), item_id, have, count))
end

local function move_to(ctx, bot, map, spot)
	if bot:map() ~= map and bot:map_move(map) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map)
	end
	if spot ~= nil and pq.move(bot, spot.x, spot.y) == false then
		return ctx:fail(string.format("%s 좌표 이동 실패: %d, %d", bot:name(), spot.x, spot.y))
	end
	return true
end

local function clear_store_bank(ctx, bot)
	if move_to(ctx, bot, FREE_MARKET) == false then
		return false
	end
	local oid = pq.npc(ctx, bot, FREDRICK)
	if oid == false then
		return false
	end
	local p = bot:request({ resp.store_bank_open, resp.store_bank_location }, req.npc_click { oid = oid }, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 프레드릭 응답 없음")
	end
	if p.map_id ~= nil then
		return check(ctx, p.map_id == NOTHING_MAP, bot:name() .. " 이전 테스트의 상점이 영업 중")
	end
	if bot:request(resp.store_bank_fee, req.store_bank { mode = STORE_BANK.withdraw }, nil, 3000) == false then
		return ctx:fail(bot:name() .. " 스토어뱅크 수수료 응답 없음")
	end
	local r = bot:request(resp.store_bank_result, req.store_bank { mode = STORE_BANK.confirm }, nil, 5000)
	return check(ctx, r ~= false and r.result == STORE_BANK.claimed, bot:name() .. " 이전 스토어뱅크 비우기 실패")
end

local function open_shop(ctx, bot, kind, permit, add_mode, price)
	local enter = resp.mini_room_enter
	local list = resp.mini_room_items
	if kind == 4 then
		enter = resp.personal_shop_enter
		list = resp.personal_shop_items
	end
	local create = req.mini_room { mode = MODE.create, type = kind, title = bot:name(), slot = bot:slot(permit) or 0, item_id = permit }
	local p = bot:request({ enter, resp.mini_room_enter_failed }, create, nil, 5000)
	if p == false or p.my_slot ~= 0 then
		return ctx:fail(bot:name() .. " 상점 개설 실패")
	end
	local add = req.mini_room {
		mode = add_mode,
		inventory_type = 2,
		slot = bot:slot(TARGET) or 0,
		bundles = 2,
		per_bundle = 5,
		price = price,
	}
	if bot:request(list, add, nil, 5000) == false then
		return ctx:fail(bot:name() .. " 상점 등록 실패")
	end
	bot:send(req.mini_room { mode = MODE.open })
	ctx:sleep(500)
	return true
end

local function ours(result, names)
	local out = {}
	for _, entry in ipairs(result.entries) do
		if names[entry.owner_name] then
			table.insert(out, entry)
		end
	end
	return out
end

local function search(ctx, bot, request, what)
	local p = bot:request(resp.shop_scanner_result, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 검색 결과 없음: " .. what)
	end
	return p
end

local function search_flow(ctx)
	local seller = ctx:bot(0)
	local merchant = ctx:bot(1)
	local finder = ctx:bot(2)
	if setup(ctx, seller, 0, ENTRUSTED_PERMIT .. ":1," .. TARGET .. ":10") == false then
		return false
	end
	if setup(ctx, merchant, 0, PERSONAL_PERMIT .. ":1," .. TARGET .. ":10") == false then
		return false
	end
	if setup(ctx, finder, 100000, OWL .. ":3," .. CASH_OWL .. ":2") == false then
		return false
	end
	if clear_store_bank(ctx, seller) == false or clear_store_bank(ctx, merchant) == false then
		return false
	end

	if move_to(ctx, seller, ENTRUSTED_ROOM, SPOT) == false or open_shop(ctx, seller, 5, ENTRUSTED_PERMIT, MODE.entrusted_add_item, 500) == false then
		return false
	end
	if move_to(ctx, merchant, PERSONAL_ROOM, SPOT) == false or open_shop(ctx, merchant, 4, PERSONAL_PERMIT, MODE.personal_add_item, 300) == false then
		return false
	end
	local names = { [seller:name()] = true, [merchant:name()] = true }

	if move_to(ctx, finder, FREE_MARKET, nil) == false then
		return false
	end
	if finder:request(resp.shop_scanner_popular, req.shop_scanner_open { mode = 5 }, nil, 5000) == false then
		return ctx:fail("부엉이 창 인기 검색어 응답 없음")
	end

	local owl = function(item_id, high_first)
		return req.use_shop_scanner { slot = finder:slot(OWL) or 0, item_id = OWL, search_id = item_id, high_first = high_first }
	end
	local p = search(ctx, finder, owl(TARGET, false), "가격 오름차순")
	if p == false then
		return false
	end
	local found = ours(p, names)
	if check(ctx, #found == 2 and found[1].price == 300 and found[2].price == 500, "오름차순 검색 결과가 다름: " .. #found) == false then
		return false
	end
	if check(ctx, found[1].owner_name == merchant:name() and found[1].map_id == PERSONAL_ROOM and found[1].bundles == 2 and found[1].per_bundle == 5, "개인상점 검색 항목이 다름") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, finder, OWL, 2) == false then
		return false
	end

	p = search(ctx, finder, owl(TARGET, true), "가격 내림차순")
	if p == false then
		return false
	end
	found = ours(p, names)
	if check(ctx, #found == 2 and found[1].price == 500, "내림차순 검색 결과가 다름") == false then
		return false
	end
	ctx:sleep(300)

	p = search(ctx, finder, owl(MISSING, false), "없는 아이템")
	if p == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, #ours(p, names) == 0, "없는 아이템 검색에 결과가 있음") == false or has(ctx, finder, OWL, 1) == false then
		return false
	end

	local cash = req.use_cash_item { slot = finder:slot(CASH_OWL) or 0, item_id = CASH_OWL, search_id = TARGET, high_first = false }
	p = search(ctx, finder, cash, "캐시 부엉이")
	if p == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, #ours(p, names) == 2, "캐시 부엉이 검색 결과가 다름") == false or has(ctx, finder, CASH_OWL, 1) == false then
		return false
	end

	local popular = finder:request(resp.shop_scanner_popular, req.shop_scanner_open { mode = 5 }, nil, 5000)
	local listed = false
	for _, item_id in ipairs(popular and popular.item_ids or {}) do
		listed = listed or item_id == TARGET
	end
	if check(ctx, listed, "인기 검색어에 검색한 아이템이 없음") == false then
		return false
	end

	local entered = finder:request(resp.personal_shop_enter, req.shop_scanner_warp { sn = found[2].sn, map_id = PERSONAL_ROOM }, nil, 10000)
	if check(ctx, entered ~= false and entered.my_slot == 1 and finder:map() == PERSONAL_ROOM, "검색 결과로 개인상점 입장 실패") == false then
		return false
	end
	finder:send(req.mini_room { mode = MODE.exit })

	merchant:send(req.mini_room { mode = MODE.exit })
	local sn = found[1].sn
	if seller:request(resp.mini_room_enter, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000) == false then
		return ctx:fail("고용상점 관리 입장 실패")
	end
	seller:send(req.mini_room { mode = MODE.close })
	ctx:sleep(500)
	if has(ctx, seller, TARGET, 10) == false or has(ctx, merchant, TARGET, 10) == false then
		return false
	end

	if move_to(ctx, finder, FREE_MARKET, nil) == false then
		return false
	end
	p = search(ctx, finder, owl(TARGET, false), "상점 종료 후")
	if p == false then
		return false
	end
	return check(ctx, #ours(p, names) == 0, "닫은 상점이 검색 결과에 남아 있음")
end

test_suite {
	name = "ShopScanner: 미네르바의 부엉이 검색·인기 검색어·상점 이동",
	bot_count = 3,

	scenarios = {
		search_flow,
	},
}
