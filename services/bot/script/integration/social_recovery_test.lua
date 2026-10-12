local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local HENESYS = 100000000
local MEGAPHONE = { sn = 10000804 }

local NOTES = { send = 0, remove = 1, list = 2 }
local NOTES_PACKET = { list = 3, sent = 4 }
local NOTE_KIND = { population = 1 }
local CASH_SHOP = { gift = 4, gift_sent = 59 }
local SETTLE_WAIT_MS = 40000

local inbox = {}

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function listen(ctx)
	ctx:hook("social_recovery", function(_, bot, p, name)
		if name ~= "notes" then
			return
		end
		local box = inbox[bot:id()]
		if box == nil then
			box = {}
			inbox[bot:id()] = box
		end
		box[#box + 1] = p
	end)
end

local function send(bot, request)
	inbox[bot:id()] = {}
	return bot:send(request)
end

local function expect(ctx, bot, action, what)
	for _ = 1, 50 do
		local box = inbox[bot:id()] or {}
		for i, p in ipairs(box) do
			if p.action == action then
				table.remove(box, i)
				return p
			end
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s 응답 없음", bot:name(), what))
end

local function setup(ctx, bot)
	if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= HENESYS and bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function list_notes(ctx, bot)
	send(bot, req.notes { action = NOTES.list })
	return expect(ctx, bot, NOTES_PACKET.list, "쪽지 목록")
end

local function clear_notes(ctx, bot)
	local list = list_notes(ctx, bot)
	if list == false then
		return false
	end
	if #list.notes == 0 then
		return true
	end

	local removals = {}
	for i, note in ipairs(list.notes) do
		removals[i] = { id = note.id, kind = note.kind }
	end
	send(bot, req.notes { action = NOTES.remove, removals = removals })
	ctx:sleep(500)
	list = list_notes(ctx, bot)
	if list == false then
		return false
	end
	return check(ctx, #list.notes == 0, string.format("%s 쪽지 %d개가 지워지지 않음", bot:name(), #list.notes))
end

-- The receiver of a cash shop gift thanks the buyer, which leaves a population note in the buyer's box.
local function thanks(ctx, buyer, receiver)
	inbox[buyer:id()] = {}
	inbox[receiver:id()] = {}
	listen(ctx)
	if setup(ctx, buyer) == false or setup(ctx, receiver) == false then
		return false
	end
	if clear_notes(ctx, buyer) == false then
		return false
	end
	if pq.command(buyer, "/캐시얻기 5000 0", "캐시 잔액: ") == false then
		return ctx:fail("캐시 지급 실패")
	end

	if buyer:enter_cash_shop() == false then
		return ctx:fail("캐시샵 입장 실패")
	end
	local gift = buyer:request(resp.cash_shop_result, req.cash_shop_operation {
		action = CASH_SHOP.gift,
		commodity_sn = MEGAPHONE.sn,
		recipient = receiver:name(),
		message = "선물",
	}, function(p)
		return p.kind == CASH_SHOP.gift_sent or p.kind == CASH_SHOP.gift_sent + 1
	end, 5000)
	if gift == false or check(ctx, gift.kind == CASH_SHOP.gift_sent, "캐시샵 선물 실패") == false then
		return false
	end
	if buyer:leave_cash_shop() == false then
		return ctx:fail("캐시샵 퇴장 실패")
	end
	if receiver:enter_cash_shop() == false or receiver:leave_cash_shop() == false then
		return ctx:fail("받는 캐릭터 캐시샵 입장·퇴장 실패")
	end

	send(receiver, req.notes { action = NOTES.send, target = buyer:name(), message = "감사합니다", kind = NOTE_KIND.population })
	if expect(ctx, receiver, NOTES_PACKET.sent, "감사 쪽지 전송") == false then
		return false
	end
	local list = list_notes(ctx, buyer)
	if list == false then
		return false
	end
	local note = list.notes[1]
	if check(ctx, #list.notes == 1 and note.kind == NOTE_KIND.population, "감사 쪽지 목록이 다름: " .. #list.notes) == false then
		return false
	end
	return note
end

-- The fault is server wide, so it is turned off right after the one request it targets.
local function remove_unknown(ctx, buyer, note, mode)
	if fault.rpc(buyer, "Notes", mode) == false then
		return ctx:fail("Notes 장애 주입 실패")
	end
	send(buyer, req.notes { action = NOTES.remove, removals = { { id = note.id, kind = note.kind } } })
	ctx:sleep(1000)
	fault.rpc(buyer, "Notes", "off")
	return true
end

local function wait_population(ctx, bot, population, what)
	for _ = 1, SETTLE_WAIT_MS / 1000 do
		if bot:population() == population then
			return true
		end
		ctx:sleep(1000)
	end
	return check(ctx, false, string.format("%s %s: 인기도 %d (기대 %d)", bot:name(), what, bot:population(), population))
end

local function settled(ctx, buyer, population)
	ctx:sleep(3000)
	if check(ctx, buyer:population() == population, string.format("정리 뒤 인기도 %d (기대 %d, 중복 지급 없음)", buyer:population(), population)) == false then
		return false
	end
	local list = list_notes(ctx, buyer)
	if list == false then
		return false
	end
	return check(ctx, #list.notes == 0, "정리 뒤 감사 쪽지가 남음: " .. #list.notes)
end

local function lost_after_remove_timer(ctx)
	local buyer, receiver = ctx:bot(0), ctx:bot(1)
	local note = thanks(ctx, buyer, receiver)
	if note == false then
		return false
	end
	local before = buyer:population()
	if remove_unknown(ctx, buyer, note, "lost") == false then
		return false
	end
	if check(ctx, buyer:population() == before, "장애 직후 인기도가 오름") == false then
		return false
	end
	if wait_population(ctx, buyer, before + 1, "주기 정리로 인기도 지급") == false then
		return false
	end
	return settled(ctx, buyer, before + 1)
end

local function lost_after_remove_relogin(ctx)
	local buyer, receiver = ctx:bot(2), ctx:bot(3)
	local note = thanks(ctx, buyer, receiver)
	if note == false then
		return false
	end
	local before = buyer:population()
	if remove_unknown(ctx, buyer, note, "lost") == false then
		return false
	end
	if buyer:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_population(ctx, buyer, before + 1, "재접속 정리로 인기도 지급") == false then
		return false
	end
	if buyer:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	return settled(ctx, buyer, before + 1)
end

local function not_sent_remove(ctx)
	local buyer, receiver = ctx:bot(0), ctx:bot(1)
	local note = thanks(ctx, buyer, receiver)
	if note == false then
		return false
	end
	local before = buyer:population()
	if remove_unknown(ctx, buyer, note, "on") == false then
		return false
	end
	if wait_population(ctx, buyer, before + 1, "주기 정리로 다시 삭제") == false then
		return false
	end
	return settled(ctx, buyer, before + 1)
end

test_suite {
	name = "Social: 감사 쪽지 삭제 응답 유실 복구",
	bot_count = 4,

	scenarios = {
		lost_after_remove_timer,
		lost_after_remove_relogin,
		not_sent_remove,
	},
}
