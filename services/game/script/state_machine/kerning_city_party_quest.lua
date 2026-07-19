-- State machine (old/scripts/event/KerningPQ.js): 커닝시티 파티 퀘스트

local config = require("script/lib/kerning_city_party_quest")

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

function on_init(group)
	group:set_property("state", "0")
	group:declare_maps(stage_maps)
end

function on_mob_kill(sm, player, mobs)
	sm:add_kill(player, #mobs)
end

function on_setup(sm)
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
	sm:start_timer(duration_ms)
end

function on_player_entry(sm, player)
	player:map(stage_maps[1])
	player:try_party_quest(ranking_quest_id)
end

function on_player_dead(sm, player)
end

function on_changed_map(sm, player, map_id)
	if map_id == 103000800
		or map_id == 103000801
		or map_id == 103000802
		or map_id == 103000803
		or map_id == 103000804
		or map_id == 103000805 then
		return
	end
	sm:unregister(player)
	if #sm:players() < config.required_party_size then
		sm:finish(exit_map_id)
		sm:group():set_property("state", "0")
	end
end

function on_player_revive(sm, player)
end

function on_player_disconnected(sm, player)
	return -3
end

function on_left_party(sm, player)
	if #sm:players() <= config.required_party_size then
		sm:finish(exit_map_id)
		sm:group():set_property("state", "0")
	else
		on_player_exit(sm, player)
	end
end

function on_disband_party(sm)
	sm:finish(exit_map_id)
	sm:group():set_property("state", "0")
end

function on_scheduled_timeout(sm)
	on_clear_party_quest(sm)
end

function on_player_exit(sm, player)
	sm:unregister(player)
	player:map(exit_map_id)
	if #sm:players() < config.required_party_size then
		sm:finish(exit_map_id)
		sm:group():set_property("state", "0")
	end
end

function on_clear_party_quest(sm)
	sm:finish(exit_map_id)
	sm:group():set_property("state", "0")
end

function on_all_monsters_dead(sm)
end

function on_cancel_schedule(group)
end
