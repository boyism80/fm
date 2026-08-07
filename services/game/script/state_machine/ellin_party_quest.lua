-- State machine (old/scripts/event/Ellin.js): 독안개의 숲

local pq = require("script/lib/party_quest")

local EXIT_MAP = 930000800
local START_MAP = 930000000
local DURATION_MS = 1800000
local RANKING_QUEST = 1206
local STAGE1_MOB = 9300172
local STAGE1_MAP = 930000100

local STAGE_MAPS = {
	930000000,
	930000010,
	930000100,
	930000200,
	930000300,
	930000400,
	930000500,
	930000600,
	930000700,
}

local function end_run(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:set_property("stage1_cleareff", "0")
		group:declare_maps(STAGE_MAPS)
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("stage1_cleareff", "0")
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
				if map_id == 930000500 then
					pq.shuffle_reactors(map)
				end
			end
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(START_MAP)
		player:try_party_quest(RANKING_QUEST)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id >= 930000000 and map_id <= 930000700 then
			return
		end
		sm:unregister(player)
	end,

	on_mob_die = function(sm, mob)
		if mob == nil or mob:id() ~= STAGE1_MOB then
			return
		end
		local group = sm:group()
		if group == nil or group:get_property("stage1_cleareff") == "1" then
			return
		end
		local map = group:map(STAGE1_MAP)
		if map == nil then
			return
		end
		if pq.mob_count(map) == 0 then
			map:clear_effect()
			group:set_property("stage1_cleareff", "1")
		end
	end,

	on_left_party = function(sm, player)
		end_run(sm)
	end,

	on_disband_party = function(sm)
		end_run(sm)
	end,

	on_scheduled_timeout = function(sm)
		end_run(sm)
	end,

	on_clear = function(sm)
		end_run(sm)
	end,

	on_finish = function(sm)
		local group = sm:group()
		group:set_property("state", "0")
		group:set_property("stage1_cleareff", "0")
	end
}
