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
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		return { BOSS_MAP }
	end,

	on_prepare = function(sm)
		sm:map(BOSS_MAP):respawn({ include_one_time = true })
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(BOSS_MAP))
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
