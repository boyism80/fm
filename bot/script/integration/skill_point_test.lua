local pq = require("script/integration/lib/party_quest")

local WARRIOR = 100
local POWER_STRIKE = 1001004
local NOT_BY_SP = {
    1004,
    1007,
    1009,
    1010,
    1011,
}

local function sp(ctx, bot)
    local p = pq.command(bot, "/봇상태", "봇상태")
    if p == false then
        ctx:fail("봇상태 응답 없음")
        return nil
    end
    return tonumber(p.message:match("sp=(%d+)"))
end

local function distribute(ctx, bot, skill_id)
    if bot:request(resp.update_stats, req.distribute_sp { skill_id = skill_id }) == false then
        ctx:fail("SP 분배 응답 없음: " .. skill_id)
        return false
    end
    return true
end

local function beginner_skills_rejected(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 30 " .. WARRIOR .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/스킬포인트 3", "스킬포인트 설정") == false then
        return ctx:fail("SP 설정 실패")
    end

    for _, skill_id in ipairs(NOT_BY_SP) do
        if distribute(ctx, bot, skill_id) == false then
            return false
        end
    end
    local left = sp(ctx, bot)
    if left ~= 3 then
        return ctx:fail("SP로 배울 수 없는 초보자 스킬에 SP가 쓰임: sp=" .. tostring(left))
    end

    if distribute(ctx, bot, POWER_STRIKE) == false then
        return false
    end
    left = sp(ctx, bot)
    if left ~= 2 then
        return ctx:fail("1차 스킬에 SP가 쓰이지 않음: sp=" .. tostring(left))
    end
    return true
end

test_suite {
    name = "SkillPoint: SP로 배울 수 없는 초보자 스킬 거절",
    bot_count = 1,

    scenarios = {
        beginner_skills_rejected,
    },
}
