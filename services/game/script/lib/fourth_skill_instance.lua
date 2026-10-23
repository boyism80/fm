local function create(config)
    local function finish(sm)
        sm:finish(config.exit_map)
    end

    return {
        on_init = function(group)
            group:set_property("started", "false")
            group:min_players(config.min_players or 1)
            group:exit_map(config.exit_map)
            if config.max_machines ~= nil then
                group:max_machines(config.max_machines)
            end
        end,

        on_create = function(sm)
            local group = sm:group()
            group:set_property("started", "true")
            return { config.map }
        end,

        on_prepare = function(sm)
            sm:map(config.map):respawn({ include_one_time = true })
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

        on_finish = function(sm)
            sm:group():set_property("started", "false")
        end
    }
end

return create
