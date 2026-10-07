local pq = require("script/integration/lib/party_quest")

local NX = 0
local MAPLE_POINT = 1
local CASH = 5

local EQUIP = 1

local ACTION = {
	buy = 3,
	gift = 4,
	wishlist = 5,
	inventory_slots = 6,
	storage_slots = 7,
	character_slots = 8,
	take_out = 12,
	put_in = 13,
	pay_back = 25,
	couple_ring = 28,
	package = 29,
	quest_item = 31,
	friendship_ring = 34,
}
local KIND = {
	wishlist_updated = 50,
	buy_failed = 53,
	coupon_redeemed = 54,
	coupon_failed = 57,
	gift_sent = 59,
	inventory_slots = 61,
	storage_slots = 63,
	character_slots = 65,
	taken_out = 67,
	take_out_failed = 68,
	put_in = 69,
	paid_back = 96,
	package_bought = 100,
	quest_item_bought = 104,
}
local NOT_ENOUGH_CASH = 122
local SAME_ACCOUNT = 125
local WRONG_NAME = 126
local COUPON_WRONG = 131
local COUPON_USED = 133
local UNKNOWN = 0
local RING = 143

local MEGAPHONE = { sn = 10000804, item_id = 5060002, price = 990 }
local SUPER_MEGAPHONE = { sn = 10001513, item_id = 5076000, price = 1150 }
local PACKAGE = { sn = 70000003, price = 3540, size = 3 }
local HAT = { sn = 20000000, item_id = 1002186, price = 1500, refund = 450 }
local QUEST_ITEM = { sn = 80000054, item_id = 4031191, price = 1 }
local CHARACTER_SLOT = { sn = 50200034, price = 6900 }
local SLOT_PRICE = 3800
local COUPLE_RING = { sn = 20900006, item_id = 1112001, price = 3500 }
local FRIENDSHIP_RING = { sn = 20900056, item_id = 1112800, price = 3500 }
local CASH_RING_SLOT = -112

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

local function paid(ctx, bot, request, kind, what)
	local p = result(ctx, bot, request, kind, what)
	if p == false then
		return false
	end
	ctx:sleep(300)
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
	local profile = bot:request(resp.character_profile, req.inspect_character { character_id = bot:id() }, nil, 3000)
	if profile == false then
		return ctx:fail("캐릭터 정보 응답 없음")
	end
	if check(ctx, profile.self and profile.guild_name == "-" and #profile.wishlist == 2 and profile.wishlist[1] == MEGAPHONE.sn and profile.wishlist[2] == SUPER_MEGAPHONE.sn, "캐릭터 정보 위시리스트가 다름: " .. #profile.wishlist) == false then
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

local function coupon(ctx, bot, code, kind, what)
	local p = bot:request(resp.cash_shop_result, req.cash_shop_coupon { code = code }, function(p)
		return p.kind == KIND.coupon_redeemed or p.kind == KIND.coupon_failed
	end, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 쿠폰 응답 없음: " .. what)
	end
	if p.kind ~= kind then
		return ctx:fail(string.format("%s 쿠폰 %s: 결과 %d (기대 %d), 실패 코드 %d", bot:name(), what, p.kind, kind, p.failure))
	end
	return p
end

local function create_coupons(ctx, bot)
	local codes = {}
	for _, spec in ipairs({ { "캐시", 100 }, { "포인트", 50 }, { "메소", 1000 }, { "아이템", MEGAPHONE.sn } }) do
		local p = pq.command(bot, string.format("/쿠폰생성 %s %d", spec[1], spec[2]), "쿠폰: ")
		if p == false then
			return ctx:fail("쿠폰 생성 실패: " .. spec[1])
		end
		codes[spec[1]] = p.message:match("쿠폰: (%w+)")
	end
	return codes
end

local function gift_flow(ctx, buyer, receiver)
	local notice = buyer:request_on(receiver, resp.notice, req.cash_shop_operation {
		action = ACTION.gift,
		commodity_sn = MEGAPHONE.sn,
		recipient = receiver:name(),
		message = "선물입니다",
	}, function(p)
		return p.message:find("캐시샵 선물", 1, true) ~= nil
	end, 5000)
	if notice == false then
		return ctx:fail(receiver:name() .. " 선물 알림을 받지 못함")
	end
	ctx:sleep(500)
	if cash_is(ctx, buyer, 30000 - MEGAPHONE.price, 5000) == false then
		return false
	end

	local p = paid(ctx, buyer, req.cash_shop_operation { action = ACTION.gift, commodity_sn = MEGAPHONE.sn, recipient = buyer:name(), message = "" }, KIND.gift_sent + 1, "같은 계정 선물")
	if p == false or check(ctx, p.failure == SAME_ACCOUNT, "같은 계정 선물 실패 코드 " .. p.failure) == false then
		return false
	end
	p = paid(ctx, buyer, req.cash_shop_operation { action = ACTION.gift, commodity_sn = MEGAPHONE.sn, recipient = "nobody_here", message = "" }, KIND.gift_sent + 1, "없는 캐릭터 선물")
	if p == false or check(ctx, p.failure == WRONG_NAME, "없는 캐릭터 선물 실패 코드 " .. p.failure) == false then
		return false
	end
	return true
end

local function purchase_flow(ctx, bot)
	local nx = 30000 - MEGAPHONE.price
	local before = #bot:locker()
	local p = paid(ctx, bot, req.cash_shop_operation { action = ACTION.package, currency = NX, commodity_sn = PACKAGE.sn }, KIND.package_bought, "패키지 구매")
	if p == false then
		return false
	end
	nx = nx - PACKAGE.price
	if check(ctx, #p.items == PACKAGE.size and #bot:locker() == before + PACKAGE.size, "패키지 구성품 수가 다름: " .. #p.items) == false then
		return false
	end

	if paid(ctx, bot, req.cash_shop_operation { action = ACTION.buy, currency = NX, commodity_sn = HAT.sn }, KIND.buy_failed - 1, "모자 구매") == false then
		return false
	end
	nx = nx - HAT.price
	p = paid(ctx, bot, req.cash_shop_operation { action = ACTION.pay_back, serial = find_locker(bot, HAT.item_id).serial }, KIND.paid_back, "환불")
	if p == false then
		return false
	end
	if check(ctx, p.maple_point == HAT.refund and find_locker(bot, HAT.item_id) == nil, "환불 포인트 " .. p.maple_point) == false then
		return false
	end
	if cash_is(ctx, bot, nx, 5000 + HAT.refund) == false then
		return false
	end

	p = result(ctx, bot, req.cash_shop_operation { action = ACTION.quest_item, commodity_sn = QUEST_ITEM.sn }, KIND.quest_item_bought, "메소 퀘스트 아이템")
	if p == false then
		return false
	end
	if check(ctx, bot:items()[QUEST_ITEM.item_id] == 1, "퀘스트 아이템이 인벤토리에 없음") == false then
		return false
	end

	p = result(ctx, bot, req.cash_shop_operation { action = ACTION.couple_ring, commodity_sn = MEGAPHONE.sn, recipient = "nobody_here", message = "" }, KIND.buy_failed, "반지가 아닌 커플링")
	if p == false or check(ctx, p.failure == UNKNOWN, "반지가 아닌 커플링 실패 코드 " .. p.failure) == false then
		return false
	end
	return nx
end

local function slot_flow(ctx, bot, nx, equip_slots)
	local p = paid(ctx, bot, req.cash_shop_operation { action = ACTION.inventory_slots, currency = NX, inventory_type = EQUIP }, KIND.inventory_slots, "인벤토리 확장")
	if p == false or check(ctx, p.slots == equip_slots, string.format("장비 슬롯 %d (기대 %d)", p.slots, equip_slots)) == false then
		return false
	end
	return nx - SLOT_PRICE
end

local function service_flow(ctx)
	local buyer = ctx:bot(1)
	local receiver = ctx:bot(0)
	if pq.command(buyer, "/캐시얻기 30000 5000", "캐시 잔액: NX 30000, 메이플포인트 5000") == false then
		return ctx:fail("캐시 지급 실패")
	end
	if pq.command(buyer, "/메소얻기 10", "메소 10 획득.") == false then
		return ctx:fail("메소 지급 실패")
	end
	local codes = create_coupons(ctx, buyer)
	if codes == false then
		return false
	end

	if buyer:enter_cash_shop() == false then
		return ctx:fail("캐시샵 입장 실패")
	end
	if check(ctx, buyer:slot_limit(EQUIP) == 32, "기본 장비 슬롯 " .. buyer:slot_limit(EQUIP)) == false then
		return false
	end
	if gift_flow(ctx, buyer, receiver) == false then
		return false
	end
	local nx = purchase_flow(ctx, buyer)
	if nx == false then
		return false
	end

	nx = slot_flow(ctx, buyer, nx, 36)
	if nx == false then
		return false
	end
	local p = paid(ctx, buyer, req.cash_shop_operation { action = ACTION.storage_slots, currency = NX }, KIND.storage_slots, "창고 확장")
	if p == false or check(ctx, p.slots == 8, "창고 슬롯 " .. p.slots) == false then
		return false
	end
	nx = nx - SLOT_PRICE
	if paid(ctx, buyer, req.cash_shop_operation { action = ACTION.character_slots, currency = NX, commodity_sn = CHARACTER_SLOT.sn }, KIND.character_slots, "캐릭터 슬롯 확장") == false then
		return false
	end
	nx = nx - CHARACTER_SLOT.price
	if cash_is(ctx, buyer, nx, 5000 + HAT.refund) == false then
		return false
	end

	if coupon(ctx, buyer, codes["캐시"], KIND.coupon_redeemed, "NX") == false then
		return false
	end
	p = coupon(ctx, buyer, codes["포인트"], KIND.coupon_redeemed, "메이플포인트")
	if p == false or check(ctx, p.maple_point == 50, "포인트 쿠폰 " .. p.maple_point) == false then
		return false
	end
	p = coupon(ctx, buyer, codes["메소"], KIND.coupon_redeemed, "메소")
	if p == false or check(ctx, p.meso == 1000, "메소 쿠폰 " .. p.meso) == false then
		return false
	end
	p = coupon(ctx, buyer, codes["아이템"], KIND.coupon_redeemed, "아이템")
	if p == false or check(ctx, #p.items == 1 and p.items[1].item_id == MEGAPHONE.item_id, "아이템 쿠폰 내용이 다름") == false then
		return false
	end
	p = coupon(ctx, buyer, codes["캐시"], KIND.coupon_failed, "재사용")
	if p == false or check(ctx, p.failure == COUPON_USED, "재사용 쿠폰 실패 코드 " .. p.failure) == false then
		return false
	end
	p = coupon(ctx, buyer, "NOSUCHCOUPON0000", KIND.coupon_failed, "없는 코드")
	if p == false or check(ctx, p.failure == COUPON_WRONG, "없는 쿠폰 실패 코드 " .. p.failure) == false then
		return false
	end
	ctx:sleep(500)
	if cash_is(ctx, buyer, nx + 100, 5000 + HAT.refund + 50) == false then
		return false
	end

	if buyer:leave_cash_shop() == false then
		return ctx:fail("캐시샵 퇴장 실패")
	end
	if check(ctx, buyer:slot_limit(EQUIP) == 36, "퇴장 후 장비 슬롯 " .. buyer:slot_limit(EQUIP)) == false then
		return false
	end
	if check(ctx, buyer:items()[QUEST_ITEM.item_id] == 1, "퀘스트 아이템이 저장되지 않음") == false then
		return false
	end
	if check(ctx, buyer:meso() == 10 - QUEST_ITEM.price + 1000, "퇴장 후 메소 " .. buyer:meso()) == false then
		return false
	end

	if receiver:enter_cash_shop() == false then
		return ctx:fail("받는 캐릭터 캐시샵 입장 실패")
	end
	local gifts = receiver:gifts()
	if check(ctx, #gifts == 1 and gifts[1].item_id == MEGAPHONE.item_id and gifts[1].sender_name == buyer:name(), "선물 목록이 다름: " .. #gifts) == false then
		return false
	end
	if receiver:leave_cash_shop() == false then
		return ctx:fail("받는 캐릭터 캐시샵 퇴장 실패")
	end

	if buyer:enter_cash_shop() == false then
		return ctx:fail("캐시샵 재입장 실패")
	end
	if slot_flow(ctx, buyer, nx, 40) == false then
		return false
	end
	if buyer:leave_cash_shop() == false then
		return ctx:fail("캐시샵 재퇴장 실패")
	end
	return true
end

local function ring_failed(ctx, bot, action, commodity, recipient, message, failure, what)
	local p = result(ctx, bot, req.cash_shop_operation { action = action, commodity_sn = commodity.sn, recipient = recipient, message = message }, KIND.buy_failed, what)
	if p == false then
		return false
	end
	return check(ctx, p.failure == failure, string.format("%s 실패 코드 %d (기대 %d)", what, p.failure, failure))
end

local function ring_flow(ctx)
	local buyer, partner = ctx:bot(1), ctx:bot(0)
	if pq.command(buyer, "/성별 0", "성별 변경: 0") == false or pq.command(partner, "/성별 0", "성별 변경: 0") == false then
		return ctx:fail("성별 변경 실패")
	end
	if partner:enter_cash_shop() == false or partner:leave_cash_shop() == false then
		return ctx:fail("받는 캐릭터 성별 저장 실패")
	end
	if pq.command(buyer, "/캐시얻기 10000 0", "캐시 잔액: ") == false then
		return ctx:fail("캐시 지급 실패")
	end

	if buyer:enter_cash_shop() == false then
		return ctx:fail("캐시샵 입장 실패")
	end
	if ring_failed(ctx, buyer, ACTION.couple_ring, COUPLE_RING, partner:name(), "", UNKNOWN, "메시지 없는 커플링") == false then
		return false
	end
	if ring_failed(ctx, buyer, ACTION.friendship_ring, FRIENDSHIP_RING, buyer:name(), "친구", WRONG_NAME, "자기 자신에게 우정링") == false then
		return false
	end
	if ring_failed(ctx, buyer, ACTION.couple_ring, COUPLE_RING, partner:name(), "사랑해", RING, "동성 커플링") == false then
		return false
	end
	if buyer:leave_cash_shop() == false then
		return ctx:fail("캐시샵 퇴장 실패")
	end

	if pq.command(partner, "/성별 1", "성별 변경: 1") == false then
		return ctx:fail("성별 변경 실패")
	end
	if partner:enter_cash_shop() == false or partner:leave_cash_shop() == false then
		return ctx:fail("받는 캐릭터 성별 저장 실패")
	end
	if buyer:enter_cash_shop() == false then
		return ctx:fail("캐시샵 재입장 실패")
	end
	local nx, mp = buyer:cash()
	local sent = buyer:request(resp.cash_shop_result, req.cash_shop_operation {
		action = ACTION.couple_ring,
		commodity_sn = COUPLE_RING.sn,
		recipient = partner:name(),
		message = "사랑해",
	}, function(p)
		return p.kind == KIND.gift_sent or p.kind == KIND.buy_failed
	end, 5000)
	if sent == false or sent.kind ~= KIND.gift_sent then
		return ctx:fail("커플링 구매 실패: " .. (sent and sent.failure or -1))
	end
	ctx:sleep(300)
	if cash_is(ctx, buyer, nx - COUPLE_RING.price, mp) == false then
		return false
	end
	local ring = find_locker(buyer, COUPLE_RING.item_id)
	if check(ctx, ring ~= nil, "산 커플링이 보관함에 없음") == false then
		return false
	end
	if result(ctx, buyer, req.cash_shop_operation { action = ACTION.take_out, serial = ring.serial, inventory_type = EQUIP }, KIND.taken_out, "커플링 꺼내기") == false then
		return false
	end
	if buyer:leave_cash_shop() == false then
		return ctx:fail("캐시샵 퇴장 실패")
	end

	if partner:map() ~= buyer:map() and partner:map_move(buyer:map()) == false then
		return ctx:fail("받는 캐릭터 맵 이동 실패")
	end
	local look = buyer:request_on(partner, resp.update_character_look, req.move_item {
		inventory_type = EQUIP,
		source = buyer:slot(COUPLE_RING.item_id),
		dest = CASH_RING_SLOT,
		count = 1,
	}, function(p)
		return p.character.id == buyer:id()
	end, 5000)
	if look == false then
		return ctx:fail("커플링 착용 외형 갱신이 상대에게 오지 않음")
	end
	if check(ctx, look.crush_ring ~= nil and look.crush_ring.item_id == COUPLE_RING.item_id, "착용한 커플링이 외형에 없음") == false then
		return false
	end

	if partner:enter_cash_shop() == false then
		return ctx:fail("받는 캐릭터 캐시샵 입장 실패")
	end
	local gifts = partner:gifts()
	if check(ctx, #gifts == 1 and gifts[1].item_id == COUPLE_RING.item_id and gifts[1].sender_name == buyer:name(), "커플링 선물이 다름: " .. #gifts) == false then
		return false
	end
	if partner:leave_cash_shop() == false then
		return ctx:fail("받는 캐릭터 캐시샵 퇴장 실패")
	end
	return true
end

test_suite {
	name = "Cash shop: 구매·보관함·선물·패키지·환불·슬롯·쿠폰·반지",
	bot_count = 2,

	scenarios = {
		shop_flow,
		service_flow,
		ring_flow,
	},
}
