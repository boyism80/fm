-- State machine (old/scripts/event/ElementThanatos.js): 속성의 타나토스

local EXIT_MAP = 220050300
local BOSS_MAP = 922020100
local DURATION_MS = 1200000

local function finish(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps({ BOSS_MAP })
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		local map = group:map(BOSS_MAP)
		if map ~= nil then
			map:respawn(true)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(BOSS_MAP)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == BOSS_MAP then
			return
		end
		sm:unregister(player)
	end,

	on_left_party = function(sm, player)
		sm:unregister(player)
		player:map(EXIT_MAP)
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
		sm:group():set_property("state", "0")
	end
}
