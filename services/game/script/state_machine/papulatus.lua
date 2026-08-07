-- State machine (old/scripts/event/Papulatus.js): 파풀라투스

local LOBBY_MAP = 220080000
local BOSS_MAP = 220080001
local GATE_REACTOR = 2208001

local function reset_gate(group)
	if group == nil then
		return
	end
	group:set_property("battle", "0")
	local lobby = group:map(LOBBY_MAP)
	if lobby == nil then
		return
	end
	local gate = lobby:reactor(GATE_REACTOR)
	if gate ~= nil then
		gate:hit(0)
	end
end

local function close_gate(group)
	if group == nil then
		return
	end
	local lobby = group:map(LOBBY_MAP)
	if lobby == nil then
		return
	end
	local gate = lobby:reactor(GATE_REACTOR)
	if gate ~= nil then
		gate:hit(1)
	end
end

return {
	on_init = function(group)
		group:set_property("battle", "0")
		group:declare_min_players(1)
		group:declare_exit_map(LOBBY_MAP)
		reset_gate(group)
	end,

	on_create = function(sm)
		local group = sm:group()
		local map = group:map(BOSS_MAP)
		if map ~= nil then
			map:reset()
		end
		reset_gate(group)
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
		player:map(BOSS_MAP)
	end,

	on_player_dead = function(sm, player)
	end,

	on_player_revive = function(sm, player)
	end,

	on_player_disconnected = function(sm, player)
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

	on_cancel_schedule = function(group)
	end
}
