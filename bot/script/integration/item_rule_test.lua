local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local HUNTING_GROUND = 104040000
local EQUIP = 1
local USE = 2
local CAP_SLOT = -1

local ONLY_CAP = 1002508
local POTION = 2000000
local HELMET_SCROLL = 2040000
local HAT = 1002008
local RETURN_SCROLL = 2030000
local LEGENDARY_SPIRIT = 1003
local EVERLAST_MAP = 209000001
local SNAIL_SHELL = 4000000
local MESO = 0x40000
local ONLY_MEDAL = 1142667
local TWO_HANDED_SWORD = 1402001
local SHIELD = 1092000
local WARRIOR = 100
local SHIELD_SLOT = -10
local WEAPON_SLOT = -11
local BOWMAN = 300
local DEX_ACCESSORY = 1012101
local BOW = 1452003
local FACE_SLOT = -2

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function setup(ctx, bot, items)
    if pq.command(bot, "/봇초기화 30 0 0 " .. items .. " -", "봇초기화 완료") == false then
        return ctx:fail(bot:name() .. " 봇 초기화 실패")
    end
    if bot:instance_move(HENESYS) == false then
        return ctx:fail(bot:name() .. " 맵 이동 실패")
    end
    return true
end

local function spawn(ctx, bot, item_id)
    local p = bot:request(resp.spawn_item, req.normal_chat { message = "/아이템드롭 " .. item_id }, function(p)
        return p.item_model ~= nil and p.item_model.id == item_id
    end, 5000)
    if p == false then
        return ctx:fail(bot:name() .. " 아이템을 떨어뜨리지 못함: " .. item_id)
    end
    return p.id
end

local function wear(ctx, bot, item_id, slot)
    local p = bot:request(resp.inventory_operation,
        req.move_item { inventory_type = EQUIP, source = bot:slot(item_id) or 0, dest = slot, count = 1 }, nil, 3000)
    return check(ctx, p ~= false, bot:name() .. " 착용 응답 없음: " .. item_id)
end

-- A worn one-of-a-kind item counts as held, so the same item on the ground cannot be picked up.
local function unique_worn(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, "-") == false then
        return false
    end
    if bot:force_instance_move(HENESYS, "다른 스위트가 치우는 바닥 아이템과 분리") == false then
        return ctx:fail("인스턴스 맵 이동 실패")
    end
    local oid = spawn(ctx, bot, ONLY_CAP)
    if oid == false then
        return false
    end
    if check(ctx, bot:loot(oid), "가지지 않은 고유 모자를 줍지 못함") == false then
        return false
    end
    if wear(ctx, bot, ONLY_CAP, CAP_SLOT) == false then
        return false
    end
    ctx:sleep(300)
    if check(ctx, (bot:items()[ONLY_CAP] or 0) == 0, "착용한 고유 모자가 인벤토리에 남음") == false then
        return false
    end

    oid = spawn(ctx, bot, ONLY_CAP)
    if oid == false then
        return false
    end
    if check(ctx, bot:loot(oid) == false, "착용 중인 고유 모자를 하나 더 주움") == false then
        return false
    end
    return check(ctx, (bot:items()[ONLY_CAP] or 0) == 0, "착용 중인 고유 모자가 인벤토리에 생김")
end

-- Only a scroll can enhance equipment; another consume item is refused and kept.
local function scroll_only(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot, HAT .. ":1," .. POTION .. ":1," .. HELMET_SCROLL .. ":1") == false then
        return false
    end
    if wear(ctx, bot, HAT, CAP_SLOT) == false then
        return false
    end

    local p = bot:request(resp.update_stats,
        req.enhance_equipment { scroll_slot = bot:slot(POTION) or 0, target_slot = CAP_SLOT }, function(p)
            return p.unlock_action
        end, 3000)
    if check(ctx, p ~= false, "물약을 주문서로 쓴 요청이 거절되지 않음") == false then
        return false
    end
    ctx:sleep(300)
    if check(ctx, (bot:items()[POTION] or 0) == 1, "주문서가 아닌 물약이 사라짐") == false then
        return false
    end

    if bot:request(resp.show_scroll_effect, req.enhance_equipment { scroll_slot = bot:slot(HELMET_SCROLL) or 0, target_slot = CAP_SLOT }, nil, 3000) == false then
        return ctx:fail("투구 주문서 응답 없음")
    end
    ctx:sleep(300)
    return check(ctx, (bot:items()[HELMET_SCROLL] or 0) == 0, "투구 주문서가 남음")
end

local function unlocked(p)
    return p.unlock_action
end

-- Only Legendary Spirit lets a scroll reach equipment that is not worn, and the effect says so.
local function legendary_spirit(ctx)
    local bot = ctx:bot(1)
    if setup(ctx, bot, HAT .. ":1," .. HELMET_SCROLL .. ":2") == false then
        return false
    end
    local scroll = req.enhance_equipment { scroll_slot = bot:slot(HELMET_SCROLL) or 0, target_slot = bot:slot(HAT) or 0 }
    if check(ctx, bot:request(resp.show_scroll_effect, scroll, nil, 1500) == false, "레전드리 스피릿 없이 인벤토리 모자가 강화됨") == false then
        return false
    end
    if check(ctx, (bot:items()[HELMET_SCROLL] or 0) == 2, "거절된 강화가 주문서를 소모함") == false then
        return false
    end

    if pq.command(bot, "/스킬레벨 " .. LEGENDARY_SPIRIT .. " 1", "스킬레벨 설정") == false then
        return ctx:fail("레전드리 스피릿 배우기 실패")
    end
    scroll = req.enhance_equipment { scroll_slot = bot:slot(HELMET_SCROLL) or 0, target_slot = bot:slot(HAT) or 0 }
    local effect = bot:request(resp.show_scroll_effect, scroll, function(p)
        return p.character_id == bot:id()
    end, 3000)
    if effect == false then
        return ctx:fail("레전드리 스피릿을 배웠는데 인벤토리 모자가 강화되지 않음")
    end
    return check(ctx, effect.legendary_spirit, "레전드리 스피릿 강화 이펙트가 일반 이펙트로 감")
end

-- Merging stacks of more than 65535 in total keeps every item.
local function merge_overflow(ctx)
    local bot = ctx:bot(1)
    local stacks = {}
    for i = 1, 7 do
        stacks[i] = POTION .. ":9999"
    end
    if setup(ctx, bot, table.concat(stacks, ",")) == false then
        return false
    end
    if check(ctx, (bot:items()[POTION] or 0) == 69993, "물약 9999개 7칸을 받지 못함: " .. tostring(bot:items()[POTION])) == false then
        return false
    end
    if bot:request(resp.update_stats, req.sort_inventory { inventory_type = USE }, unlocked, 3000) == false then
        return ctx:fail("소비 탭 정렬 응답 없음")
    end
    ctx:sleep(300)
    return check(ctx, (bot:items()[POTION] or 0) == 69993, "정렬 뒤 물약 개수가 바뀜: " .. tostring(bot:items()[POTION]))
end

-- A dead character can neither drop meso or items nor read a return scroll.
local function dead_refusals(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 1000 " .. POTION .. ":1," .. RETURN_SCROLL .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if pq.command(bot, "/체력바꾸기 0", "체력 설정") == false then
        return ctx:fail("체력 설정 실패")
    end

    if bot:request(resp.update_stats, req.drop_meso { count = 100 }, unlocked, 3000) == false then
        return ctx:fail("죽은 채 메소 버리기 응답 없음")
    end
    ctx:sleep(300)
    if check(ctx, bot:meso() == 1000, "죽은 채 메소를 버림: " .. bot:meso()) == false then
        return false
    end

    local drop = req.move_item { inventory_type = USE, source = bot:slot(POTION) or 0, dest = 0, count = 1 }
    if bot:request(resp.update_stats, drop, unlocked, 3000) == false then
        return ctx:fail("죽은 채 아이템 버리기 응답 없음")
    end
    ctx:sleep(300)
    if check(ctx, (bot:items()[POTION] or 0) == 1, "죽은 채 물약을 버림") == false then
        return false
    end

    local scroll = req.use_return_scroll { slot = bot:slot(RETURN_SCROLL) or 0, item_id = RETURN_SCROLL }
    if check(ctx, bot:request(resp.warp, scroll, nil, 1500) == false, "죽은 채 귀환 주문서로 이동함") == false then
        return false
    end
    if check(ctx, (bot:items()[RETURN_SCROLL] or 0) == 1 and bot:map() == HUNTING_GROUND, "죽은 채 귀환 주문서를 씀") == false then
        return false
    end
    if pq.command(bot, "/체력바꾸기 3000", "체력 설정") == false then
        return ctx:fail("체력 복구 실패")
    end
    return true
end

local function leave_on_everlast(ctx, bot)
    if bot:force_instance_move(EVERLAST_MAP, "버린 아이템이 다른 스위트와 섞이지 않게") == false then
        return ctx:fail("에버라스트 맵 이동 실패")
    end
    if bot:drop(SNAIL_SHELL, 1) == nil then
        return ctx:fail("달팽이 껍질을 버리지 못함")
    end
    local dropped = bot:request(resp.update_stats, req.drop_meso { count = 100 }, function(p)
        return p.stats[MESO] ~= nil
    end, 3000)
    if dropped == false then
        return ctx:fail("메소를 버리지 못함")
    end
    return true
end

local function returned(ctx, bot, what)
    if check(ctx, (bot:items()[SNAIL_SHELL] or 0) == 2, what .. " 뒤 버린 달팽이 껍질이 돌아오지 않음: " .. tostring(bot:items()[SNAIL_SHELL])) == false then
        return false
    end
    return check(ctx, bot:meso() == 1000, what .. " 뒤 버린 메소가 돌아오지 않음: " .. bot:meso())
end

-- What a character dropped on an everlast map comes back before it is saved on logout or a channel change.
local function everlast_drops(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 1000 " .. SNAIL_SHELL .. ":2 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if leave_on_everlast(ctx, bot) == false then
        return false
    end
    if bot:relog({ passwords = { "" } }) == false then
        return ctx:fail("재접속 실패")
    end
    if returned(ctx, bot, "재접속") == false then
        return false
    end

    if leave_on_everlast(ctx, bot) == false then
        return false
    end
    if bot:transfer(1) == false then
        return ctx:fail("채널 이동 실패")
    end
    if returned(ctx, bot, "채널 이동") == false then
        return false
    end
    if bot:transfer(0) == false then
        return ctx:fail("원래 채널 복귀 실패")
    end
    return true
end

-- A two-handed weapon and a shield never stay on together: wearing one takes the other off.
local function two_handed_shield(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 " .. WARRIOR .. " 0 " .. TWO_HANDED_SWORD .. ":1," .. SHIELD .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/힘바꾸기 40", "힘 설정") == false then
        return ctx:fail("힘 설정 실패")
    end
    if wear(ctx, bot, TWO_HANDED_SWORD, WEAPON_SLOT) == false then
        return false
    end
    if wear(ctx, bot, SHIELD, SHIELD_SLOT) == false then
        return false
    end
    ctx:sleep(300)
    if check(ctx, bot:slot(SHIELD) == nil, "방패를 착용하지 못함") == false then
        return false
    end
    if check(ctx, bot:slot(TWO_HANDED_SWORD) ~= nil, "방패를 끼웠는데 두손검이 벗겨지지 않음") == false then
        return false
    end

    if wear(ctx, bot, TWO_HANDED_SWORD, WEAPON_SLOT) == false then
        return false
    end
    ctx:sleep(300)
    if check(ctx, bot:slot(TWO_HANDED_SWORD) == nil, "두손검을 착용하지 못함") == false then
        return false
    end
    return check(ctx, bot:slot(SHIELD) ~= nil, "두손검을 끼웠는데 방패가 벗겨지지 않음")
end

-- An equipped item's DEX counts toward a weapon's DEX requirement.
local function equip_with_item_stats(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 " .. BOWMAN .. " 0 " .. DEX_ACCESSORY .. ":1," .. BOW .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/힘바꾸기 40", "힘 설정") == false then
        return ctx:fail("힘 설정 실패")
    end
    if pq.command(bot, "/덱스바꾸기 37", "덱스 설정") == false then
        return ctx:fail("덱스 설정 실패")
    end
    if wear(ctx, bot, DEX_ACCESSORY, FACE_SLOT) == false then
        return false
    end
    if wear(ctx, bot, BOW, WEAPON_SLOT) == false then
        return false
    end
    ctx:sleep(300)
    return check(ctx, bot:slot(BOW) == nil, "DEX 37에 장신구 DEX 5로 요구 DEX 40을 채웠는데 활을 착용하지 못함")
end

-- A one-of-a-kind item left on an everlast map stays there when the character already holds another.
local function everlast_unique(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 0 " .. ONLY_MEDAL .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:force_instance_move(EVERLAST_MAP, "버린 아이템이 다른 스위트와 섞이지 않게") == false then
        return ctx:fail("에버라스트 맵 이동 실패")
    end
    if bot:drop(ONLY_MEDAL, 1) == nil then
        return ctx:fail("고유 훈장을 버리지 못함")
    end
    if pq.command(bot, "/봇초기화 30 0 0 " .. ONLY_MEDAL .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("고유 훈장을 다시 받지 못함")
    end
    if bot:relog({ passwords = { "" } }) == false then
        return ctx:fail("재접속 실패")
    end
    return check(ctx, (bot:items()[ONLY_MEDAL] or 0) == 1,
        "바닥의 고유 훈장까지 돌아와 두 개가 됨: " .. tostring(bot:items()[ONLY_MEDAL]))
end

test_suite {
    name = "Item: 착용한 고유 장비 줍기·주문서 아닌 강화·레전드리 스피릿·정렬 묶음 합계·사망 중 버리기와 귀환 주문서·에버라스트 드롭·양손 무기와 방패·장비 스탯으로 채운 착용 요구치",
    bot_count = 2,

    scenarios = {
        unique_worn,
        scroll_only,
        legendary_spirit,
        merge_overflow,
        dead_refusals,
        everlast_drops,
        two_handed_shield,
        everlast_unique,
        equip_with_item_stats,
    },
}
