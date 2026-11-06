local pq = require("script/integration/lib/party_quest")

local HUNTING_GROUND = 104040000
local HIT_ENV = -3
local HP = 0x400
local WARRIOR = 100
local SAFETY_CHARM = 5130000
local START_EXP = 10000
-- Level 30 needs 95700 EXP, and a warrior with 20 LUK loses 0.2 / 20 + 5% of it on a field.
local FIELD_LOSS = 5742

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

local function setup(ctx, bot, items)
    local commands = {
        { "/봇초기화 30 " .. WARRIOR .. " 0 " .. items, "봇초기화 완료" },
        { "/럭바꾸기 20", "럭 설정" },
        { "/경험치얻기 " .. START_EXP, "경험치 " .. START_EXP .. " 획득" },
    }
    for _, c in ipairs(commands) do
        if pq.command(bot, c[1], c[2]) == false then
            return ctx:fail(c[1] .. " 실패")
        end
    end
    if bot:instance_move(HUNTING_GROUND) == false then
        return ctx:fail("사냥터 이동 실패")
    end
    return true
end

local function die(ctx, bot)
    local dead = bot:request(resp.update_stats, req.damaged { type = HIT_ENV, damage = 30000 }, function(p)
        return p.stats[HP] == 0
    end, 5000)
    if dead == false then
        return ctx:fail("죽을 만큼 맞았는데 HP가 0이 되지 않음")
    end
    ctx:sleep(300)
    return true
end

local function revive(ctx, bot)
    if bot:relog { passwords = { "" } } == false then
        return ctx:fail("재접속 실패")
    end
    return true
end

-- Dying on a field takes the job's share of the level's EXP.
local function field_loss(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, "-") == false or die(ctx, bot) == false then
        return false
    end
    local exp = bot:exp()
    if check(ctx, exp == START_EXP - FIELD_LOSS, string.format("사망 경험치가 다름: %d (기대 %d)", exp, START_EXP - FIELD_LOSS)) == false then
        return false
    end
    return revive(ctx, bot)
end

-- A safety charm is used up instead of the EXP.
local function safety_charm(ctx)
    local bot = ctx:bot(0)
    if setup(ctx, bot, SAFETY_CHARM .. ":2") == false or die(ctx, bot) == false then
        return false
    end
    local exp, charms = bot:exp(), bot:items()[SAFETY_CHARM] or 0
    if check(ctx, exp == START_EXP and charms == 1, string.format("호신부적 사망 결과가 다름: 경험치 %d 부적 %d", exp, charms)) == false then
        return false
    end
    return revive(ctx, bot)
end

test_suite {
    name = "Death: 사망 경험치·호신부적",
    bot_count = 1,

    scenarios = {
        field_loss,
        safety_charm,
    },
}
