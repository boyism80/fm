-- State machine (old/scripts/event/AirStrike.js): 바트의 방

local pq = require("script/lib/party_quest")
local EXIT_MAP = 120000102
local QUEST_MAP = 912020000
local DURATION_MS = 300000

local function finish(sm)
    local players = sm:players()
    sm:finish(0)
    for _, player in ipairs(players) do
        player:map(EXIT_MAP, 1)
    end
end

return {
    on_init = function(group)
        group:min_players(1)
    end,

    on_create = function(sm)
        return { QUEST_MAP }
    end,

    on_prepare = function(sm)
        local map = sm:map(QUEST_MAP)
        map:respawn({ include_one_time = true })
        pq.shuffle_reactors(map)
    end,

    on_start = function(sm)
        sm:start_timer(DURATION_MS)
    end,

    on_player_enter = function(sm, player)
        player:map(sm:map(QUEST_MAP))
    end,

    on_scheduled_timeout = function(sm)
        finish(sm)
    end,

}
