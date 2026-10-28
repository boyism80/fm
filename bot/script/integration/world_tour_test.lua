local pq = require("script/integration/lib/party_quest")

local HENESYS = 100000000
local MUSHROOM_SHRINE = 800000000
local SHOWA = 801000000
local BATHHOUSE = 801000100
local SPINEL = 9000020
local SMILING_CAT = 9120013
local SPINEL_JAPAN = 2
local DIALOG_YES_NO = 1
local DIALOG_LIST = 4

local function travel(ctx, bot, npc, selected, map_id)
    local oid = pq.npc(ctx, bot, npc)
    if oid == false then
        return false
    end
    local dlg = bot:npc_click(oid)
    if dlg == false or dlg.selections == nil then
        return ctx:fail("선택지가 오지 않음: " .. npc)
    end
    if bot:request(resp.dialog_yes_no, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = selected }, nil,
            5000) == false then
        return ctx:fail("이동 확인이 오지 않음: " .. npc)
    end
    local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
        return p.character.map == map_id
    end, 10000)
    if warp == false then
        return ctx:fail("맵 이동 실패: " .. map_id)
    end
    return dlg
end

-- A paid move inside Showa keeps the town Spinel's world tour left from.
local function showa_keeps_departure(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 10000 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:instance_move(HENESYS) == false then
        return ctx:fail("헤네시스 이동 실패")
    end
    if travel(ctx, bot, SPINEL, SPINEL_JAPAN, MUSHROOM_SHRINE) == false then
        return false
    end
    if bot:instance_move(SHOWA) == false then
        return ctx:fail("쇼와 마을 이동 실패")
    end
    if travel(ctx, bot, SMILING_CAT, 0, BATHHOUSE) == false then
        return false
    end
    if bot:instance_move(MUSHROOM_SHRINE) == false then
        return ctx:fail("버섯신사 이동 실패")
    end
    local back = travel(ctx, bot, SPINEL, 0, HENESYS)
    if back == false then
        return false
    end
    if back.selections[1]:find("#m" .. HENESYS .. "#", 1, true) == nil then
        return ctx:fail("스피넬 귀환지가 출발 마을이 아님: " .. back.selections[1])
    end
    return true
end

test_suite {
    name = "World tour: 쇼와 마을 이동 뒤 스피넬 귀환",
    bot_count = 1,

    scenarios = {
        showa_keeps_departure,
    },
}
