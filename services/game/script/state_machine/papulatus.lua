-- State machine (old/scripts/event/Papulatus.js): 파풀라투스

local LOBBY_MAP = 220080000
local BOSS_MAP = 220080001
local GATE_REACTOR = 2208001
local SCRIPT = "script/state_machine/papulatus.lua"

local function reset_gate(group)
	if group == nil then
		return
	end
	group:set_property("battle", "0")
	run_on_map(LOBBY_MAP, SCRIPT, "hit_gate", 0)
end

local function close_gate(group)
	if group == nil then
		return
	end
	run_on_map(LOBBY_MAP, SCRIPT, "hit_gate", 1)
end

return {
	hit_gate = function(map, state)
		local gate = map:reactor(GATE_REACTOR)
		if gate ~= nil then
			gate:hit(state)
		end
	end,

	on_init = function(group)
		group:set_property("battle", "0")
		group:min_players(1)
		group:exit_map(LOBBY_MAP)
	end,

	on_create = function(sm)
		reset_gate(sm:group())
		return { BOSS_MAP }
	end,

	on_start = function(sm)
		local group = sm:group()
		if group:get_property("battle") == "1" then
			return
		end
		group:set_property("battle", "1")
		close_gate(group)
	end,

	on_player_enter = function(sm, player)
		player:map(sm:map(BOSS_MAP))
	end,

	on_player_dead = function(sm, player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_left_party = function(sm, player)
	end,

	on_disband_party = function(sm)
	end,

	on_scheduled_timeout = function(sm)
		sm:finish(LOBBY_MAP)
	end,

	on_finish = function(sm)
		reset_gate(sm:group())
	end,

	on_all_monsters_dead = function(sm)
	end,
}
