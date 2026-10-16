-- State machine (old/scripts/event/FireDemon.js): 용암의 심장부

local EXIT_MAP = 211042300
local QUEST_MAP = 921100000
local DURATION_MS = 300000

local function finish(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { QUEST_MAP }
	end,

	on_prepare = function(sm)
		sm:map(QUEST_MAP):reload_reactors()
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(QUEST_MAP))
	end,

	on_scheduled_timeout = function(sm)
		finish(sm)
	end,

}
