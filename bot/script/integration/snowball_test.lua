local pq = require("script/integration/lib/party_quest")
local event = require("script/integration/lib/channel_event")

local SNOWBALL_MAP = 109060000
local STATE_MOVE = 4
local BOTTOM = 0
-- The bottom snowball starts between x 175 and 450 on the track at y 155.
local BESIDE_BALL = { x = 300, y = 155 }
local AWAY_FROM_BALL = { x = 800, y = 155 }

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

-- A hit beside the snowball rolls it one step; a hit away from it does nothing; running into it knocks the bot back.
local function snowball(ctx)
    local bot = ctx:bot(0)
    if event.join(ctx, bot, "SNOWBALL", "눈덩이 굴리기", SNOWBALL_MAP) == false then
        return false
    end
    if event.start(ctx, bot, "눈덩이 굴리기", resp.snowball_state, nil) == false then
        return false
    end

    if pq.move(bot, AWAY_FROM_BALL.x, AWAY_FROM_BALL.y) == false then
        return ctx:fail("눈덩이에서 먼 곳으로 이동 실패")
    end
    local missed = bot:request(resp.snowball_hit, req.snowball_hit { target = BOTTOM, damage = 10, delay = 0 }, nil, 3000)
    -- One hit in a hundred rolls the snowball wherever it lands, so only the answer is checked here.
    if check(ctx, missed ~= false, "눈덩이에서 먼 공격에 응답이 없음") == false then
        return false
    end

    if pq.move(bot, BESIDE_BALL.x, BESIDE_BALL.y) == false then
        return ctx:fail("눈덩이 옆으로 이동 실패")
    end
    local rolled = bot:request(resp.snowball_state, req.snowball_hit { target = BOTTOM, damage = 10, delay = 0 },
        function(p)
            return p.state == STATE_MOVE and p.position[1] >= 1
        end, 3000)
    if check(ctx, rolled ~= false, "눈덩이 옆에서 공격했는데 굴러가지 않음") == false then
        return false
    end

    local knocked = bot:request(resp.snowball_knockback, req.snowball_touch {}, nil, 3000)
    event.cancel(bot)
    return check(ctx, knocked ~= false, "눈덩이에 닿았는데 밀려나지 않음")
end

test_suite {
    name = "Snowball: 눈덩이 굴리기 채널 이벤트 (시작·타격 판정·밀려남)",
    bot_count = 1,

    on_finished = function(ctx)
        event.cancel(ctx:bot(0))
    end,

    scenarios = {
        snowball,
    },
}
