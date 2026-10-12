local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local MODE = { send = 3, receive = 5, delete = 6, close = 8 }
local RESULT = { unknown = 20 }

local SHELL = 4000000
local MESO = 100000
local COST = 1000 + 5000
local SETTLE_WAIT_MS = 40000

local function check(ctx, ok, message)
	if ok == false then
		return ctx:fail(message)
	end
	return true
end

local function setup(ctx, bot, items)
	if pq.command(bot, string.format("/봇초기화 30 0 %d %s -", MESO, items or (SHELL .. ":5")), "봇초기화 완료") == false then
		return ctx:fail(bot:name() .. " 봇 초기화 실패")
	end
	return true
end

local function holds(ctx, bot, shell, meso, what)
	local have = bot:items()[SHELL] or 0
	return check(ctx, have == shell and bot:meso() == meso, string.format("%s %s: 조개껍질 %d 메소 %d (기대 %d, %d)", bot:name(), what, have, bot:meso(), shell, meso))
end

local function open(ctx, bot)
	bot:send(req.duey { mode = MODE.close })
	local p = bot:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
	if p == false then
		return ctx:fail(bot:name() .. " 택배 보관함이 열리지 않음")
	end
	return p
end

local function parcels_from(ctx, recipient, sender)
	local p = open(ctx, recipient)
	if p == false then
		return false
	end
	local found = {}
	for _, parcel in ipairs(p.parcels) do
		if parcel.sender == sender:name() then
			found[#found + 1] = parcel
		end
	end
	return found
end

local function clear(ctx, bot)
	local p = open(ctx, bot)
	if p == false then
		return false
	end
	for _, parcel in ipairs(p.parcels) do
		bot:request(resp.duey_removed, req.duey { mode = MODE.delete, parcel_id = parcel.id }, nil, 3000)
	end
	bot:send(req.duey { mode = MODE.close })
	return true
end

local function send_request(sender, recipient)
	return req.duey {
		mode = MODE.send,
		inventory_type = 4,
		slot = sender:slot(SHELL) or 0,
		count = 5,
		meso = 1000,
		recipient = recipient:name(),
		quick = false,
	}
end

-- The fault is server wide, so it is turned off right after the one request it targets.
local function send_unknown(ctx, sender, recipient, mode)
	if open(ctx, sender) == false then
		return false
	end
	if fault.rpc(sender, "SendParcel", mode) == false then
		return ctx:fail("SendParcel 장애 주입 실패")
	end
	local p = sender:request(resp.duey, send_request(sender, recipient), nil, 5000)
	fault.rpc(sender, "SendParcel", "off")
	if p == false then
		return ctx:fail("장애 중 배송 응답 없음")
	end
	if p.result ~= RESULT.unknown then
		return ctx:fail("장애 중 배송 결과 " .. tostring(p.result) .. " (기대 알 수 없음)")
	end
	sender:send(req.duey { mode = MODE.close })
	return holds(ctx, sender, 0, MESO - COST, "장애 직후 보류")
end

local function wait_holds(ctx, bot, shell, meso, what)
	for _ = 1, SETTLE_WAIT_MS / 1000 do
		if (bot:items()[SHELL] or 0) == shell and bot:meso() == meso then
			return true
		end
		ctx:sleep(1000)
	end
	return holds(ctx, bot, shell, meso, what)
end

local function wait_refund(ctx, bot)
	return wait_holds(ctx, bot, 5, MESO, "환불 대기")
end

local function prepare(ctx, sender, recipient)
	if clear(ctx, sender) == false or clear(ctx, recipient) == false then
		return false
	end
	return setup(ctx, sender)
end

local function lost_after_send_timer(ctx)
	local sender = ctx:bot(0)
	local recipient = ctx:bot(1)
	if prepare(ctx, sender, recipient) == false then
		return false
	end
	if send_unknown(ctx, sender, recipient, "lost") == false then
		return false
	end

	ctx:sleep(SETTLE_WAIT_MS)
	if holds(ctx, sender, 0, MESO - COST, "주기 정리 뒤 (이미 배송됨)") == false then
		return false
	end
	local found = parcels_from(ctx, recipient, sender)
	if found == false then
		return false
	end
	if check(ctx, #found == 1, "배송된 택배 수 " .. #found .. " (기대 1)") == false then
		return false
	end
	return clear(ctx, recipient)
end

local function lost_after_send_relogin(ctx)
	local sender = ctx:bot(2)
	local recipient = ctx:bot(3)
	if prepare(ctx, sender, recipient) == false then
		return false
	end
	if send_unknown(ctx, sender, recipient, "lost") == false then
		return false
	end

	if sender:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	ctx:sleep(2000)
	if holds(ctx, sender, 0, MESO - COST, "재접속 정리 뒤 (이미 배송됨)") == false then
		return false
	end
	if sender:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, sender, 0, MESO - COST, "복귀 뒤") == false then
		return false
	end
	local found = parcels_from(ctx, recipient, sender)
	if found == false then
		return false
	end
	if check(ctx, #found == 1, "배송된 택배 수 " .. #found .. " (기대 1)") == false then
		return false
	end
	return clear(ctx, recipient)
end

local function not_sent_refund(ctx)
	local sender = ctx:bot(4)
	local recipient = ctx:bot(5)
	if prepare(ctx, sender, recipient) == false then
		return false
	end
	if send_unknown(ctx, sender, recipient, "on") == false then
		return false
	end

	if wait_refund(ctx, sender) == false then
		return false
	end
	local found = parcels_from(ctx, recipient, sender)
	if found == false then
		return false
	end
	return check(ctx, #found == 0, "보내지 않은 택배가 도착함: " .. #found)
end

local function late_send_after_abort(ctx)
	local sender = ctx:bot(0)
	local recipient = ctx:bot(1)
	if prepare(ctx, sender, recipient) == false then
		return false
	end
	if open(ctx, sender) == false then
		return false
	end
	if fault.rpc_delay(sender, "SendParcel", 6000) == false then
		return ctx:fail("SendParcel 지연 주입 실패")
	end

	sender:send(send_request(sender, recipient))
	local moved = sender:transfer(1)
	fault.rpc_delay(recipient, "SendParcel", 0)
	if moved == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_refund(ctx, sender) == false then
		return false
	end

	ctx:sleep(8000)
	local found = parcels_from(ctx, recipient, sender)
	if found == false then
		return false
	end
	if check(ctx, #found == 0, "중단한 뒤 늦게 도착한 배송이 반영됨: " .. #found) == false then
		return false
	end
	if holds(ctx, sender, 5, MESO, "늦은 응답 뒤") == false then
		return false
	end
	if sender:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	return holds(ctx, sender, 5, MESO, "복귀 뒤 (중복 환불 없음)")
end

local SYSTEM_SENDER = "[복구]"

local function box_empty(ctx, bot, what)
	local p = open(ctx, bot)
	if p == false then
		return false
	end
	bot:send(req.duey { mode = MODE.close })
	return check(ctx, #p.parcels == 0, string.format("%s %s: 보관함에 택배 %d개", bot:name(), what, #p.parcels))
end

-- The reply to the claim is lost or the request never leaves; the parcel stays held by internal or in the box.
local function claim_unknown(ctx, bot, mode)
	if clear(ctx, bot) == false or setup(ctx, bot, "-") == false then
		return false
	end
	if pq.command(bot, string.format("/택배 %s %s %d 5 1000", bot:name(), SYSTEM_SENDER, SHELL), "택배 발송: ") == false then
		return ctx:fail("시스템 택배 발송 실패")
	end
	ctx:sleep(500)
	local p = open(ctx, bot)
	if p == false then
		return false
	end
	local parcel = p.parcels[1]
	if check(ctx, parcel ~= nil and parcel.sender == SYSTEM_SENDER, "시스템 택배가 보관함에 없음") == false then
		return false
	end

	local failed
	if mode == "delay" then
		if fault.rpc_delay(bot, "ClaimParcel", 6000) == false then
			return ctx:fail("ClaimParcel 지연 주입 실패")
		end
		bot:send(req.duey { mode = MODE.receive, parcel_id = parcel.id })
		return true
	end
	if fault.rpc(bot, "ClaimParcel", mode) == false then
		return ctx:fail("ClaimParcel 장애 주입 실패")
	end
	failed = bot:request(resp.duey, req.duey { mode = MODE.receive, parcel_id = parcel.id }, nil, 5000)
	fault.rpc(bot, "ClaimParcel", "off")
	if failed == false then
		return ctx:fail("장애 중 수령 응답 없음")
	end
	if failed.result ~= RESULT.unknown then
		return ctx:fail("장애 중 수령 결과 " .. tostring(failed.result) .. " (기대 알 수 없음)")
	end
	bot:send(req.duey { mode = MODE.close })
	return holds(ctx, bot, 0, MESO, "장애 직후 미지급")
end

local function claim_lost_timer(ctx)
	local bot = ctx:bot(0)
	if claim_unknown(ctx, bot, "lost") == false then
		return false
	end
	if wait_holds(ctx, bot, 5, MESO + 1000, "주기 정리로 수령") == false then
		return false
	end
	ctx:sleep(3000)
	if holds(ctx, bot, 5, MESO + 1000, "수령 뒤 (중복 지급 없음)") == false then
		return false
	end
	return box_empty(ctx, bot, "수령 뒤")
end

local function claim_lost_relogin(ctx)
	local bot = ctx:bot(1)
	if claim_unknown(ctx, bot, "lost") == false then
		return false
	end
	if bot:transfer(1) == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, bot, 5, MESO + 1000, "재접속 정리로 수령") == false then
		return false
	end
	if bot:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, bot, 5, MESO + 1000, "복귀 뒤 (중복 지급 없음)") == false then
		return false
	end
	return box_empty(ctx, bot, "수령 뒤")
end

local function claim_not_sent(ctx)
	local bot = ctx:bot(2)
	if claim_unknown(ctx, bot, "on") == false then
		return false
	end
	if wait_holds(ctx, bot, 5, MESO + 1000, "주기 정리로 다시 수령") == false then
		return false
	end
	return box_empty(ctx, bot, "다시 수령한 뒤")
end

local function late_claim_after_resend(ctx)
	local bot = ctx:bot(3)
	local other = ctx:bot(4)
	if claim_unknown(ctx, bot, "delay") == false then
		return false
	end
	local moved = bot:transfer(1)
	fault.rpc_delay(other, "ClaimParcel", 0)
	if moved == false then
		return ctx:fail("채널 이동 실패")
	end
	if wait_holds(ctx, bot, 5, MESO + 1000, "재접속 재요청으로 수령") == false then
		return false
	end

	ctx:sleep(8000)
	if holds(ctx, bot, 5, MESO + 1000, "늦은 원래 요청 뒤 (중복 지급 없음)") == false then
		return false
	end
	if bot:transfer(0) == false then
		return ctx:fail("원래 채널 복귀 실패")
	end
	ctx:sleep(1000)
	if holds(ctx, bot, 5, MESO + 1000, "복귀 뒤") == false then
		return false
	end
	return box_empty(ctx, bot, "수령 뒤")
end

test_suite {
	name = "Duey: 택배 발송·수령 응답 유실 복구",
	bot_count = 6,

	scenarios = {
		lost_after_send_timer,
		lost_after_send_relogin,
		not_sent_refund,
		late_send_after_abort,
		claim_lost_timer,
		claim_lost_relogin,
		claim_not_sent,
		late_claim_after_resend,
	},
}
