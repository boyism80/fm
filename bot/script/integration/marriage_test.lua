local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")
local fault = require("script/integration/lib/fault")
local settle = require("script/integration/lib/settle")

local VILLAGE = 680000000
local LOBBY = 680000200
local CATHEDRAL = 680000210
local EXIT_MAP = 680000500
local HENESYS = 100000000

local CLARANCE = 9201005
local VALENTINA = 9201006
local ANNA = 9201007
local MARGARET = 9201013
local PILA = 9201014
local ANGELIQUE = 9201036
local ALLEGRO = 9201000
local MAYA = 1012101
local MAYA_HOUSE = 100000001

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
local GIFT_MODE = { give = 6, receive = 7, close = 8, open_give = 9, open_receive = 10, given = 11, give_failed = 13, received = 15, receive_fail = 16 }
local SETTLE_WAIT_MS = 40000
local REPLY_DELAY_MS = 4000
local DUEY = { send = 3, receive = 5, delete = 6, close = 8 }
local ETC = 4
local TOKEN_QUEST = 1100
local MAYA_TOKEN_QUEST = 1115
local QUEST_START = 1

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
    if bot:map() ~= VILLAGE and route.travel(ctx, bot, VILLAGE) == false then
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
    if route.travel(ctx, bot, VILLAGE) == false then
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
            local p = bot:request(resp.engage_result, req.dialog { dialog_type = DIALOG.yes_no, next = true },
                function(p)
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

    local proposal = groom:request_on(bride, resp.engage_request,
        req.ring_action { mode = RING.propose, name = bride:name(), item_id = RING_BOX }, nil, 5000)
    if proposal == false then
        return ctx:fail("청혼 요청이 신부에게 오지 않음")
    end
    if check(ctx, proposal.name == groom:name() and proposal.character_id == groom:id(), "청혼 요청 내용이 다름: " .. tostring(proposal.name)) == false then
        return false
    end

    local engaged = bride:request_on(groom, resp.engage_result,
        req.ring_action { mode = RING.answer, accepted = true, name = groom:name(), character_id = groom:id() },
        function(p)
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
    local reserved = bride:request(resp.engage_result, req.ring_action { mode = RING.wishlist, wishes = { "파란 달팽이 껍질" } },
        function(p)
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
    local p, name = bot:request({ resp.warp, resp.dialog },
        req.dialog { dialog_type = DIALOG.list, next = true, selected = selected }, function(p, name)
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
    local open = guest:request(resp.wedding_gift, req.dialog { dialog_type = DIALOG.list, next = true, selected = 0 },
        function(p)
            return p.mode == GIFT_MODE.open_give
        end, 5000)
    if open == false then
        return ctx:fail("하객 선물 창이 열리지 않음")
    end
    if check(ctx, #open.wishes == 2 and open.wishes[1] == "모래 주머니", "신랑 위시리스트가 다름") == false then
        return false
    end
    local given = guest:request(resp.wedding_gift,
        req.wedding_present { mode = GIFT_MODE.give, slot = guest:slot(GIFT), item_id = GIFT, count = 2 }, function(p)
            return p.mode == GIFT_MODE.given
        end, 5000)
    if given == false then
        return ctx:fail("결혼 선물 전달 결과가 오지 않음")
    end
    for _, mode in ipairs({ "on", "lost" }) do
        if fault.rpc(guest, "GiveWeddingGift", mode) == false then
            return ctx:fail("GiveWeddingGift 장애 주입 실패")
        end
        local failed = guest:request(resp.wedding_gift,
            req.wedding_present { mode = GIFT_MODE.give, slot = guest:slot(GIFT), item_id = GIFT, count = 1 },
            function(p)
                return p.mode == GIFT_MODE.give_failed
            end, 5000)
        fault.rpc(guest, "GiveWeddingGift", "off")
        if failed == false then
            return ctx:fail("장애 중 결혼 선물 실패 결과가 오지 않음 (" .. mode .. ")")
        end
    end
    guest:send(req.wedding_present { mode = GIFT_MODE.close })
    ctx:sleep(500)
    if has(ctx, guest, GIFT, 1) == false then
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
    local couple = groom:request_on(guest, resp.wedding_couple, req.dialog { dialog_type = DIALOG.yes_no, next = true },
        nil, 150000)
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
        if route.travel(ctx, bot, VILLAGE) == false then
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

    local spouse = groom:request_on(bride, resp.spouse_map, req.normal_chat { message = "/맵이동 " .. HENESYS }, function(p)
        return p.character_id == groom:id() and p.map_id == HENESYS
    end, 5000)
    if spouse == false then
        return ctx:fail("신랑이 맵을 옮겼는데 신부에게 배우자 위치가 오지 않음")
    end
    if route.travel(ctx, groom, VILLAGE) == false then
        return ctx:fail("신랑 웨딩빌리지 복귀 실패")
    end
    return true
end

local function wait_count(ctx, bot, item_id, count, what)
    if settle.wait(ctx, function()
            return (bot:items()[item_id] or 0) == count
        end) then
        return true
    end
    return ctx:fail(string.format("%s %s: %d 개수 %d (기대 %d)", bot:name(), what, item_id, bot:items()[item_id] or 0, count))
end

local function resume(ctx, bot)
    if settle.resume(bot) == false then
        return ctx:fail("정산 타이머 재개 실패")
    end
    return true
end

local function pause(ctx, bot)
    if settle.pause(bot) == false then
        return ctx:fail("정산 타이머 정지 실패")
    end
    return true
end

-- The guest's unsent gift is refunded and the lost one stays given; the groom's lost claim is granted on settle.
local function claim_refunds(ctx, bot)
    for _ = 1, SETTLE_WAIT_MS / 500 do
        bot:send(req.duey { mode = DUEY.close })
        local p = bot:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
        if p == false then
            return ctx:fail(bot:name() .. " 택배 보관함이 열리지 않음")
        end
        local found = {}
        for _, parcel in ipairs(p.parcels) do
            if parcel.sender == "택배" then
                found[#found + 1] = parcel
            end
        end
        if #found > 0 then
            for _, parcel in ipairs(found) do
                if bot:request(resp.duey_removed, req.duey { mode = DUEY.receive, parcel_id = parcel.id }, function(got)
                        return got.reason == 4
                    end, 3000) == false then
                    return ctx:fail(bot:name() .. " 환불 택배 수령 실패")
                end
            end
            bot:send(req.duey { mode = DUEY.close })
            return true
        end
        bot:send(req.duey { mode = DUEY.close })
        ctx:sleep(500)
    end
    return ctx:fail(bot:name() .. " 환불 택배가 오지 않음")
end

local function gifts(ctx)
    local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
    if claim_refunds(ctx, guest) == false then
        return false
    end
    if wait_count(ctx, guest, GIFT, 2, "보내지 않은 선물 환불") == false then
        return false
    end
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
    local received = groom:request(resp.wedding_gift,
        req.wedding_present { mode = GIFT_MODE.receive, inventory_type = ETC, index = 0 }, function(p)
            return p.mode == GIFT_MODE.received
        end, 5000)
    if received == false then
        return ctx:fail("결혼 선물 수령 결과가 오지 않음")
    end
    ctx:sleep(500)
    if has(ctx, groom, GIFT, 2) == false then
        return false
    end

    if fault.rpc(groom, "ClaimWeddingGift", "lost") == false then
        return ctx:fail("ClaimWeddingGift 장애 주입 실패")
    end
    local failed = groom:request(resp.wedding_gift,
        req.wedding_present { mode = GIFT_MODE.receive, inventory_type = ETC, index = 0 }, function(p)
            return p.mode == GIFT_MODE.receive_fail
        end, 5000)
    fault.rpc(groom, "ClaimWeddingGift", "off")
    if failed == false then
        return ctx:fail("장애 중 결혼 선물 수령 실패 결과가 오지 않음")
    end
    groom:send(req.wedding_present { mode = GIFT_MODE.close })
    if wait_count(ctx, groom, GIFT, 3, "유실된 선물 수령 정리") == false then
        return false
    end
    if groom:transfer(1) == false or groom:transfer(0) == false then
        return ctx:fail("신랑 채널 이동 실패")
    end
    ctx:sleep(1000)
    if has(ctx, groom, GIFT, 3) == false then
        return false
    end
    return talk(ctx, groom, ANGELIQUE, "받으실 수 있는 선물이 없군요", "선물을 모두 받은 신랑")
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
    local hidden = {}
    ctx:hook("spouse_hidden", function(_, bot, p)
        if p.map_id == 0 and p.character_id == 0 then
            hidden[bot:id()] = true
        end
    end)
    local divorced = groom:request(resp.engage_result, req.dialog { dialog_type = DIALOG.yes_no, next = true },
        function(p)
            return p.result == ENGAGE.divorced
        end, 5000)
    if divorced == false then
        ctx:unhook("spouse_hidden")
        return ctx:fail("이혼 결과가 오지 않음")
    end
    ctx:sleep(1000)
    ctx:unhook("spouse_hidden")
    if check(ctx, hidden[groom:id()] == true and hidden[bride:id()] == true, "이혼 후 배우자 위치가 지워지지 않음") == false then
        return false
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

local function holds(ctx, bot, item_id, expected, what)
    local have = bot:items()[item_id] or 0
    return check(ctx, have == expected,
        string.format("%s %s: 아이템 %d 개수 %d (기대 %d)", bot:name(), what, item_id, have, expected))
end

-- /파혼 also clears the engagement for the next scenario.
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
    return check(ctx, was == expected,
        string.format("%s %s: 약혼 상태 %s (기대 %s)", bot:name(), what, tostring(was), tostring(expected)))
end

local function prepare_couple(ctx)
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
    local proposal = groom:request_on(bride, resp.engage_request,
        req.ring_action { mode = RING.propose, name = bride:name(), item_id = RING_BOX }, nil, 5000)
    return check(ctx, proposal ~= false, "청혼 요청이 신부에게 오지 않음")
end

local function engage_couple(ctx, groom, bride)
    if propose(ctx, groom, bride) == false then
        return false
    end
    local p = bride:request_on(groom, resp.engage_result, answer(groom), function(p)
        return p.result == ENGAGE.engaged
    end, 5000)
    if p == false then
        return ctx:fail("약혼 결과가 오지 않음")
    end
    state.marriage_id = p.marriage.id
    return true
end

-- The fault is server wide, so it is turned off right after the one answer it targets.
-- The settle timer stays paused so the held state can be seen; the caller resumes it.
local function proposal_unknown(ctx, mode)
    local groom, bride = ctx:bot(0), ctx:bot(1)
    if prepare_couple(ctx) == false or propose(ctx, groom, bride) == false then
        return false
    end
    if pause(ctx, bride) == false then
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
    if wait_count(ctx, groom, ENGAGEMENT_RING, expected, what) == false then
        return false
    end
    return wait_count(ctx, bride, ENGAGEMENT_RING, expected, what)
end

local function proposal_lost_timer(ctx)
    local groom, bride = ctx:bot(0), ctx:bot(1)
    if proposal_unknown(ctx, "lost") == false or resume(ctx, groom) == false then
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
    if wait_count(ctx, groom, ENGAGEMENT_RING, 1, "재접속 정리로 약혼반지") == false then
        return false
    end
    if groom:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    if resume(ctx, groom) == false then
        return false
    end
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
    if proposal_unknown(ctx, "on") == false or resume(ctx, groom) == false then
        return false
    end
    if wait_count(ctx, groom, RING_BOX, 1, "주기 정리로 상자 환불") == false then
        return false
    end
    settle.tick(ctx)
    if holds(ctx, groom, ENGAGEMENT_RING, 0, "환불 뒤") == false or holds(ctx, bride, ENGAGEMENT_RING, 0, "환불 뒤") == false then
        return false
    end
    return engaged(ctx, groom, false, "환불 뒤")
end

local function late_proposal_after_abort(ctx)
    local groom, bride = ctx:bot(0), ctx:bot(1)
    if prepare_couple(ctx) == false or propose(ctx, groom, bride) == false then
        return false
    end
    if fault.rpc_delay(bride, "CreateMarriage", REPLY_DELAY_MS) == false then
        return ctx:fail("CreateMarriage 지연 주입 실패")
    end
    bride:send(answer(groom))
    ctx:sleep(500)
    local moved = groom:transfer(1)
    fault.rpc_delay(bride, "CreateMarriage", 0)
    if moved == false then
        return ctx:fail("채널 이동 실패")
    end
    if wait_count(ctx, groom, RING_BOX, 1, "재접속 중단으로 상자 환불") == false then
        return false
    end

    ctx:sleep(REPLY_DELAY_MS + 1000)
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
    local dlg = click(ctx, groom, MARGARET)
    if dlg == false then
        return false
    end
    if dlg.selections == nil then
        groom:dialog(false)
        return ctx:fail("마가렛 예약 선택지가 오지 않음")
    end
    return said(ctx, groom, groom:dialog(true, 3), "약혼했군요", "마가렛 예약 확인") ~= false
end

-- Reserved couples see the wishlist step at Margaret instead of the reservation list.
local function reserved(ctx, groom, expected, what)
    local dlg = click(ctx, groom, MARGARET)
    groom:dialog(false)
    if dlg == false then
        return false
    end
    local was = dlg.text ~= nil and dlg.text:find("결혼식 예약접수를 완료하기 전에", 1, true) ~= nil
    return check(ctx, was == expected, string.format("%s: 예약 상태 %s (기대 %s)", what, tostring(was), tostring(expected)))
end

local function reservation_unknown(ctx, mode)
    local groom, bride = ctx:bot(0), ctx:bot(1)
    if prepare_couple(ctx) == false or engage_couple(ctx, groom, bride) == false then
        return false
    end
    if open_reservation(ctx, groom) == false then
        return false
    end
    if mode == "delay" then
        if fault.rpc_delay(groom, "ReserveWedding", REPLY_DELAY_MS) == false then
            return ctx:fail("ReserveWedding 지연 주입 실패")
        end
        groom:send(req.dialog { dialog_type = DIALOG.normal, next = true })
        ctx:sleep(500)
        return true
    end

    if pause(ctx, groom) == false then
        return false
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
    if reservation_unknown(ctx, "lost") == false or resume(ctx, groom) == false then
        return false
    end
    settle.tick(ctx)
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
    if resume(ctx, groom) == false then
        return false
    end
    if holds(ctx, groom, TICKET, 0, "복귀 뒤") == false then
        return false
    end
    return reserved(ctx, groom, true, "재접속 정리 뒤")
end

local function reservation_not_sent(ctx)
    local groom = ctx:bot(0)
    if reservation_unknown(ctx, "on") == false or resume(ctx, groom) == false then
        return false
    end
    if wait_count(ctx, groom, TICKET, 1, "주기 정리로 티켓 환불") == false then
        return false
    end
    return reserved(ctx, groom, false, "환불 뒤")
end

-- A script still waiting for the reply keeps the character on this channel, so the move waits for it.
local function reservation_delayed(ctx)
    local groom = ctx:bot(0)
    if reservation_unknown(ctx, "delay") == false then
        return false
    end
    ctx:sleep(REPLY_DELAY_MS + 1500)
    groom:dialog(false)
    fault.rpc_delay(groom, "ReserveWedding", 0)
    if groom:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    ctx:sleep(1000)
    if holds(ctx, groom, TICKET, 0, "늦은 응답으로 예약") == false then
        return false
    end
    if groom:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    if holds(ctx, groom, TICKET, 0, "채널 이동 후 (중복 환불 없음)") == false then
        return false
    end
    return reserved(ctx, groom, true, "늦은 응답 후")
end

local function wishlist_unknown(ctx, mode)
    local groom, bride = ctx:bot(0), ctx:bot(1)
    if prepare_couple(ctx) == false then
        return false
    end
    if groom:transfer(1) == false or bride:transfer(1) == false or groom:transfer(0) == false or bride:transfer(0) == false then
        return ctx:fail("위시리스트 테스트 전 약혼 상태 새로고침 실패")
    end
    if engage_couple(ctx, groom, bride) == false then
        return false
    end
    if open_reservation(ctx, groom) == false then
        return false
    end
    local input = groom:request(resp.engage_request, req.dialog { dialog_type = DIALOG.normal, next = true }, function(p)
        return p.mode == RING.wishlist
    end, 5000)
    if input == false then
        return ctx:fail("위시리스트 입력창이 오지 않음")
    end
    if pause(ctx, groom) == false then
        return false
    end
    if fault.rpc(groom, "SetWeddingWishlist", mode) == false then
        return ctx:fail("SetWeddingWishlist 장애 주입 실패")
    end
    groom:send(req.ring_action { mode = RING.wishlist, wishes = { "모래 주머니" } })
    ctx:sleep(1000)
    fault.rpc(groom, "SetWeddingWishlist", "off")
    return holds(ctx, groom, INVITATION, 0, "응답 유실 직후 청첩장 보류")
end

local function wishlist_lost_timer(ctx)
    local groom = ctx:bot(0)
    if wishlist_unknown(ctx, "lost") == false or resume(ctx, groom) == false then
        return false
    end
    if wait_count(ctx, groom, INVITATION, 5, "주기 정리 뒤 청첩장") == false then
        return false
    end
    settle.tick(ctx)
    return has(ctx, groom, INVITATION, 5)
end

local function wishlist_lost_relogin(ctx)
    local groom = ctx:bot(0)
    if wishlist_unknown(ctx, "lost") == false then
        return false
    end
    if groom:transfer(1) == false then
        return ctx:fail("위시리스트 유실 뒤 채널 이동 실패")
    end
    if wait_count(ctx, groom, INVITATION, 5, "재접속 뒤 청첩장") == false then
        return false
    end
    if groom:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    if resume(ctx, groom) == false then
        return false
    end
    return has(ctx, groom, INVITATION, 5)
end

local function wishlist_not_sent(ctx)
    local groom = ctx:bot(0)
    if wishlist_unknown(ctx, "on") == false or resume(ctx, groom) == false then
        return false
    end
    if wait_count(ctx, groom, INVITATION, 5, "미전송 요청 재시도 뒤 청첩장") == false then
        return false
    end
    return has(ctx, groom, INVITATION, 5)
end

local function invite_ready(ctx)
    local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
    if prepare_couple(ctx) == false then
        return false
    end
    if groom:transfer(1) == false or bride:transfer(1) == false or groom:transfer(0) == false or bride:transfer(0) == false then
        return ctx:fail("초대 테스트 전 약혼 상태 새로고침 실패")
    end
    if engage_couple(ctx, groom, bride) == false or open_reservation(ctx, groom) == false then
        return false
    end
    local input = groom:request(resp.engage_request, req.dialog { dialog_type = DIALOG.normal, next = true }, function(p)
        return p.mode == RING.wishlist
    end, 5000)
    if input == false then
        return ctx:fail("위시리스트 입력창이 오지 않음")
    end
    groom:send(req.ring_action { mode = RING.wishlist, wishes = { "모래 주머니" } })
    if wait_count(ctx, groom, INVITATION, 5, "초대 준비 청첩장") == false then
        return false
    end
    return clear_duey(ctx, guest)
end

local function invite(groom, guest)
    return req.ring_action { mode = RING.invite, name = guest:name(), marriage_id = state.marriage_id, slot = groom:slot(INVITATION) }
end

local function invited(ctx, guest, expected, what)
    guest:send(req.duey { mode = DUEY.close })
    local p = guest:request(resp.duey_open, req.normal_chat { message = "/듀이" }, nil, 5000)
    guest:send(req.duey { mode = DUEY.close })
    if p == false then
        return ctx:fail(guest:name() .. " 택배 보관함이 열리지 않음")
    end
    return check(ctx, #p.parcels == expected, string.format("%s: 하객 택배 %d개 (기대 %d)", what, #p.parcels, expected))
end

local function invite_unknown(ctx, mode)
    local groom, guest = ctx:bot(0), ctx:bot(2)
    if invite_ready(ctx) == false then
        return false
    end
    if pause(ctx, groom) == false then
        return false
    end
    if fault.rpc(groom, "InviteWeddingGuest", mode) == false then
        return ctx:fail("InviteWeddingGuest 장애 주입 실패")
    end
    groom:send(invite(groom, guest))
    ctx:sleep(1000)
    fault.rpc(groom, "InviteWeddingGuest", "off")
    return holds(ctx, groom, INVITATION, 4, "장애 직후 청첩장 보류")
end

local function invite_lost_timer(ctx)
    local groom, guest = ctx:bot(0), ctx:bot(2)
    if invite_unknown(ctx, "lost") == false or resume(ctx, groom) == false then
        return false
    end
    settle.tick(ctx)
    if holds(ctx, groom, INVITATION, 4, "주기 정리 뒤 (이미 초대됨)") == false or invited(ctx, guest, 1, "주기 정리 뒤") == false then
        return false
    end
    local again = groom:request(resp.notice, invite(groom, guest), function(p)
        return p.message:find("이미 결혼식에 초대", 1, true) ~= nil
    end, 5000)
    if again == false then
        return ctx:fail("재초대가 이미 초대됨으로 거절되지 않음")
    end
    ctx:sleep(500)
    return holds(ctx, groom, INVITATION, 4, "거절된 재초대 뒤 환불")
end

local function invite_lost_relogin(ctx)
    local groom, guest = ctx:bot(0), ctx:bot(2)
    if invite_unknown(ctx, "lost") == false then
        return false
    end
    if groom:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    ctx:sleep(2000)
    if holds(ctx, groom, INVITATION, 4, "재접속 정리 뒤 (이미 초대됨)") == false then
        return false
    end
    if groom:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    if resume(ctx, groom) == false then
        return false
    end
    if holds(ctx, groom, INVITATION, 4, "복귀 뒤") == false then
        return false
    end
    return invited(ctx, guest, 1, "재접속 정리 뒤")
end

-- Settle resends the same request, so an invite that never reached internal goes through on settle.
local function invite_not_sent(ctx)
    local groom, guest = ctx:bot(0), ctx:bot(2)
    if invite_unknown(ctx, "on") == false or resume(ctx, groom) == false then
        return false
    end
    settle.tick(ctx)
    if holds(ctx, groom, INVITATION, 4, "재전송 뒤") == false then
        return false
    end
    return invited(ctx, guest, 1, "재전송 뒤")
end

local function invite_late(ctx)
    local groom, bride, guest = ctx:bot(0), ctx:bot(1), ctx:bot(2)
    if invite_ready(ctx) == false then
        return false
    end
    if fault.rpc_delay(groom, "InviteWeddingGuest", REPLY_DELAY_MS) == false then
        return ctx:fail("InviteWeddingGuest 지연 주입 실패")
    end
    groom:send(invite(groom, guest))
    ctx:sleep(500)
    local moved = groom:transfer(1)
    fault.rpc_delay(bride, "InviteWeddingGuest", 0)
    if moved == false then
        return ctx:fail("채널 이동 실패")
    end
    ctx:sleep(REPLY_DELAY_MS + 2000)
    if holds(ctx, groom, INVITATION, 4, "늦은 요청 뒤") == false or invited(ctx, guest, 1, "늦은 요청 뒤 (중복 발송 없음)") == false then
        return false
    end
    if groom:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    return holds(ctx, groom, INVITATION, 4, "복귀 뒤 (중복 환불 없음)")
end

local function invite_guest_full(ctx)
    local groom, guest = ctx:bot(0), ctx:bot(2)
    if invite_ready(ctx) == false then
        return false
    end
    for _ = 1, 50 do
        if command(ctx, groom, string.format("/택배 %s 시스템 - 1 1", guest:name()), "택배 발송") == false then
            return false
        end
    end
    ctx:sleep(2000)
    if invited(ctx, guest, 50, "택배함 채운 뒤") == false then
        return false
    end

    local sent = groom:request(resp.notice, invite(groom, guest), function(p)
        return p.message:find("청첩장을 보냈습니다", 1, true) ~= nil
    end, 5000)
    if sent == false then
        return ctx:fail("택배함이 가득 찬 하객에게 청첩장이 가지 않음")
    end
    ctx:sleep(500)
    if holds(ctx, groom, INVITATION, 4, "가득 찬 택배함으로 초대한 뒤") == false then
        return false
    end
    if invited(ctx, guest, 51, "가득 찬 택배함으로 초대한 뒤") == false then
        return false
    end
    return clear_duey(ctx, guest)
end

local function started(p, quest)
    return p.quest_status.quest_id == quest and p.quest_status.status == 1
end

-- The token quests start only while Allegro's quest is in progress.
local function token_quest(ctx)
    local guest = ctx:bot(2)
    if command(ctx, guest, "/성별 0", "성별 변경: 0") == false or setup(ctx, guest, "-") == false then
        return false
    end
    local dlg = click(ctx, guest, ALLEGRO)
    if dlg == false or said(ctx, guest, dlg, "사랑의 증표 4개", "알레그로 안내") == false then
        return false
    end
    local p = guest:request(resp.update_quest, req.dialog { dialog_type = DIALOG.normal, next = true }, function(p)
        return started(p, TOKEN_QUEST)
    end, 5000)
    guest:dialog(false)
    if p == false then
        return ctx:fail("알레그로가 퀘스트 1100을 시작하지 않음")
    end

    if guest:force_instance_move(MAYA_HOUSE, "증표 퀘스트") == false then
        return ctx:fail("마야의 집 이동 실패")
    end
    if pq.npc(ctx, guest, MAYA) == false then
        return false
    end
    p = guest:request(resp.update_quest,
        req.quest_action { mode = QUEST_START, quest_id = MAYA_TOKEN_QUEST, npcid = MAYA },
        function(p)
            return started(p, MAYA_TOKEN_QUEST)
        end, 5000)
    if p == false then
        return ctx:fail("마야의 사랑의 증표 퀘스트를 시작하지 못함")
    end
    return true
end

local function cleanup(ctx)
    break_off(ctx, ctx:bot(0))
    ctx:sleep(1000)
    break_off(ctx, ctx:bot(1))
    return command(ctx, ctx:bot(1), "/성별 0", "성별 변경: 0")
end

test_suite {
    name = "Marriage: 약혼·예약·청첩장·예식·선물·이혼·사랑의 증표 퀘스트·응답 유실 복구",
    bot_count = 3,

    on_scenario_finished = function(ctx)
        settle.resume_all(ctx, 3)
    end,

    scenarios = {
        prepare,
        engage,
        reserve,
        ceremony,
        gifts,
        divorce,
        proposal_lost_timer,
        proposal_lost_relogin,
        proposal_not_sent,
        late_proposal_after_abort,
        reservation_lost_timer,
        reservation_lost_relogin,
        reservation_not_sent,
        reservation_delayed,
        wishlist_lost_timer,
        wishlist_lost_relogin,
        wishlist_not_sent,
        invite_lost_timer,
        invite_lost_relogin,
        invite_not_sent,
        invite_late,
        invite_guest_full,
        token_quest,
        cleanup,
    },
}
