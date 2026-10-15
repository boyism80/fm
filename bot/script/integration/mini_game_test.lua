local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local OMOK_SET = 4080000
local MATCH_CARD_SET = 4080100

local TYPE = { omok = 1, match_card = 2 }
local MODE = {
	create = 0x00,
	visit = 0x04,
	exit = 0x0A,
	request_tie = 0x2A,
	answer_tie = 0x2B,
	give_up = 0x2C,
	ready = 0x32,
	expel = 0x34,
	start = 0x35,
	move_omok = 0x38,
	select_card = 0x3C,
}
local OUTCOME = { give_up = 0, tie = 1, win = 2 }
local LEAVE = { kicked = 5, closed = 3 }

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, items)
	if pq.command(bot, string.format("/봇초기화 30 0 0 %s -", items), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= HENESYS and bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 헤네시스 이동 실패")
	end
	return true
end

local function open_room(ctx, owner, watcher, kind, title, password)
	local balloon = owner:request_on(watcher, resp.user_mini_room_balloon, req.mini_room {
		mode = MODE.create,
		type = kind,
		title = title,
		private = password ~= nil,
		password = password,
		piece = 0,
	}, function(p)
		return p.character_id == owner:id() and p.balloon ~= nil
	end, 5000)
	if balloon == false then
		return ctx:fail("미니게임 풍선이 다른 캐릭터에게 보이지 않음")
	end
	if check(ctx, balloon.balloon.type == kind and balloon.balloon.title == title and balloon.balloon.users == 1, "미니게임 풍선 정보가 다름") == false then
		return false
	end
	return balloon.balloon
end

local function visit(ctx, bot, sn, password)
	local p = bot:request(resp.mini_game_enter, req.mini_room { mode = MODE.visit, sn = sn, private = password ~= nil, password = password }, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 미니게임 입장 응답 없음")
	end
	if check(ctx, p.my_slot == 1 and #p.members == 2, "미니게임 입장 정보가 다름") == false then
		return false
	end
	return p
end

local function start(ctx, owner, visitor, loser)
	local ready = visitor:request_on(owner, resp.mini_game_ready, req.mini_room { mode = MODE.ready }, function(p)
		return p.ready
	end, 5000)
	if ready == false then
		return ctx:fail("준비 알림이 주인에게 없음")
	end
	local started = owner:request_on(visitor, resp.mini_game_start, req.mini_room { mode = MODE.start }, nil, 5000)
	if started == false then
		return ctx:fail("게임 시작 알림이 손님에게 없음")
	end
	if check(ctx, started.loser == loser, string.format("시작 패킷 순서 %d (기대 %d)", started.loser, loser)) == false then
		return false
	end
	return started
end

local function move(ctx, bot, other, x, y, stone)
	local p = bot:request_on(other, resp.mini_game_move_omok, req.mini_room { mode = MODE.move_omok, x = x, y = y, stone = stone }, function(p)
		return p.x == x and p.y == y
	end, 5000)
	if p == false then
		return ctx:fail(string.format("%s 오목 착수 (%d, %d) 알림 없음", bot:name(), x, y))
	end
	return true
end

local function game_over(ctx, p, outcome, winner, what)
	if p == false then
		return ctx:fail(what .. ": 게임 결과 알림 없음")
	end
	if check(ctx, p.outcome == outcome, string.format("%s: 결과 %d (기대 %d)", what, p.outcome, outcome)) == false then
		return false
	end
	if outcome ~= OUTCOME.tie and check(ctx, p.winner == winner, string.format("%s: 승자 %d (기대 %d)", what, p.winner, winner)) == false then
		return false
	end
	return true
end

local function record_is(ctx, record, wins, ties, losses, what)
	return check(ctx, record.wins == wins and record.ties == ties and record.losses == losses,
		string.format("%s 전적 %d승 %d무 %d패 (기대 %d승 %d무 %d패)", what, record.wins, record.ties, record.losses, wins, ties, losses))
end

local function omok_flow(ctx)
	local owner = ctx:bot(0)
	local visitor = ctx:bot(1)
	if setup(ctx, owner, OMOK_SET .. ":1") == false or setup(ctx, visitor, "-") == false then
		return false
	end

	local balloon = open_room(ctx, owner, visitor, TYPE.omok, "봇 오목", "1234")
	if balloon == false then
		return false
	end
	if check(ctx, balloon.private, "비밀방 표시가 없음") == false then
		return false
	end
	local wrong = visitor:request(resp.mini_game_enter, req.mini_room { mode = MODE.visit, sn = balloon.sn, private = true, password = "0000" }, nil, 1500)
	if check(ctx, wrong == false, "틀린 비밀번호로 입장됨") == false then
		return false
	end
	if visit(ctx, visitor, balloon.sn, "1234") == false then
		return false
	end

	local early = owner:request_on(visitor, resp.mini_game_start, req.mini_room { mode = MODE.start }, nil, 1500)
	if check(ctx, early == false, "준비 전에 게임이 시작됨") == false then
		return false
	end
	if start(ctx, owner, visitor, 1) == false then
		return false
	end
	local out_of_turn = visitor:request_on(owner, resp.mini_game_move_omok, req.mini_room { mode = MODE.move_omok, x = 7, y = 7, stone = 2 }, nil, 1500)
	if check(ctx, out_of_turn == false, "차례가 아닌 착수가 받아들여짐") == false then
		return false
	end
	for x = 0, 3 do
		if move(ctx, owner, visitor, x, 0, 1) == false or move(ctx, visitor, owner, x, 1, 2) == false then
			return false
		end
	end
	local won = owner:request_on(visitor, resp.mini_game_over, req.mini_room { mode = MODE.move_omok, x = 4, y = 0, stone = 1 }, nil, 5000)
	if game_over(ctx, won, OUTCOME.win, 0, "오목 5목") == false then
		return false
	end
	if record_is(ctx, won.records[1], 1, 0, 0, "주인") == false or record_is(ctx, won.records[2], 0, 0, 1, "손님") == false then
		return false
	end

	if start(ctx, owner, visitor, 0) == false then
		return false
	end
	local gave_up = visitor:request_on(owner, resp.mini_game_over, req.mini_room { mode = MODE.give_up }, nil, 5000)
	if game_over(ctx, gave_up, OUTCOME.give_up, 0, "기권") == false then
		return false
	end
	if record_is(ctx, gave_up.records[1], 2, 0, 0, "주인") == false then
		return false
	end

	if start(ctx, owner, visitor, 0) == false then
		return false
	end
	local asked = owner:request_on(visitor, resp.mini_game_request_tie, req.mini_room { mode = MODE.request_tie }, nil, 5000)
	if asked == false then
		return ctx:fail("무승부 요청이 손님에게 없음")
	end
	local tied = visitor:request_on(owner, resp.mini_game_over, req.mini_room { mode = MODE.answer_tie, accept = true }, nil, 5000)
	if game_over(ctx, tied, OUTCOME.tie, 0, "무승부") == false then
		return false
	end
	if record_is(ctx, tied.records[2], 0, 1, 2, "손님") == false then
		return false
	end

	local kicked = owner:request_on(visitor, resp.mini_room_leave, req.mini_room { mode = MODE.expel }, function(p)
		return p.reason == LEAVE.kicked
	end, 5000)
	if kicked == false then
		return ctx:fail("강퇴한 손님이 나가지 않음")
	end
	local removed = owner:request_on(visitor, resp.user_mini_room_balloon, req.mini_room { mode = MODE.exit }, function(p)
		return p.character_id == owner:id() and p.balloon == nil
	end, 5000)
	if removed == false then
		return ctx:fail("방을 닫아도 풍선이 사라지지 않음")
	end
	return true
end

local function pick(ctx, bot, other, first_pick, card)
	local p = bot:request_on(other, resp.mini_game_select_card, req.mini_room { mode = MODE.select_card, first_pick = first_pick, card = card }, function(p)
		return p.card == card
	end, 5000)
	if p == false then
		return ctx:fail(string.format("%s 카드 %d 선택 알림 없음", bot:name(), card))
	end
	return p
end

local function match_card_flow(ctx)
	local owner = ctx:bot(0)
	local visitor = ctx:bot(1)
	if setup(ctx, owner, MATCH_CARD_SET .. ":1") == false or setup(ctx, visitor, "-") == false then
		return false
	end

	local balloon = open_room(ctx, owner, visitor, TYPE.match_card, "봇 짝맞추기", nil)
	if balloon == false then
		return false
	end
	if visit(ctx, visitor, balloon.sn, nil) == false then
		return false
	end
	local started = start(ctx, owner, visitor, 1)
	if started == false then
		return false
	end
	local cards = started.cards
	if check(ctx, cards ~= nil and #cards == 12, "카드 수가 12장이 아님") == false then
		return false
	end

	local pairs_of = {}
	local miss = nil
	for i, card in ipairs(cards) do
		pairs_of[card] = pairs_of[card] or {}
		table.insert(pairs_of[card], i - 1)
	end
	for i = 2, #cards do
		if cards[i] ~= cards[1] then
			miss = i - 1
			break
		end
	end

	if pick(ctx, owner, visitor, true, 0) == false then
		return false
	end
	local wrong = pick(ctx, owner, visitor, false, miss)
	if wrong == false then
		return false
	end
	if check(ctx, wrong.result == 0 and wrong.first_card == 0, "틀린 짝 결과가 다름") == false then
		return false
	end

	local matched = 0
	for card = 0, 5 do
		local slots = pairs_of[card]
		if pick(ctx, visitor, owner, true, slots[1]) == false then
			return false
		end
		matched = matched + 1
		if matched == 6 then
			local over = visitor:request_on(owner, resp.mini_game_over, req.mini_room { mode = MODE.select_card, first_pick = false, card = slots[2] }, nil, 5000)
			if game_over(ctx, over, OUTCOME.win, 1, "짝맞추기") == false then
				return false
			end
			if record_is(ctx, over.records[1], 0, 0, 1, "주인") == false or record_is(ctx, over.records[2], 1, 0, 0, "손님") == false then
				return false
			end
		else
			local right = pick(ctx, visitor, owner, false, slots[2])
			if right == false then
				return false
			end
			if check(ctx, right.result == 3, "맞춘 짝 결과가 다름") == false then
				return false
			end
		end
	end

	local closed = owner:request_on(visitor, resp.mini_room_leave, req.mini_room { mode = MODE.exit }, function(p)
		return p.reason == LEAVE.closed
	end, 5000)
	if closed == false then
		return ctx:fail("주인이 나가도 손님 창이 닫히지 않음")
	end
	return true
end

test_suite {
	name = "MiniGame: 오목·짝맞추기",
	bot_count = 2,

	scenarios = {
		omok_flow,
		match_card_flow,
	},
}
