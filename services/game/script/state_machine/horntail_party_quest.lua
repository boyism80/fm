-- State machine (old/scripts/event/HorntailPQ.js): 생명의 동굴 파티 퀘스트

local stage_maps = {
	240050100,
	240050101,
	240050102,
	240050103,
	240050104,
	240050105,
	240050200,
	240050300,
	240050310,
}

local FAIL_EXIT = 240050500
local CLEAR_EXIT = 240050400
local DURATION_MS = 1800000

local function clear(sm)
	local exit_map = FAIL_EXIT
	local allfinish = sm:get_property("allfinish")
	if allfinish ~= nil and allfinish ~= "" then
		exit_map = CLEAR_EXIT
	end
	sm:finish(exit_map)
end

local function in_stage(map_id)
	return map_id >= 240050100 and map_id <= 240050310
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps(stage_maps)
		group:declare_min_players(5)
		group:declare_exit_map(FAIL_EXIT)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		sm:set_property("stage1progress", "0")
		sm:set_property("allfinish", "")
		for _, map_id in ipairs(stage_maps) do
			local map = group:map(map_id)
			map:reset()
			map:respawn(true)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(stage_maps[1])
	end,

	on_mob_kill = function(sm, player, mobs)
	end,

	on_changed_map = function(sm, player, map_id)
		if in_stage(map_id) then
			return
		end
		sm:unregister(player)
	end,

	on_player_dead = function(sm, player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_player_disconnected = function(sm, player)
	end,

	on_left_party = function(sm, player)
	end,

	on_disband_party = function(sm)
		clear(sm)
	end,

	on_scheduled_timeout = function(sm)
		clear(sm)
	end,

	on_clear = function(sm)
		clear(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end,

	on_all_monsters_dead = function(sm)
	end,

	on_cancel_schedule = function(group)
	end
}
