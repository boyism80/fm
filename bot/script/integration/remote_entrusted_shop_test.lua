local pq = require("script/integration/lib/party_quest")

local FREE_MARKET = 910000000
local MARKET_ROOM = 910000008
local FREDRICK = 9030000

local PERMIT = 5030000
local REMOTE = 5470000
local POTION = 2000000
local SWORD = 1302000

local MODE = {
	create = 0x00,
	visit = 0x04,
	exit = 0x0A,
	open = 0x0B,
	add_item = 0x1D,
	buy = 0x1E,
	remove_item = 0x22,
	close = 0x25,
	withdraw_meso = 0x27,
}
local CHECK = { remote_location = 16, remote_visit = 17 }
local ENTER_ERROR = { organizing = 16 }
local LEAVE = { organizing = 13 }
local STORE_BANK = { withdraw = 0x19, confirm = 0x1A, claimed = 0x1D }
local NOTHING_MAP = 999999999
local NOTHING_CHANNEL = 255

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

local function remote(ctx, bot, result, what)
	local p = bot:request(resp.entrusted_shop_check_result, req.remote_entrusted_shop { slot = bot:slot(REMOTE) or 0 }, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 원격관리기 응답 없음: " .. what)
	end
	if p.result ~= result then
		return ctx:fail(string.format("%s %s: 결과 %d (기대 %d, 위치 %s/%s)", bot:name(), what, p.result, result, tostring(p.map_id), tostring(p.channel)))
	end
	return p
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

local function remote_flow(ctx)
	local owner = ctx:bot(0)
	local buyer = ctx:bot(1)
	if setup(ctx, owner, 0, PERMIT .. ":1," .. REMOTE .. ":1," .. POTION .. ":100," .. SWORD .. ":1") == false then
		return false
	end
	if setup(ctx, buyer, 100000, REMOTE .. ":1") == false then
		return false
	end
	if clear_store_bank(ctx, owner) == false then
		return false
	end

	local missing = remote(ctx, buyer, CHECK.remote_location, "상점 없이 원격관리")
	if missing == false then
		return false
	end
	if check(ctx, missing.map_id == NOTHING_MAP and missing.channel == NOTHING_CHANNEL, "상점 없음 위치가 다름") == false then
		return false
	end

	if move_to(ctx, buyer, MARKET_ROOM, GUEST_SPOT) == false or move_to(ctx, owner, MARKET_ROOM, SPOT) == false then
		return false
	end
	local create = req.mini_room { mode = MODE.create, type = 5, title = "원격 상점", slot = owner:slot(PERMIT) or 0, item_id = PERMIT }
	if entered(ctx, owner, create, 0, "개설") == false then
		return false
	end
	if items(ctx, owner, add_item(owner, POTION, 10, 5, 1000), 1, "물약 등록") == false then
		return false
	end
	local spawn = owner:request_on(buyer, resp.spawn_entrusted_shop, req.mini_room { mode = MODE.open }, function(s)
		return s.employer_id == owner:id()
	end, 5000)
	if spawn == false then
		return ctx:fail("상점 개설이 다른 캐릭터에게 보이지 않음")
	end
	local sn = spawn.entrusted_shop_balloon.sn

	if move_to(ctx, owner, FREE_MARKET, nil) == false then
		return false
	end
	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = sn }, 1, "손님 방문") == false then
		return false
	end
	if items(ctx, buyer, req.mini_room { mode = MODE.buy, index = 0, bundles = 3 }, 1, "물약 3묶음 구매") == false then
		return false
	end

	local p = remote(ctx, owner, CHECK.remote_visit, "같은 채널 원격관리")
	if p == false then
		return false
	end
	if check(ctx, p.sn == sn, string.format("원격관리 상점 번호 %d (기대 %d)", p.sn, sn)) == false then
		return false
	end
	local kicked = owner:request_on(buyer, resp.mini_room_leave, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000)
	if check(ctx, kicked ~= false and kicked.reason == LEAVE.organizing, "원격관리 시작에 손님이 정리 중으로 나가지 않음") == false then
		return false
	end
	local fail = buyer:request({ resp.mini_room_enter, resp.mini_room_enter_failed }, req.mini_room { mode = MODE.visit, sn = sn }, nil, 5000)
	if check(ctx, fail ~= false and fail.error == ENTER_ERROR.organizing, "원격관리 중 방문이 거절되지 않음") == false then
		return false
	end

	if items(ctx, owner, add_item(owner, SWORD, 1, 1, 50000), 2, "원격 검 등록") == false then
		return false
	end
	local removed = owner:request(resp.mini_room_item_removed, req.mini_room { mode = MODE.remove_item, index = 1 }, nil, 5000)
	if check(ctx, removed ~= false and removed.count == 1 and removed.index == 1, "원격 검 회수 응답이 다름") == false then
		return false
	end
	if owner:request(resp.mini_room_meso_withdrawn, req.mini_room { mode = MODE.withdraw_meso }, nil, 3000) == false then
		return ctx:fail("원격 판매 대금 회수 응답 없음")
	end
	ctx:sleep(300)
	if meso_is(ctx, owner, 3000) == false or has(ctx, owner, SWORD, 1) == false then
		return false
	end
	owner:send(req.mini_room { mode = MODE.exit })
	ctx:sleep(500)
	if entered(ctx, buyer, req.mini_room { mode = MODE.visit, sn = sn }, 1, "원격관리 종료 후 방문") == false then
		return false
	end
	buyer:send(req.mini_room { mode = MODE.exit })

	if owner:transfer(1) == false then
		return ctx:fail("1채널 이동 실패")
	end
	p = remote(ctx, owner, CHECK.remote_location, "다른 채널 원격관리")
	if p == false then
		return false
	end
	if check(ctx, p.map_id == MARKET_ROOM and p.channel == 0, string.format("다른 채널 상점 위치 %d/%d", p.map_id, p.channel)) == false then
		return false
	end
	if owner:transfer(0) == false then
		return ctx:fail("0채널 복귀 실패")
	end

	if remote(ctx, owner, CHECK.remote_visit, "닫기 전 원격관리") == false then
		return false
	end
	if entered(ctx, owner, req.mini_room { mode = MODE.visit, sn = sn }, 0, "닫기 전 원격 입장") == false then
		return false
	end
	local destroyed = owner:request_on(buyer, resp.destroy_entrusted_shop, req.mini_room { mode = MODE.close }, nil, 5000)
	if destroyed == false then
		return ctx:fail("원격으로 닫아도 맵에서 사라지지 않음")
	end
	ctx:sleep(300)
	if has(ctx, owner, POTION, 85) == false or has(ctx, owner, REMOTE, 1) == false then
		return false
	end
	return remote(ctx, owner, CHECK.remote_location, "닫은 뒤 원격관리") ~= false
end

test_suite {
	name = "RemoteEntrustedShop: 원격관리기 관리·채널 안내",
	bot_count = 2,

	scenarios = {
		remote_flow,
	},
}
