local pq = require("script/integration/lib/party_quest")

local VILLAGE = 680000000
local LOBBY = 680000200
local CATHEDRAL = 680000210
local EXIT_MAP = 680000500

local CLARANCE = 9201005
local VALENTINA = 9201006
local ANNA = 9201007
local MARGARET = 9201013
local PILA = 9201014
local ANGELIQUE = 9201036

local RING_BOX = 2240004
local ENGAGEMENT_RING = 4210000
local TICKET = 5251004
local INVITATION = 4211000
local INVITED = 4212000
local PERMIT = 4213001
local GIFT = 4000000

local DIALOG = { normal = 0, yes_no = 1, list = 4 }
local RING = { propose = 0, answer = 2, invite = 5, wishlist = 9 }
local ENGAGE = { engaged = 11, married = 12, divorced = 14, reserved = 16 }
local GIFT_MODE = { give = 6, receive = 7, close = 8, open_give = 9, open_receive = 10, given = 11, received = 15 }
local DUEY = { send = 3, receive = 5, delete = 6, close = 8 }
local ETC = 4

local state = {}

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, items)
	if pq.command(bot, string.format("/봇초기화 30 0 1000000 %s -", items), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map_move(VILLAGE) == false then
		return ctx:fail(bot:name() .. " 웨딩빌리지 이동 실패")
	end
	return true
end

local function has(ctx, bot, item_id, count)
	local have = bot:items()[item_id] or 0
	return check(ctx, have == count, string.format("%s 인벤토리 %d 개수 %d (기대 %d)", bot:name(), item_id, have, count))
end

local function click(ctx, bot, template)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false or dlg == nil then
		return ctx:fail(bot:name() .. " NPC 대화 응답 없음: " .. template)
	end
	return dlg
end

local function said(ctx, bot, dlg, expected, what)
	if dlg == nil or dlg == false or dlg.text == nil or dlg.text:find(expected, 1, true) == nil then
		bot:dialog(false)
		return ctx:fail(string.format("%s %s: 예상과 다른 대화 (%s)", bot:name(), what, dlg and tostring(dlg.text) or "nil"))
	end
	return dlg
end

local function talk(ctx, bot, template, expected, what)
	local dlg = click(ctx, bot, template)
	if dlg == false then
		return false
	end
	if said(ctx, bot, dlg, expected, what) == false then
		return false
	end
	bot:dialog(false)
	return true
end

local function command(ctx, bot, text, expected)
	if pq.command(bot, text, expected) == false then
		return ctx:fail(bot:name() .. " 명령 실패: " .. text)
	end
	return true
end

local function shift_clock(ctx, bot, text, expected)
	if command(ctx, bot, text, expected) == false then
		return false
	end
	ctx:sleep(2000)
	return true
end

local function settle_divorce(ctx, bot)
	if bot:map_move(VILLAGE) == false then
		return ctx:fail(bot:name() .. " 웨딩빌리지 이동 실패")
	end
	for _ = 1, 3 do
		local dlg = click(ctx, bot, PILA)
		if dlg == false then
			return false
		end
		if dlg.selections == nil then
			bot:dialog(false)
			return true
		end
		dlg = bot:dialog(true, 0)
		if dlg == nil then
			return ctx:fail(bot:name() .. " 이혼 상담 응답 없음")
		end
		if dlg.text:find("지나셨군요", 1, true) ~= nil then
			if bot:dialog(true) == nil then
				return ctx:fail(bot:name() .. " 이혼 확인 질문 없음")
			end
			local p = bot:request(resp.engage_result, req.dialog { dialog_type = DIALOG.yes_no, next = true }, function(p)
				return p.result == ENGAGE.divorced
			end, 5000)
			return check(ctx, p ~= false, bot:name() .. " 이혼 결과가 오지 않음")
		end
		if dlg.text:find("안타깝군요", 1, true) ~= nil then
			bot:dialog(true)
			bot:dialog(true)
		end
		bot:dialog(false)
		if shift_clock(ctx, bot, "/시간가속 3.01:00:00", "시간가속 적용 요청") == false then
			return false
		end
	end
	return ctx:fail(bot:name() .. " 이전 결혼 정리 실패")
end

local function reset(ctx, bot)
	local p = bot:request(resp.notice, req.normal_chat { message = "/파혼" }, function(p)
		return p.message:find("파혼 완료", 1, true) ~= nil or p.message:find("약혼 상태가 아닙니다", 1, true) ~= nil
	end)
	if p == false then
		return ctx:fail(bot:name() .. " 파혼 명령 응답 없음")
	end
	return settle_divorce(ctx, bot)
end

local function clear_duey(ctx, bot)
	bot:send(req.duey { mode = DUEY.close })
	local p = bot:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 택배 보관함이 열리지 않음")
	end
	for _, parcel in ipairs(p.parcels) do
		if bot:request(resp.duey_removed, req.duey { mode = DUEY.delete, parcel_id = parcel.id }, nil, 3000) == false then
			return ctx:fail(bot:name() .. " 택배 정리 실패")
		end
	end
	bot:send(req.duey { mode = DUEY.close })
	return true
end

local function prepare(ctx)
	local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
	if command(ctx, bride, "/성별 1", "성별 변경: 1") == false or command(ctx, groom, "/성별 0", "성별 변경: 0") == false then
		return false
	end
	if reset(ctx, groom) == false or reset(ctx, bride) == false then
		return false
	end
	if shift_clock(ctx, groom, "/현재시간초기화", "현재 시간 보정 초기화") == false then
		return false
	end
	if clear_duey(ctx, guest) == false then
		return false
	end
	if setup(ctx, groom, string.format("%d:1,%d:1", RING_BOX, TICKET)) == false then
		return false
	end
	if setup(ctx, bride, "-") == false then
		return false
	end
	return setup(ctx, guest, GIFT .. ":5")
end

local function engage(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	local x, y = bride:position()
	if pq.move(groom, x, y) == false then
		return ctx:fail("신부 앞으로 이동 실패")
	end

	local proposal = groom:request_on(bride, resp.engage_request, req.ring_action { mode = RING.propose, name = bride:name(), item_id = RING_BOX }, nil, 5000)
	if proposal == false then
		return ctx:fail("청혼 요청이 신부에게 오지 않음")
	end
	if check(ctx, proposal.name == groom:name() and proposal.character_id == groom:id(), "청혼 요청 내용이 다름: " .. tostring(proposal.name)) == false then
		return false
	end

	local engaged = bride:request_on(groom, resp.engage_result, req.ring_action { mode = RING.answer, accepted = true, name = groom:name(), character_id = groom:id() }, function(p)
		return p.result == ENGAGE.engaged
	end, 5000)
	if engaged == false then
		return ctx:fail("약혼 결과가 신랑에게 오지 않음")
	end
	local marriage = engaged.marriage
	if check(ctx, marriage ~= nil and marriage.groom_id == groom:id() and marriage.bride_id == bride:id(), "약혼 레코드의 신랑/신부가 다름") == false then
		return false
	end
	ctx:sleep(500)
	if has(ctx, groom, RING_BOX, 0) == false or has(ctx, groom, ENGAGEMENT_RING, 1) == false or has(ctx, bride, ENGAGEMENT_RING, 1) == false then
		return false
	end
	state.marriage_id = marriage.id
	return talk(ctx, groom, PILA, "근심이 가득", "약혼 상태 필라")
end

local function reserve(ctx)
	local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
	local dlg = click(ctx, groom, MARGARET)
	if dlg == false then
		return false
	end
	if check(ctx, dlg.selections ~= nil, "마가렛 예약 선택지가 오지 않음") == false then
		groom:dialog(false)
		return false
	end
	if said(ctx, groom, groom:dialog(true, 3), "약혼했군요", "예약 확인") == false then
		return false
	end
	local input = groom:request(resp.engage_request, req.dialog { dialog_type = DIALOG.normal, next = true }, function(p)
		return p.mode == RING.wishlist
	end, 5000)
	if input == false then
		return ctx:fail("예약 후 위시리스트 입력창이 오지 않음")
	end
	ctx:sleep(500)
	if has(ctx, groom, TICKET, 0) == false then
		return false
	end
	if talk(ctx, groom, MARGARET, "결혼식 예약접수를 완료하기 전에", "위시리스트 다시 열기") == false then
		return false
	end

	groom:send(req.ring_action { mode = RING.wishlist, wishes = { "모래 주머니", "달팽이 껍질" } })
	ctx:sleep(1000)
	if has(ctx, groom, INVITATION, 5) == false then
		return false
	end
	if talk(ctx, groom, MARGARET, "이미 위시리스트를 등록", "중복 위시리스트") == false then
		return false
	end
	local reserved = bride:request(resp.engage_result, req.ring_action { mode = RING.wishlist, wishes = { "파란 달팽이 껍질" } }, function(p)
		return p.result == ENGAGE.reserved
	end, 5000)
	if reserved == false then
		return ctx:fail("두 사람 위시리스트 등록 후 예약 완료가 오지 않음")
	end
	if has(ctx, bride, INVITATION, 5) == false then
		return false
	end

	local arrival = groom:request_on(guest, resp.duey_arrival, req.ring_action {
		mode = RING.invite,
		name = guest:name(),
		marriage_id = state.marriage_id,
		slot = groom:slot(INVITATION),
	}, nil, 5000)
	if arrival == false then
		return ctx:fail("청첩장 퀵배송이 하객에게 오지 않음")
	end
	ctx:sleep(500)
	if has(ctx, groom, INVITATION, 4) == false then
		return false
	end
	local box = guest:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
	if box == false or #box.parcels ~= 1 then
		return ctx:fail("하객 보관함에 청첩장이 없음")
	end
	if guest:request(resp.duey_removed, req.duey { mode = DUEY.receive, parcel_id = box.parcels[1].id }, nil, 3000) == false then
		return ctx:fail("청첩장 수령 실패")
	end
	guest:send(req.duey { mode = DUEY.close })
	ctx:sleep(500)
	if has(ctx, guest, INVITED, 1) == false then
		return false
	end

	if talk(ctx, groom, MARGARET, "이미 결혼식이 예약", "신랑 승낙서 요청") == false then
		return false
	end
	if talk(ctx, bride, MARGARET, "주례 승낙서#k를 드렸어요", "신부 승낙서 요청") == false then
		return false
	end
	ctx:sleep(500)
	return has(ctx, bride, PERMIT, 1)
end

local function enter_lobby(ctx, bot, selected)
	local dlg = click(ctx, bot, CLARANCE)
	if dlg == false then
		return false
	end
	local p, name = bot:request({ resp.warp, resp.dialog }, req.dialog { dialog_type = DIALOG.list, next = true, selected = selected }, function(p, name)
		return name == resp.dialog or p.character.map == LOBBY
	end, 5000)
	if p == false or name == resp.dialog then
		bot:dialog(false)
		return ctx:fail(bot:name() .. " 웨딩홀 로비로 이동하지 않음: " .. (p and tostring(p.text) or "응답 없음"))
	end
	return true
end

local function ceremony(ctx)
	local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
	local dlg = click(ctx, guest, CLARANCE)
	if dlg == false then
		return false
	end
	if said(ctx, guest, guest:dialog(true, 1), "현재 시작되어 있는 결혼식이 없군요", "식 전 하객 입장") == false then
		return false
	end
	guest:dialog(false)

	if enter_lobby(ctx, bride, 0) == false or pq.wait_map(ctx, groom, LOBBY, 5000) == false then
		return false
	end
	if has(ctx, bride, PERMIT, 0) == false then
		return false
	end
	if enter_lobby(ctx, guest, 1) == false then
		return false
	end

	dlg = click(ctx, guest, ANGELIQUE)
	if dlg == false then
		return false
	end
	local open = guest:request(resp.wedding_gift, req.dialog { dialog_type = DIALOG.list, next = true, selected = 0 }, function(p)
		return p.mode == GIFT_MODE.open_give
	end, 5000)
	if open == false then
		return ctx:fail("하객 선물 창이 열리지 않음")
	end
	if check(ctx, #open.wishes == 2 and open.wishes[1] == "모래 주머니", "신랑 위시리스트가 다름") == false then
		return false
	end
	local given = guest:request(resp.wedding_gift, req.wedding_present { mode = GIFT_MODE.give, slot = guest:slot(GIFT), item_id = GIFT, count = 2 }, function(p)
		return p.mode == GIFT_MODE.given
	end, 5000)
	if given == false then
		return ctx:fail("결혼 선물 전달 결과가 오지 않음")
	end
	guest:send(req.wedding_present { mode = GIFT_MODE.close })
	ctx:sleep(500)
	if has(ctx, guest, GIFT, 3) == false then
		return false
	end

	dlg = click(ctx, groom, VALENTINA)
	if dlg == false then
		return false
	end
	local warp = groom:request(resp.warp, req.dialog { dialog_type = DIALOG.yes_no, next = true }, function(p)
		return p.character.map == CATHEDRAL
	end, 5000)
	if warp == false then
		return ctx:fail("신랑이 웨딩홀로 이동하지 않음")
	end
	if pq.wait_map(ctx, bride, CATHEDRAL, 5000) == false or pq.wait_map(ctx, guest, CATHEDRAL, 5000) == false then
		return false
	end

	dlg = click(ctx, groom, ANNA)
	if dlg == false then
		return false
	end
	local couple = groom:request_on(guest, resp.wedding_couple, req.dialog { dialog_type = DIALOG.yes_no, next = true }, nil, 150000)
	if couple == false then
		return ctx:fail("키스 연출이 하객에게 오지 않음")
	end
	if check(ctx, couple.groom_id == groom:id() and couple.bride_id == bride:id(), "키스 연출의 신랑/신부가 다름") == false then
		return false
	end
	ctx:sleep(2000)

	dlg = click(ctx, groom, VALENTINA)
	if dlg == false then
		return false
	end
	if said(ctx, groom, dlg, "다정한 한쌍의 달팽이", "피날레 질문") == false then
		return false
	end
	warp = groom:request(resp.warp, req.dialog { dialog_type = DIALOG.yes_no, next = true }, function(p)
		return p.character.map == EXIT_MAP
	end, 5000)
	if warp == false then
		return ctx:fail("조촐한 예식 후 퇴장 맵으로 이동하지 않음")
	end
	if pq.wait_map(ctx, bride, EXIT_MAP, 5000) == false or pq.wait_map(ctx, guest, EXIT_MAP, 5000) == false then
		return false
	end

	for _, bot in ipairs({ groom, bride }) do
		if bot:map_move(VILLAGE) == false then
			return ctx:fail(bot:name() .. " 웨딩빌리지 이동 실패")
		end
		if talk(ctx, bot, MARGARET, "이미 결혼을 하신 커플", "결혼 후 마가렛") == false then
			return false
		end
	end
	local profile = groom:request(resp.character_profile, req.inspect_character { character_id = bride:id() }, nil, 3000)
	if profile == false then
		return ctx:fail("신부 캐릭터 정보 응답 없음")
	end
	if check(ctx, profile.character_id == bride:id() and profile.married and profile.self == false, "신부 캐릭터 정보의 결혼 여부가 다름") == false then
		return false
	end
	return true
end

local function gifts(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if talk(ctx, bride, ANGELIQUE, "받으실 수 있는 선물이 없군요", "선물 없는 신부") == false then
		return false
	end
	local oid = pq.npc(ctx, groom, ANGELIQUE)
	if oid == false then
		return false
	end
	local open = groom:request(resp.wedding_gift, req.npc_click { oid = oid }, function(p)
		return p.mode == GIFT_MODE.open_receive
	end, 5000)
	if open == false then
		return ctx:fail("신랑 선물 보관함이 열리지 않음")
	end
	local received = groom:request(resp.wedding_gift, req.wedding_present { mode = GIFT_MODE.receive, inventory_type = ETC, index = 0 }, function(p)
		return p.mode == GIFT_MODE.received
	end, 5000)
	if received == false then
		return ctx:fail("결혼 선물 수령 결과가 오지 않음")
	end
	groom:send(req.wedding_present { mode = GIFT_MODE.close })
	ctx:sleep(500)
	return has(ctx, groom, GIFT, 2)
end

local function request_divorce(ctx, bot)
	local dlg = click(ctx, bot, PILA)
	if dlg == false then
		return false
	end
	if check(ctx, #dlg.selections == 1, bot:name() .. " 이혼 신청 전 취소 메뉴가 보임") == false then
		bot:dialog(false)
		return false
	end
	if said(ctx, bot, bot:dialog(true, 0), "안타깝군요", "이혼 상담") == false then
		return false
	end
	if bot:dialog(true) == nil then
		return ctx:fail("이혼 유예 질문이 오지 않음")
	end
	if said(ctx, bot, bot:dialog(true), "이혼 신청이 접수되었습니다", "이혼 신청") == false then
		return false
	end
	bot:dialog(false)
	return true
end

local function cancel_divorce(ctx, bot)
	ctx:sleep(1000)
	local dlg = click(ctx, bot, PILA)
	if dlg == false then
		return false
	end
	if check(ctx, #dlg.selections == 2, bot:name() .. " 이혼 신청 취소 메뉴가 없음") == false then
		bot:dialog(false)
		return false
	end
	if said(ctx, bot, bot:dialog(true, 1), "취소하러 오셨군요", "이혼 신청 취소") == false then
		return false
	end
	if said(ctx, bot, bot:dialog(true), "이혼 신청이 취소되었습니다", "이혼 신청 취소 결과") == false then
		return false
	end
	bot:dialog(false)
	ctx:sleep(1000)
	return true
end

local function divorce(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if request_divorce(ctx, groom) == false then
		return false
	end
	if cancel_divorce(ctx, bride) == false then
		return false
	end
	if request_divorce(ctx, groom) == false then
		return false
	end

	local dlg = click(ctx, groom, PILA)
	if dlg == false then
		return false
	end
	if said(ctx, groom, groom:dialog(true, 0), "72시간이 경과하지", "유예 기간 중 이혼") == false then
		return false
	end
	groom:dialog(false)

	if shift_clock(ctx, groom, "/시간가속 3.01:00:00", "시간가속 적용 요청") == false then
		return false
	end
	dlg = click(ctx, groom, PILA)
	if dlg == false then
		return false
	end
	if said(ctx, groom, groom:dialog(true, 0), "지나셨군요", "유예 기간 후 이혼") == false then
		return false
	end
	if groom:dialog(true) == nil then
		return ctx:fail("이혼 확인 질문이 오지 않음")
	end
	local divorced = groom:request(resp.engage_result, req.dialog { dialog_type = DIALOG.yes_no, next = true }, function(p)
		return p.result == ENGAGE.divorced
	end, 5000)
	if divorced == false then
		return ctx:fail("이혼 결과가 오지 않음")
	end
	if shift_clock(ctx, groom, "/현재시간초기화", "현재 시간 보정 초기화") == false then
		return false
	end
	ctx:sleep(1000)

	for _, bot in ipairs({ groom, bride }) do
		if talk(ctx, bot, MARGARET, "먼저 약혼을 해야합니다", "이혼 후 마가렛") == false then
			return false
		end
	end
	return command(ctx, bride, "/성별 0", "성별 변경: 0")
end

test_suite {
	name = "Marriage: 약혼·예약·청첩장·예식·선물·이혼",
	bot_count = 3,

	scenarios = {
		prepare,
		engage,
		reserve,
		ceremony,
		gifts,
		divorce,
	},
}
