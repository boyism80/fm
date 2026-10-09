local pq = require("script/integration/lib/party_quest")

local FREE_MARKET = 910000000
local MARKET_ROOM = 910000007
local FREDRICK = 9030000

local PERMIT = 5030000
local POTION = 2000000
local SWORD = 1302000

local CONSUME = 2

local MODE = {
	create = 0x00,
	visit = 0x04,
	open = 0x0B,
	add_item = 0x1D,
	buy = 0x1E,
	remove_item = 0x22,
	close = 0x25,
	withdraw_meso = 0x27,
}
local CHECK = { title = 7, already_open = 8, store_bank_full = 9, cannot_open = 11 }
local ENTER_ERROR = { near_portal = 10, organizing = 16 }
local LEAVE = { organizing = 13 }
local BUY = { not_enough_item = 1 }
local STORE_BANK = { withdraw = 0x19, confirm = 0x1A, claimed = 0x1D }
local NOTHING_MAP = 999999999

local SPOT = { x = -1000, y = 102 }
local PORTAL_SPOT = { x = -360, y = 102 }

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

local function meso_is(ctx, bot, meso)
	return check(ctx, bot:meso() == meso, string.format("%s 메소 %d (기대 %d)", bot:name(), bot:meso(), meso))
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

local function shop_check(ctx, bot, result, what)
	local p = bot:request(resp.entrusted_shop_check_result, req.use_entrusted_shop {}, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 개설 문의 응답 없음: " .. what)
	end
	if p.result ~= result then
		return ctx:fail(string.format("%s %s: 결과 %d (기대 %d)", bot:name(), what, p.result, result))
	end
	return p
end

local function create_request(bot, title)
	return req.mini_room {
		mode = MODE.create,
		type = 5,
		title = title,
		slot = bot:slot(PERMIT) or 0,
		item_id = PERMIT,
	}
end

local function enter_failed(ctx, bot, request, code, what)
	local p = bot:request({ resp.mini_room_enter, resp.mini_room_enter_failed }, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 응답 없음: " .. what)
	end
	if p.error ~= code then
		return ctx:fail(string.format("%s %s: 입장 에러 %s (기대 %d)", bot:name(), what, tostring(p.error), code))
	end
	return true
end

local function entered(ctx, bot, request, slot, what)
	local p = bot:request({ resp.mini_room_enter, resp.mini_room_enter_failed }, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 응답 없음: " .. what)
	end
	if p.my_slot == nil or p.my_slot ~= slot then
		return ctx:fail(string.format("%s %s: 입장 실패 (에러 %s, 슬롯 %s)", bot:name(), what, tostring(p.error), tostring(p.my_slot)))
	end
	return p
end

local function add_item(bot, item_id, bundles, per_bundle, price)
	return req.mini_room {
		mode = MODE.add_item,
		inventory_type = math.floor(item_id / 1000000),
		slot = bot:slot(item_id) or 0,
		bundles = bundles,
		per_bundle = per_bundle,
		price = price,
	}
end

local function items(ctx, bot, request, count, what)
	local p = bot:request(resp.mini_room_items, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 상점 목록 갱신 없음: " .. what)
	end
	if #p.items ~= count then
		return ctx:fail(string.format("%s %s: 등록 수 %d (기대 %d)", bot:name(), what, #p.items, count))
	end
	return p
end

local function claim_store_bank(ctx, bot)
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
	return p
end

local function clear_store_bank(ctx, bot)
	local p = claim_store_bank(ctx, bot)
	if p == false then
		return false
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

local function shop_flow(ctx)
	local owner = ctx:bot(0)
	local buyer = ctx:bot(1)
	if setup(ctx, owner, 0, PERMIT .. ":1," .. POTION .. ":100," .. SWORD .. ":1") == false then
		return false
	end
	if setup(ctx, buyer, 100000, "-") == false then
		return false
	end
	if clear_store_bank(ctx, owner) == false then
		return false
	end

	if shop_check(ctx, owner, CHECK.cannot_open, "자유시장 입구에서 개설") == false then
		return false
	end

	if move_to(ctx, buyer, MARKET_ROOM, { x = -1300, y = 102 }) == false then
		return false
	end
	if move_to(ctx, owner, MARKET_ROOM, PORTAL_SPOT) == false then
		return false
	end
	if shop_check(ctx, owner, CHECK.title, "방에서 개설 문의") == false then
		return false
	end
	if enter_failed(ctx, owner, create_request(owner, "포털 앞"), ENTER_ERROR.near_portal, "포털 앞 개설") == false then
		return false
	end
	if move_to(ctx, owner, MARKET_ROOM, SPOT) == false then
		return false
	end
	local p = entered(ctx, owner, create_request(owner, "봇 상점"), 0, "개설")
	if p == false then
		return false
	end
	if check(ctx, p.title == "봇 상점" and p.max_items == 10 and #p.mini_room_items.items == 0, "개설 입장 정보가 다름") == false then
		return false
	end

	if items(ctx, owner, add_item(owner, POTION, 10, 5, 1000), 1, "물약 등록") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, SWORD, 1, 1, 50000), 2, "검 등록") == false then
		return false
	end
	if has(ctx, owner, POTION, 50) == false or has(ctx, owner, SWORD, 0) == false then
		return false
	end
	if items(ctx, owner, req.mini_room { mode = MODE.remove_item, index = 1 }, 1, "검 회수") == false then
		return false
	end
	if has(ctx, owner, SWORD, 1) == false then
		return false
	end

	local spawn = owner:request_on(buyer, resp.spawn_entrusted_shop, req.mini_room { mode = MODE.open }, function(s)
		return s.employer_id == owner:id()
	end, 5000)
	if spawn == false then
		return ctx:fail("상점 개설이 다른 캐릭터에게 보이지 않음")
	end
	local sn = spawn.entrusted_shop_balloon.sn
	if check(ctx, spawn.owner_name == owner:name() and spawn.entrusted_shop_balloon.title == "봇 상점", "상점 말풍선 정보가 다름") == false then
		return false
	end

	p = shop_check(ctx, owner, CHECK.already_open, "영업 중 다시 개설 문의")
	if p == false then
		return false
	end
	if check(ctx, p.map_id == MARKET_ROOM, "영업 중 상점 위치가 다름: " .. tostring(p.map_id)) == false then
		return false
	end

	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = sn }, 1, "손님 방문") == false then
		return false
	end
	if items(ctx, buyer, req.mini_room { mode = MODE.buy, index = 0, bundles = 3 }, 1, "물약 3묶음 구매") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, buyer, POTION, 15) == false or meso_is(ctx, buyer, 97000) == false then
		return false
	end
	local failed = buyer:request({ resp.mini_room_buy_failed, resp.mini_room_items }, req.mini_room { mode = MODE.buy, index = 0, bundles = 100 }, nil, 3000)
	if check(ctx, failed ~= false and failed.result == BUY.not_enough_item, "남은 묶음보다 많이 구매가 거절되지 않음") == false then
		return false
	end

	local kicked = owner:request_on(buyer, resp.mini_room_leave, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000)
	if check(ctx, kicked ~= false and kicked.slot == 1 and kicked.reason == LEAVE.organizing, "주인 관리 시작에 손님이 정리 중으로 나가지 않음") == false then
		return false
	end
	if enter_failed(ctx, buyer, req.mini_room { mode = MODE.visit, sn = sn }, ENTER_ERROR.organizing, "관리 중 방문") == false then
		return false
	end

	if owner:request(resp.mini_room_meso_withdrawn, req.mini_room { mode = MODE.withdraw_meso }, nil, 3000) == false then
		return ctx:fail("판매 대금 회수 응답 없음")
	end
	ctx:sleep(300)
	if meso_is(ctx, owner, 3000) == false then
		return false
	end

	local destroyed = owner:request_on(buyer, resp.destroy_entrusted_shop, req.mini_room { mode = MODE.close }, nil, 5000)
	if destroyed == false then
		return ctx:fail("상점을 닫아도 맵에서 사라지지 않음")
	end
	ctx:sleep(300)
	if has(ctx, owner, POTION, 85) == false then
		return false
	end
	return shop_check(ctx, owner, CHECK.title, "닫은 뒤 개설 문의") ~= false
end

local function store_bank_flow(ctx)
	local owner = ctx:bot(0)
	if setup(ctx, owner, 1000, PERMIT .. ":1," .. POTION .. ":10") == false then
		return false
	end
	if move_to(ctx, owner, MARKET_ROOM, SPOT) == false then
		return false
	end
	if entered(ctx, owner, create_request(owner, "보관 상점"), 0, "개설") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, POTION, 2, 5, 100), 1, "물약 등록") == false then
		return false
	end
	if owner:map_move(FREE_MARKET) == false then
		return ctx:fail("개설 중 맵 이동 실패")
	end
	ctx:sleep(500)
	if has(ctx, owner, POTION, 0) == false then
		return false
	end

	if move_to(ctx, owner, MARKET_ROOM, nil) == false then
		return false
	end
	if shop_check(ctx, owner, CHECK.store_bank_full, "스토어뱅크가 찬 채 개설 문의") == false then
		return false
	end

	local p = claim_store_bank(ctx, owner)
	if p == false then
		return false
	end
	if check(ctx, p.npc_id == FREDRICK and p.storage ~= nil, "스토어뱅크 창이 열리지 않음") == false then
		return false
	end
	local consume = p.storage.tabs[CONSUME] or {}
	if check(ctx, #consume == 1 and consume[1].count == 10, "스토어뱅크 물약 목록이 다름") == false then
		return false
	end
	local fee = owner:request(resp.store_bank_fee, req.store_bank { mode = STORE_BANK.withdraw }, nil, 3000)
	if check(ctx, fee ~= false and fee.fee == 0, "보관 첫날 수수료가 0이 아님") == false then
		return false
	end
	local r = owner:request(resp.store_bank_result, req.store_bank { mode = STORE_BANK.confirm }, nil, 5000)
	if check(ctx, r ~= false and r.result == STORE_BANK.claimed, "스토어뱅크 찾기 실패") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, owner, POTION, 10) == false then
		return false
	end

	p = claim_store_bank(ctx, owner)
	if p == false then
		return false
	end
	return check(ctx, p.map_id == NOTHING_MAP, "찾은 뒤에도 스토어뱅크가 비지 않음")
end

local function overflow_flow(ctx)
	local owner = ctx:bot(0)
	local buyer = ctx:bot(1)
	if setup(ctx, owner, 0, PERMIT .. ":1," .. POTION .. ":10") == false then
		return false
	end
	if setup(ctx, buyer, 100000, "-") == false then
		return false
	end
	if clear_store_bank(ctx, owner) == false then
		return false
	end
	if setup(ctx, owner, 2147483000, PERMIT .. ":1," .. POTION .. ":10") == false then
		return false
	end

	if move_to(ctx, buyer, MARKET_ROOM, { x = -1300, y = 102 }) == false then
		return false
	end
	if move_to(ctx, owner, MARKET_ROOM, SPOT) == false then
		return false
	end
	if entered(ctx, owner, create_request(owner, "넘침 상점"), 0, "개설") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, POTION, 2, 5, 5000), 1, "물약 등록") == false then
		return false
	end
	local spawn = owner:request_on(buyer, resp.spawn_entrusted_shop, req.mini_room { mode = MODE.open }, function(s)
		return s.employer_id == owner:id()
	end, 5000)
	if spawn == false then
		return ctx:fail("상점 개설이 다른 캐릭터에게 보이지 않음")
	end
	local sn = spawn.entrusted_shop_balloon.sn
	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = sn }, 1, "손님 방문") == false then
		return false
	end
	if items(ctx, buyer, req.mini_room { mode = MODE.buy, index = 0, bundles = 1 }, 1, "물약 1묶음 구매") == false then
		return false
	end

	if owner:request_on(buyer, resp.mini_room_leave, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000) == false then
		return ctx:fail("주인 관리 시작에 손님이 나가지 않음")
	end
	local popup = owner:request(resp.notice, req.mini_room { mode = MODE.withdraw_meso }, function(p)
		return p.message:find("프레드릭", 1, true) ~= nil
	end, 5000)
	if popup == false then
		return ctx:fail("메소 한도 초과 시 프레드릭 보관 팝업이 오지 않음")
	end
	ctx:sleep(300)
	if meso_is(ctx, owner, 2147483000) == false then
		return false
	end

	if owner:request_on(buyer, resp.destroy_entrusted_shop, req.mini_room { mode = MODE.close }, nil, 5000) == false then
		return ctx:fail("상점을 닫아도 맵에서 사라지지 않음")
	end
	ctx:sleep(300)
	if has(ctx, owner, POTION, 5) == false then
		return false
	end

	local p = claim_store_bank(ctx, owner)
	if p == false then
		return false
	end
	if check(ctx, p.storage ~= nil and p.storage.meso == 5000, "프레드릭 보관 메소가 다름: " .. tostring(p.storage and p.storage.meso)) == false then
		return false
	end
	if setup(ctx, owner, 0, "-") == false then
		return false
	end
	return clear_store_bank(ctx, owner)
end

test_suite {
	name = "EntrustedShop: 고용상인 개설·판매·관리·스토어뱅크",
	bot_count = 2,

	scenarios = {
		shop_flow,
		store_bank_flow,
		overflow_flow,
	},
}
