local pq = require("script/integration/lib/party_quest")

local HENESYS_PARK = 100000200
local SNAIL_CARD = 2380000
local SNAIL = 100100
local MANO_CARD = 2388000
local NOT_OWNED_CARD = 2381000
local REGISTER_CARD_EFFECT = 13

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot)
	if pq.command(bot, "/봇초기화 30 0 0 - -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if pq.command(bot, "/몬스터북초기화", "몬스터북 초기화 완료") == false then
		return ctx:fail(bot:name() .. " 몬스터북 초기화 실패")
	end
	if bot:map_move(HENESYS_PARK) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function drop_card(ctx, bot)
	local spawn = bot:request(resp.spawn_item, req.normal_chat { message = "/아이템드롭 " .. SNAIL_CARD }, function(p)
		return p.item_model ~= nil and p.item_model.id == SNAIL_CARD
	end, 3000)
	if spawn == false then
		return ctx:fail("카드가 떨어지지 않음")
	end
	local drop = bot:drops(SNAIL_CARD)[1]
	if drop == nil then
		return ctx:fail("떨어진 카드를 찾지 못함")
	end
	local x, y = bot:position()
	return req.item_loot { position = { x = x, y = y }, oid = drop.oid }
end

local function profile(ctx, bot, what)
	local p = bot:request(resp.character_profile, req.inspect_character { character_id = bot:id() }, nil, 3000)
	if p == false then
		return ctx:fail("캐릭터 정보 응답 없음: " .. what)
	end
	return p
end

local function register(ctx, bot, card_id)
	local p = bot:request(resp.monster_book_set_card, req.normal_chat { message = "/몬스터북카드 " .. card_id }, function(p)
		return p.card_id == card_id
	end, 3000)
	if p == false then
		return ctx:fail("GM 명령 카드 등록 응답 없음: " .. card_id)
	end
	return p
end

local function collect(ctx)
	local bot = ctx:bot(0)
	local observer = ctx:bot(1)
	if setup(ctx, bot) == false or setup(ctx, observer) == false then
		return false
	end

	local loot = drop_card(ctx, bot)
	if loot == false then
		return false
	end
	local effect = bot:request_on(observer, resp.show_effect, loot, function(p)
		return p.character_id == bot:id()
	end, 3000)
	if effect == false then
		return ctx:fail("다른 캐릭터에게 카드 획득 이펙트가 보이지 않음")
	end
	if check(ctx, effect.type == REGISTER_CARD_EFFECT, "카드 획득 이펙트 종류가 다름: " .. tostring(effect.type)) == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, (bot:items()[SNAIL_CARD] or 0) == 0, "카드가 인벤토리에 들어감") == false then
		return false
	end

	for count = 2, 5 do
		loot = drop_card(ctx, bot)
		if loot == false then
			return false
		end
		local p = bot:request(resp.monster_book_set_card, loot, nil, 3000)
		if p == false then
			return ctx:fail("카드 등록 응답 없음: " .. count)
		end
		if check(ctx, p.success and p.card_id == SNAIL_CARD and p.count == count, string.format("카드 등록 결과가 다름: %s %d/%d", tostring(p.success), p.count, count)) == false then
			return false
		end
	end

	loot = drop_card(ctx, bot)
	if loot == false then
		return false
	end
	local full = bot:request(resp.monster_book_set_card, loot, nil, 3000)
	if full == false then
		return ctx:fail("가득 찬 카드 응답 없음")
	end
	if check(ctx, full.success == false, "6장째 카드가 등록됨") == false then
		return false
	end
	ctx:sleep(300)
	return check(ctx, (bot:items()[SNAIL_CARD] or 0) == 0, "가득 찬 카드가 인벤토리에 남음")
end

local function level_and_cover(ctx)
	local bot = ctx:bot(0)
	if register(ctx, bot, MANO_CARD) == false then
		return false
	end
	for card_id = SNAIL_CARD + 1, SNAIL_CARD + 10 do
		if register(ctx, bot, card_id) == false then
			return false
		end
	end

	local p = profile(ctx, bot, "등록 후")
	if p == false then
		return false
	end
	if check(ctx, p.book_level == 2 and p.book_normal_cards == 11 and p.book_special_cards == 1, string.format("몬스터북 집계가 다름: 레벨 %d 일반 %d 특수 %d", p.book_level, p.book_normal_cards, p.book_special_cards)) == false then
		return false
	end

	local cover = bot:request(resp.monster_book_set_cover, req.monster_book_cover { card_id = SNAIL_CARD }, nil, 3000)
	if cover == false then
		return ctx:fail("표지 변경 응답 없음")
	end
	if check(ctx, cover.card_id == SNAIL_CARD, "표지 카드가 다름: " .. cover.card_id) == false then
		return false
	end
	if bot:request(resp.monster_book_set_cover, req.monster_book_cover { card_id = NOT_OWNED_CARD }, nil, 1000) ~= false then
		return ctx:fail("가지지 않은 카드로 표지가 바뀜")
	end

	p = profile(ctx, bot, "표지 변경 후")
	if p == false then
		return false
	end
	return check(ctx, p.book_cover_mob_id == SNAIL, "표지 몬스터가 다름: " .. p.book_cover_mob_id)
end

local function persistence(ctx)
	local bot = ctx:bot(0)
	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	local p = profile(ctx, bot, "채널 이동 후")
	if p == false then
		return false
	end
	return check(ctx, p.book_level == 2 and p.book_normal_cards == 11 and p.book_special_cards == 1 and p.book_cover_mob_id == SNAIL,
		string.format("채널 이동 후 몬스터북이 다름: 레벨 %d 일반 %d 특수 %d 표지 %d", p.book_level, p.book_normal_cards, p.book_special_cards, p.book_cover_mob_id))
end

test_suite {
	name = "Monster book: 카드 등록·표지·저장",
	bot_count = 2,

	scenarios = {
		collect,
		level_and_cover,
		persistence,
	},
}
