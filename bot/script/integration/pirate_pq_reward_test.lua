local pq = require("script/integration/lib/party_quest")

local REWARD = 925100600
local WYANG = 2094001
local SHURIKEN = 2070005
local GLASSES = 1022079
local DIALOG_NEXT = 0
local DIALOG_LIST = 4
local MENU_COUNT = 0

local function init(ctx, bot, cmp, items, stage)
    local command = string.format("/봇초기화 55 0 0 %s - 1204.cmp=%d pirate_pq.reward=%d", items, cmp, stage)
    if pq.command(bot, command, "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:force_instance_move(REWARD, "해적 PQ 보상 맵") == false then
        return ctx:fail("보상 맵 이동 실패")
    end
    return true
end

local function ask_count(ctx, bot)
    local oid = pq.npc(ctx, bot, WYANG)
    if oid == false then
        return false
    end
    local dlg = bot:npc_click(oid)
    if dlg == false or dlg.selections == nil then
        return ctx:fail("우양 선택지가 오지 않음")
    end
    local reply = bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = MENU_COUNT },
        nil, 5000)
    if reply == false then
        return ctx:fail("우양의 횟수 안내가 오지 않음")
    end
    return reply
end

local function claim(ctx, bot, expected)
    local offer = ask_count(ctx, bot)
    if offer == false then
        return false
    end
    if offer.text:find("감사하는 마음으로", 1, true) == nil then
        bot:dialog(false)
        return ctx:fail("보상 안내가 아님: " .. offer.text)
    end
    local result = bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_NEXT, next = true }, nil, 5000)
    bot:dialog(false)
    if result == false or result.text:find(expected, 1, true) == nil then
        return ctx:fail("보상 결과 대사가 다름: " .. (result and result.text or "없음"))
    end
    ctx:sleep(500)
    return true
end

local function told(ctx, bot, expected)
    local reply = ask_count(ctx, bot)
    if reply == false then
        return false
    end
    bot:dialog(false)
    if reply.text:find(expected, 1, true) == nil then
        return ctx:fail("우양 대사가 다름: " .. reply.text)
    end
    return true
end

test_suite {
    name = "Party Quest: 해적 PQ 클리어 횟수 보상",
    bot_count = 1,

    scenarios = {
        function(ctx)
            local bot = ctx:bot(0)
            if init(ctx, bot, 8, "-", 0) == false then
                return false
            end
            if claim(ctx, bot, "잘 받으셨나요") == false then
                return false
            end
            if (bot:items()[SHURIKEN] or 0) ~= 1 then
                return ctx:fail("8회 보상 뇌전 수리검이 없음")
            end
            if told(ctx, bot, "#b15#k 번") == false then
                return false
            end
            if (bot:items()[SHURIKEN] or 0) ~= 1 then
                return ctx:fail("뇌전 수리검이 다시 지급됨")
            end
            return true
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if init(ctx, bot, 30, "-", 2) == false then
                return false
            end
            local offer = ask_count(ctx, bot)
            if offer == false then
                return false
            end
            local refused = bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_NEXT, next = true }, nil, 5000)
            bot:dialog(false)
            if refused == false or refused.text:find("착용을 해제하시고", 1, true) == nil then
                return ctx:fail("안경 없이 30회 보상을 받으려는데 거절 대사가 없음")
            end
            if (bot:items()[GLASSES] or 0) ~= 0 then
                return ctx:fail("안경 없이 강화 안경이 지급됨")
            end
            return true
        end,
        function(ctx)
            local bot = ctx:bot(0)
            if init(ctx, bot, 30, GLASSES .. ":1", 2) == false then
                return false
            end
            if claim(ctx, bot, "잘 받으셨나요") == false then
                return false
            end
            if (bot:items()[GLASSES] or 0) ~= 1 then
                return ctx:fail("30회 보상 뒤 투명 안경 개수: " .. (bot:items()[GLASSES] or 0))
            end
            return told(ctx, bot, "수고하셨습니다")
        end,
    },
}
