local event = require("script/integration/lib/channel_event")

local COCONUT_MAP = 109080000
local COKE_MAP = 109080010
local STATE_RESET = 0
local STATE_HIT = 1
-- Coconuts struck one after another, so none waits on its one-second cooldown.
local TARGETS = 40
local HITS_EACH = 6

local function check(ctx, ok, message)
    if ok == false then
        return ctx:fail(message)
    end
    return true
end

-- Every coconut is reset at the start, a hit shakes one, and a coconut that falls scores for the hitter's team.
local function play(ctx, key, name, map_id)
    local bot = ctx:bot(0)
    if event.join(ctx, bot, key, name, map_id) == false then
        return false
    end
    local reset = event.start(ctx, bot, name, resp.coconut_hit, function(p)
        return p.id < 0 and p.state == STATE_RESET
    end)
    if reset == false then
        return false
    end

    local shaken = bot:request(resp.coconut_hit, req.coconut_hit { id = 0, delay = 0 }, function(p)
        return p.id == 0 and p.state == STATE_HIT
    end, 3000)
    if check(ctx, shaken ~= false, "때린 코코넛이 흔들리지 않음") == false then
        return false
    end

    local scored = false
    for _ = 1, HITS_EACH do
        for id = 1, TARGETS do
            bot:send(req.coconut_hit { id = id, delay = 0 })
        end
        local p = bot:request(resp.coconut_score, nil, function(p)
            return p.maple + p.story > 0
        end, 1200)
        if p ~= false then
            scored = true
            break
        end
    end
    event.cancel(bot)
    return check(ctx, scored, "코코넛이 떨어졌는데 점수가 오르지 않음")
end

test_suite {
    name = "Coconut: 코코넛 시즌·코크 플레이 채널 이벤트 (시작·타격·점수)",
    bot_count = 1,

    on_finished = function(ctx)
        event.cancel(ctx:bot(0))
    end,

    scenarios = {
        function(ctx)
            return play(ctx, "COCONUT", "코코넛 시즌", COCONUT_MAP)
        end,
        function(ctx)
            return play(ctx, "COKE", "코크 플레이", COKE_MAP)
        end,
    },
}
