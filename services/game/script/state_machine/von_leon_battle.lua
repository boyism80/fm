-- State machine (old/scripts/event/VonLeonBattle.js): 반 레온 원정대

local ex = require("script/lib/expedition")

local EXIT_MAP = 211061001
local AUDIENCE_HALL = 211070100
local CORRIDOR = 211070110
local VON_LEON = 2161000
local DURATION_MS = 3600000

return {
	on_init = function(group)
		group:min_players(1)
		group:max_machines(nil)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { AUDIENCE_HALL, CORRIDOR }
	end,

	on_prepare = function(sm)
		sm:map(AUDIENCE_HALL):spawn_npc(VON_LEON, 0, -181)
	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		local function on_arrive(player)
			ex.stamp(player, "von_leon")
		end
		player:map(sm:map(AUDIENCE_HALL), { callback = on_arrive })
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(EXIT_MAP)
	end,
}
