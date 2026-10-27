local pq = require("script/integration/lib/party_quest")

local FIELD = 100000000
local SNAIL = 100100
local SWORD = 1302000
local WEAPON_SLOT = -11

local function prepare(ctx, bot)
    if pq.command(bot, "/봇초기화 30 100 0 " .. SWORD .. ":1 -", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:request(resp.inventory_operation, req.move_item { inventory_type = 1, source = bot:slot(SWORD) or 0, dest = WEAPON_SLOT, count = 1 }, nil, 3000) == false then
        return ctx:fail("검 착용 실패")
    end
    if pq.command(bot, "/즉사 off", "즉사 상태: disabled") == false then
        return ctx:fail("즉사 끄기 실패")
    end
    return true
end

local function attack(oid)
    return req.attack {
        close_attack_info = {
            attack_header = { targets = 1, hits = 1 },
            damages = { { oid = oid, damage_pairs = { { damage = 1 } } } },
        },
    }
end

local function control_of(oid, aggro)
    return function(p)
        return p.mob.oid == oid and p.aggro == aggro
    end
end

local function stop_of(oid)
    return function(p)
        return p.oid == oid
    end
end

-- The attacker takes control of the mob with aggro; a non-controller's move is refused unless the mob chases it.
local function attacker_takes_control(ctx)
    local first = ctx:bot(0)
    local second = ctx:bot(1)
    if prepare(ctx, first) == false or prepare(ctx, second) == false then
        return false
    end
    if first:instance_move(FIELD) == false then
        return ctx:fail("첫 봇 맵 이동 실패")
    end
    local started = first:request(resp.start_control_mob, req.normal_chat { message = "/몬스터생성 " .. SNAIL }, nil, 3000)
    if started == false then
        return ctx:fail("생성한 달팽이의 컨트롤이 첫 봇에게 오지 않음")
    end
    local oid = started.mob.oid
    if started.aggro then
        return ctx:fail("먼저 공격하지 않는 달팽이가 어그로로 배정됨")
    end
    if second:instance_move(FIELD) == false then
        return ctx:fail("두 번째 봇 맵 이동 실패")
    end

    if second:request_on(first, resp.stop_control_mob, attack(oid), stop_of(oid), 3000) == false then
        return ctx:fail("다른 봇이 때렸는데 컨트롤러가 컨트롤을 잃지 않음")
    end
    if first:request(resp.stop_control_mob, req.move_mob { oid = oid, center_split = -1 }, stop_of(oid), 3000) == false then
        return ctx:fail("컨트롤러가 아닌 봇의 이동에 컨트롤 중지가 오지 않음")
    end
    if first:request(resp.start_control_mob, req.move_mob { oid = oid, center_split = -1, aggro = true }, control_of(oid, true), 3000) == false then
        return ctx:fail("몹이 쫓는 봇의 이동인데 컨트롤이 넘어오지 않음")
    end

    if second:request(resp.start_control_mob, req.aggro { oid = oid }, nil, 1000) ~= false then
        return ctx:fail("컨트롤러가 맵에 있는데 자동 어그로로 컨트롤이 넘어감")
    end
    if second:request(resp.start_control_mob, attack(oid), control_of(oid, true), 3000) == false then
        return ctx:fail("때린 봇에게 어그로 컨트롤이 오지 않음")
    end
    return true
end

test_suite {
    name = "Mob aggro: 공격자 컨트롤 전환·어그로 이동·자동 어그로",
    bot_count = 2,

    on_finished = function(ctx)
        pq.command(ctx:bot(0), "/즉사 on", "즉사 상태: enabled")
        pq.command(ctx:bot(1), "/즉사 on", "즉사 상태: enabled")
    end,

    scenarios = {
        attacker_takes_control,
    },
}
