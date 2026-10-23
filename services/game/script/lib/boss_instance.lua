local M = {}

-- config: map, exit_map, duration_ms, bosses = { { id, x, y }, ... } spawned one after another,
-- on_cleared(sm, map) after the last boss falls.
function M.create(config)
    local function summon(sm, phase)
        local boss = config.bosses[phase]
        sm:set_property("phase", tostring(phase))
        sm:map(config.map):spawn_mob(boss.id, boss.x, boss.y)
    end

    return {
        on_init = function(group)
            group:min_players(1)
            group:exit_map(config.exit_map)
            group:max_machines(1)
        end,

        on_create = function(sm)
            return { config.map }
        end,

        on_prepare = function(sm)
            sm:map(config.map):respawn({ include_one_time = true })
            summon(sm, 1)
        end,

        on_start = function(sm)
            sm:start_timer(config.duration_ms)
        end,

        on_player_enter = function(sm, player)
            player:map(sm:map(config.map), 0)
        end,

        on_mob_die = function(sm, mob)
            local phase = tonumber(sm:get_property("phase")) or 0
            local boss = config.bosses[phase]
            if boss == nil or mob:id() ~= boss.id then
                return
            end
            if phase < #config.bosses then
                summon(sm, phase + 1)
                return
            end
            sm:set_property("phase", tostring(phase + 1))
            if config.on_cleared ~= nil then
                config.on_cleared(sm, sm:map(config.map))
            end
        end,

        on_scheduled_timeout = function(sm)
            sm:finish(config.exit_map)
        end,
    }
end

return M
