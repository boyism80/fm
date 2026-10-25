-- State machine (old/scripts/event/DollHouse.js): 인형의 집

local pq = require("script/lib/party_quest")
local EXIT_MAP = 221024400
local QUEST_MAP = 922000010
local DURATION_MS = 600000

return {
    on_init = function(group)
        group:min_players(1)
        group:exit_map(EXIT_MAP)
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
        sm:finish(EXIT_MAP, 4)
    end,

}
