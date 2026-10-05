-- State machine (old/scripts/event/tokyo_2095.js): 2095 tokyo

local EXIT_MAP = 802000312
local TEMPLATE_MAP = 802000311
local DURATION_MS = 1200000

local function finish(sm, exit_map)
	sm:finish(exit_map or EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:set_property("leader", "true")
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("leader", "true")
		return { TEMPLATE_MAP }
	end,

	on_prepare = function(sm)
		sm:map(TEMPLATE_MAP):kill_all_mobs()
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(TEMPLATE_MAP), 0)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm, EXIT_MAP)
	end,

	on_clear = function(sm)
		finish(sm, 0)
	end,

	on_finish = function(sm)
		local group = sm:group()
		group:set_property("state", "0")
		group:set_property("leader", "true")
	end,
}
