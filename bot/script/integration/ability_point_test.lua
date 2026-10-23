local pq = require("script/integration/lib/party_quest")

local STR = 0x40
local DEX = 0x80
local AVAILABLE_AP = 0x4000
local MAX_STAT = 999

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function assign(bot, entries)
    return bot:request(resp.update_stats, req.auto_assign_ap { entries = entries }, function(p)
        return p.stats[AVAILABLE_AP] ~= nil
    end, 3000)
end

-- Two entries for the same stat count together against the stat cap before any AP is spent.
local function split_entries(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 0 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/힘바꾸기 990", "힘 설정: 990") == false or pq.command(bot, "/덱스바꾸기 4", "덱스 설정: 4") == false then
        return ctx:fail("능력치 설정 실패")
    end
    if pq.command(bot, "/AP바꾸기 12", "AP 설정: 12") == false then
        return ctx:fail("AP 설정 실패")
    end

    if assign(bot, { { stat = STR, amount = 6 }, { stat = STR, amount = 6 } }) ~= false then
        return ctx:fail("같은 능력치로 나눈 배분이 한도를 넘었는데 반영됨")
    end

    local p = assign(bot, { { stat = STR, amount = 9 }, { stat = DEX, amount = 3 } })
    if p == false then
        return ctx:fail("한도 안의 배분이 반영되지 않음")
    end
    return check(ctx, p.stats[STR] == MAX_STAT and p.stats[DEX] == 7 and p.stats[AVAILABLE_AP] == 0,
        string.format("배분 결과 힘 %s 덱스 %s AP %s (기대 %d, 7, 0)", tostring(p.stats[STR]), tostring(p.stats[DEX]),
            tostring(p.stats[AVAILABLE_AP]), MAX_STAT))
end

test_suite {
    name = "AbilityPoint: 같은 능력치 나눠 배분 한도",
    bot_count = 1,

    scenarios = {
        split_entries,
    },
}
