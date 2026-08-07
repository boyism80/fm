local EXIT_MAP = 120000101
local DURATION_MS = 600000

local function create(map_id)
	local function finish(sm)
		local players = sm:players()
		sm:finish(0)
		for _, player in ipairs(players) do
			player:map(EXIT_MAP)
		end
	end

	return {
		on_init = function(group)
			group:set_property("state", "0")
			group:declare_maps({ map_id })
			group:declare_min_players(1)
		end,

		on_create = function(sm)
			local group = sm:group()
			group:set_property("state", "1")
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
			end
		end,

		on_start = function(sm)
			sm:start_timer(DURATION_MS)
		end,

		on_player_enter = function(sm, player)
			player:map(map_id)
		end,

		on_changed_map = function(sm, player, entered_map_id)
			if entered_map_id == map_id then
				return
			end
			sm:finish(0)
		end,

		on_player_disconnected = function(sm, player)
			sm:finish(0)
		end,

		on_scheduled_timeout = function(sm)
			finish(sm)
		end,

		on_finish = function(sm)
			sm:group():set_property("state", "0")
		end
	}
end

return create
