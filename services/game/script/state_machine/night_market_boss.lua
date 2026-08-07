-- State machine (old/scripts/event/NightMarketBoss.js): 포장마차

local EXIT_MAP = 741020100
local BOSS_MAP = 741020101
local DURATION_MS = 1900000

local function end_run(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		local map = group:map(BOSS_MAP)
		if map ~= nil then
			map:reset()
			map:respawn(true)
		end
		return { BOSS_MAP }

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(BOSS_MAP)
	end,

	on_scheduled_timeout = function(sm)
		end_run(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}
