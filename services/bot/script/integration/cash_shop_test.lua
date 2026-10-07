local pq = require("script/integration/lib/party_quest")

local NX = 0
local MAPLE_POINT = 1
local CASH = 5

local ACTION = { buy = 3, wishlist = 5, take_out = 12, put_in = 13 }
local KIND = { wishlist_updated = 50, buy_failed = 51, taken_out = 67, take_out_failed = 68, put_in = 69 }
local NOT_ENOUGH_CASH = 122

local MEGAPHONE = { sn = 10000804, item_id = 5060002, price = 990 }
local SUPER_MEGAPHONE = { sn = 10001513, item_id = 5076000, price = 1150 }

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function cash_is(ctx, bot, nx, mp)
	local have_nx, have_mp = bot:cash()
	return check(ctx, have_nx == nx and have_mp == mp, string.format("%s 캐시 NX %d / MP %d (기대 %d / %d)", bot:name(), have_nx, have_mp, nx, mp))
end

local function find_locker(bot, item_id)
	for _, item in ipairs(bot:locker()) do
		if item.item_id == item_id then
			return item
		end
	end
	return nil
end

local function buy(ctx, bot, commodity, currency)
	local p = bot:request(resp.cash_shop_balance, req.cash_shop_operation { action = ACTION.buy, currency = currency, commodity_sn = commodity.sn }, nil, 5000)
	if p == false then
		return ctx:fail(string.format("%s 구매 응답 없음: %d", bot:name(), commodity.sn))
	end
	return check(ctx, find_locker(bot, commodity.item_id) ~= nil, string.format("%s 보관함에 %d 없음", bot:name(), commodity.item_id))
end

local function result(ctx, bot, request, kind, what)
	local p = bot:request(resp.cash_shop_result, request, function(p)
		return p.kind == kind or p.kind == kind + 1
	end, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 응답 없음: " .. what)
	end
	if p.kind ~= kind then
		return ctx:fail(string.format("%s %s: 결과 %d (기대 %d), 실패 코드 %d", bot:name(), what, p.kind, kind, p.failure))
	end
	return p
end

local function shop_flow(ctx)
	local bot = ctx:bot(0)
	local map = bot:map()
	if pq.command(bot, "/캐시얻기 5000 2000", "캐시 잔액: NX 5000, 메이플포인트 2000") == false then
		return ctx:fail("캐시 지급 실패")
	end

	if bot:enter_cash_shop() == false then
		return ctx:fail("캐시샵 입장 실패")
	end
	if cash_is(ctx, bot, 5000, 2000) == false then
		return false
	end
	if check(ctx, #bot:locker() == 0, "새 계정 보관함이 비어 있지 않음") == false then
		return false
	end

	if buy(ctx, bot, MEGAPHONE, NX) == false or cash_is(ctx, bot, 5000 - MEGAPHONE.price, 2000) == false then
		return false
	end
	if buy(ctx, bot, SUPER_MEGAPHONE, MAPLE_POINT) == false or cash_is(ctx, bot, 4010, 2000 - SUPER_MEGAPHONE.price) == false then
		return false
	end
	local p = bot:request(resp.cash_shop_result, req.cash_shop_operation { action = ACTION.buy, currency = MAPLE_POINT, commodity_sn = SUPER_MEGAPHONE.sn }, function(p)
		return p.kind == KIND.buy_failed
	end, 5000)
	if p == false or p.failure ~= NOT_ENOUGH_CASH then
		return ctx:fail("메이플포인트 부족 구매가 거절되지 않음")
	end
	if cash_is(ctx, bot, 4010, 850) == false then
		return false
	end

	local megaphone = find_locker(bot, MEGAPHONE.item_id)
	if result(ctx, bot, req.cash_shop_operation { action = ACTION.take_out, serial = megaphone.serial, inventory_type = CASH }, KIND.taken_out, "보관함 꺼내기") == false then
		return false
	end
	if check(ctx, bot:items()[MEGAPHONE.item_id] == 1 and find_locker(bot, MEGAPHONE.item_id) == nil, "꺼낸 아이템이 인벤토리에 없음") == false then
		return false
	end
	if result(ctx, bot, req.cash_shop_operation { action = ACTION.take_out, serial = megaphone.serial, inventory_type = CASH }, KIND.taken_out + 1, "꺼낸 아이템 다시 꺼내기") == false then
		return false
	end

	if result(ctx, bot, req.cash_shop_operation { action = ACTION.put_in, serial = megaphone.serial, inventory_type = CASH }, KIND.put_in, "보관함 넣기") == false then
		return false
	end
	if check(ctx, bot:items()[MEGAPHONE.item_id] == nil and find_locker(bot, MEGAPHONE.item_id) ~= nil, "넣은 아이템이 보관함에 없음") == false then
		return false
	end
	if result(ctx, bot, req.cash_shop_operation { action = ACTION.put_in, serial = megaphone.serial, inventory_type = CASH }, KIND.put_in + 1, "넣은 아이템 다시 넣기") == false then
		return false
	end

	p = result(ctx, bot, req.cash_shop_operation { action = ACTION.wishlist, wishlist = { MEGAPHONE.sn, SUPER_MEGAPHONE.sn } }, KIND.wishlist_updated, "위시리스트")
	if p == false then
		return false
	end
	if check(ctx, p.wishlist[1] == MEGAPHONE.sn and p.wishlist[2] == SUPER_MEGAPHONE.sn and p.wishlist[3] == 0, "위시리스트 내용이 다름") == false then
		return false
	end

	local super = find_locker(bot, SUPER_MEGAPHONE.item_id)
	if result(ctx, bot, req.cash_shop_operation { action = ACTION.take_out, serial = super.serial, inventory_type = CASH }, KIND.taken_out, "확성기 꺼내기") == false then
		return false
	end

	if bot:leave_cash_shop() == false then
		return ctx:fail("캐시샵 퇴장 실패")
	end
	if check(ctx, bot:map() == map, string.format("퇴장 후 맵이 바뀜 %d -> %d", map, bot:map())) == false then
		return false
	end
	if check(ctx, bot:items()[SUPER_MEGAPHONE.item_id] == 1, "꺼낸 아이템이 게임 인벤토리에 저장되지 않음") == false then
		return false
	end

	if bot:enter_cash_shop() == false then
		return ctx:fail("캐시샵 재입장 실패")
	end
	if cash_is(ctx, bot, 4010, 850) == false then
		return false
	end
	if check(ctx, #bot:locker() == 1 and find_locker(bot, MEGAPHONE.item_id) ~= nil, "재입장 후 보관함이 다름: " .. #bot:locker()) == false then
		return false
	end
	if bot:leave_cash_shop() == false then
		return ctx:fail("캐시샵 재퇴장 실패")
	end
	return true
end

test_suite {
	name = "Cash shop: 입장·구매·보관함·위시리스트·퇴장",
	bot_count = 1,

	scenarios = {
		shop_flow,
	},
}
