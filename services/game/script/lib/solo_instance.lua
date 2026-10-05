local M = {}

function M.create(config)
	return {
		on_init = function(group)
			group:min_players(1)
			group:exit_map(config.exit_map)
			group:max_machines(nil)
		end,

		on_create = function(sm)
			return config.maps
		end,

		on_prepare = function(sm)
			if config.setup == nil then
				return
			end
			local maps = {}
			for i, template_id in ipairs(config.maps) do
				maps[i] = sm:map(template_id)
			end
			config.setup(sm, maps)
		end,

		on_start = function(sm)
			sm:start_timer(config.duration_ms)
		end,

		on_player_enter = function(sm, player)
			player:map(sm:map(config.maps[1]), 0)
		end,

		on_scheduled_timeout = function(sm)
			if config.timeout_message ~= nil then
				sm:message(config.timeout_message, Msg.LightBlueText)
			end
			sm:finish(config.exit_map, config.exit_portal)
		end,
	}
end

return M
