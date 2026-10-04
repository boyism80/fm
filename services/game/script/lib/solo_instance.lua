local M = {}

function M.find_map(sm, template_id)
	for _, map in ipairs(sm:maps()) do
		if map:template_id() == template_id then
			return map
		end
	end
	return nil
end

function M.create(config)
	local links = config.links or {}

	return {
		on_init = function(group)
			group:declare_min_players(1)
			group:declare_exit_map(config.exit_map)
		end,

		on_create = function(sm)
			local maps = {}
			for _, template_id in ipairs(config.maps) do
				local map, err = id2map(template_id):create_instance()
				if map == nil then
					log(sm:group():name() .. " create_instance:", err)
					for _, created in ipairs(maps) do
						created:destroy()
					end
					sm:finish(0)
					return
				end
				maps[#maps + 1] = map
			end
			for _, map in ipairs(maps) do
				for _, name in ipairs(links[map:template_id()] or {}) do
					map:portal(name):script("instance_link")
				end
			end
			if config.setup ~= nil then
				config.setup(sm, maps)
			end
			return maps
		end,

		on_start = function(sm)
			sm:start_timer(config.duration_ms)
		end,

		on_player_enter = function(sm, player)
			player:map(M.find_map(sm, config.maps[1]), 0)
		end,

		on_scheduled_timeout = function(sm)
			if config.timeout_message ~= nil then
				sm:message(config.timeout_message, Msg.LightBlueText)
			end
			sm:finish(config.exit_map, config.exit_portal)
		end,

		on_finish = function(sm)
			for _, map in ipairs(sm:maps()) do
				map:destroy()
			end
		end,
	}
end

return M
