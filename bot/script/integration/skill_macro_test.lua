local pq = require("script/integration/lib/party_quest")

local WARRIOR = 100
local POWER_STRIKE = 1001004
local SLASH_BLAST = 1001005

local function macros(ctx, bot, channel, label)
    if bot:transfer(channel) == false then
        ctx:fail(label .. ": 채널 이동 실패")
        return nil
    end
    return bot:skill_macros()
end

local function describe(list)
    local parts = {}
    for _, macro in ipairs(list) do
        table.insert(parts,
            string.format("%s/%s/%d,%d,%d", macro.name, tostring(macro.shout), macro.skills[1], macro.skills[2],
                macro.skills[3]))
    end
    return "[" .. table.concat(parts, " ") .. "]"
end

local function save_and_reload(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 10 " .. WARRIOR .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(bot, "/스킬레벨 " .. POWER_STRIKE .. " 1", "스킬레벨 설정") == false then
        return ctx:fail("파워 스트라이크 배우기 실패")
    end

    local list = bot:skill_macros()
    if #list ~= 0 then
        return ctx:fail("새 캐릭터에 매크로가 있음: " .. describe(list))
    end

    if bot:send(req.skill_macro {
            count = 2,
            entries = {
                { name = "공격", shout = 1, skills = { POWER_STRIKE } },
                { name = "empty", shout = 0, skills = {} },
            },
        }) == false then
        return ctx:fail("매크로 저장 요청 실패")
    end
    list = macros(ctx, bot, 1, "저장 후")
    if list == nil then
        return false
    end
    if #list ~= 2
        or list[1].name ~= "공격" or list[1].shout ~= true
        or list[1].skills[1] ~= POWER_STRIKE or list[1].skills[2] ~= 0 or list[1].skills[3] ~= 0
        or list[2].name ~= "empty" or list[2].shout ~= false or list[2].skills[1] ~= 0 then
        return ctx:fail("채널 이동 후 매크로가 다름: " .. describe(list))
    end

    if bot:send(req.skill_macro {
            count = 1,
            entries = { { name = "bad", shout = 0, skills = { SLASH_BLAST } } },
        }) == false then
        return ctx:fail("잘못된 매크로 요청 실패")
    end
    list = macros(ctx, bot, 0, "배우지 않은 스킬 저장 후")
    if list == nil then
        return false
    end
    if #list ~= 2 or list[1].name ~= "공격" then
        return ctx:fail("배우지 않은 스킬로 매크로가 바뀜: " .. describe(list))
    end

    if bot:send(req.skill_macro {
            count = 1,
            entries = { { name = "a\0b", shout = 0, skills = { POWER_STRIKE } } },
        }) == false then
        return ctx:fail("제어 문자 이름 매크로 요청 실패")
    end
    list = macros(ctx, bot, 1, "제어 문자 이름 저장 후")
    if list == nil then
        return false
    end
    if #list ~= 2 or list[1].name ~= "공격" then
        return ctx:fail("제어 문자 이름으로 매크로가 바뀜: " .. describe(list))
    end

    if bot:send(req.skill_macro { count = 0, entries = {} }) == false then
        return ctx:fail("매크로 비우기 요청 실패")
    end
    list = macros(ctx, bot, 0, "비운 후")
    if list == nil then
        return false
    end
    if #list ~= 0 then
        return ctx:fail("매크로가 비워지지 않음: " .. describe(list))
    end
    return true
end

test_suite {
    name = "SkillMacro: 저장·불러오기·배우지 않은 스킬 거절",
    bot_count = 1,

    scenarios = {
        save_and_reload,
    },
}
