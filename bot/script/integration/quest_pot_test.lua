local pq = require("script/integration/lib/party_quest")
local check = require("script/integration/lib/script_check")

local QUEST = 3250
local POT_QUEST = 7067
local NPC = 2041026
local POT = 4220046
local FEED = 4031992
local OTHER_ETC = 4000000
local GROWTH = 100
local FULL = 3000
local SCRIPTED_START = 4
local SCRIPTED_END = 5

local REPLIES = { resp.update_quest }
for _, name in ipairs(check.REPLIES) do
    REPLIES[#REPLIES + 1] = name
end

local function talk(bot, mode, status)
    local p, name = bot:request(REPLIES, req.quest_action { mode = mode, quest_id = QUEST, npcid = NPC }, check.replied,
        5000)
    for _ = 1, 20 do
        if p == false then
            return false
        end
        if name == resp.update_quest and p.quest_status.quest_id == QUEST and p.quest_status.status == status then
            return true
        end
        local next = nil
        if name == resp.dialog_yes_no then
            next = req.dialog { dialog_type = 1, next = true }
        elseif name == resp.dialog then
            next = req.dialog { dialog_type = 0, next = true }
        end
        p, name = bot:request(REPLIES, next, check.replied, 5000)
    end
    return false
end

local function feed(bot, item, pred, timeout)
    return bot:request({ resp.update_quest, resp.show_quest_completion }, req.quest_pot_feed {
        slot = bot:slot(item) or 0,
        item_id = item,
        quest_id = POT_QUEST,
        exp = GROWTH,
    }, pred, timeout)
end

local function growth(p, name)
    return name == resp.update_quest and p.quest_status.quest_id == POT_QUEST
end

local function timer(ctx)
    local bot = ctx:bot(0)
    local items = string.format("%d:31,%d:1", FEED, OTHER_ETC)
    if pq.command(bot, "/봇초기화 50 0 0 " .. items .. " 3249:2 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    local p, name = bot:request({ resp.warp, resp.notice }, req.normal_chat { message = "/퀘스트시작준비 " .. QUEST },
        function(p, name)
            return name == resp.warp or p.message:find("NPC", 1, true) ~= nil
        end, 5000)
    if p == false then
        return ctx:fail("퀘스트 시작 준비 실패")
    end
    if talk(bot, SCRIPTED_START, 1) == false then
        return ctx:fail("타이머 퀘스트 시작 실패")
    end
    if (bot:items()[POT] or 0) ~= 1 then
        return ctx:fail("타이머의 알을 받지 못함")
    end

    if feed(bot, OTHER_ETC, growth, 500) ~= false or (bot:items()[OTHER_ETC] or 0) ~= 1 then
        return ctx:fail("먹이가 아닌 아이템이 먹여짐")
    end

    for i = 1, FULL / GROWTH - 1 do
        p = feed(bot, FEED, growth, 3000)
        if p == false or p.quest_status.status_record ~= tostring(i * GROWTH) then
            return ctx:fail(string.format("%d번째 먹이 기록이 다름: %s", i, p and p.quest_status.status_record or "응답 없음"))
        end
    end
    p = feed(bot, FEED, function(p, name)
        return name == resp.show_quest_completion and p.quest_id == QUEST
    end, 3000)
    if p == false then
        return ctx:fail("다 자란 뒤 3250 완료 표시가 오지 않음")
    end
    if (bot:items()[FEED] or 0) ~= 1 then
        return ctx:fail("먹이가 30개 소비되지 않음: " .. tostring(bot:items()[FEED]))
    end

    if feed(bot, FEED, growth, 500) ~= false or (bot:items()[FEED] or 0) ~= 1 then
        return ctx:fail("다 자란 타이머가 먹이를 더 먹음")
    end

    p = bot:request(resp.update_quest, req.quest_pot_open { quest_id = POT_QUEST }, function(p)
        return p.quest_status.quest_id == POT_QUEST
    end, 3000)
    if p == false or p.quest_status.status_record ~= tostring(FULL) then
        return ctx:fail("화분 열기 기록이 다름")
    end

    if talk(bot, SCRIPTED_END, 2) == false then
        return ctx:fail("타이머 퀘스트 완료 실패")
    end
    if (bot:items()[POT] or 0) ~= 0 then
        return ctx:fail("완료 뒤 타이머의 알이 남음")
    end
    return true
end

test_suite {
    name = "Quest pot",
    bot_count = 1,
    scenarios = { timer },
}
