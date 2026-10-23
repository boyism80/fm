local pq = require("script/integration/lib/party_quest")
local fault = require("script/integration/lib/fault")

local GUILD_HQ = 200000301
local HERACLE = 2010007
local LENARIO = 2010009
local GUILD_SHOW_INFO = 0x1A
local GUILD_EMBLEM_CHANGE = 0x42
local DIALOG = { yes_no = 1 }
local HERACLE_MENU = { create = 0, capacity = 2 }
local GUILD_NAME_INPUT = 1
local LENARIO_MENU = { create = 2, capacity = 3, disband = 4 }

local MESO = 10000000
local CREATE_COST = 1500000
local EMBLEM_COST = 5000000
local CAPACITY_COST = 500000
local ALLIANCE_COST = 5000000
local ALLIANCE_CAPACITY_COST = 5000000
local ALLIANCE_CAPACITY = 2
local GP = 100
local SETTLE_WAIT_MS = 45000
local REPLY_DELAY_MS = 6000

local WATCHED = {
    guild_capacity_change = true,
    guild_update_gp = true,
    alliance_create = true,
    alliance_update_info = true,
    alliance_disband = true,
}

local seen = {}

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function code(value)
    return function(p)
        return p.code == value
    end
end

local function setup(ctx, bot)
    if pq.command(bot, string.format("/봇초기화 50 0 %d - -", MESO), "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if bot:instance_move(GUILD_HQ) == false then
        return ctx:fail(bot:name() .. " 길드 본부 이동 실패")
    end
    return true
end

local function holds(ctx, bot, meso, what)
    return check(ctx, bot:meso() == meso, string.format("%s %s: 메소 %d (기대 %d)", bot:name(), what, bot:meso(), meso))
end

local function create_request(bot)
    return req.guild_operation { operation = GUILD.Create, guild_name = bot:name() }
end

local function emblem_request()
    return req.guild_operation { operation = GUILD.ChangeEmblem, logo_bg = 1001, logo_bg_color = 2, logo = 2001, logo_color = 3 }
end

-- The fault is server wide, so it is turned off right after the one request it targets.
local function send_unknown(ctx, bot, method, request, mode)
    if fault.rpc(bot, method, mode) == false then
        return ctx:fail(method .. " 장애 주입 실패")
    end
    bot:send(request)
    ctx:sleep(1500)
    fault.rpc(bot, method, "off")
    return true
end

local function change_emblem(ctx, bot)
    if bot:request(resp.guild_message, emblem_request(), code(GUILD_EMBLEM_CHANGE), 10000) == false then
        return ctx:fail(bot:name() .. " 길드 마크 변경 알림 없음")
    end
    return true
end

local function create_lost_timer(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot) == false then
        return false
    end
    if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "lost") == false then
        return false
    end
    if holds(ctx, bot, MESO - CREATE_COST, "장애 직후 보류") == false then
        return false
    end
    if bot:request(resp.guild_message, nil, code(GUILD_SHOW_INFO), SETTLE_WAIT_MS) == false then
        return ctx:fail("주기 정리 뒤 길드 정보가 오지 않음")
    end
    return holds(ctx, bot, MESO - CREATE_COST, "주기 정리 뒤 (환불 없음)")
end

local function create_lost_relogin(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot) == false then
        return false
    end
    if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "lost") == false then
        return false
    end
    if bot:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    ctx:sleep(3000)
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    if holds(ctx, bot, MESO - CREATE_COST, "재접속 정리 뒤 (환불 없음)") == false then
        return false
    end
    if change_emblem(ctx, bot) == false then
        return false
    end
    ctx:sleep(500)
    return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "길드 마크 변경 뒤")
end

local function create_not_sent(ctx)
    local bot = ctx:bot(2)
    if setup(ctx, bot) == false then
        return false
    end
    if send_unknown(ctx, bot, "CreateGuild", create_request(bot), "on") == false then
        return false
    end
    if bot:request(resp.guild_message, nil, code(GUILD_SHOW_INFO), SETTLE_WAIT_MS) == false then
        return ctx:fail("주기 정리로 다시 보낸 길드 생성 정보가 오지 않음")
    end
    return holds(ctx, bot, MESO - CREATE_COST, "다시 보낸 생성 뒤")
end

local function late_create_after_resend(ctx)
    local bot = ctx:bot(3)
    local other = ctx:bot(4)
    if setup(ctx, bot) == false then
        return false
    end
    if fault.rpc_delay(bot, "CreateGuild", 6000) == false then
        return ctx:fail("CreateGuild 지연 주입 실패")
    end
    bot:send(create_request(bot))
    local moved = bot:transfer(1)
    fault.rpc_delay(other, "CreateGuild", 0)
    if moved == false then
        return ctx:fail("채널 이동 실패")
    end

    ctx:sleep(10000)
    if holds(ctx, bot, MESO - CREATE_COST, "재접속 재요청 뒤") == false then
        return false
    end
    ctx:sleep(8000)
    if holds(ctx, bot, MESO - CREATE_COST, "늦은 원래 요청 뒤 (환불·중복 차감 없음)") == false then
        return false
    end
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    return change_emblem(ctx, bot)
end

local function emblem_lost(ctx)
    local bot = ctx:bot(0)
    if send_unknown(ctx, bot, "ChangeGuildEmblem", emblem_request(), "lost") == false then
        return false
    end
    if holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "장애 직후 보류") == false then
        return false
    end
    ctx:sleep(SETTLE_WAIT_MS)
    return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "주기 정리 뒤 (환불 없음)")
end

local function emblem_not_sent(ctx)
    local bot = ctx:bot(2)
    if send_unknown(ctx, bot, "ChangeGuildEmblem", emblem_request(), "on") == false then
        return false
    end
    if bot:request(resp.guild_message, nil, code(GUILD_EMBLEM_CHANGE), SETTLE_WAIT_MS) == false then
        return ctx:fail("주기 정리로 다시 보낸 길드 마크 변경 알림 없음")
    end
    ctx:sleep(500)
    return holds(ctx, bot, MESO - CREATE_COST - EMBLEM_COST, "다시 보낸 마크 변경 뒤")
end

-- Internal announces each change it applies, so a duplicated effect shows up as a second packet.
local function watch(ctx)
    seen = {}
    ctx:hook("guild", function(_, bot, p, name)
        if WATCHED[name] == nil then
            return
        end
        local key = bot:id() .. name
        seen[key] = seen[key] or {}
        table.insert(seen[key], p)
    end)
end

local function packets(bot, name)
    return seen[bot:id() .. name] or {}
end

local function arrived(ctx, bot, name, what)
    for _ = 1, SETTLE_WAIT_MS / 500 do
        if #packets(bot, name) > 0 then
            return true
        end
        ctx:sleep(500)
    end
    return ctx:fail(string.format("%s %s: %s 없음", bot:name(), what, name))
end

local function once(ctx, bot, name, what)
    local count = #packets(bot, name)
    return check(ctx, count == 1, string.format("%s %s: %s %d개 (기대 1)", bot:name(), what, name, count))
end

local function relogin(ctx, bot)
    if bot:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    ctx:sleep(3000)
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    return true
end

local function yes(bot)
    bot:send(req.dialog { dialog_type = DIALOG.yes_no, next = true })
end

local function ask(ctx, bot, npc, selected, expected)
    local oid = pq.npc(ctx, bot, npc)
    if oid == false then
        return false
    end
    local dlg = bot:npc_click(oid)
    if dlg == false or dlg == nil or dlg.selections == nil then
        bot:dialog(false)
        return ctx:fail(bot:name() .. " NPC 선택지가 오지 않음: " .. npc)
    end
    dlg = bot:dialog(true, selected)
    if dlg == nil or dlg.text == nil or dlg.text:find(expected, 1, true) == nil then
        bot:dialog(false)
        return ctx:fail(string.format("%s NPC %d 확인 질문이 다름: %s", bot:name(), npc, dlg and tostring(dlg.text) or "nil"))
    end
    return true
end

local function ask_guild_create(ctx, bot)
    return ask(ctx, bot, HERACLE, HERACLE_MENU.create, "길드 제작 수수료")
end

local function name_input(bot)
    return bot:request(resp.guild_message, req.dialog { dialog_type = DIALOG.yes_no, next = true },
        code(GUILD_NAME_INPUT), 3000)
end

local function create_cancelled(ctx)
    local bot = ctx:bot(5)
    if setup(ctx, bot) == false then
        return false
    end
    if ask_guild_create(ctx, bot) == false then
        return false
    end
    if name_input(bot) == false then
        return ctx:fail("길드 이름 입력창이 오지 않음")
    end

    if ask_guild_create(ctx, bot) == false then
        return false
    end
    bot:dialog(false)
    return holds(ctx, bot, MESO, "이름 입력 취소 뒤")
end

local function create_during_question(ctx)
    local bot = ctx:bot(5)
    if ask_guild_create(ctx, bot) == false then
        return false
    end
    if bot:request(resp.guild_message, create_request(bot), code(GUILD_SHOW_INFO), 10000) == false then
        bot:dialog(false)
        return ctx:fail("길드 생성 정보가 오지 않음")
    end
    if bot:request(resp.guild_message, nil, code(GUILD_NAME_INPUT), 1500) ~= false then
        return ctx:fail("길드 생성 패킷이 열린 질문을 진행시킴")
    end
    if name_input(bot) == false then
        return ctx:fail("질문에 답한 뒤 대화가 이어지지 않음")
    end
    return holds(ctx, bot, MESO - CREATE_COST, "질문 중 길드 생성 뒤")
end

local function ask_guild_capacity(ctx, bot)
    return ask(ctx, bot, HERACLE, HERACLE_MENU.capacity, "길드 최대 인원 추가 비용")
end

local function ask_alliance_create(ctx, leader)
    if ask(ctx, leader, LENARIO, LENARIO_MENU.create, "정말 길드 연합을 만들고") == false then
        return false
    end
    if leader:dialog(true) == nil then
        return ctx:fail("연합 이름 입력창이 오지 않음")
    end
    if leader:dialog(true, 0, leader:name()) == nil then
        return ctx:fail("연합 이름 확인 질문이 오지 않음")
    end
    return true
end

local function ask_alliance_capacity(ctx, leader)
    return ask(ctx, leader, LENARIO, LENARIO_MENU.capacity, "길드 수를 늘리는데에는")
end

local function disband_alliance(ctx, leader)
    if ask(ctx, leader, LENARIO, LENARIO_MENU.disband, "해체하고 싶으세요") == false then
        return false
    end
    local done = leader:dialog(true)
    leader:dialog(false)
    if check(ctx, done ~= nil and done.text ~= nil and done.text:find("해체되었습니다", 1, true) ~= nil, "길드 연합 해체 실패: " .. tostring(done and done.text)) == false then
        return false
    end
    -- The guild leaves the alliance only after the disband packet; the dialog answers first.
    return arrived(ctx, leader, "alliance_disband", "연합 해체")
end

-- The fault is server wide, so it is turned off right after the one answer it targets.
local function answer_unknown(ctx, bot, method, mode, ask_fn)
    if fault.rpc(bot, method, mode) == false then
        return ctx:fail(method .. " 장애 주입 실패")
    end
    local asked = ask_fn(ctx, bot)
    if asked ~= false then
        yes(bot)
        ctx:sleep(1500)
        bot:dialog(false)
    end
    fault.rpc(bot, method, "off")
    return asked
end

-- A script still waiting for the reply keeps the character on this channel, so the move waits for it.
local function answer_delayed(ctx, bot, method, ask_fn)
    if fault.rpc_delay(bot, method, REPLY_DELAY_MS) == false then
        return ctx:fail(method .. " 지연 주입 실패")
    end
    local asked = ask_fn(ctx, bot)
    if asked ~= false then
        yes(bot)
        ctx:sleep(REPLY_DELAY_MS + 1500)
        bot:dialog(false)
    end
    fault.rpc_delay(bot, method, 0)
    if asked == false then
        return false
    end
    if bot:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    return true
end

local function capacity_unknown(ctx, bot, mode)
    watch(ctx)
    if setup(ctx, bot) == false then
        return false
    end
    if answer_unknown(ctx, bot, "IncreaseGuildCapacity", mode, ask_guild_capacity) == false then
        return false
    end
    return holds(ctx, bot, MESO - CAPACITY_COST, "장애 직후 보류")
end

local function capacity_lost_timer(ctx)
    local bot = ctx:bot(0)
    if capacity_unknown(ctx, bot, "lost") == false then
        return false
    end
    if arrived(ctx, bot, "guild_capacity_change", "반영된 인원 확장") == false then
        return false
    end
    ctx:sleep(SETTLE_WAIT_MS)
    if once(ctx, bot, "guild_capacity_change", "주기 정리 뒤 (중복 확장 없음)") == false then
        return false
    end
    return holds(ctx, bot, MESO - CAPACITY_COST, "주기 정리 뒤 (환불 없음)")
end

local function capacity_lost_relogin(ctx)
    local bot = ctx:bot(1)
    if capacity_unknown(ctx, bot, "lost") == false then
        return false
    end
    if arrived(ctx, bot, "guild_capacity_change", "반영된 인원 확장") == false or relogin(ctx, bot) == false then
        return false
    end
    if once(ctx, bot, "guild_capacity_change", "재접속 정리 뒤 (중복 확장 없음)") == false then
        return false
    end
    return holds(ctx, bot, MESO - CAPACITY_COST, "재접속 정리 뒤 (환불 없음)")
end

local function capacity_not_sent(ctx)
    local bot = ctx:bot(2)
    if capacity_unknown(ctx, bot, "on") == false then
        return false
    end
    if arrived(ctx, bot, "guild_capacity_change", "주기 정리로 다시 보낸 인원 확장") == false then
        return false
    end
    return holds(ctx, bot, MESO - CAPACITY_COST, "다시 보낸 인원 확장 뒤")
end

local function capacity_delayed(ctx)
    local bot = ctx:bot(3)
    watch(ctx)
    if setup(ctx, bot) == false then
        return false
    end
    if answer_delayed(ctx, bot, "IncreaseGuildCapacity", ask_guild_capacity) == false then
        return false
    end
    if arrived(ctx, bot, "guild_capacity_change", "늦은 응답으로 인원 확장") == false then
        return false
    end
    ctx:sleep(1000)
    if once(ctx, bot, "guild_capacity_change", "채널 이동 뒤 (중복 확장 없음)") == false then
        return false
    end
    if holds(ctx, bot, MESO - CAPACITY_COST, "채널 이동 뒤 (환불·중복 차감 없음)") == false then
        return false
    end
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

local function gain_gp(ctx, bot)
    if pq.command(bot, "/길드GP얻기 " .. GP, "길드 GP " .. GP .. " 획득 요청") == false then
        return ctx:fail(bot:name() .. " 길드 GP 명령 실패")
    end
    return true
end

local function gp_unknown(ctx, bot, mode)
    watch(ctx)
    if fault.rpc(bot, "GainGuildGP", mode) == false then
        return ctx:fail("GainGuildGP 장애 주입 실패")
    end
    local gained = gain_gp(ctx, bot)
    ctx:sleep(1500)
    fault.rpc(bot, "GainGuildGP", "off")
    return gained
end

local function gp_lost_timer(ctx)
    local bot = ctx:bot(0)
    if gp_unknown(ctx, bot, "lost") == false then
        return false
    end
    if arrived(ctx, bot, "guild_update_gp", "반영된 GP") == false then
        return false
    end
    ctx:sleep(SETTLE_WAIT_MS)
    return once(ctx, bot, "guild_update_gp", "주기 정리 뒤 (중복 지급 없음)")
end

local function gp_lost_relogin(ctx)
    local bot = ctx:bot(1)
    if gp_unknown(ctx, bot, "lost") == false then
        return false
    end
    if arrived(ctx, bot, "guild_update_gp", "반영된 GP") == false or relogin(ctx, bot) == false then
        return false
    end
    ctx:sleep(3000)
    return once(ctx, bot, "guild_update_gp", "재접속 정리 뒤 (중복 지급 없음)")
end

local function gp_not_sent(ctx)
    local bot = ctx:bot(2)
    if gp_unknown(ctx, bot, "on") == false then
        return false
    end
    if check(ctx, #packets(bot, "guild_update_gp") == 0, "보내지 않은 GP가 반영됨") == false then
        return false
    end
    if arrived(ctx, bot, "guild_update_gp", "주기 정리로 다시 보낸 GP") == false then
        return false
    end
    ctx:sleep(3000)
    return once(ctx, bot, "guild_update_gp", "다시 보낸 GP 뒤")
end

local function late_gp_after_resend(ctx)
    local bot, other = ctx:bot(3), ctx:bot(4)
    watch(ctx)
    if fault.rpc_delay(bot, "GainGuildGP", 6000) == false then
        return ctx:fail("GainGuildGP 지연 주입 실패")
    end
    local gained = gain_gp(ctx, bot)
    local moved = gained ~= false and bot:transfer(1)
    fault.rpc_delay(other, "GainGuildGP", 0)
    if gained == false then
        return false
    end
    if moved == false then
        return ctx:fail("채널 이동 실패")
    end
    if arrived(ctx, bot, "guild_update_gp", "재접속 재요청으로 GP") == false then
        return false
    end
    ctx:sleep(8000)
    if once(ctx, bot, "guild_update_gp", "늦은 원래 요청 뒤 (중복 지급 없음)") == false then
        return false
    end
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

-- Two guild masters in one party on Lenario's map; the party leader becomes the alliance leader.
local function form_party(ctx, leader, partner)
    if leader:request(resp.party_created, req.party_operation { operation = PARTY.Create }) == false then
        return ctx:fail("파티 생성 실패")
    end
    local invite = leader:request_on(partner, resp.party_invite,
        req.party_operation { operation = PARTY.Invite, target_name = partner:name() })
    if invite == false then
        return ctx:fail("파티 초대가 오지 않음")
    end
    if partner:request(resp.party_update_join, req.party_operation { operation = PARTY.AcceptInvite, party_id = invite.party_id }) == false then
        return ctx:fail("파티 가입 실패")
    end
    return true
end

local function created(ctx, bot, what)
    if arrived(ctx, bot, "alliance_create", what) == false then
        return false
    end
    local list = packets(bot, "alliance_create")
    local id = list[1].info.alliance_id
    for _, p in ipairs(list) do
        if p.info.alliance_id ~= id then
            return ctx:fail(string.format("%s %s: 연합이 둘 생김 (%d, %d)", bot:name(), what, id, p.info.alliance_id))
        end
    end
    return true
end

local function alliance_capacity(bot)
    local top = 0
    for _, p in ipairs(packets(bot, "alliance_update_info")) do
        if p.info ~= nil and p.info.capacity > top then
            top = p.info.capacity
        end
    end
    return top
end

local function capacity_is(ctx, bot, what)
    for _ = 1, SETTLE_WAIT_MS / 500 do
        if alliance_capacity(bot) > ALLIANCE_CAPACITY then
            break
        end
        ctx:sleep(500)
    end
    local have = alliance_capacity(bot)
    return check(ctx, have == ALLIANCE_CAPACITY + 1,
        string.format("%s %s: 연합 정원 %d (기대 %d)", bot:name(), what, have, ALLIANCE_CAPACITY + 1))
end

local function alliance_unknown(ctx, mode)
    local leader, partner = ctx:bot(0), ctx:bot(1)
    watch(ctx)
    if setup(ctx, leader) == false or setup(ctx, partner) == false then
        return false
    end
    if answer_unknown(ctx, leader, "CreateAlliance", mode, ask_alliance_create) == false then
        return false
    end
    return holds(ctx, leader, MESO - ALLIANCE_COST, "장애 직후 보류")
end

local function alliance_capacity_unknown(ctx, mode)
    local leader = ctx:bot(0)
    watch(ctx)
    if setup(ctx, leader) == false then
        return false
    end
    if answer_unknown(ctx, leader, "IncreaseAllianceCapacity", mode, ask_alliance_capacity) == false then
        return false
    end
    return holds(ctx, leader, MESO - ALLIANCE_CAPACITY_COST, "장애 직후 보류")
end

local function alliance_create_lost_timer(ctx)
    local leader = ctx:bot(0)
    if form_party(ctx, leader, ctx:bot(1)) == false or alliance_unknown(ctx, "lost") == false then
        return false
    end
    if created(ctx, leader, "반영된 연합 생성") == false then
        return false
    end
    ctx:sleep(SETTLE_WAIT_MS)
    if created(ctx, leader, "주기 정리 뒤") == false then
        return false
    end
    return holds(ctx, leader, MESO - ALLIANCE_COST, "주기 정리 뒤 (환불 없음)")
end

local function alliance_create_lost_relogin(ctx)
    local leader = ctx:bot(0)
    if alliance_unknown(ctx, "lost") == false then
        return false
    end
    if created(ctx, leader, "반영된 연합 생성") == false or relogin(ctx, leader) == false then
        return false
    end
    if created(ctx, leader, "재접속 정리 뒤") == false then
        return false
    end
    return holds(ctx, leader, MESO - ALLIANCE_COST, "재접속 정리 뒤 (환불 없음)")
end

local function alliance_create_not_sent(ctx)
    local leader = ctx:bot(0)
    if alliance_unknown(ctx, "on") == false then
        return false
    end
    if created(ctx, leader, "주기 정리로 다시 보낸 연합 생성") == false then
        return false
    end
    return holds(ctx, leader, MESO - ALLIANCE_COST, "다시 보낸 연합 생성 뒤")
end

local function alliance_create_delayed(ctx)
    local leader, partner = ctx:bot(0), ctx:bot(1)
    watch(ctx)
    if setup(ctx, leader) == false or setup(ctx, partner) == false then
        return false
    end
    if answer_delayed(ctx, leader, "CreateAlliance", ask_alliance_create) == false then
        return false
    end
    if created(ctx, leader, "늦은 응답으로 연합 생성") == false then
        return false
    end
    ctx:sleep(1000)
    if created(ctx, leader, "채널 이동 뒤") == false then
        return false
    end
    if holds(ctx, leader, MESO - ALLIANCE_COST, "채널 이동 뒤 (환불·중복 차감 없음)") == false then
        return false
    end
    if leader:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

local function alliance_capacity_lost_timer(ctx)
    local leader = ctx:bot(0)
    if alliance_capacity_unknown(ctx, "lost") == false then
        return false
    end
    if capacity_is(ctx, leader, "반영된 연합 정원") == false then
        return false
    end
    ctx:sleep(SETTLE_WAIT_MS)
    if capacity_is(ctx, leader, "주기 정리 뒤 (중복 확장 없음)") == false then
        return false
    end
    if holds(ctx, leader, MESO - ALLIANCE_CAPACITY_COST, "주기 정리 뒤 (환불 없음)") == false then
        return false
    end
    return disband_alliance(ctx, leader)
end

local function alliance_capacity_lost_relogin(ctx)
    local leader = ctx:bot(0)
    if alliance_capacity_unknown(ctx, "lost") == false then
        return false
    end
    if capacity_is(ctx, leader, "반영된 연합 정원") == false or relogin(ctx, leader) == false then
        return false
    end
    ctx:sleep(3000)
    if capacity_is(ctx, leader, "재접속 정리 뒤 (중복 확장 없음)") == false then
        return false
    end
    if holds(ctx, leader, MESO - ALLIANCE_CAPACITY_COST, "재접속 정리 뒤 (환불 없음)") == false then
        return false
    end
    return disband_alliance(ctx, leader)
end

local function alliance_capacity_not_sent(ctx)
    local leader = ctx:bot(0)
    if alliance_capacity_unknown(ctx, "on") == false then
        return false
    end
    if capacity_is(ctx, leader, "주기 정리로 다시 보낸 연합 정원") == false then
        return false
    end
    if holds(ctx, leader, MESO - ALLIANCE_CAPACITY_COST, "다시 보낸 연합 정원 뒤") == false then
        return false
    end
    return disband_alliance(ctx, leader)
end

local function alliance_capacity_delayed(ctx)
    local leader = ctx:bot(0)
    watch(ctx)
    if setup(ctx, leader) == false then
        return false
    end
    if answer_delayed(ctx, leader, "IncreaseAllianceCapacity", ask_alliance_capacity) == false then
        return false
    end
    if capacity_is(ctx, leader, "늦은 응답으로 연합 정원") == false then
        return false
    end
    ctx:sleep(1000)
    if capacity_is(ctx, leader, "채널 이동 뒤 (중복 확장 없음)") == false then
        return false
    end
    if holds(ctx, leader, MESO - ALLIANCE_CAPACITY_COST, "채널 이동 뒤 (환불·중복 차감 없음)") == false then
        return false
    end
    if leader:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    ctx:sleep(1000)
    return disband_alliance(ctx, leader)
end

test_suite {
    name = "Guild: 길드 생성·마크·인원·GP·연합 응답 유실 복구",
    bot_count = 6,

    scenarios = {
        create_cancelled,
        create_during_question,
        create_lost_timer,
        create_lost_relogin,
        create_not_sent,
        late_create_after_resend,
        emblem_lost,
        emblem_not_sent,
        capacity_lost_timer,
        capacity_lost_relogin,
        capacity_not_sent,
        capacity_delayed,
        gp_lost_timer,
        gp_lost_relogin,
        gp_not_sent,
        late_gp_after_resend,
        alliance_create_lost_timer,
        alliance_capacity_lost_timer,
        alliance_create_lost_relogin,
        alliance_capacity_lost_relogin,
        alliance_create_not_sent,
        alliance_capacity_not_sent,
        alliance_create_delayed,
        alliance_capacity_delayed,
    },
}
