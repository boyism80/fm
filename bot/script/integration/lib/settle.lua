local pq = require("script/integration/lib/party_quest")

local M = {}

-- Matches settle_interval_seconds in config/game.yaml, which the bot servers are generated from.
local INTERVAL_MS = 1000
local MARGIN_MS   = 2000
local POLL_MS     = 500
local DEADLINE_MS = 30000

-- Sleeps through at least one run of the settle timer.
function M.tick(ctx)
    ctx:sleep(INTERVAL_MS + MARGIN_MS)
end

function M.wait(ctx, check)
    local waited = 0
    while check() ~= true do
        if waited >= DEADLINE_MS then
            return false
        end
        ctx:sleep(POLL_MS)
        waited = waited + POLL_MS
    end
    return true
end

-- The timer is per game server, so pause and resume through a bot on the server that holds the request.
function M.pause(bot)
    return pq.command(bot, "/정산타이머 off", "정산 타이머: OFF") ~= false
end

function M.resume(bot)
    return pq.command(bot, "/정산타이머 on", "정산 타이머: ON") ~= false
end

-- For on_scenario_finished: a failed scenario must not leave the timer off on any server its bots are on.
function M.resume_all(ctx, count)
    for i = 0, count - 1 do
        ctx:bot(i):send(req.normal_chat { message = "/정산타이머 on" })
    end
end

return M
