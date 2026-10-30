local pq = require("script/integration/lib/party_quest")

local LEVEL = 0x10
local EXP = 0x10000
local WARRIOR = 100

-- Exp past the max level is dropped, not refused, so a large gain at 199 reaches 200.
local function reach_max_level(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 199 " .. WARRIOR .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    local p = bot:request(resp.update_stats, req.normal_chat { message = "/경험치얻기 4000000000" }, function(p)
        return p.stats[LEVEL] ~= nil
    end, 3000)
    if p == false then
        return ctx:fail("남은 양보다 큰 경험치로 레벨이 오르지 않음")
    end
    if p.stats[LEVEL] ~= 200 or p.stats[EXP] ~= 0 then
        return ctx:fail(string.format("레벨 %s 경험치 %s (기대 200, 0)", tostring(p.stats[LEVEL]), tostring(p.stats[EXP])))
    end
    return true
end

test_suite {
    name = "MaxLevel: 레벨 199 넘치는 경험치",
    bot_count = 1,

    scenarios = {
        reach_max_level,
    },
}
