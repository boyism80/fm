local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local SUMMON_MAP = 100000000
local SNAIL = 100100
local SWORD = 1302000
local WEAPON_SLOT = -11
local FORGED_DAMAGE = 2000000000
local POINTS = 10
local THRESHOLD = 100
local OTHER_CHANNEL = 1
local OUTLAW = 521
local OCTOPUS = 5211001
local LAG_MS = 700
local PING_CYCLE_MS = 8000

local function command(bot, name, text)
    return bot:request(name, req.normal_chat { message = text })
end

local function forged_attack(oid)
    return req.attack {
        close_attack_info = {
            attack_header = { targets = 1, hits = 1 },
            damages = { { oid = oid, damage_pairs = { { damage = FORGED_DAMAGE } } } },
        },
    }
end

local function prepare(ctx, bot)
    if pq.command(bot, "/봇초기화 62 110 0 " .. SWORD .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:request(resp.inventory_operation, req.move_item { inventory_type = 1, source = bot:slot(SWORD) or 0, dest = WEAPON_SLOT, count = 1 }, nil, 3000) == false then
        return ctx:fail("검 착용 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    if command(bot, resp.spawn_mob, "/몬스터생성 " .. SNAIL) == false then
        return ctx:fail("달팽이 생성 응답 없음")
    end
    -- The spawn packet can land after the command's reply.
    for _ = 1, 10 do
        local mobs = bot:mobs(SNAIL)
        if #mobs > 0 then
            return mobs[#mobs].oid
        end
        ctx:sleep(300)
    end
    return ctx:fail("생성한 달팽이를 찾지 못함")
end

local function score_and_notice(ctx)
    local cheater = ctx:bot(0)
    local gm = ctx:bot(1)
    if gm:transfer(OTHER_CHANNEL) == false then
        return ctx:fail("GM 채널 이동 실패")
    end
    local oid = prepare(ctx, cheater)
    if oid == false then
        return false
    end

    for _ = 1, THRESHOLD / POINTS - 1 do
        cheater:request(resp.update_stats, forged_attack(oid), nil, 1000)
    end
    local notice = cheater:request_on(gm, resp.notice, forged_attack(oid), function(p)
        return p.message:find("[치트 의심]", 1, true) ~= nil and p.message:find(cheater:name(), 1, true) ~= nil
    end, 5000)
    if notice == false then
        return ctx:fail("기준에 닿았는데 GM에게 알림이 오지 않음")
    end

    local expected = string.format("clean=%d lagged=0 damage_exceeds=%d", THRESHOLD, THRESHOLD / POINTS)
    if pq.command(cheater, "/치트점수", expected) == false then
        return ctx:fail("조작한 피해가 점수로 쌓이지 않음 (기대: " .. expected .. ")")
    end
    return true
end

local function lagged_score(ctx)
    local bot = ctx:bot(2)
    local oid = prepare(ctx, bot)
    if oid == false then
        return false
    end

    bot:set_pong_delay(LAG_MS)
    ctx:sleep(PING_CYCLE_MS)
    for _ = 1, 3 do
        bot:request(resp.update_stats, forged_attack(oid), nil, 2000)
    end
    bot:set_pong_delay(0)

    local expected = "clean=0 lagged=30 damage_exceeds=3"
    if pq.command(bot, "/치트점수", expected) == false then
        return ctx:fail("느린 연결의 위반이 지연 의심으로 분류되지 않음 (기대: " .. expected .. ")")
    end
    return true
end

local function summon_pace(ctx)
    local outlaw = ctx:bot(3)
    if pq.command(outlaw, "/봇초기화 120 " .. OUTLAW .. " 0 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if pq.command(outlaw, "/스킬레벨 " .. OCTOPUS .. " 1", "스킬레벨 설정") == false or pq.command(outlaw, "/마력바꾸기 5000", "마력 설정: 5000") == false then
        return ctx:fail("스킬 준비 실패")
    end
    if outlaw:instance_move(SUMMON_MAP) == false then
        return ctx:fail("맵 이동 실패")
    end
    local summon = outlaw:request(resp.spawn_summon, req.active_skill { skill_id = OCTOPUS, skill_level = 1 },
        function(p)
            return p.owner_id == outlaw:id() and p.skill_id == OCTOPUS
        end, 5000)
    if summon == false then
        return ctx:fail("소환수가 나타나지 않음")
    end
    local spawned = outlaw:request(resp.spawn_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }, function(p)
        return p.mob.mob_id == SNAIL
    end, 5000)
    if spawned == false then
        return ctx:fail("달팽이 생성 응답 없음")
    end

    for _ = 1, 40 do
        outlaw:request(resp.summon_attack, req.summon_attack {
            summon_oid = summon.oid,
            damages = { { oid = spawned.mob.oid, damage_pairs = { { damage = 1 } } } },
        }, nil, 100)
    end
    if pq.command(outlaw, "/치트점수", "summon_rate=") == false then
        return ctx:fail("소환수 공격이 너무 빨라도 위반으로 쌓이지 않음")
    end
    return true
end

test_suite {
    name = "CheatTracker: 조작한 피해 점수·다른 채널 GM 알림·지연 의심 분류·소환수 공격 속도",
    bot_count = 4,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/즉사 on", "즉사 상태: enabled")
        pq.command(ctx:bot(2), "/즉사 on", "즉사 상태: enabled")
    end,

    scenarios = {
        score_and_notice,
        lagged_score,
        summon_pace,
    },
}
