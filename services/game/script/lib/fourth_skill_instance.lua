local function create(config)
	local function finish(sm)
		sm:finish(config.exit_map)
	end

	return {
		on_init = function(group)
			group:set_property("started", "false")
			group:declare_min_players(1)
			group:declare_exit_map(config.exit_map)
		end,

		on_create = function(sm)
			local group = sm:group()
			group:set_property("started", "true")
			local map = group:map(config.map)
			if map ~= nil then
				map:reset()
				map:respawn(true)
			end
			return { config.map }
		end,

		on_start = function(sm)
			sm:start_timer(config.duration_ms)
		end,

		on_player_enter = function(sm, player)
			player:map(config.map)
		end,

		on_disband_party = function(sm)
			finish(sm)
		end,

		on_scheduled_timeout = function(sm)
			finish(sm)
		end,

		on_clear = function(sm)
			finish(sm)
		end,

		on_finish = function(sm)
			sm:group():set_property("started", "false")
		end
	}
end

return create
