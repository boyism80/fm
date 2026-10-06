-- State machine (old/scripts/event/PinkBeanBattle.js): 핑크빈 원정대

local ex = require("script/lib/expedition")

local EXIT_MAP = 270050300
local TWILIGHT = 270050100
local DURATION_MS = 86400000

return {
	on_init = function(group)
		group:min_players(1)
		group:max_machines(nil)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { TWILIGHT }
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		ex.stamp(player, "pink_bean")
		player:map(sm:map(TWILIGHT))
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(EXIT_MAP)
	end,
}
