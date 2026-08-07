-- State machine (old/scripts/event/ZakumPQ.js): 아도비스의 임무1

local pq = require("script/lib/party_quest")

local EXIT_MAP = 280090000
local START_MAP = 280010000
local DURATION_MS = 1800000
local STAGE1_QUEST = 100001
local SHUFFLE_REACTOR_MIN = 2110000
local SHUFFLE_REACTOR_MAX = 2112013

local STAGE_MAPS = {
	280010000,
	280010010,
	280010011,
	280010020,
	280010030,
	280010031,
	280010040,
	280010041,
	280010050,
	280010060,
	280010070,
	280010071,
	280010080,
	280010081,
	280010090,
	280010091,
	280010100,
	280010101,
	280010110,
	280010120,
	280010130,
	280010140,
	280010150,
	280011000,
	280011001,
	280011002,
	280011003,
	280011004,
	280011005,
	280011006,
}

local function is_stage_map(map_id)
	for _, id in ipairs(STAGE_MAPS) do
		if id == map_id then
			return true
		end
	end
	return false
end

local function ensure_quest_started(me, quest_id, record)
	local q = me:quest(quest_id)
	if q == nil then
		return nil
	end
	if q:started() or q:completed() then
		return q
	end
	if record == nil then
		record = ""
	end
	if q:wz() == nil then
		q:start(record)
	else
		q:start(0, true)
	end
	return me:quest(quest_id)
end

local function end_run(sm)
	sm:finish(EXIT_MAP)
end

local function exit_player(sm, player)
	sm:unregister(player)
	if player ~= nil then
		player:map(EXIT_MAP)
	end
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps(STAGE_MAPS)
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		sm:set_property("clear", "")
		sm:set_property("paper", "")
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
				pq.shuffle_reactors(map, SHUFFLE_REACTOR_MIN, SHUFFLE_REACTOR_MAX)
			end
		end
		local hub = group:map(START_MAP)
		if hub ~= nil then
			pq.shuffle_reactors(hub)
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		ensure_quest_started(player, STAGE1_QUEST, "")
		player:map(START_MAP)
	end,

	on_changed_map = function(sm, player, map_id)
		if is_stage_map(map_id) then
			return
		end
		sm:unregister(player)
	end,

	on_left_party = function(sm, player)
		exit_player(sm, player)
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
	end
}
