-- State machine (old/scripts/event/KerningPQ.js): 커닝시티 파티 퀘스트

local stage_maps = {
	103000800,
	103000801,
	103000802,
	103000803,
	103000804,
	103000805,
}

local exit_map_id = 103000890
local ranking_quest_id = 1201
local duration_ms = 1800000

local function clear(sm)
	sm:finish(exit_map_id)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps(stage_maps)
		group:declare_min_players(2)
		group:declare_exit_map(exit_map_id)
	end,

	on_mob_kill = function(sm, player, mobs)
		sm:add_kill(player, #mobs)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		for _, map_id in ipairs(stage_maps) do
			local map = group:map(map_id)
			map:reset()
			map:respawn(true)
			local portal = map:portal("next00")
			if portal ~= nil then
				portal:script("enter_kpq")
			end
		end
	end,

	on_start = function(sm)
		sm:start_timer(duration_ms)
	end,

	on_player_enter = function(sm, player)
		player:map(stage_maps[1])
		player:try_party_quest(ranking_quest_id)
	end,

	on_player_dead = function(sm, player)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == 103000800
			or map_id == 103000801
			or map_id == 103000802
			or map_id == 103000803
			or map_id == 103000804
			or map_id == 103000805 then
			return
		end
		sm:unregister(player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_player_disconnected = function(sm, player)
	end,

	on_left_party = function(sm, player)
	end,

	on_disband_party = function(sm)
		sm:finish(exit_map_id)
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
