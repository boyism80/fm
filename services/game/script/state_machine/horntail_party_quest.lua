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

local function in_stage(map_id)
	return map_id >= 240050100 and map_id <= 240050310
end

function on_init(group)
	group:set_property("state", "0")
	group:declare_maps(stage_maps)
	group:declare_min_players(5)
	group:declare_exit_map(FAIL_EXIT)
end

function on_setup(sm)
	local group = sm:group()
	group:set_property("state", "1")
	sm:set_property("stage1progress", "0")
	sm:set_property("allfinish", "")
	for _, map_id in ipairs(stage_maps) do
		local map = group:map(map_id)
		map:reset()
		map:respawn(true)
	end
	sm:start_timer(DURATION_MS)
end

function on_player_entry(sm, player)
	player:map(stage_maps[1])
end

function on_mob_kill(sm, player, mobs)
end

function on_changed_map(sm, player, map_id)
	if in_stage(map_id) then
		return
	end
	sm:unregister(player)
end

function on_player_dead(sm, player)
end

function on_player_revive(sm, player)
end

function on_player_disconnected(sm, player)
end

function on_left_party(sm, player)
end

function on_disband_party(sm)
	on_clear(sm)
end

function on_scheduled_timeout(sm)
	on_clear(sm)
end

function on_clear(sm)
	local exit_map = FAIL_EXIT
	local allfinish = sm:get_property("allfinish")
	if allfinish ~= nil and allfinish ~= "" then
		exit_map = CLEAR_EXIT
	end
	sm:finish(exit_map)
end

function on_finish(sm)
	sm:group():set_property("state", "0")
end

function on_all_monsters_dead(sm)
end

function on_cancel_schedule(group)
end
