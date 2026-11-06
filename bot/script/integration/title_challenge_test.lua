local check = require("script/integration/lib/script_check")
local pq = require("script/integration/lib/party_quest")

local SCRIPTED_START = 4
local DIALOG_DEFAULT = 0
local CHALLENGE_NPC = 9000040
local NPC_X = 527
local NPC_Y = 454
local POPULARITY = 29002
local DILIGENT = 29003
local MAX_NOTICES = 5

local function status(bot, quest)
    local p = bot:request(resp.notice, req.normal_chat { message = "/퀘스트상태 " .. quest }, function(p)
        return p.message:find("퀘스트 " .. quest .. "[ :]") ~= nil
    end)
    if p == false then
        return nil
    end
    return p.message
end

-- Closing the notices of the popularity title challenge leaves it unstarted.
local function close_notices(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/퀘스트초기화 " .. POPULARITY, "퀘스트 " .. POPULARITY .. " ") == false then
        return ctx:fail("초기화 실패")
    end
    if bot:request(resp.warp, req.normal_chat { message = "/퀘스트시작준비 " .. POPULARITY }) == false then
        return ctx:fail("시작 준비 실패")
    end
    bot:move(NPC_X, NPC_Y)

    local dlg, name = bot:request(check.REPLIES,
        req.quest_action { mode = SCRIPTED_START, quest_id = POPULARITY, npcid = CHALLENGE_NPC }, check.replied, 5000)
    if dlg == false or name ~= resp.dialog then
        return ctx:fail("시작 안내 대화가 오지 않음: " .. tostring(name))
    end
    for _ = 1, MAX_NOTICES do
        if bot:request(resp.dialog, req.dialog { dialog_type = DIALOG_DEFAULT, next = false }, nil, 1000) == false then
            break
        end
    end

    local m = status(bot, POPULARITY)
    if m == nil then
        return ctx:fail("상태 응답 없음")
    end
    if m:sub(- #"진행") == "진행" then
        return ctx:fail("안내 대화를 닫았는데 도전이 시작됨: " .. m)
    end
    return true
end

-- The diligent explorer counts a day once today's play reaches 3 hours, sends today's seconds and the days, and its
-- 27th day shows the challenge complete.
local function diligent_explorer(ctx)
    local bot = ctx:bot(0)
    local records = string.format("quest.%d.play_time.today=10795,quest.%d.play_time.days=26", DILIGENT, DILIGENT)
    if pq.command(bot, string.format("/봇초기화 30 0 0 - %d:1 - %s", DILIGENT, records), "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    local completed = false
    ctx:hook("diligent_explorer", function(_, b, p, name)
        if b:id() == bot:id() and name == resp.show_quest_completion and p.quest_id == DILIGENT then
            completed = true
        end
    end)
    local p = bot:request(resp.update_quest_play_time, nil, function(p)
        return p.quest_id == DILIGENT
    end, 15000)
    if p == false then
        return ctx:fail("성실한 모험가 접속 시간 기록이 오지 않음")
    end
    if p.today < 10800 or p.days ~= 27 then
        return ctx:fail(string.format("성실한 모험가 기록이 다름: 오늘 %d초, %d일", p.today, p.days))
    end
    ctx:sleep(500)
    if not completed then
        return ctx:fail("27일째에 도전 완료 알림이 오지 않음")
    end
    return true
end

test_suite {
    name = "Title challenge: 안내 대화를 닫으면 도전이 시작되지 않음·성실한 모험가 접속 시간",
    bot_count = 1,

    scenarios = {
        close_notices,
        diligent_explorer,
    },
}
