-- State machine (old/scripts/event/ZakumBattle.js): 자쿰 원정대

local ex = require("script/lib/expedition")

local GROUP = "zakum_battle"
local EXIT_MAP = 211042300
local ALTAR_MAP = 280030000
local GATE_REACTOR = 2118002
local SCRIPT = "script/state_machine/zakum_battle.lua"

local function gate_closed()
	local group = state_machine(GROUP)
	if group == nil then
		return false
	end
	local summoned = 0
	for _, sm in ipairs(group:machines()) do
		if sm:get_property("summoned") == "1" then
			summoned = summoned + 1
		end
	end
	return summoned >= ex.BOSS.zakum.max_battles
end

return {
	gate_closed = gate_closed,

	sync_gate = function(map)
		local gate = map:reactor(GATE_REACTOR)
		if gate == nil then
			return
		end
		if gate_closed() then
			gate:hit(1)
		else
			gate:hit(0)
		end
	end,

	summon = function(sm)
		sm:set_property("summoned", "1")
		run_on_map(EXIT_MAP, SCRIPT, "sync_gate")
	end,

	on_init = function(group)
		group:min_players(1)
		group:max_machines(nil)
		group:exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		return { ALTAR_MAP }
	end,

	on_player_enter = function(sm, player)
		local function on_arrive(player)
			ex.stamp(player, "zakum")
		end
		player:map(sm:map(ALTAR_MAP), { callback = on_arrive })
	end,

	on_finish = function(sm)
		sm:set_property("summoned", "0")
		run_on_map(EXIT_MAP, SCRIPT, "sync_gate")
	end,
}
