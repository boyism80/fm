-- State machine (old/scripts/event/HorntailBattle.js): 혼테일 원정대

local ex = require("script/lib/expedition")

local EXIT_MAP = 240050400
local CAVE1 = 240060000
local CAVE2 = 240060100
local CAVE3 = 240060200
local DURATION_MS = 43200000

return {
	on_init = function(group)
		group:min_players(1)
		group:max_machines(nil)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { CAVE1, CAVE2, CAVE3 }
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		ex.stamp(player, "horntail")
		player:map(sm:map(CAVE1))
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(EXIT_MAP)
	end,
}
