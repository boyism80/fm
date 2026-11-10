local pq = require("script/integration/lib/party_quest")
local route = require("script/integration/lib/route")

local HENESYS_PARK = 100000200
local KEEPER = 1012009

local ETC = 4
local TRADE = 3
local UNIQUE = 4031187
-- A one-of-a-kind item of a party quest (Magatia's suspicious liquid).
local PARTY_QUEST_ONLY = 4001132

local MODE = {
    create = 0x00,
    invite = 0x02,
    visit = 0x04,
    exit = 0x0A,
    put_item = 0x0D,
    confirm = 0x0F,
}
local LEAVE = { only = 8 }
local STORAGE = { take_out = 4, store = 5 }
local RESULT = { store = 0x0D, open = 0x16 }

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function has(ctx, bot, count)
    local have = bot:items()[UNIQUE] or 0
    return check(ctx, have == count, string.format("%s 고유 아이템 개수 %d (기대 %d)", bot:name(), have, count))
end

local function setup(ctx, bot, items)
    if pq.command(bot, "/봇초기화 30 0 100000 " .. items .. " -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if route.travel_instance(ctx, bot, HENESYS_PARK) == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패")
    end
    return true
end

local function spawn(ctx, bot, item_id)
    local p = bot:request(resp.spawn_item, req.normal_chat { message = "/아이템드롭 " .. item_id }, function(p)
        return p.item_model ~= nil and p.item_model.id == item_id
    end, 5000)
    if p == false then
        return ctx:fail(bot:name() .. " 아이템을 떨어뜨리지 못함")
    end
    return p.id
end

local function open_storage(ctx, bot)
    local oid = pq.npc(ctx, bot, KEEPER)
    if oid == false then
        return false
    end
    local p = bot:request(resp.storage, req.npc_click { oid = oid }, function(p)
        return p.result == RESULT.open
    end, 5000)
    return check(ctx, p ~= false, bot:name() .. " 창고가 열리지 않음")
end

local function open_trade(ctx, a, b)
    local enter = a:request(resp.trade_enter, req.mini_room { mode = MODE.create, type = TRADE }, nil, 5000)
    if check(ctx, enter ~= false, "교환창 생성 응답 없음") == false then
        return false
    end
    local invite = a:request_on(b, resp.trade_invite, req.mini_room { mode = MODE.invite, target_id = b:id() }, nil, 5000)
    if check(ctx, invite ~= false, "교환 초대가 상대에게 오지 않음") == false then
        return false
    end
    enter = b:request(resp.trade_enter, req.mini_room { mode = MODE.visit, sn = invite.sn }, nil, 5000)
    if check(ctx, enter ~= false, "초대 수락 입장 응답 없음") == false then
        return false
    end
    ctx:sleep(300)
    return true
end

-- A unique item put into a trade still counts as held, so neither the storage nor the ground can give a second one.
local function offered_unique(ctx)
    local a = ctx:bot(0)
    local b = ctx:bot(1)
    if setup(ctx, a, UNIQUE .. ":1") == false or setup(ctx, b, "-") == false then
        return false
    end

    if open_storage(ctx, a) == false then
        return false
    end
    local stored = a:request({ resp.storage, resp.storage_error },
        req.storage { mode = STORAGE.store, slot = a:slot(UNIQUE) or 0, item_id = UNIQUE, count = 1 }, nil, 3000)
    if check(ctx, stored ~= false and stored.result == RESULT.store, "고유 아이템 보관 실패") == false then
        return false
    end
    local oid = spawn(ctx, a, UNIQUE)
    if oid == false or check(ctx, a:loot(oid), "창고에 맡긴 고유 아이템을 줍지 못함") == false then
        return false
    end
    if has(ctx, a, 1) == false then
        return false
    end

    local slot = a:slot(UNIQUE) or 0
    if open_trade(ctx, a, b) == false then
        return false
    end
    local put = a:request_on(b, resp.trade_item,
        req.mini_room { mode = MODE.put_item, inventory_type = ETC, slot = slot, count = 1, trade_slot = 1 }, nil, 5000)
    if check(ctx, put ~= false, "고유 아이템 올리기가 상대에게 보이지 않음") == false then
        return false
    end
    ctx:sleep(300)
    if has(ctx, a, 0) == false then
        return false
    end

    local refused = a:request(resp.update_stats, req.storage { mode = STORAGE.take_out, inventory_type = ETC, index = 0 },
        function(p)
            return p.unlock_action
        end, 3000)
    if check(ctx, refused ~= false, "교환에 올린 고유 아이템을 창고에서 또 꺼냄") == false then
        return false
    end
    ctx:sleep(300)
    if has(ctx, a, 0) == false then
        return false
    end

    oid = spawn(ctx, a, UNIQUE)
    if oid == false or check(ctx, a:loot(oid) == false, "교환에 올린 고유 아이템을 또 주움") == false then
        return false
    end
    if has(ctx, a, 0) == false then
        return false
    end

    local left = b:request_on(a, resp.mini_room_leave, req.mini_room { mode = MODE.exit }, nil, 5000)
    if check(ctx, left ~= false, "상대 취소 알림이 오지 않음") == false then
        return false
    end
    ctx:sleep(500)
    return has(ctx, a, 1)
end

-- A trade that would give a character a second one-of-a-kind item closes as a one-of-a-kind refusal, and both keep theirs.
local function held_unique_trade(ctx)
    local a = ctx:bot(0)
    local b = ctx:bot(1)
    if setup(ctx, a, UNIQUE .. ":1") == false or setup(ctx, b, UNIQUE .. ":1") == false then
        return false
    end
    if open_trade(ctx, a, b) == false then
        return false
    end
    local put = a:request_on(b, resp.trade_item,
        req.mini_room { mode = MODE.put_item, inventory_type = ETC, slot = a:slot(UNIQUE) or 0, count = 1, trade_slot = 1 },
        nil, 5000)
    if check(ctx, put ~= false, "고유 아이템 올리기가 상대에게 보이지 않음") == false then
        return false
    end
    if a:request_on(b, resp.trade_confirm, req.mini_room { mode = MODE.confirm }, nil, 5000) == false then
        return ctx:fail("교환 확정이 상대에게 보이지 않음")
    end
    local closed = b:request(resp.mini_room_leave, req.mini_room { mode = MODE.confirm }, nil, 5000)
    if check(ctx, closed ~= false and closed.reason == LEAVE.only, "이미 가진 고유 아이템 교환이 고유 아이템 거절로 끝나지 않음") == false then
        return false
    end
    ctx:sleep(500)
    return has(ctx, a, 1) and has(ctx, b, 1)
end

-- A party quest item is never one of a kind, so a character holding one picks up another.
local function party_quest_only(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, PARTY_QUEST_ONLY .. ":1") == false then
        return false
    end
    local oid = spawn(ctx, bot, PARTY_QUEST_ONLY)
    if oid == false or check(ctx, bot:loot(oid), "가진 파티 퀘스트 아이템을 하나 더 줍지 못함") == false then
        return false
    end
    ctx:sleep(300)
    local have = bot:items()[PARTY_QUEST_ONLY] or 0
    return check(ctx, have == 2, "파티 퀘스트 아이템 개수 " .. have .. " (기대 2)")
end

test_suite {
    name = "Unique: 교환에 올린 고유 아이템은 보유로 침·이미 가진 고유 아이템 교환 거절·파티 퀘스트 아이템 줍기",
    bot_count = 2,

    scenarios = {
        offered_unique,
        held_unique_trade,
        party_quest_only,
    },
}
