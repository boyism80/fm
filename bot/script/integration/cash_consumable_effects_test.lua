local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local TOWN = 100000000
local SONG = 5100000
local MEGAPHONE = 5070000
local ITEM_MEGAPHONE = 5076000
local TRIPLE_MEGAPHONE = 5077000
local KITE = 5080000
local NOTE = 5090000
local CHALKBOARD = 5370000
local GREEN_LENS = 5152103
local BLACK_LENS = 5152100
local FACE = 0x2
local ATTACHMENT = 4000000
local ETC = 4
local MSG_ITEM_MEGAPHONE = 8
local MSG_TRIPLE_MEGAPHONE = 10
local NOTES_SENT = 4
local NOTES_FAILED = 5
local QUIET_MS = 1500

local function setup(ctx, bot, items)
    if pq.command(bot, "/봇초기화 30 0 0 " .. items, "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    return true
end

local function meet(ctx, actor, observer)
    if actor:instance_move(HUNTING_GROUND) == false or observer:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    return true
end

local function count(bot, item_id)
    return bot:items()[item_id] or 0
end

local function consumed(ctx, bot, item_id, before, label)
    if count(bot, item_id) ~= before - 1 then
        bot:request(resp.inventory_operation, nil, nil, 3000)
    end
    if count(bot, item_id) ~= before - 1 then
        return ctx:fail(label .. " 사용 후 아이템이 한 개 줄지 않음")
    end
    return true
end

local function kept(ctx, bot, item_id, before, label)
    bot:request(resp.inventory_operation, nil, nil, 500)
    if count(bot, item_id) ~= before then
        return ctx:fail(label .. " 실패했는데 아이템이 소모됨")
    end
    return true
end

local function use(bot, item_id, fields)
    fields = fields or {}
    fields.slot = bot:slot(item_id) or 0
    fields.item_id = item_id
    return req.use_cash_item(fields)
end

local function chalkboard(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, CHALKBOARD .. ":3") == false or meet(ctx, actor, observer) == false then
        return false
    end

    local before = count(actor, CHALKBOARD)
    local long = string.rep("가", 21)
    if actor:request_on(observer, resp.chalkboard, use(actor, CHALKBOARD, { text = long }), nil, QUIET_MS) ~= false then
        return ctx:fail("40바이트를 넘는 칠판이 보임")
    end
    if kept(ctx, actor, CHALKBOARD, before, "긴 칠판") == false then
        return false
    end

    local shown = actor:request_on(observer, resp.chalkboard, use(actor, CHALKBOARD, { text = "칠판 테스트" }), function(p)
        return p.character_id == actor:id() and p.text == "칠판 테스트"
    end, 5000)
    if shown == false then
        return ctx:fail("칠판이 다른 플레이어에게 보이지 않음")
    end
    if kept(ctx, actor, CHALKBOARD, before, "기간제 칠판") == false then
        return false
    end

    if actor:request_on(observer, resp.chalkboard, use(actor, CHALKBOARD, { text = "칠판 테스트" }), nil, QUIET_MS) ~= false then
        return ctx:fail("같은 문장으로 칠판을 다시 알림")
    end
    if kept(ctx, actor, CHALKBOARD, before, "같은 칠판") == false then
        return false
    end

    if actor:request_on(observer, resp.chalkboard, req.close_chalkboard {}, function(p)
            return p.character_id == actor:id() and p.text == ""
        end, 5000) == false then
        return ctx:fail("칠판 닫기가 다른 플레이어에게 보이지 않음")
    end
    return true
end

local function note(ctx)
    local sender = ctx:bot(0)
    local receiver = ctx:bot(1)
    if setup(ctx, sender, NOTE .. ":2") == false then
        return false
    end

    local before = count(sender, NOTE)
    local failed = sender:request(resp.notes, use(sender, NOTE, { recipient = "없는캐릭터명", text = "쪽지 테스트" }), function(p)
        return p.action == NOTES_FAILED
    end, 5000)
    if failed == false then
        return ctx:fail("없는 수신자에게 쪽지 실패가 오지 않음")
    end
    if kept(ctx, sender, NOTE, before, "없는 수신자 쪽지") == false then
        return false
    end

    local sent = sender:request(resp.notes, use(sender, NOTE, { recipient = receiver:name(), text = "쪽지 테스트" }),
        function(p)
            return p.action == NOTES_SENT
        end, 5000)
    if sent == false then
        return ctx:fail("쪽지 전송 성공이 오지 않음")
    end
    return consumed(ctx, sender, NOTE, before, "쪽지")
end

local function message_box(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, KITE .. ":3") == false or meet(ctx, actor, observer) == false then
        return false
    end

    local before = count(actor, KITE)
    if actor:request_on(observer, resp.spawn_message_box, use(actor, KITE, { text = "한 줄 방패연" }), nil, QUIET_MS) ~= false then
        return ctx:fail("세 줄 형식이 아닌 메시지 박스가 보임")
    end
    local long = "\n" .. string.rep("가", 31) .. "\n"
    if actor:request_on(observer, resp.spawn_message_box, use(actor, KITE, { text = long }), nil, QUIET_MS) ~= false then
        return ctx:fail("한 줄이 60바이트를 넘는 메시지 박스가 보임")
    end
    if kept(ctx, actor, KITE, before, "잘못된 방패연") == false then
        return false
    end

    local text = "방패연 테스트\n둘째 줄\n"
    local spawned = actor:request_on(observer, resp.spawn_message_box, use(actor, KITE, { text = text }), function(p)
        return p.item_id == KITE and p.text == text and p.owner == actor:name()
    end, 5000)
    if spawned == false then
        return ctx:fail("메시지 박스가 다른 플레이어에게 보이지 않음")
    end
    if consumed(ctx, actor, KITE, before, "방패연") == false then
        return false
    end

    before = count(actor, KITE)
    if actor:request(resp.message_box_failed, use(actor, KITE, { text = "겹친 방패연\n\n" }), nil, 5000) == false then
        return ctx:fail("같은 자리에 메시지 박스를 다시 놓았는데 실패가 오지 않음")
    end
    return kept(ctx, actor, KITE, before, "겹친 방패연")
end

local function triple_megaphone(ctx)
    local sender = ctx:bot(0)
    local receiver = ctx:bot(1)
    if setup(ctx, sender, TRIPLE_MEGAPHONE .. ":2") == false then
        return false
    end

    local before = count(sender, TRIPLE_MEGAPHONE)
    if sender:request_on(receiver, resp.megaphone, use(sender, TRIPLE_MEGAPHONE, { lines = {} }), nil, QUIET_MS) ~= false then
        return ctx:fail("줄이 없는 세 줄 확성기가 전달됨")
    end
    if kept(ctx, sender, TRIPLE_MEGAPHONE, before, "빈 세 줄 확성기") == false then
        return false
    end

    local heard = sender:request_on(receiver, resp.megaphone,
        use(sender, TRIPLE_MEGAPHONE, { lines = { "첫째 줄", "둘째 줄", "셋째 줄" }, ear = true }), function(p)
            return p.type == MSG_TRIPLE_MEGAPHONE and p.lines[1]:find("첫째 줄", 1, true) ~= nil
        end, 5000)
    if heard == false then
        return ctx:fail("세 줄 확성기가 전달되지 않음")
    end
    if #heard.lines ~= 3 or heard.lines[3]:find("셋째 줄", 1, true) == nil then
        return ctx:fail("세 줄 확성기의 줄이 모두 전달되지 않음")
    end
    if heard.lines[2]:find(sender:name(), 1, true) == nil or heard.ear ~= true then
        return ctx:fail("세 줄 확성기의 이름이나 귀 표시가 빠짐")
    end
    return consumed(ctx, sender, TRIPLE_MEGAPHONE, before, "세 줄 확성기")
end

local function item_megaphone(ctx)
    local sender = ctx:bot(0)
    local receiver = ctx:bot(1)
    if setup(ctx, sender, ITEM_MEGAPHONE .. ":3," .. ATTACHMENT .. ":1") == false then
        return false
    end

    local before = count(sender, ITEM_MEGAPHONE)
    local missing = use(sender, ITEM_MEGAPHONE,
        { text = "빈 칸 첨부", attached = true, attached_type = ETC, attached_slot = 99 })
    if sender:request_on(receiver, resp.megaphone, missing, nil, QUIET_MS) ~= false then
        return ctx:fail("없는 아이템을 첨부한 확성기가 전달됨")
    end
    if kept(ctx, sender, ITEM_MEGAPHONE, before, "빈 칸 첨부 확성기") == false then
        return false
    end

    local attached = use(sender, ITEM_MEGAPHONE, {
        text = "아이템 확성기",
        attached = true,
        attached_type = ETC,
        attached_slot = sender:slot(ATTACHMENT),
    })
    local heard = sender:request_on(receiver, resp.megaphone, attached, function(p)
        return p.type == MSG_ITEM_MEGAPHONE and p.lines[1]:find("아이템 확성기", 1, true) ~= nil
    end, 5000)
    if heard == false then
        return ctx:fail("아이템 확성기가 전달되지 않음")
    end
    if heard.item == nil or heard.item.item_id ~= ATTACHMENT then
        return ctx:fail("아이템 확성기에 첨부 아이템이 없음")
    end
    if consumed(ctx, sender, ITEM_MEGAPHONE, before, "아이템 확성기") == false then
        return false
    end
    if count(sender, ATTACHMENT) ~= 1 then
        return ctx:fail("첨부한 아이템이 바뀜")
    end

    before = count(sender, ITEM_MEGAPHONE)
    heard = sender:request_on(receiver, resp.megaphone, use(sender, ITEM_MEGAPHONE, { text = "첨부 없음" }), function(p)
        return p.type == MSG_ITEM_MEGAPHONE and p.lines[1]:find("첨부 없음", 1, true) ~= nil
    end, 5000)
    if heard == false or heard.item ~= nil then
        return ctx:fail("첨부 없는 아이템 확성기가 올바르게 전달되지 않음")
    end
    return consumed(ctx, sender, ITEM_MEGAPHONE, before, "첨부 없는 아이템 확성기")
end

local function color_lens(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, GREEN_LENS .. ":1," .. BLACK_LENS .. ":1") == false then
        return false
    end

    local face_of = function(p)
        return p.stats[FACE] ~= nil
    end
    local green = bot:request(resp.update_stats, use(bot, GREEN_LENS), face_of, 5000)
    if green == false then
        return ctx:fail("초록색 렌즈를 썼는데 얼굴이 바뀌지 않음")
    end
    if math.floor(green.stats[FACE] / 100) % 10 ~= 3 then
        return ctx:fail("초록색 렌즈 뒤 얼굴 " .. green.stats[FACE] .. "의 눈 색이 3이 아님")
    end
    if consumed(ctx, bot, GREEN_LENS, 1, "초록색 렌즈") == false then
        return false
    end

    local black = bot:request(resp.update_stats, use(bot, BLACK_LENS), face_of, 5000)
    if black == false then
        return ctx:fail("검은색 렌즈를 썼는데 얼굴이 바뀌지 않음")
    end
    if black.stats[FACE] ~= green.stats[FACE] - 300 then
        return ctx:fail("검은색 렌즈가 눈 색 말고 다른 자리도 바꿈: " .. green.stats[FACE] .. " -> " .. black.stats[FACE])
    end
    return consumed(ctx, bot, BLACK_LENS, 1, "검은색 렌즈")
end

local function song(ctx)
    local actor = ctx:bot(0)
    local observer = ctx:bot(1)
    if setup(ctx, actor, SONG .. ":2") == false or meet(ctx, actor, observer) == false then
        return false
    end
    actor:request(resp.notice, req.normal_chat { message = "/맵리셋" }, nil, 300)

    local played = function(p)
        return p.item_id == SONG and p.name == actor:name()
    end
    local before = count(actor, SONG)
    if actor:request_on(observer, resp.cash_song, use(actor, SONG), played, 5000) == false then
        return ctx:fail("축하곡이 다른 플레이어에게 들리지 않음")
    end
    if consumed(ctx, actor, SONG, before, "축하곡") == false then
        return false
    end

    before = count(actor, SONG)
    if actor:request_on(observer, resp.cash_song, use(actor, SONG), nil, QUIET_MS) ~= false then
        return ctx:fail("곡이 끝나기 전에 다른 곡이 시작됨")
    end
    if kept(ctx, actor, SONG, before, "재생 중 축하곡") == false then
        return false
    end

    if observer:force_instance_move(TOWN, "축하곡 재입장") == false then
        return ctx:fail("마을 이동 실패")
    end
    local back = req.normal_chat { message = string.format("/인스턴스이동 %d %d", HUNTING_GROUND, ctx:seat()) }
    if observer:request(resp.cash_song, back, played, 5000) == false then
        return ctx:fail("다시 들어온 플레이어에게 재생 중인 곡이 오지 않음")
    end
    return true
end

local function unsupported(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, MEGAPHONE .. ":1") == false then
        return false
    end
    local before = count(bot, MEGAPHONE)
    if bot:request(resp.notice, use(bot, MEGAPHONE, { text = "맵 확성기" }), function(p)
            return p.message:find("사용할 수 없는 아이템", 1, true) ~= nil
        end, 5000) == false then
        return ctx:fail("지원하지 않는 캐시 아이템에 안내가 오지 않음")
    end
    return kept(ctx, bot, MEGAPHONE, before, "지원하지 않는 확성기")
end

test_suite {
    name = "Cash consumable effects: 칠판·쪽지·메시지 박스·세 줄/아이템 확성기·컬러렌즈·축하곡",
    bot_count = 2,

    scenarios = {
        unsupported,
        color_lens,
        song,
        chalkboard,
        note,
        message_box,
        triple_megaphone,
        item_megaphone,
    },
}
