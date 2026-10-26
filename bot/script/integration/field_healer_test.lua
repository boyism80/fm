local pq = require("script/integration/lib/party_quest")

local THRONE = 990000900
local TOP = -564
local FALL = 10
local MAX_HP = 5000
local HIT_ENV = -3
local DAMAGE = 3000
local PULL_INTERVAL_MS = 550

local function pull(bot)
    return bot:request(resp.healer_move, req.pulley_pull {}, nil, 3000)
end

-- Each pull lowers Ergoth's healer, it climbs back on its own, and at the bottom it heals the field and returns to the top.
local function healer(ctx)
    local bot = ctx:bot(0)
    if pq.command(bot, "/봇초기화 70 0 0", "봇초기화 완료") == false then
        return ctx:fail("봇 초기화 실패")
    end
    if bot:force_instance_move(THRONE, "에레고스의 왕좌 힐러") == false then
        return ctx:fail("에레고스의 왕좌 이동 실패")
    end

    local first = pull(bot)
    if first == false or first.y ~= TOP + FALL then
        return ctx:fail("도르래를 당겼는데 힐러가 내려오지 않음")
    end
    for _ = 1, 22 do
        bot:send(req.pulley_pull {})
    end
    local rushed = bot:request(resp.healer_move, nil, function(p)
        return p.y > first.y
    end, 800)
    if rushed ~= false then
        return ctx:fail("간격 없이 연타한 도르래에 힐러가 더 내려옴: " .. rushed.y)
    end
    local rose = bot:request(resp.healer_move, nil, nil, 3000)
    if rose == false or rose.y >= first.y then
        return ctx:fail("내려온 힐러가 다시 올라가지 않음")
    end

    if pq.command(bot, "/체력바꾸기 " .. MAX_HP, "체력 설정") == false then
        return ctx:fail("체력 설정 실패")
    end
    if bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = DAMAGE }, nil, 3000) == false then
        return ctx:fail("체력을 깎지 못함")
    end
    local hurt = bot:hp()
    for _ = 1, 60 do
        ctx:sleep(PULL_INTERVAL_MS)
        if pull(bot) == false then
            return ctx:fail("도르래 응답이 오지 않음")
        end
        if bot:hp() > hurt then
            return true
        end
    end
    return ctx:fail("간격을 두고 도르래를 계속 당겼는데 힐러가 회복시키지 않음: " .. bot:hp())
end

test_suite {
    name = "Field healer: 에레고스의 왕좌 도르래·힐러",
    bot_count = 1,

    scenarios = {
        healer,
    },
}
