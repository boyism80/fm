-- State machine (old/scripts/event/FireDemon.js): 용암의 심장부

local EXIT_MAP = 211042300
local QUEST_MAP = 921100000
local DURATION_MS = 300000

local function finish(sm)
	sm:finish(EXIT_MAP)
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_maps({ QUEST_MAP })
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		local map = group:map(QUEST_MAP)
		if map ~= nil then
			map:reload_reactors()
		end
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(QUEST_MAP)
	end,

	on_changed_map = function(sm, player, map_id)
		if map_id == QUEST_MAP then
			return
		end
		sm:unregister(player)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm)
	end,

	on_clear = function(sm)
		finish(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end
}
