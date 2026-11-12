local function create(config)
    return {
        on_init = function(group)
            group:leave_on_party()
            group:min_players(1)
        end,

        on_create = function(sm)
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
            sm:message("카이린의 공격으로부터 2분 이상 버티세요.", Msg.PinkText)
        end,

        on_scheduled_timeout = function(sm)
            sm:finish(config.exit_map)
        end,
    }
end

return create
