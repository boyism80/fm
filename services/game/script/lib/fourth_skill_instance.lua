local pq = require("script/lib/party_quest")

local function create(config)
    local function finish(sm)
        sm:finish(config.exit_map)
    end

    return {
        on_init = function(group)
            if config.leave_on_party then
                group:leave_on_party()
            end
            group:min_players(config.min_players or 1)
            group:exit_map(config.exit_map)
        end,

        on_create = function(sm)
            return { config.map }
        end,

        on_prepare = function(sm)
            local map = sm:map(config.map)
            map:respawn({ include_one_time = true })
            if config.shuffle_reactors then
                pq.shuffle_reactors(map)
            end
        end,

        on_start = function(sm)
            sm:start_timer(config.duration_ms)
        end,

        on_player_enter = function(sm, player)
            player:map(sm:map(config.map))
        end,

        on_disband_party = function(sm)
            finish(sm)
        end,

        on_scheduled_timeout = function(sm)
            finish(sm)
        end,
    }
end

return create
