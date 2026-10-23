local pq = require("script/integration/lib/party_quest")

local LITH_WEAPON_SHOP = 104000003
local SHOP_NPC = 1001000
local STAR = 2070000
local ROGUE = 400
local SHOP_LEAVE = 3
local EL_NATH_SHOP_MAP = 221022000
local EXPENSIVE_SHOP_NPC = 2041016
local EXPENSIVE_ITEM = 4031348
local EXPENSIVE_PRICE = 10000000
local OVERFLOWING_QUANTITY = 430
local OTHER_CHANNEL = 1
local HOME_CHANNEL = 0

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

-- Subi throwing stars: price 250 for a full stack of 500, unitPrice 0.3 per star.
local function recharge_and_sell(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 " .. ROGUE .. " 200000 " .. STAR .. ":100", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(LITH_WEAPON_SHOP) == false then
        return ctx:fail("무기 상점 이동 실패")
    end
    local oid = pq.npc(ctx, bot, SHOP_NPC)
    if oid == false then
        return false
    end
    if bot:request(resp.open_npc_shop, req.npc_click { oid = oid }, nil, 5000) == false then
        return ctx:fail("상점이 열리지 않음")
    end

    local slot = bot:slot(STAR)
    if check(ctx, slot ~= nil, "표창 칸을 찾지 못함") == false then
        return false
    end
    bot:shop_recharge(slot)
    ctx:sleep(500)
    local stars = bot:items()[STAR] or 0
    if check(ctx, stars == 500 and bot:meso() == 199880, string.format("400개 충전 결과가 다름: 표창 %d, 메소 %d (기대 500, 199880)", stars, bot:meso())) == false then
        return false
    end

    bot:shop_sell(slot, STAR, 500)
    ctx:sleep(500)
    stars = bot:items()[STAR] or 0
    return check(ctx, stars == 0 and bot:meso() == 200130,
        string.format("표창 500개 판매 결과가 다름: 표창 %d, 메소 %d (기대 0, 200130)", stars, bot:meso()))
end

-- 10,000,000 x 430 wraps to a few million as an int32, so the purchase has to be refused.
local function price_overflow(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 20000000 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(EL_NATH_SHOP_MAP) == false then
        return ctx:fail("상점 맵 이동 실패")
    end
    local oid = pq.npc(ctx, bot, EXPENSIVE_SHOP_NPC)
    if oid == false then
        return false
    end
    if bot:request(resp.open_npc_shop, req.npc_click { oid = oid }, nil, 5000) == false then
        return ctx:fail("상점이 열리지 않음")
    end

    bot:shop_buy(EXPENSIVE_ITEM, 1)
    ctx:sleep(500)
    local meso = 20000000 - EXPENSIVE_PRICE
    if check(ctx, bot:meso() == meso and (bot:items()[EXPENSIVE_ITEM] or 0) == 1, string.format("한 개 구매 결과가 다름: 메소 %d, 아이템 %d", bot:meso(), bot:items()[EXPENSIVE_ITEM] or 0)) == false then
        return false
    end

    bot:shop_buy(EXPENSIVE_ITEM, OVERFLOWING_QUANTITY)
    ctx:sleep(500)
    return check(ctx, bot:meso() == meso and (bot:items()[EXPENSIVE_ITEM] or 0) == 1,
        string.format("가격이 넘치는 구매가 거절되지 않음: 메소 %d, 아이템 %d", bot:meso(), bot:items()[EXPENSIVE_ITEM] or 0))
end

local function busy_while_open(ctx)
    local bot = ctx:bot(0)
    local oid = pq.npc(ctx, bot, SHOP_NPC)
    if oid == false then
        return false
    end
    if bot:request(resp.open_npc_shop, req.npc_click { oid = oid }, nil, 1500) ~= false then
        return ctx:fail("상점이 열린 채 NPC가 다시 열림")
    end
    if bot:request(resp.server_blocked, req.switch_channel { channel = OTHER_CHANNEL }, nil, 5000) == false then
        return ctx:fail("상점이 열린 채 채널 이동이 거절되지 않음")
    end

    bot:send(req.npc_shop { mode = SHOP_LEAVE })
    ctx:sleep(500)
    if bot:request(resp.open_npc_shop, req.npc_click { oid = oid }, nil, 5000) == false then
        return ctx:fail("상점을 닫은 뒤 다시 열리지 않음")
    end
    bot:send(req.npc_shop { mode = SHOP_LEAVE })
    ctx:sleep(500)
    if bot:transfer(OTHER_CHANNEL) == false then
        return ctx:fail("상점을 닫은 뒤 채널 이동 실패")
    end
    if bot:transfer(HOME_CHANNEL) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

test_suite {
    name = "NpcShop: 표창 충전·판매 가격·가격 넘침 구매 거절",
    bot_count = 1,

    scenarios = {
        recharge_and_sell,
        busy_while_open,
        price_overflow,
    },
}
