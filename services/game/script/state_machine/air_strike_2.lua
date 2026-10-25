-- State machine (old/scripts/event/AirStrike2.js): 샤레니안 성문

local EXIT_MAP = 120000102
local QUEST_MAP = 910210000
local DURATION_MS = 900000

return {
    on_init = function(group)
        group:min_players(1)
    end,

    on_create = function(sm)
        return { QUEST_MAP }
    end,

    on_start = function(sm)
        sm:start_timer(DURATION_MS)
    end,

    on_player_enter = function(sm, player)
        player:map(sm:map(QUEST_MAP))
    end,

    on_scheduled_timeout = function(sm)
        sm:finish(EXIT_MAP, 1)
    end,

}
