local pq = require("script/integration/lib/party_quest")

local HENESYS_PARK = 100000200
local KEEPER = 1012009
local FREE_MARKET = 910000000
local SCROOGE = 9030100

local MODE = { take_out = 4, store = 5, arrange = 6, meso = 7, close = 8 }
local RESULT = {
	take_out = 0x09,
	inventory_full = 0x0A,
	not_enough_meso = 0x0B,
	store = 0x0D,
	arrange = 0x0F,
	full = 0x11,
	meso = 0x13,
	open = 0x16,
}
local EQUIP = 1
local CONSUME = 2
local ETC = 4

local POTION = 2000000
local STAR = 2070000
local SWORD = 1302000
local SHELL = 4000000
local PET = 5000000

local function setup(ctx, bot, meso, items)
	local text = string.format("/봇초기화 30 0 %d %s -", meso, items)
	if pq.command(bot, text, "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	return true
end

local function open(ctx, bot, map, keeper)
	if bot:map_move(map) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map)
	end
	local oid = pq.npc(ctx, bot, keeper)
	if oid == false then
		return false
	end
	local p = bot:request(resp.storage, req.npc_click { oid = oid }, function(p)
		return p.result == RESULT.open
	end, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 창고가 열리지 않음: " .. keeper)
	end
	if p.npc_id ~= keeper then
		return ctx:fail(string.format("%s 창고 NPC가 다름: %d", bot:name(), p.npc_id))
	end
	return p
end

local function call(bot, request)
	return bot:request({ resp.storage, resp.storage_error }, request, nil, 3000)
end

local function expect(ctx, bot, request, result, what)
	local p = call(bot, request)
	if p == false then
		return ctx:fail(bot:name() .. " 응답 없음: " .. what)
	end
	if p.result ~= result then
		return ctx:fail(string.format("%s %s: 결과 0x%02X (기대 0x%02X)", bot:name(), what, p.result, result))
	end
	return p
end

local function ignored(ctx, bot, request, what)
	local p = bot:request({ resp.storage, resp.storage_error }, request, nil, 1500)
	if p ~= false then
		return ctx:fail(string.format("%s 거절되어야 할 요청에 창고 응답: %s (0x%02X)", bot:name(), what, p.result))
	end
	return true
end

local function store_request(bot, item_id, count)
	return req.storage { mode = MODE.store, slot = bot:slot(item_id) or 0, item_id = item_id, count = count }
end

local function take_out_request(inventory_type, index)
	return req.storage { mode = MODE.take_out, inventory_type = inventory_type, index = index }
end

local function meso_request(amount)
	return req.storage { mode = MODE.meso, meso = amount }
end

local function tab_count(p, inventory_type, item_id)
	local count = 0
	for _, item in ipairs(p.tabs[inventory_type] or {}) do
		if item.item_id == item_id then
			count = count + (item.count or 1)
		end
	end
	return count
end

local function tab_size(p, inventory_type)
	return #(p.tabs[inventory_type] or {})
end

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
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

local function normal_flow(ctx)
	local bot = ctx:bot(0)
	if setup(ctx, bot, 10000, POTION .. ":50," .. SWORD .. ":1," .. SHELL .. ":10") == false then
		return false
	end
	local p = open(ctx, bot, HENESYS_PARK, KEEPER)
	if p == false then
		return false
	end
	if check(ctx, p.slots == 4 and p.meso == 0, string.format("빈 창고가 아님: 칸 %d 메소 %s", p.slots, tostring(p.meso))) == false then
		return false
	end
	for inventory_type = 1, 5 do
		if check(ctx, tab_size(p, inventory_type) == 0, "빈 창고에 아이템이 있음: 탭 " .. inventory_type) == false then
			return false
		end
	end

	p = expect(ctx, bot, store_request(bot, POTION, 20), RESULT.store, "포션 보관")
	if p == false then
		return false
	end
	if check(ctx, tab_count(p, CONSUME, POTION) == 20, "보관한 포션 개수가 다름") == false then
		return false
	end
	if has(ctx, bot, POTION, 30) == false or meso_is(ctx, bot, 9900) == false then
		return false
	end

	p = expect(ctx, bot, store_request(bot, SWORD, 1), RESULT.store, "장비 보관")
	if p == false then
		return false
	end
	if check(ctx, tab_count(p, EQUIP, SWORD) == 1, "보관한 장비가 창고에 없음") == false then
		return false
	end
	if has(ctx, bot, SWORD, 0) == false or meso_is(ctx, bot, 9800) == false then
		return false
	end

	p = expect(ctx, bot, take_out_request(EQUIP, 0), RESULT.take_out, "장비 꺼내기")
	if p == false then
		return false
	end
	if check(ctx, tab_size(p, EQUIP) == 0, "꺼낸 장비가 창고에 남음") == false then
		return false
	end
	if has(ctx, bot, SWORD, 1) == false or meso_is(ctx, bot, 9800) == false then
		return false
	end

	p = expect(ctx, bot, meso_request(-5000), RESULT.meso, "메소 맡기기")
	if p == false then
		return false
	end
	if check(ctx, p.meso == 5000, "창고 메소가 5000이 아님: " .. tostring(p.meso)) == false or meso_is(ctx, bot, 4800) == false then
		return false
	end

	p = expect(ctx, bot, meso_request(2000), RESULT.meso, "메소 찾기")
	if p == false then
		return false
	end
	if check(ctx, p.meso == 3000, "창고 메소가 3000이 아님: " .. tostring(p.meso)) == false or meso_is(ctx, bot, 6800) == false then
		return false
	end

	if expect(ctx, bot, req.storage { mode = MODE.arrange }, RESULT.arrange, "정렬") == false then
		return false
	end
	return ignored(ctx, bot, req.storage { mode = MODE.close }, "닫기")
end

local function persistence(ctx)
	local bot = ctx:bot(1)
	if setup(ctx, bot, 10000, SHELL .. ":10") == false then
		return false
	end
	if open(ctx, bot, HENESYS_PARK, KEEPER) == false then
		return false
	end
	if expect(ctx, bot, store_request(bot, SHELL, 10), RESULT.store, "기타 보관") == false then
		return false
	end
	if expect(ctx, bot, meso_request(-1000), RESULT.meso, "메소 맡기기") == false then
		return false
	end

	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	local p = open(ctx, bot, HENESYS_PARK, KEEPER)
	if p == false then
		return false
	end
	if check(ctx, tab_count(p, ETC, SHELL) == 10, "채널 이동 후 보관 아이템이 사라짐") == false then
		return false
	end
	if check(ctx, p.meso == 1000, "채널 이동 후 창고 메소가 다름: " .. tostring(p.meso)) == false then
		return false
	end
	if has(ctx, bot, SHELL, 0) == false or meso_is(ctx, bot, 8900) == false then
		return false
	end
	if expect(ctx, bot, take_out_request(ETC, 0), RESULT.take_out, "채널 이동 후 꺼내기") == false then
		return false
	end
	return has(ctx, bot, SHELL, 10)
end

local function limits(ctx)
	local bot = ctx:bot(2)
	local items = {}
	for i = 0, 31 do
		items[#items + 1] = (SHELL + i) .. ":1"
	end
	if setup(ctx, bot, 10000, table.concat(items, ",")) == false then
		return false
	end
	if open(ctx, bot, HENESYS_PARK, KEEPER) == false then
		return false
	end
	for i = 0, 3 do
		if expect(ctx, bot, store_request(bot, SHELL + i, 1), RESULT.store, "보관 " .. (SHELL + i)) == false then
			return false
		end
	end
	if expect(ctx, bot, store_request(bot, SHELL + 4, 1), RESULT.full, "다섯 번째 보관") == false then
		return false
	end
	if has(ctx, bot, SHELL + 4, 1) == false or meso_is(ctx, bot, 9600) == false then
		return false
	end

	for i = 32, 35 do
		if pq.command(bot, "/아이템생성 " .. (SHELL + i) .. " 1", "아이템 생성: ") == false then
			return ctx:fail("기타 탭 채우기 실패: " .. (SHELL + i))
		end
	end
	if expect(ctx, bot, take_out_request(ETC, 0), RESULT.inventory_full, "가득 찬 인벤토리로 꺼내기") == false then
		return false
	end
	local p = open(ctx, bot, HENESYS_PARK, KEEPER)
	if p == false then
		return false
	end
	if check(ctx, tab_size(p, ETC) == 4, "꺼내기 실패 후 창고 아이템 수가 다름: " .. tab_size(p, ETC)) == false then
		return false
	end
	return has(ctx, bot, SHELL, 0)
end

local function not_enough_fee(ctx)
	local bot = ctx:bot(3)
	if setup(ctx, bot, 50, SHELL .. ":1") == false then
		return false
	end
	if open(ctx, bot, HENESYS_PARK, KEEPER) == false then
		return false
	end
	if expect(ctx, bot, store_request(bot, SHELL, 1), RESULT.not_enough_meso, "보관료 부족") == false then
		return false
	end
	if has(ctx, bot, SHELL, 1) == false or meso_is(ctx, bot, 50) == false then
		return false
	end

	if pq.command(bot, "/봇초기화 30 0 1000 - -", "봇초기화 완료") == false then
		return ctx:fail("메소 설정 실패")
	end
	local cases = {
		{ amount = -2000, what = "가진 메소보다 많이 맡기기" },
		{ amount = 1, what = "빈 창고에서 찾기" },
		{ amount = 0, what = "0 메소" },
		{ amount = -2147483648, what = "음수 최솟값" },
	}
	for _, case in ipairs(cases) do
		if expect(ctx, bot, meso_request(case.amount), RESULT.not_enough_meso, case.what) == false then
			return false
		end
	end
	if meso_is(ctx, bot, 1000) == false then
		return false
	end
	if expect(ctx, bot, meso_request(-500), RESULT.meso, "메소 맡기기") == false then
		return false
	end
	if expect(ctx, bot, meso_request(501), RESULT.not_enough_meso, "맡긴 것보다 많이 찾기") == false then
		return false
	end
	local p = open(ctx, bot, HENESYS_PARK, KEEPER)
	if p == false then
		return false
	end
	if check(ctx, p.meso == 500, "창고 메소가 500이 아님: " .. tostring(p.meso)) == false then
		return false
	end
	return meso_is(ctx, bot, 500)
end

local function forged_requests(ctx)
	local bot = ctx:bot(4)
	local items = string.format("%d:10,%d:500,%d:1,%d:1", POTION, STAR, SWORD, PET)
	if setup(ctx, bot, 10000, items) == false then
		return false
	end
	if has(ctx, bot, PET, 1) == false then
		return false
	end
	if ignored(ctx, bot, take_out_request(CONSUME, 0), "창고를 열기 전 꺼내기") == false then
		return false
	end
	if open(ctx, bot, HENESYS_PARK, KEEPER) == false then
		return false
	end

	local slot = bot:slot(POTION)
	local cases = {
		{ request = req.storage { mode = MODE.store, slot = slot, item_id = POTION + 1, count = 1 }, what = "칸과 다른 아이템 ID" },
		{ request = store_request(bot, POTION, 11), what = "가진 것보다 많은 개수" },
		{ request = store_request(bot, POTION, 0), what = "0개 보관" },
		{ request = req.storage { mode = MODE.store, slot = -11, item_id = SWORD, count = 1 }, what = "장착 칸" },
		{ request = req.storage { mode = MODE.store, slot = 99, item_id = POTION, count = 1 }, what = "범위 밖 칸" },
		{ request = store_request(bot, PET, 1), what = "펫 보관" },
		{ request = take_out_request(CONSUME, 5), what = "없는 순번 꺼내기" },
		{ request = take_out_request(9, 0), what = "없는 탭 꺼내기" },
	}
	for _, case in ipairs(cases) do
		if ignored(ctx, bot, case.request, case.what) == false then
			return false
		end
	end
	if has(ctx, bot, POTION, 10) == false or has(ctx, bot, PET, 1) == false or meso_is(ctx, bot, 10000) == false then
		return false
	end

	local p = expect(ctx, bot, store_request(bot, STAR, 1), RESULT.store, "표창 보관")
	if p == false then
		return false
	end
	if check(ctx, tab_count(p, CONSUME, STAR) == 500, "표창이 묶음째 보관되지 않음: " .. tab_count(p, CONSUME, STAR)) == false then
		return false
	end
	if has(ctx, bot, STAR, 0) == false then
		return false
	end

	if expect(ctx, bot, store_request(bot, POTION, 5), RESULT.store, "포션 보관") == false then
		return false
	end
	bot:send(take_out_request(CONSUME, 1))
	bot:send(take_out_request(CONSUME, 1))
	bot:request(resp.storage, nil, nil, 1500)
	bot:request(resp.storage, nil, nil, 1500)
	p = open(ctx, bot, HENESYS_PARK, KEEPER)
	if p == false then
		return false
	end
	if check(ctx, tab_size(p, CONSUME) == 1 and tab_count(p, CONSUME, STAR) == 500, "연속 꺼내기 후 창고 상태가 다름") == false then
		return false
	end
	if has(ctx, bot, POTION, 10) == false or has(ctx, bot, STAR, 0) == false then
		return false
	end

	if ignored(ctx, bot, req.storage { mode = MODE.close }, "닫기") == false then
		return false
	end
	if ignored(ctx, bot, store_request(bot, POTION, 1), "닫은 뒤 보관") == false then
		return false
	end
	return has(ctx, bot, POTION, 10)
end

local function free_market_fee(ctx)
	local bot = ctx:bot(5)
	if setup(ctx, bot, 2000, SHELL .. ":1") == false then
		return false
	end
	if open(ctx, bot, FREE_MARKET, SCROOGE) == false then
		return false
	end
	if expect(ctx, bot, store_request(bot, SHELL, 1), RESULT.store, "자유시장 보관") == false then
		return false
	end
	if meso_is(ctx, bot, 1500) == false then
		return false
	end
	if expect(ctx, bot, take_out_request(ETC, 0), RESULT.take_out, "자유시장 꺼내기") == false then
		return false
	end
	if has(ctx, bot, SHELL, 1) == false or meso_is(ctx, bot, 500) == false then
		return false
	end
	if expect(ctx, bot, store_request(bot, SHELL, 1), RESULT.store, "자유시장 다시 보관") == false then
		return false
	end
	if meso_is(ctx, bot, 0) == false then
		return false
	end
	local p = expect(ctx, bot, take_out_request(ETC, 0), RESULT.not_enough_meso, "수수료 없이 꺼내기")
	if p == false then
		return false
	end
	p = open(ctx, bot, FREE_MARKET, SCROOGE)
	if p == false then
		return false
	end
	if check(ctx, tab_count(p, ETC, SHELL) == 1, "수수료 부족 후 창고 아이템이 사라짐") == false then
		return false
	end
	return has(ctx, bot, SHELL, 0)
end

test_suite {
	name = "Storage: 창고 보관·꺼내기·메소·어뷰징",
	bot_count = 6,

	scenarios = {
		normal_flow,
		persistence,
		limits,
		not_enough_fee,
		forged_requests,
		free_market_fee,
	},
}
