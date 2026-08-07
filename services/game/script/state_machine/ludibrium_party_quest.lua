-- State machine (old/scripts/event/LudiPQ.js): 루디브리엄 파티 퀘스트

local MIN_PARTY_SIZE = 3
local DURATION_MS = 3600000
local BONUS_DURATION_MS = 60000
local RANKING_QUEST = 1202
local LOBBY_MAP = 922010000
local START_MAP = 922010100
local BONUS_MAP = 922011000
local REWARD_MAP = 922011100

local STAGE_MAPS = {
	922010100,
	922010200,
	922010201,
	922010300,
	922010400,
	922010401,
	922010402,
	922010403,
	922010404,
	922010405,
	922010500,
	922010501,
	922010502,
	922010503,
	922010504,
	922010505,
	922010506,
	922010600,
	922010700,
	922010800,
	922010900,
	922011000,
}

local function end_run(sm)
	local group = sm:group()
	if group ~= nil and group:get_property("allfinish") == "1" then
		sm:finish(REWARD_MAP)
	else
		sm:finish(LOBBY_MAP)
	end
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:set_property("allfinish", "")
		group:declare_min_players(MIN_PARTY_SIZE)
		group:declare_exit_map(LOBBY_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("allfinish", "")
		sm:set_property("stage", "1")
		sm:set_property("guideRead", "0")
		for _, map_id in ipairs(STAGE_MAPS) do
			local map = group:map(map_id)
			if map ~= nil then
				map:reset()
				map:respawn(true)
				local portal = map:portal("next00")
				if portal ~= nil then
					portal:script("enter_lpq")
				end
			end
		end
		return STAGE_MAPS

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(START_MAP)
		player:try_party_quest(RANKING_QUEST)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == BONUS_MAP then
			local group = sm:group()
			if group ~= nil and group:get_property("allfinish") ~= "1" then
				group:set_property("allfinish", "1")
				sm:restart_timer(BONUS_DURATION_MS)
			end
		end
	end,

	on_left_party = function(sm, player)
		sm:finish(LOBBY_MAP)
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
		group:set_property("allfinish", "")
	end
}
