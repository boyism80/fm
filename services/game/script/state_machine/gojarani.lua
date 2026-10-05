-- State machine (old/scripts/event/Gojarani.js): Gojarani

local EXIT_MAP = 200000204
local TEMPLATE_MAP = 103000004
local NPC_ID = 1052004
local NPC_X = -70
local NPC_Y = 95

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:min_players(1)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		sm:group():set_property("state", "1")
		return {
			{ id = TEMPLATE_MAP, npcs = false, reactors = false, respawns = false },
		}
	end,

	on_prepare = function(sm)
		local map = sm:map(TEMPLATE_MAP)
		map:spawn_npc(NPC_ID, NPC_X, NPC_Y)
		map:portal(2):script("goja_out")
		map:portal(3):script("goja_out")
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(TEMPLATE_MAP), 0)
	end,

	on_finish = function(sm)
		sm:group():set_property("state", "0")
	end,
}
