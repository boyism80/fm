local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local VILLAGE = 680000000
local MARGARET = 9201013

local RING_BOX = 2240004
local ENGAGEMENT_RING = 4210000
local TICKET = 5251004

local DIALOG = { normal = 0 }
local RING = { propose = 0, answer = 2, wishlist = 9 }
local ENGAGE = { engaged = 11 }
local SETTLE_WAIT_MS = 40000

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function count(bot, item_id)
	return bot:items()[item_id] or 0
end

local function holds(ctx, bot, item_id, expected, what)
	local have = count(bot, item_id)
	return check(ctx, have == expected, string.format("%s %s: 아이템 %d 개수 %d (기대 %d)", bot:name(), what, item_id, have, expected))
end

local function wait_holds(ctx, bot, item_id, expected, what)
	for _ = 1, SETTLE_WAIT_MS / 1000 do
		if count(bot, item_id) == expected then
			return true
		end
		ctx:sleep(1000)
	end
	return holds(ctx, bot, item_id, expected, what)
end

local function command(ctx, bot, text, expected)
	if pq.command(bot, text, expected) == false then
		return ctx:fail(bot:name() .. " 명령 실패: " .. text)
	end
	return true
end

-- Returns whether the bot was engaged; /파혼 also clears the engagement for the next scenario.
local function break_off(ctx, bot)
	local p = bot:request(resp.notice, req.normal_chat { message = "/파혼" }, function(p)
		return p.message:find("파혼 완료", 1, true) ~= nil or p.message:find("약혼 상태가 아닙니다", 1, true) ~= nil
	end)
	if p == false then
		ctx:fail(bot:name() .. " 파혼 명령 응답 없음")
		return nil
	end
	return p.message:find("파혼 완료", 1, true) ~= nil
end

local function engaged(ctx, bot, expected, what)
	local was = break_off(ctx, bot)
	if was == nil then
		return false
	end
	return check(ctx, was == expected, string.format("%s %s: 약혼 상태 %s (기대 %s)", bot:name(), what, tostring(was), tostring(expected)))
end

local function setup(ctx, bot, items)
	if pq.command(bot, string.format("/봇초기화 30 0 1000000 %s -", items), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	if bot:map() ~= VILLAGE and bot:map_move(VILLAGE) == false then
		return ctx:fail(bot:name() .. " 웨딩빌리지 이동 실패")
	end
	return true
end

local function prepare(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if command(ctx, bride, "/성별 1", "성별 변경: 1") == false or command(ctx, groom, "/성별 0", "성별 변경: 0") == false then
		return false
	end
	break_off(ctx, groom)
	ctx:sleep(1000)
	break_off(ctx, bride)
	if setup(ctx, groom, string.format("%d:1,%d:1", RING_BOX, TICKET)) == false then
		return false
	end
	return setup(ctx, bride, "-")
end

local function answer(groom)
	return req.ring_action { mode = RING.answer, accepted = true, name = groom:name(), character_id = groom:id() }
end

local function propose(ctx, groom, bride)
	local x, y = bride:position()
	if pq.move(groom, x, y) == false then
		return ctx:fail("신부 앞으로 이동 실패")
	end
	local proposal = groom:request_on(bride, resp.engage_request, req.ring_action { mode = RING.propose, name = bride:name(), item_id = RING_BOX }, nil, 5000)
	return check(ctx, proposal ~= false, "청혼 요청이 신부에게 오지 않음")
end

local function engage(ctx, groom, bride)
	if propose(ctx, groom, bride) == false then
		return false
	end
	local p = bride:request_on(groom, resp.engage_result, answer(groom), function(p)
		return p.result == ENGAGE.engaged
	end, 5000)
	return check(ctx, p ~= false, "약혼 결과가 오지 않음")
end

-- The fault is server wide, so it is turned off right after the one answer it targets.
local function proposal_unknown(ctx, mode)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if prepare(ctx) == false or propose(ctx, groom, bride) == false then
		return false
	end
	if fault.rpc(bride, "CreateMarriage", mode) == false then
		return ctx:fail("CreateMarriage 장애 주입 실패")
	end
	bride:send(answer(groom))
	ctx:sleep(1000)
	fault.rpc(bride, "CreateMarriage", "off")
	return holds(ctx, groom, RING_BOX, 0, "장애 직후 상자 보류")
end

local function rings(ctx, groom, bride, expected, what)
	if wait_holds(ctx, groom, ENGAGEMENT_RING, expected, what) == false then
		return false
	end
	return wait_holds(ctx, bride, ENGAGEMENT_RING, expected, what)
end

local function proposal_lost_timer(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if proposal_unknown(ctx, "lost") == false then
		return false
	end
	if rings(ctx, groom, bride, 1, "주기 정리로 약혼반지") == false then
		return false
	end
	if holds(ctx, groom, RING_BOX, 0, "약혼 뒤 상자") == false then
		return false
	end
	return engaged(ctx, groom, true, "주기 정리 뒤")
end

local function proposal_lost_relogin(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if proposal_unknown(ctx, "lost") == false then
		return false
	end
	if groom:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, groom, ENGAGEMENT_RING, 1, "재접속 정리로 약혼반지") == false then
		return false
	end
	if groom:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if rings(ctx, groom, bride, 1, "복귀 뒤 (중복 지급 없음)") == false then
		return false
	end
	if holds(ctx, groom, RING_BOX, 0, "복귀 뒤 상자") == false then
		return false
	end
	return engaged(ctx, groom, true, "재접속 정리 뒤")
end

local function proposal_not_sent(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if proposal_unknown(ctx, "on") == false then
		return false
	end
	if wait_holds(ctx, groom, RING_BOX, 1, "주기 정리로 상자 환불") == false then
		return false
	end
	ctx:sleep(3000)
	if holds(ctx, groom, ENGAGEMENT_RING, 0, "환불 뒤") == false or holds(ctx, bride, ENGAGEMENT_RING, 0, "환불 뒤") == false then
		return false
	end
	return engaged(ctx, groom, false, "환불 뒤")
end

local function late_proposal_after_abort(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if prepare(ctx) == false or propose(ctx, groom, bride) == false then
		return false
	end
	if fault.rpc_delay(bride, "CreateMarriage", 6000) == false then
		return ctx:fail("CreateMarriage 지연 주입 실패")
	end
	bride:send(answer(groom))
	ctx:sleep(500)
	local moved = groom:transfer(1)
	fault.rpc_delay(bride, "CreateMarriage", 0)
	if moved == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, groom, RING_BOX, 1, "재접속 중단으로 상자 환불") == false then
		return false
	end

	ctx:sleep(8000)
	if holds(ctx, bride, ENGAGEMENT_RING, 0, "늦은 요청 뒤") == false or holds(ctx, groom, ENGAGEMENT_RING, 0, "늦은 요청 뒤") == false then
		return false
	end
	if groom:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, groom, RING_BOX, 1, "복귀 뒤 (중복 환불 없음)") == false then
		return false
	end
	return engaged(ctx, bride, false, "중단한 뒤")
end

local function open_reservation(ctx, groom)
	local oid = pq.npc(ctx, groom, MARGARET)
	if oid == false then
		return false
	end
	local dlg = groom:npc_click(oid)
	if dlg == false or dlg == nil or dlg.selections == nil then
		groom:dialog(false)
		return ctx:fail("마가렛 예약 선택지가 오지 않음")
	end
	dlg = groom:dialog(true, 3)
	if dlg == nil or dlg.text == nil or dlg.text:find("약혼했군요", 1, true) == nil then
		groom:dialog(false)
		return ctx:fail("마가렛 예약 확인 대화가 다름")
	end
	return true
end

-- Reserved couples see the wishlist step at Margaret instead of the reservation list.
local function reserved(ctx, groom, expected, what)
	local oid = pq.npc(ctx, groom, MARGARET)
	if oid == false then
		return false
	end
	local dlg = groom:npc_click(oid)
	groom:dialog(false)
	if dlg == false or dlg == nil then
		return ctx:fail("마가렛 대화 응답 없음")
	end
	local was = dlg.text ~= nil and dlg.text:find("결혼식 예약접수를 완료하기 전에", 1, true) ~= nil
	return check(ctx, was == expected, string.format("%s: 예약 상태 %s (기대 %s)", what, tostring(was), tostring(expected)))
end

local function reservation_unknown(ctx, mode)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if prepare(ctx) == false or engage(ctx, groom, bride) == false then
		return false
	end
	if open_reservation(ctx, groom) == false then
		return false
	end
	if mode == "delay" then
		if fault.rpc_delay(groom, "ReserveWedding", 6000) == false then
			return ctx:fail("ReserveWedding 지연 주입 실패")
		end
		groom:send(req.dialog { dialog_type = DIALOG.normal, next = true })
		ctx:sleep(500)
		return true
	end

	if fault.rpc(groom, "ReserveWedding", mode) == false then
		return ctx:fail("ReserveWedding 장애 주입 실패")
	end
	local failed = groom:request(resp.dialog, req.dialog { dialog_type = DIALOG.normal, next = true }, function(p)
		return p.text ~= nil and p.text:find("문제가", 1, true) ~= nil
	end, 15000)
	fault.rpc(groom, "ReserveWedding", "off")
	groom:dialog(false)
	if failed == false then
		return ctx:fail("장애 중 예약 실패 대화가 오지 않음")
	end
	return holds(ctx, groom, TICKET, 0, "장애 직후 티켓 보류")
end

local function reservation_lost_timer(ctx)
	local groom = ctx:bot(0)
	if reservation_unknown(ctx, "lost") == false then
		return false
	end
	ctx:sleep(SETTLE_WAIT_MS)
	if holds(ctx, groom, TICKET, 0, "주기 정리 뒤 (이미 예약됨)") == false then
		return false
	end
	return reserved(ctx, groom, true, "주기 정리 뒤")
end

local function reservation_lost_relogin(ctx)
	local groom = ctx:bot(0)
	if reservation_unknown(ctx, "lost") == false then
		return false
	end
	if groom:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	ctx:sleep(2000)
	if holds(ctx, groom, TICKET, 0, "재접속 정리 뒤 (이미 예약됨)") == false then
		return false
	end
	if groom:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, groom, TICKET, 0, "복귀 뒤") == false then
		return false
	end
	return reserved(ctx, groom, true, "재접속 정리 뒤")
end

local function reservation_not_sent(ctx)
	local groom = ctx:bot(0)
	if reservation_unknown(ctx, "on") == false then
		return false
	end
	if wait_holds(ctx, groom, TICKET, 1, "주기 정리로 티켓 환불") == false then
		return false
	end
	return reserved(ctx, groom, false, "환불 뒤")
end

local function late_reservation_after_abort(ctx)
	local groom, bride = ctx:bot(0), ctx:bot(1)
	if reservation_unknown(ctx, "delay") == false then
		return false
	end
	local moved = groom:transfer(1)
	fault.rpc_delay(bride, "ReserveWedding", 0)
	if moved == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, groom, TICKET, 1, "재접속 중단으로 티켓 환불") == false then
		return false
	end

	ctx:sleep(8000)
	if groom:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, groom, TICKET, 1, "복귀 뒤 (중복 환불 없음)") == false then
		return false
	end
	return reserved(ctx, groom, false, "중단한 뒤")
end

local function cleanup(ctx)
	break_off(ctx, ctx:bot(0))
	ctx:sleep(1000)
	break_off(ctx, ctx:bot(1))
	return true
end

test_suite {
	name = "Marriage: 약혼·예약 응답 유실 복구",
	bot_count = 2,

	scenarios = {
		proposal_lost_timer,
		proposal_lost_relogin,
		proposal_not_sent,
		late_proposal_after_abort,
		reservation_lost_timer,
		reservation_lost_relogin,
		reservation_not_sent,
		late_reservation_after_abort,
		cleanup,
	},
}
