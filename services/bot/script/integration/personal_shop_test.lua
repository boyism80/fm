local pq = require("script/integration/lib/party_quest")

local FREE_MARKET = 910000000
local SHOP_ROOM = 910000001

local PERMIT = 5140000
local POTION = 2000000
local SWORD = 1302000

local MODE = {
	create = 0x00,
	visit = 0x04,
	exit = 0x0A,
	open = 0x0B,
	add_item = 0x12,
	buy = 0x13,
	remove_item = 0x17,
	kick = 0x18,
	blacklist = 0x1A,
}
local ENTER_ERROR = { free_market = 13, blocked = 15 }
local LEAVE = { kicked = 5, sold_out = 10 }
local BUY = { not_enough_item = 1 }

local SPOT = { x = -1000, y = 102 }
local GUEST_SPOT = { x = -1300, y = 102 }

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

local function create_request(bot, title)
	return req.mini_room {
		mode = MODE.create,
		type = 4,
		title = title,
		slot = bot:slot(PERMIT) or 0,
		item_id = PERMIT,
	}
end

local function enter_failed(ctx, bot, request, code, what)
	local p = bot:request({ resp.personal_shop_enter, resp.mini_room_enter_failed }, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 응답 없음: " .. what)
	end
	if p.error ~= code then
		return ctx:fail(string.format("%s %s: 입장 에러 %s (기대 %d)", bot:name(), what, tostring(p.error), code))
	end
	return true
end

local function entered(ctx, bot, request, slot, what)
	local p = bot:request({ resp.personal_shop_enter, resp.mini_room_enter_failed }, request, nil, 5000)
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
	local p = bot:request(resp.personal_shop_items, request, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 상점 목록 갱신 없음: " .. what)
	end
	if #p.items ~= count then
		return ctx:fail(string.format("%s %s: 등록 수 %d (기대 %d)", bot:name(), what, #p.items, count))
	end
	return p
end

local function open_shop(ctx, owner, watcher)
	local balloon = owner:request_on(watcher, resp.user_mini_room_balloon, req.mini_room { mode = MODE.open }, function(p)
		return p.character_id == owner:id() and p.balloon ~= nil
	end, 5000)
	if balloon == false then
		return ctx:fail("개인상점 풍선이 다른 캐릭터에게 보이지 않음")
	end
	return balloon.balloon
end

local function sell_flow(ctx)
	local owner = ctx:bot(0)
	local buyer = ctx:bot(1)
	local other = ctx:bot(2)
	if setup(ctx, owner, 0, PERMIT .. ":1," .. POTION .. ":100," .. SWORD .. ":1") == false then
		return false
	end
	if setup(ctx, buyer, 100000, "-") == false or setup(ctx, other, 100000, "-") == false then
		return false
	end

	if move_to(ctx, owner, FREE_MARKET, nil) == false then
		return false
	end
	if enter_failed(ctx, owner, create_request(owner, "입구"), ENTER_ERROR.free_market, "자유시장 입구 개설") == false then
		return false
	end

	if move_to(ctx, buyer, SHOP_ROOM, GUEST_SPOT) == false or move_to(ctx, other, SHOP_ROOM, GUEST_SPOT) == false then
		return false
	end
	if move_to(ctx, owner, SHOP_ROOM, SPOT) == false then
		return false
	end
	local p = entered(ctx, owner, create_request(owner, "봇 개인상점"), 0, "개설")
	if p == false then
		return false
	end
	if check(ctx, p.title == "봇 개인상점" and p.max_items == 16 and #p.items == 0, "개설 입장 정보가 다름") == false then
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
	local removed = owner:request(resp.mini_room_item_removed, req.mini_room { mode = MODE.remove_item, index = 1 }, nil, 5000)
	if check(ctx, removed ~= false and removed.count == 1 and removed.index == 1, "검 회수 응답이 다름") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, owner, SWORD, 1) == false then
		return false
	end

	local balloon = open_shop(ctx, owner, buyer)
	if balloon == false then
		return false
	end
	if check(ctx, balloon.title == "봇 개인상점" and balloon.users == 1, "개인상점 풍선 정보가 다름") == false then
		return false
	end
	local late = owner:request(resp.personal_shop_items, add_item(owner, SWORD, 1, 1, 50000), nil, 1500)
	if check(ctx, late == false, "개설 후 등록이 거절되지 않음") == false then
		return false
	end

	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = balloon.sn }, 1, "손님 방문") == false then
		return false
	end
	local sold = buyer:request_on(owner, resp.mini_room_sold, req.mini_room { mode = MODE.buy, index = 0, bundles = 3 }, nil, 5000)
	if check(ctx, sold ~= false and sold.index == 0 and sold.bundles == 3 and sold.buyer == buyer:name(), "주인에게 판매 알림이 없음") == false then
		return false
	end
	ctx:sleep(300)
	if has(ctx, buyer, POTION, 15) == false or meso_is(ctx, buyer, 97000) == false or meso_is(ctx, owner, 3000) == false then
		return false
	end
	local failed = buyer:request({ resp.mini_room_buy_failed, resp.personal_shop_items }, req.mini_room { mode = MODE.buy, index = 0, bundles = 100 }, nil, 3000)
	if check(ctx, failed ~= false and failed.result == BUY.not_enough_item, "남은 묶음보다 많이 구매가 거절되지 않음") == false then
		return false
	end

	if entered(ctx, other, req.mini_room { mode = MODE.visit, sn = balloon.sn }, 2, "두 번째 손님 방문") == false then
		return false
	end
	local closed = other:request_on(buyer, resp.mini_room_leave, req.mini_room { mode = MODE.buy, index = 0, bundles = 7 }, function(leave)
		return leave.reason == LEAVE.sold_out
	end, 5000)
	if closed == false then
		return ctx:fail("모두 팔린 뒤 손님이 품절로 나가지 않음")
	end
	ctx:sleep(300)
	if has(ctx, other, POTION, 35) == false or meso_is(ctx, owner, 10000) == false then
		return false
	end
	return true
end

local function kick_flow(ctx)
	local owner = ctx:bot(0)
	local buyer = ctx:bot(1)
	local other = ctx:bot(2)
	if move_to(ctx, owner, SHOP_ROOM, SPOT) == false then
		return false
	end
	if entered(ctx, owner, create_request(owner, "추방 상점"), 0, "다시 개설") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, SWORD, 1, 1, 50000), 1, "검 등록") == false then
		return false
	end
	local balloon = open_shop(ctx, owner, buyer)
	if balloon == false then
		return false
	end
	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = balloon.sn }, 1, "손님 방문") == false then
		return false
	end

	local kicked = owner:request_on(buyer, resp.mini_room_leave, req.mini_room { mode = MODE.kick, slot = 1, name = buyer:name() }, function(leave)
		return leave.reason == LEAVE.kicked
	end, 5000)
	if kicked == false then
		return ctx:fail("추방한 손님이 강제 퇴장되지 않음")
	end
	if enter_failed(ctx, buyer, req.mini_room { mode = MODE.visit, sn = balloon.sn }, ENTER_ERROR.blocked, "추방 후 재방문") == false then
		return false
	end
	owner:send(req.mini_room { mode = MODE.blacklist, names = { other:name() } })
	ctx:sleep(300)
	if enter_failed(ctx, other, req.mini_room { mode = MODE.visit, sn = balloon.sn }, ENTER_ERROR.blocked, "블랙리스트 방문") == false then
		return false
	end

	local removed = owner:request_on(buyer, resp.user_mini_room_balloon, req.mini_room { mode = MODE.exit }, function(p)
		return p.character_id == owner:id() and p.balloon == nil
	end, 5000)
	if removed == false then
		return ctx:fail("상점을 닫아도 풍선이 사라지지 않음")
	end
	ctx:sleep(300)
	if has(ctx, owner, SWORD, 1) == false then
		return false
	end

	if entered(ctx, owner, create_request(owner, "이동 상점"), 0, "이동 전 개설") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, SWORD, 1, 1, 50000), 1, "검 등록") == false then
		return false
	end
	if owner:map_move(FREE_MARKET) == false then
		return ctx:fail("개설 중 맵 이동 실패")
	end
	ctx:sleep(500)
	return has(ctx, owner, SWORD, 1)
end

test_suite {
	name = "PersonalShop: 개인상점 개설·판매·추방·종료",
	bot_count = 3,

	scenarios = {
		sell_flow,
		kick_flow,
	},
}
