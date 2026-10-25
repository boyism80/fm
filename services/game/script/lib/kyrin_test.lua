local EXIT_MAP = 120000101
local DURATION_MS = 600000

local function create(map_id)
    return {
        on_init = function(group)
            group:min_players(1)
        end,

        on_create = function(sm)
            return { map_id }
        end,

        on_prepare = function(sm)
            sm:map(map_id):respawn({ include_one_time = true })
        end,

        on_start = function(sm)
            sm:start_timer(DURATION_MS)
        end,

        on_player_enter = function(sm, player)
            player:map(sm:map(map_id))
        end,

        on_scheduled_timeout = function(sm)
            sm:finish(EXIT_MAP)
        end,

    }
end

return create
