-- State machine (old/scripts/event/ShanghaiBoss.js): 대왕지네

local EXIT_MAP = 701010321
local ENTRY_MAP = 701010322
local BOSS_MAP = 701010323
local SIDE_MAP = 701010324
local BOSS_MOB = 9600009
local BOSS_SPAWN_X = 2176
local BOSS_SPAWN_Y = 823
local DURATION_MS = 600000000

local STAGE_MAPS = {
	ENTRY_MAP,
	BOSS_MAP,
	SIDE_MAP,
}

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
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
			end
		end
		local boss_map = group:map(BOSS_MAP)
		if boss_map ~= nil then
			boss_map:kill_all_mobs()
			boss_map:spawn_mob(BOSS_MOB, BOSS_SPAWN_X, BOSS_SPAWN_Y)
		end
		return STAGE_MAPS

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(ENTRY_MAP)
	end,

	on_scheduled_timeout = function(sm)
		end_run(sm)
	end,

	on_clear = function(sm)
		end_run(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}
