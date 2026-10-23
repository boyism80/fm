local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")
local settle = require("script/integration/lib/settle")

local HENESYS = 100000000
local DUEY = 9010009

local MODE = { from_arrival = 0, identity = 1, send = 3, receive = 5, delete = 6, close = 8 }
local RESULT = {
    not_enough_meso = 12,
    invalid = 13,
    recipient_not_found = 14,
    same_account = 15,
    recipient_full = 16,
    cannot_receive = 17,
    only_in_box = 18,
    sent = 19,
    unknown = 20,
    inventory_full = 22,
    only_held = 23,
    quick_open = 27,
}
local REMOVED = { deleted = 3, received = 4 }
local MAX_PARCELS = 50

local EQUIP = 1
local CONSUME = 2
local ETC = 4

local POTION = 2000000
local SWORD = 1302000
local SHELL = 4000000
local ONLY = 4001025
local BLOCK = 4001078
local PET = 5000000
local COUPON = 5330000

local MESO = 100000
local COST = 1000 + 5000
local SETTLE_WAIT_MS = 40000
local SYSTEM_SENDER = "[복구]"
local REFUND_SENDER = "택배"

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function setup(ctx, bot, meso, items)
    local text = string.format("/봇초기화 30 0 %d %s -", meso, items)
    if pq.command(bot, text, "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
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

local function open(ctx, bot)
    bot:send(req.duey { mode = MODE.close })
    local p = bot:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
    if p == false then
        return ctx:fail(bot:name() .. " 택배 보관함이 열리지 않음")
    end
    return p
end

local function call(bot, request)
    return bot:request({ resp.duey, resp.duey_removed }, request, nil, 3000)
end

local function expect(ctx, bot, request, result, what)
    local p = call(bot, request)
    if p == false then
        return ctx:fail(bot:name() .. " 응답 없음: " .. what)
    end
    if p.result ~= result then
        return ctx:fail(string.format("%s %s: 결과 %s (기대 %d)", bot:name(), what, tostring(p.result), result))
    end
    return p
end

local function ignored(ctx, bot, request, what)
    local p = bot:request({ resp.duey, resp.duey_removed }, request, nil, 1500)
    if p ~= false then
        return ctx:fail(string.format("%s 거절되어야 할 요청에 택배 응답: %s (%s)", bot:name(), what, tostring(p.result)))
    end
    return true
end

local function removed(ctx, bot, request, reason, what)
    local p = bot:request({ resp.duey, resp.duey_removed }, request, nil, 3000)
    if p == false then
        return ctx:fail(bot:name() .. " 응답 없음: " .. what)
    end
    if p.reason ~= reason then
        return ctx:fail(string.format("%s %s: 택배 제거 응답이 아님 (결과 %s, 사유 %s)", bot:name(), what, tostring(p.result),
            tostring(p.reason)))
    end
    return p
end

local function clear(ctx, bot)
    local p = open(ctx, bot)
    if p == false then
        return false
    end
    for _, parcel in ipairs(p.parcels) do
        if removed(ctx, bot, req.duey { mode = MODE.delete, parcel_id = parcel.id }, REMOVED.deleted, "정리") == false then
            return false
        end
    end
    bot:send(req.duey { mode = MODE.close })
    return true
end

local function send_fields(bot, item_id, count, meso, recipient)
    local slot = 0
    local inventory_type = 0
    if item_id ~= nil then
        slot = bot:slot(item_id) or 0
        inventory_type = math.floor(item_id / 1000000)
    end
    return {
        mode = MODE.send,
        inventory_type = inventory_type,
        slot = slot,
        count = count or 0,
        meso = meso or 0,
        recipient = recipient,
        quick = false,
    }
end

local function send_request(bot, item_id, count, meso, recipient)
    return req.duey(send_fields(bot, item_id, count, meso, recipient))
end

local function quick_request(bot, item_id, count, meso, recipient, message)
    local fields = send_fields(bot, item_id, count, meso, recipient)
    fields.quick = true
    fields.message = message or ""
    fields.coupon = 1
    return req.duey(fields)
end

local function find(p, sender)
    local found = {}
    for _, parcel in ipairs(p.parcels) do
        if sender == nil or parcel.sender == sender then
            found[#found + 1] = parcel
        end
    end
    return found
end

local function system_parcel(ctx, bot, recipient, sender, item_id, count, meso)
    local text = string.format("/택배 %s %s %s %d %d", recipient, sender, item_id and tostring(item_id) or "-", count or 1,
        meso or 0)
    if pq.command(bot, text, "택배 발송: ") == false then
        return ctx:fail(bot:name() .. " 시스템 택배 발송 실패: " .. text)
    end
    ctx:sleep(300)
    return true
end

local function normal_flow(ctx)
    local sender = ctx:bot(0)
    local recipient = ctx:bot(1)
    if clear(ctx, sender) == false or clear(ctx, recipient) == false then
        return false
    end
    if setup(ctx, sender, 200000, SHELL .. ":10," .. SWORD .. ":1") == false then
        return false
    end

    if sender:instance_move(HENESYS) == false then
        return ctx:fail("헤네시스 이동 실패")
    end
    local oid = pq.npc(ctx, sender, DUEY)
    if oid == false then
        return false
    end
    local p = sender:request(resp.duey_open, req.npc_click { oid = oid }, nil, 5000)
    if p == false then
        return ctx:fail("듀이 NPC로 보관함이 열리지 않음")
    end
    if check(ctx, p.from_arrival == false and #p.parcels == 0, "빈 보관함이 아니거나 알림 열기로 표시됨") == false then
        return false
    end

    if expect(ctx, sender, send_request(sender, SHELL, 5, 1000, recipient:name()), RESULT.sent, "일반 배송") == false then
        return false
    end
    if has(ctx, sender, SHELL, 5) == false or meso_is(ctx, sender, 200000 - 1000 - 5000) == false then
        return false
    end

    if expect(ctx, sender, send_request(sender, nil, 0, 100000, recipient:name()), RESULT.sent, "메소만 배송 (세금 1%)") == false then
        return false
    end
    if meso_is(ctx, sender, 194000 - 100000 - 6000) == false then
        return false
    end

    if expect(ctx, sender, send_request(sender, SWORD, 1, 0, recipient:name()), RESULT.sent, "장비 배송") == false then
        return false
    end
    if has(ctx, sender, SWORD, 0) == false or meso_is(ctx, sender, 88000 - 5000) == false then
        return false
    end

    p = open(ctx, recipient)
    if p == false then
        return false
    end
    local parcels = find(p, sender:name())
    if check(ctx, #parcels == 3, "받은 택배 수가 3이 아님: " .. #parcels) == false then
        return false
    end
    local shell = parcels[1]
    if check(ctx, shell.meso == 1000 and shell.quick == false and shell.item ~= nil and shell.item.item_id == SHELL and shell.item.count == 5, "첫 택배 내용이 다름") == false then
        return false
    end
    if check(ctx, parcels[2].meso == 100000 and parcels[2].item == nil, "메소 택배 내용이 다름") == false then
        return false
    end
    if check(ctx, parcels[3].item ~= nil and parcels[3].item.item_id == SWORD, "장비 택배 내용이 다름") == false then
        return false
    end

    local meso = recipient:meso()
    if expect(ctx, recipient, req.duey { mode = MODE.receive, parcel_id = shell.id }, RESULT.cannot_receive, "배송 중 수령") == false then
        return false
    end
    if has(ctx, recipient, SHELL, 0) == false or meso_is(ctx, recipient, meso) == false then
        return false
    end

    for _, parcel in ipairs(parcels) do
        if removed(ctx, recipient, req.duey { mode = MODE.delete, parcel_id = parcel.id }, REMOVED.deleted, "택배 삭제") == false then
            return false
        end
    end
    p = open(ctx, recipient)
    if p == false then
        return false
    end
    return check(ctx, #p.parcels == 0, "삭제한 택배가 남음: " .. #p.parcels)
end

local function quick_flow(ctx)
    local sender = ctx:bot(2)
    local recipient = ctx:bot(3)
    if clear(ctx, sender) == false or clear(ctx, recipient) == false then
        return false
    end
    if setup(ctx, sender, 100000, COUPON .. ":2," .. SWORD .. ":1," .. POTION .. ":10") == false then
        return false
    end
    if setup(ctx, recipient, 1000, "-") == false then
        return false
    end

    if expect(ctx, sender, req.normal_chat { message = "/퀵배송" }, RESULT.quick_open, "퀵배송 창 열기") == false then
        return false
    end
    if pq.command(sender, "/퀵배송", "퀵배송 창을 열 수 없습니다.") == false then
        return ctx:fail("이미 열린 퀵배송 창이 다시 열림")
    end

    local arrival = sender:request_on(recipient, resp.duey_arrival,
        quick_request(sender, SWORD, 1, 2000, recipient:name(), "선물입니다"), nil, 5000)
    if arrival == false then
        return ctx:fail("퀵배송 도착 알림이 오지 않음")
    end
    if check(ctx, arrival.sender == sender:name() and arrival.quick == true, "도착 알림 내용이 다름: " .. tostring(arrival.sender)) == false then
        return false
    end
    ctx:sleep(500)
    if has(ctx, sender, SWORD, 0) == false or has(ctx, sender, COUPON, 1) == false or meso_is(ctx, sender, 100000 - 2000) == false then
        return false
    end
    sender:send(req.duey { mode = MODE.close })

    local p = recipient:request(resp.duey_open, req.duey { mode = MODE.from_arrival, value = -1, kind = 2 }, nil, 5000)
    if p == false then
        return ctx:fail("도착 알림으로 보관함이 열리지 않음")
    end
    local parcels = find(p, sender:name())
    if check(ctx, p.from_arrival == true and #parcels == 1, "알림으로 연 보관함 상태가 다름") == false then
        return false
    end
    local parcel = parcels[1]
    if check(ctx, parcel.quick == true and parcel.message == "선물입니다" and parcel.meso == 2000, "퀵배송 택배 내용이 다름") == false then
        return false
    end
    if removed(ctx, recipient, req.duey { mode = MODE.receive, parcel_id = parcel.id }, REMOVED.received, "퀵배송 수령") == false then
        return false
    end
    if has(ctx, recipient, SWORD, 1) == false or meso_is(ctx, recipient, 3000) == false then
        return false
    end
    if expect(ctx, recipient, req.duey { mode = MODE.receive, parcel_id = parcel.id }, RESULT.invalid, "이미 받은 택배 다시 수령") == false then
        return false
    end
    if has(ctx, recipient, SWORD, 1) == false then
        return false
    end

    if recipient:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    if open(ctx, sender) == false then
        return false
    end
    arrival = sender:request_on(recipient, resp.duey_arrival, quick_request(sender, POTION, 3, 0, recipient:name(), ""),
        nil, 5000)
    if arrival == false then
        return ctx:fail("다른 채널로 퀵배송 도착 알림이 오지 않음")
    end
    ctx:sleep(500)
    if has(ctx, sender, COUPON, 0) == false or has(ctx, sender, POTION, 7) == false then
        return false
    end
    if expect(ctx, sender, quick_request(sender, POTION, 1, 0, recipient:name(), ""), RESULT.invalid, "이용권 없이 퀵배송") == false then
        return false
    end
    if has(ctx, sender, POTION, 7) == false then
        return false
    end

    p = open(ctx, recipient)
    if p == false then
        return false
    end
    parcels = find(p, sender:name())
    if check(ctx, #parcels == 1 and parcels[1].item.item_id == POTION, "다른 채널에서 받은 택배가 없음") == false then
        return false
    end
    if removed(ctx, recipient, req.duey { mode = MODE.receive, parcel_id = parcels[1].id }, REMOVED.received, "다른 채널 수령") == false then
        return false
    end
    if has(ctx, recipient, POTION, 3) == false then
        return false
    end
    if recipient:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

local function claim_refunds(ctx, bot, wait_ms)
    local steps = math.max(1, math.floor((wait_ms or 10000) / 500))
    for _ = 1, steps do
        local p = open(ctx, bot)
        if p == false then
            return false
        end
        local found = find(p, REFUND_SENDER)
        if #found > 0 then
            for _, parcel in ipairs(found) do
                if removed(ctx, bot, req.duey { mode = MODE.receive, parcel_id = parcel.id }, REMOVED.received, "환불 택배 수령") == false then
                    return false
                end
            end
            return true
        end
        bot:send(req.duey { mode = MODE.close })
        ctx:sleep(500)
    end
    return ctx:fail(bot:name() .. " 환불 택배가 오지 않음")
end

local function system_parcels(ctx)
    local bot = ctx:bot(4)
    local gm = ctx:bot(0)
    if clear(ctx, bot) == false then
        return false
    end
    if setup(ctx, bot, 1000, ONLY .. ":1") == false then
        return false
    end

    local arrival = gm:request_on(bot, resp.duey_arrival,
        req.normal_chat { message = string.format("/택배 %s [이벤트] %d 5 500", bot:name(), SHELL) }, nil, 5000)
    if arrival == false then
        return ctx:fail("시스템 택배 도착 알림이 오지 않음")
    end
    if check(ctx, arrival.sender == "[이벤트]", "시스템 택배 보낸 이름이 다름: " .. tostring(arrival.sender)) == false then
        return false
    end
    if system_parcel(ctx, gm, bot:name(), "[이벤트]", ONLY, 1, 0) == false then
        return false
    end

    local p = open(ctx, bot)
    if p == false then
        return false
    end
    local parcels = find(p, "[이벤트]")
    if check(ctx, #parcels == 2, "시스템 택배 수가 2가 아님: " .. #parcels) == false then
        return false
    end
    local shell, only = parcels[1], parcels[2]
    if check(ctx, shell.item.item_id == SHELL and shell.item.count == 5 and shell.meso == 500, "기타 택배 내용이 다름") == false then
        return false
    end

    if expect(ctx, bot, req.duey { mode = MODE.receive, parcel_id = only.id }, RESULT.only_held, "가진 고유 아이템 수령") == false then
        return false
    end
    if expect(ctx, gm, send_request(gm, nil, 0, 0, bot:name()), RESULT.invalid, "빈 택배") == false then
        return false
    end

    if setup(ctx, gm, 100000, ONLY .. ":1") == false then
        return false
    end
    if open(ctx, gm) == false then
        return false
    end
    if expect(ctx, gm, send_request(gm, ONLY, 1, 0, bot:name()), RESULT.only_in_box, "보관함에 있는 고유 아이템 배송") == false then
        return false
    end
    if claim_refunds(ctx, gm) == false then
        return false
    end
    if has(ctx, gm, ONLY, 1) == false or meso_is(ctx, gm, 100000) == false then
        return false
    end

    local items = {}
    for i = 1, 32 do
        items[#items + 1] = (SHELL + i) .. ":1"
    end
    if pq.command(bot, "/봇초기화 30 0 1000 " .. table.concat(items, ",") .. " -", "봇초기화 완료") == false then
        return ctx:fail("기타 탭 채우기 실패")
    end
    p = open(ctx, bot)
    if p == false then
        return false
    end
    if expect(ctx, bot, req.duey { mode = MODE.receive, parcel_id = shell.id }, RESULT.inventory_full, "가득 찬 인벤토리로 수령") == false then
        return false
    end
    if pq.command(bot, "/봇초기화 30 0 2147483600 - -", "봇초기화 완료") == false then
        return ctx:fail("메소 설정 실패")
    end
    p = open(ctx, bot)
    if p == false then
        return false
    end
    if expect(ctx, bot, req.duey { mode = MODE.receive, parcel_id = shell.id }, RESULT.cannot_receive, "메소 한도 초과 수령") == false then
        return false
    end

    if pq.command(bot, "/봇초기화 30 0 1000 - -", "봇초기화 완료") == false then
        return ctx:fail("인벤토리 비우기 실패")
    end
    p = open(ctx, bot)
    if p == false then
        return false
    end
    bot:send(req.duey { mode = MODE.receive, parcel_id = shell.id })
    bot:send(req.duey { mode = MODE.receive, parcel_id = shell.id })
    if removed(ctx, bot, nil, REMOVED.received, "시스템 택배 수령") == false then
        return false
    end
    ctx:sleep(1000)
    if has(ctx, bot, SHELL, 5) == false or meso_is(ctx, bot, 1500) == false then
        return false
    end
    if removed(ctx, bot, req.duey { mode = MODE.receive, parcel_id = only.id }, REMOVED.received, "고유 아이템 수령") == false then
        return false
    end
    return has(ctx, bot, ONLY, 1)
end

local function box_limit(ctx)
    local bot = ctx:bot(5)
    local gm = ctx:bot(0)
    if clear(ctx, bot) == false then
        return false
    end
    for i = 1, MAX_PARCELS do
        if system_parcel(ctx, gm, bot:name(), "[이벤트]", nil, 0, i) == false then
            return false
        end
    end
    if setup(ctx, gm, 100000, SHELL .. ":1") == false then
        return false
    end
    if open(ctx, gm) == false then
        return false
    end
    if expect(ctx, gm, send_request(gm, SHELL, 1, 0, bot:name()), RESULT.recipient_full, "가득 찬 보관함으로 배송") == false then
        return false
    end
    if claim_refunds(ctx, gm) == false then
        return false
    end
    if has(ctx, gm, SHELL, 1) == false or meso_is(ctx, gm, 100000) == false then
        return false
    end

    local p = open(ctx, bot)
    if p == false then
        return false
    end
    if check(ctx, #p.parcels == MAX_PARCELS, "보관함 택배 수가 다름: " .. #p.parcels) == false then
        return false
    end
    return clear(ctx, bot)
end

local function forged_requests(ctx)
    local bot = ctx:bot(1)
    local other = ctx:bot(0)
    local items = string.format("%d:10,%d:1,%d:1,%d:1", POTION, BLOCK, PET, COUPON)
    if setup(ctx, bot, 10000, items) == false then
        return false
    end
    bot:send(req.duey { mode = MODE.close })
    if ignored(ctx, bot, send_request(bot, POTION, 1, 0, other:name()), "보관함을 열기 전 배송") == false then
        return false
    end
    if open(ctx, bot) == false then
        return false
    end

    local slot = bot:slot(POTION)
    local long = string.rep("가", 51)
    local cases = {
        { request = send_request(bot, POTION, 1, 0, "없는캐릭터이름"), result = RESULT.recipient_not_found, what = "없는 받는 사람" },
        { request = send_request(bot, POTION, 1, 0, bot:name()), result = RESULT.same_account, what = "자기 자신에게" },
        { request = send_request(bot, POTION, 11, 0, other:name()), result = RESULT.invalid, what = "가진 것보다 많은 개수" },
        { request = send_request(bot, POTION, 0, 0, other:name()), result = RESULT.invalid, what = "0개 배송" },
        { request = send_request(bot, POTION, -1, 0, other:name()), result = RESULT.invalid, what = "음수 개수" },
        { request = req.duey { mode = MODE.send, inventory_type = CONSUME, slot = 99, count = 1, recipient = other:name() }, result = RESULT.invalid, what = "범위 밖 칸" },
        { request = req.duey { mode = MODE.send, inventory_type = EQUIP, slot = slot, count = 1, recipient = other:name() }, result = RESULT.invalid, what = "다른 탭 칸" },
        { request = send_request(bot, nil, 0, -1, other:name()), result = RESULT.invalid, what = "음수 메소" },
        { request = send_request(bot, nil, 0, 6000, other:name()), result = RESULT.not_enough_meso, what = "수수료 포함 메소 부족" },
        { request = send_request(bot, nil, 0, 2147483647, other:name()), result = RESULT.not_enough_meso, what = "메소 최댓값" },
        { request = send_request(bot, BLOCK, 1, 0, other:name()), result = RESULT.invalid, what = "교환 불가 아이템" },
        { request = send_request(bot, PET, 1, 0, other:name()), result = RESULT.invalid, what = "펫" },
        { request = quick_request(bot, POTION, 1, 0, other:name(), long), result = RESULT.invalid, what = "100바이트 넘는 메시지" },
        { request = req.duey { mode = MODE.send, inventory_type = CONSUME, slot = slot, count = 1, recipient = other:name(), quick = false }, result = RESULT.sent, what = "기준 배송" },
    }
    for _, case in ipairs(cases) do
        if expect(ctx, bot, case.request, case.result, case.what) == false then
            return false
        end
        if case.result == RESULT.recipient_not_found or case.result == RESULT.same_account or case.result == RESULT.recipient_full or case.result == RESULT.only_in_box then
            if claim_refunds(ctx, bot) == false then
                return false
            end
        end
    end
    if has(ctx, bot, POTION, 9) == false or has(ctx, bot, BLOCK, 1) == false or has(ctx, bot, PET, 1) == false then
        return false
    end
    if has(ctx, bot, COUPON, 1) == false or meso_is(ctx, bot, 5000) == false then
        return false
    end

    local p = open(ctx, other)
    if p == false then
        return false
    end
    local parcels = find(p, bot:name())
    if check(ctx, #parcels == 1, "기준 배송 택배가 하나가 아님: " .. #parcels) == false then
        return false
    end
    local foreign = parcels[1].id

    p = open(ctx, bot)
    if p == false then
        return false
    end
    local deny = {
        { request = req.duey { mode = MODE.receive, parcel_id = foreign }, what = "남의 택배 수령" },
        { request = req.duey { mode = MODE.delete, parcel_id = foreign }, what = "남의 택배 삭제" },
        { request = req.duey { mode = MODE.receive, parcel_id = 0 }, what = "없는 택배 수령" },
    }
    for _, case in ipairs(deny) do
        if expect(ctx, bot, case.request, RESULT.invalid, case.what) == false then
            return false
        end
    end
    if ignored(ctx, bot, req.duey { mode = MODE.identity, value = 1234567, kind = 1 }, "요청하지 않은 주민번호 응답") == false then
        return false
    end

    bot:send(req.duey { mode = MODE.close })
    if ignored(ctx, bot, req.duey { mode = MODE.delete, parcel_id = foreign }, "닫은 뒤 삭제") == false then
        return false
    end
    if ignored(ctx, bot, send_request(bot, POTION, 1, 0, other:name()), "닫은 뒤 배송") == false then
        return false
    end
    if has(ctx, bot, POTION, 9) == false then
        return false
    end
    return clear(ctx, other)
end

local function holds(ctx, bot, shell, meso, what)
    local have = bot:items()[SHELL] or 0
    return check(ctx, have == shell and bot:meso() == meso,
        string.format("%s %s: 조개껍질 %d 메소 %d (기대 %d, %d)", bot:name(), what, have, bot:meso(), shell, meso))
end

local function wait_holds(ctx, bot, shell, meso, what)
    if settle.wait(ctx, function()
            return (bot:items()[SHELL] or 0) == shell and bot:meso() == meso
        end) then
        return true
    end
    return holds(ctx, bot, shell, meso, what)
end

local function wait_refund(ctx, bot)
    if claim_refunds(ctx, bot, SETTLE_WAIT_MS) == false then
        return false
    end
    return holds(ctx, bot, 5, MESO, "환불")
end

local function parcels_from(ctx, recipient, sender)
    local p = open(ctx, recipient)
    if p == false then
        return false
    end
    return find(p, sender:name())
end

local function box_empty(ctx, bot, what)
    local p = open(ctx, bot)
    if p == false then
        return false
    end
    bot:send(req.duey { mode = MODE.close })
    return check(ctx, #p.parcels == 0, string.format("%s %s: 보관함에 택배 %d개", bot:name(), what, #p.parcels))
end

local function prepare(ctx, sender, recipient)
    if clear(ctx, sender) == false or clear(ctx, recipient) == false then
        return false
    end
    return setup(ctx, sender, MESO, SHELL .. ":5")
end

local function recovery_request(sender, recipient)
    return send_request(sender, SHELL, 5, 1000, recipient:name())
end

-- The fault is server wide, so it is turned off right after the one request it targets.
-- The settle timer stays paused so the held state can be seen; the caller resumes it.
local function send_unknown(ctx, sender, recipient, mode)
    if open(ctx, sender) == false then
        return false
    end
    if settle.pause(sender) == false then
        return ctx:fail("정산 타이머 정지 실패")
    end
    if fault.rpc(sender, "SendParcel", mode) == false then
        return ctx:fail("SendParcel 장애 주입 실패")
    end
    local p = sender:request(resp.duey, recovery_request(sender, recipient), nil, 5000)
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

local function lost_after_send_timer(ctx)
    local sender = ctx:bot(0)
    local recipient = ctx:bot(1)
    if prepare(ctx, sender, recipient) == false then
        return false
    end
    if send_unknown(ctx, sender, recipient, "lost") == false then
        return false
    end

    if settle.resume(sender) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
    settle.tick(ctx)
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
    if settle.resume(sender) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
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
    if settle.resume(sender) == false then
        return ctx:fail("정산 타이머 재개 실패")
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
    if fault.rpc_delay(sender, "SendParcel", 4000) == false then
        return ctx:fail("SendParcel 지연 주입 실패")
    end

    sender:send(recovery_request(sender, recipient))
    local moved = sender:transfer(1)
    fault.rpc_delay(recipient, "SendParcel", 0)
    if moved == false then
        return ctx:fail("채널 이동 실패")
    end
    if wait_refund(ctx, sender) == false then
        return false
    end

    ctx:sleep(5000)
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

local function claim_unknown(ctx, bot, mode)
    if clear(ctx, bot) == false or setup(ctx, bot, MESO, "-") == false then
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

    if mode == "delay" then
        if fault.rpc_delay(bot, "ClaimParcel", 4000) == false then
            return ctx:fail("ClaimParcel 지연 주입 실패")
        end
        bot:send(req.duey { mode = MODE.receive, parcel_id = parcel.id })
        return true
    end
    if settle.pause(bot) == false then
        return ctx:fail("정산 타이머 정지 실패")
    end
    if fault.rpc(bot, "ClaimParcel", mode) == false then
        return ctx:fail("ClaimParcel 장애 주입 실패")
    end
    local failed = bot:request(resp.duey, req.duey { mode = MODE.receive, parcel_id = parcel.id }, nil, 5000)
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
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
    if wait_holds(ctx, bot, 5, MESO + 1000, "주기 정리로 수령") == false then
        return false
    end
    settle.tick(ctx)
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
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
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
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
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

    ctx:sleep(5000)
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
    name = "Duey: 택배 배송·퀵배송·수령·어뷰징·응답 유실 복구",
    bot_count = 6,

    on_scenario_finished = function(ctx)
        settle.resume_all(ctx, 6)
    end,

    scenarios = {
        normal_flow,
        quick_flow,
        system_parcels,
        box_limit,
        forged_requests,
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
