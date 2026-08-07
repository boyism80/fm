-- State machine (old/scripts/event/AirStrike2.js): 샤레니안 성문

local EXIT_MAP = 120000102
local QUEST_MAP = 910210000
local DURATION_MS = 900000

local function finish(sm)
	local players = sm:players()
	sm:finish(0)
	for _, player in ipairs(players) do
		player:map(EXIT_MAP, 1)
	end
end

return {
	on_init = function(group)
		group:set_property("noEntry", "false")
		group:declare_min_players(1)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("noEntry", "true")
		local map = group:map(QUEST_MAP)
		if map ~= nil then
			map:reset()
		end
		return { QUEST_MAP }

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		player:map(QUEST_MAP)
	end,

	on_player_disconnected = function(sm, player)
		sm:finish(0)
	end,

	on_scheduled_timeout = function(sm)
		finish(sm)
	end,

	on_finish = function(sm)
		sm:group():set_property("noEntry", "false")
	end
}
