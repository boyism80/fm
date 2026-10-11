local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local MEGAPHONE = { sn = 10000804 }

local WHISPER = { find = 5, send = 6, find_result = 9, send_result = 10, receive = 18 }
local WHISPER_LOCATION = { map = 1 }
local MESSENGER = { enter = 0, self_enter = 1, leave = 2, invite = 3, invite_result = 4, chat = 6 }
local NOTES = { send = 0, remove = 1, list = 2 }
local NOTES_PACKET = { list = 3, sent = 4, failed = 5, refresh = 7 }
local NOTE_KIND = { normal = 0, population = 1 }
local NOTE_RESULT = { wrong_name = 1, inbox_full = 2 }
local POPULATION = { down = 0, up = 1 }
local POPULATION_RESULT = { given = 0, wrong_name = 1, level_too_low = 2, already_today = 3, received = 5 }
local FAMILY_RESULT = { failed = 5, junior_level_too_high = 8 }
local CASH_SHOP = { gift = 4, gift_sent = 59 }

local SOCIAL = {
	whisper = true,
	messenger = true,
	notes = true,
	population = true,
	family_result = true,
	family_invite = true,
	family_invite_result = true,
	family_joined = true,
	family_info = true,
	family_chart = true,
}

local inbox = {}

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function listen(ctx)
	ctx:hook("social", function(_, bot, p, name)
		if SOCIAL[name] == nil then
			return
		end
		local box = inbox[bot:id()]
		if box == nil then
			box = {}
			inbox[bot:id()] = box
		end
		box[#box + 1] = { name = name, packet = p }
	end)
end

-- Each queue runs on its own bots, so a flow only resets its own inboxes.
local function flow(fn, ...)
	local ids = { ... }
	return function(ctx)
		local bots = {}
		for i, id in ipairs(ids) do
			bots[i] = ctx:bot(id)
			inbox[bots[i]:id()] = {}
		end
		listen(ctx)
		return fn(ctx, unpack(bots))
	end
end

local function send(bot, request)
	inbox[bot:id()] = {}
	return bot:send(request)
end

local function expect(ctx, bot, name, match, what)
	for _ = 1, 50 do
		local box = inbox[bot:id()] or {}
		for i, item in ipairs(box) do
			if item.name == name and (match == nil or match(item.packet)) then
				table.remove(box, i)
				return item.packet
			end
		end
		ctx:sleep(100)
	end
	return ctx:fail(string.format("%s %s 응답 없음", bot:name(), what))
end

local function setup(ctx, bot, level)
	if pq.command(bot, string.format("/봇초기화 %d 0 0 -", level), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= HENESYS and bot:map_move(HENESYS) == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패")
	end
	return true
end

local function list_notes(ctx, bot)
	send(bot, req.notes { action = NOTES.list })
	return expect(ctx, bot, "notes", function(p)
		return p.action == NOTES_PACKET.list
	end, "쪽지 목록")
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

local function whisper_flow(ctx, a, b)
	if setup(ctx, a, 30) == false or setup(ctx, b, 30) == false then
		return false
	end

	send(a, req.whisper { mode = WHISPER.find, target = b:name() })
	local found = expect(ctx, a, "whisper", function(p)
		return p.mode == WHISPER.find_result
	end, "귓속말 찾기")
	if found == false then
		return false
	end
	if check(ctx, found.name == b:name() and found.location == WHISPER_LOCATION.map and found.map_id == HENESYS,
		string.format("찾기 결과가 다름: location=%d map=%d", found.location, found.map_id)) == false then
		return false
	end

	send(a, req.whisper { mode = WHISPER.send, target = b:name(), message = "귓속말 테스트" })
	local received = expect(ctx, b, "whisper", function(p)
		return p.mode == WHISPER.receive
	end, "귓속말 수신")
	if received == false then
		return false
	end
	if check(ctx, received.name == a:name() and received.message == "귓속말 테스트", "받은 귓속말 내용이 다름") == false then
		return false
	end
	local sent = expect(ctx, a, "whisper", function(p)
		return p.mode == WHISPER.send_result
	end, "귓속말 전송 결과")
	if sent == false or check(ctx, sent.success, "귓속말 전송이 실패로 응답됨") == false then
		return false
	end

	send(a, req.whisper { mode = WHISPER.send, target = "nobody_here", message = "x" })
	local missing = expect(ctx, a, "whisper", function(p)
		return p.mode == WHISPER.send_result
	end, "없는 대상 귓속말 결과")
	if missing == false then
		return false
	end
	return check(ctx, missing.success == false, "없는 대상에게 귓속말이 성공함")
end

local function messenger_flow(ctx, a, b)
	if setup(ctx, a, 30) == false or setup(ctx, b, 30) == false then
		return false
	end

	send(a, req.messenger { action = MESSENGER.enter, room_id = 0 })
	local opened = expect(ctx, a, "messenger", function(p)
		return p.action == MESSENGER.self_enter
	end, "메신저 방 생성")
	if opened == false or check(ctx, opened.position == 0, "방장 자리가 0이 아님: " .. opened.position) == false then
		return false
	end

	send(a, req.messenger { action = MESSENGER.invite, text = b:name() })
	local result = expect(ctx, a, "messenger", function(p)
		return p.action == MESSENGER.invite_result
	end, "메신저 초대 결과")
	if result == false or check(ctx, result.success and result.name == b:name(), "메신저 초대가 실패함") == false then
		return false
	end
	local invite = expect(ctx, b, "messenger", function(p)
		return p.action == MESSENGER.invite
	end, "메신저 초대 수신")
	if invite == false or check(ctx, invite.name == a:name() and invite.room_id ~= 0, "받은 초대 내용이 다름") == false then
		return false
	end

	send(b, req.messenger { action = MESSENGER.enter, room_id = invite.room_id })
	local joined = expect(ctx, b, "messenger", function(p)
		return p.action == MESSENGER.self_enter
	end, "메신저 입장")
	if joined == false or check(ctx, joined.position == 1, "입장 자리가 1이 아님: " .. joined.position) == false then
		return false
	end
	local host = expect(ctx, b, "messenger", function(p)
		return p.action == MESSENGER.enter and p.name == a:name()
	end, "기존 참가자 정보")
	if host == false or check(ctx, host.position == 0 and #host.look > 0, "기존 참가자 정보가 다름") == false then
		return false
	end
	local guest = expect(ctx, a, "messenger", function(p)
		return p.action == MESSENGER.enter and p.name == b:name()
	end, "새 참가자 알림")
	if guest == false or check(ctx, guest.position == 1, "새 참가자 자리가 다름: " .. guest.position) == false then
		return false
	end

	send(a, req.messenger { action = MESSENGER.chat, text = a:name() .. " : 메신저 테스트" })
	local chat = expect(ctx, b, "messenger", function(p)
		return p.action == MESSENGER.chat
	end, "메신저 채팅")
	if chat == false or check(ctx, chat.text:find("메신저 테스트", 1, true) ~= nil, "메신저 채팅 내용이 다름") == false then
		return false
	end

	send(b, req.messenger { action = MESSENGER.leave })
	local left = expect(ctx, a, "messenger", function(p)
		return p.action == MESSENGER.leave
	end, "메신저 퇴장 알림")
	if left == false or check(ctx, left.position == 1, "퇴장한 자리가 다름: " .. left.position) == false then
		return false
	end
	send(a, req.messenger { action = MESSENGER.leave })
	ctx:sleep(300)
	return true
end

local function notes_flow(ctx, a, b)
	if setup(ctx, a, 30) == false or setup(ctx, b, 30) == false then
		return false
	end
	if clear_notes(ctx, b) == false then
		return false
	end

	send(a, req.notes { action = NOTES.send, target = b:name(), message = "쪽지 테스트", kind = NOTE_KIND.normal })
	if expect(ctx, a, "notes", function(p)
		return p.action == NOTES_PACKET.sent
	end, "쪽지 전송") == false then
		return false
	end
	if expect(ctx, b, "notes", function(p)
		return p.action == NOTES_PACKET.refresh
	end, "쪽지 도착 알림") == false then
		return false
	end

	local list = list_notes(ctx, b)
	if list == false then
		return false
	end
	local note = list.notes[1]
	if check(ctx, #list.notes == 1 and note.sender == a:name() and note.message == "쪽지 테스트" and note.kind == NOTE_KIND.normal,
		"받은 쪽지 목록이 다름: " .. #list.notes) == false then
		return false
	end

	send(a, req.notes { action = NOTES.send, target = b:name(), message = "감사합니다", kind = NOTE_KIND.population })
	local thanks = expect(ctx, a, "notes", function(p)
		return p.action == NOTES_PACKET.failed
	end, "선물 없는 감사 쪽지")
	if thanks == false or check(ctx, thanks.status == NOTE_RESULT.inbox_full, "선물 없는 감사 쪽지 결과 " .. thanks.status) == false then
		return false
	end
	send(a, req.notes { action = NOTES.send, target = "nobody_here", message = "x", kind = NOTE_KIND.normal })
	local missing = expect(ctx, a, "notes", function(p)
		return p.action == NOTES_PACKET.failed
	end, "없는 대상 쪽지")
	if missing == false or check(ctx, missing.status == NOTE_RESULT.wrong_name, "없는 대상 쪽지 결과 " .. missing.status) == false then
		return false
	end

	return clear_notes(ctx, b)
end

local function thanks_flow(ctx, buyer, receiver)
	if setup(ctx, buyer, 30) == false or setup(ctx, receiver, 30) == false then
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
	inbox[buyer:id()] = {}
	inbox[receiver:id()] = {}

	send(receiver, req.notes { action = NOTES.send, target = buyer:name(), message = "감사합니다", kind = NOTE_KIND.population })
	if expect(ctx, receiver, "notes", function(p)
		return p.action == NOTES_PACKET.sent
	end, "감사 쪽지 전송") == false then
		return false
	end
	send(receiver, req.notes { action = NOTES.send, target = buyer:name(), message = "또 감사", kind = NOTE_KIND.population })
	local again = expect(ctx, receiver, "notes", function(p)
		return p.action == NOTES_PACKET.failed
	end, "두 번째 감사 쪽지")
	if again == false or check(ctx, again.status == NOTE_RESULT.inbox_full, "두 번째 감사 쪽지 결과 " .. again.status) == false then
		return false
	end

	local list = list_notes(ctx, buyer)
	if list == false then
		return false
	end
	local note = list.notes[1]
	if check(ctx, #list.notes == 1 and note.sender == receiver:name() and note.kind == NOTE_KIND.population, "감사 쪽지 목록이 다름: " .. #list.notes) == false then
		return false
	end

	local before = buyer:population()
	if buyer:request(resp.update_stats, req.notes { action = NOTES.remove, removals = { { id = note.id, kind = note.kind } } }, nil, 5000) == false then
		return ctx:fail("감사 쪽지 삭제 후 인기도 갱신 없음")
	end
	if check(ctx, buyer:population() == before + 1, string.format("감사 쪽지 인기도 %d (기대 %d)", buyer:population(), before + 1)) == false then
		return false
	end
	return clear_notes(ctx, buyer)
end

local function give(ctx, giver, target, mode)
	send(giver, req.give_population { target_id = target:id(), mode = mode })
	return expect(ctx, giver, "population", function(p)
		return p.status ~= POPULATION_RESULT.received
	end, "인기도 결과")
end

local function population_flow(ctx, a, b, c)
	if setup(ctx, a, 30) == false or setup(ctx, b, 30) == false or setup(ctx, c, 10) == false then
		return false
	end

	local p = give(ctx, a, b, POPULATION.up)
	if p == false or check(ctx, p.status == POPULATION_RESULT.given and p.name == b:name() and p.value == 1, "인기도 올리기 결과가 다름: " .. p.status) == false then
		return false
	end
	local received = expect(ctx, b, "population", function(r)
		return r.status == POPULATION_RESULT.received
	end, "인기도 받음 알림")
	if received == false or check(ctx, received.name == a:name() and received.mode == POPULATION.up, "받은 인기도 알림이 다름") == false then
		return false
	end
	ctx:sleep(300)
	if check(ctx, b:population() == 1, "받은 캐릭터 인기도 " .. b:population()) == false then
		return false
	end

	p = give(ctx, a, c, POPULATION.up)
	if p == false or check(ctx, p.status == POPULATION_RESULT.already_today, "하루 제한 결과가 다름: " .. p.status) == false then
		return false
	end
	p = give(ctx, c, a, POPULATION.up)
	if p == false or check(ctx, p.status == POPULATION_RESULT.level_too_low, "레벨 제한 결과가 다름: " .. p.status) == false then
		return false
	end
	p = give(ctx, b, b, POPULATION.up)
	if p == false or check(ctx, p.status == POPULATION_RESULT.wrong_name, "자기 자신 결과가 다름: " .. p.status) == false then
		return false
	end

	p = give(ctx, b, a, POPULATION.down)
	if p == false or check(ctx, p.status == POPULATION_RESULT.given and p.value == -1, "인기도 내리기 결과가 다름: " .. p.status) == false then
		return false
	end
	ctx:sleep(300)
	return check(ctx, a:population() == -1, "내려간 인기도 " .. a:population())
end

local function family_flow(ctx, senior, junior, low)
	if setup(ctx, senior, 50) == false or setup(ctx, junior, 40) == false or setup(ctx, low, 30) == false then
		return false
	end

	send(senior, req.family_remove_senior {})
	send(junior, req.family_remove_senior {})
	ctx:sleep(1000)

	send(low, req.family_invite { family = { name = junior:name() } })
	local rejected = expect(ctx, low, "family_result", nil, "레벨 높은 주니어 초대")
	if rejected == false or check(ctx, rejected.code == FAMILY_RESULT.junior_level_too_high, "레벨 높은 주니어 초대 결과 " .. rejected.code) == false then
		return false
	end

	send(senior, req.family_invite { family = { name = junior:name() } })
	local invite = expect(ctx, junior, "family_invite", nil, "패밀리 초대 수신")
	if invite == false or check(ctx, invite.inviter_id == senior:id() and invite.name == senior:name(), "받은 패밀리 초대가 다름") == false then
		return false
	end

	send(junior, req.family_accept { family = { target_id = senior:id(), name = senior:name(), accepted = true } })
	local result = expect(ctx, senior, "family_invite_result", nil, "패밀리 초대 결과")
	if result == false or check(ctx, result.accepted and result.name == junior:name(), "패밀리 초대 결과가 다름") == false then
		return false
	end
	local joined = expect(ctx, junior, "family_joined", nil, "패밀리 가입 알림")
	if joined == false or check(ctx, joined.name == senior:name(), "가입 알림의 선배 이름이 다름") == false then
		return false
	end

	local function pedigree(viewer)
		send(viewer, req.family_pedigree { family = { name = viewer:name() } })
		local chart = expect(ctx, viewer, "family_chart", nil, viewer:name() .. " 가계도")
		if chart == false then
			return nil
		end
		local ids = {}
		for _, member in ipairs(chart.state.members) do
			ids[#ids + 1] = tostring(member.character_id)
		end
		return table.concat(ids, ","), chart.state.member_count
	end
	local tree = senior:id() .. "," .. senior:id() .. "," .. junior:id()
	for _, viewer in ipairs({ junior, senior }) do
		local nodes, count = pedigree(viewer)
		if nodes == nil then
			return false
		end
		if check(ctx, nodes == tree and count == 2,
			string.format("%s 가계도가 다름: nodes=%s count=%s (기대 %s, 2)", viewer:name(), nodes, tostring(count), tree)) == false then
			return false
		end
	end

	send(junior, req.family_remove_senior {})
	local info = expect(ctx, junior, "family_info", function(p)
		return p.state.leader_id == 0
	end, "패밀리 탈퇴 후 정보")
	if info == false then
		return false
	end

	send(junior, req.family_remove_senior {})
	local again = expect(ctx, junior, "family_result", nil, "선배 없는 탈퇴")
	if again == false then
		return false
	end
	return check(ctx, again.code == FAMILY_RESULT.failed, "선배 없는 탈퇴 결과 " .. again.code)
end

test_suite {
	name = "Social: 귓속말·메신저·쪽지·감사 쪽지·인기도·패밀리",
	bot_count = 14,

	scenarios = {
		{
			parallel = {
				flow(whisper_flow, 0, 1),
				flow(messenger_flow, 2, 3),
				flow(notes_flow, 4, 5),
				flow(thanks_flow, 6, 7),
				flow(population_flow, 8, 9, 10),
				flow(family_flow, 11, 12, 13),
			},
		},
	},
}
