-- State machine (old/scripts/event/elevator.js): 핼리오스탑 엘리베이터

local FLOOR_2 = 222020100
local FLOOR_99 = 222020200
local WAIT_DOWN = 222020210
local RUN_DOWN = 222020211
local WAIT_UP = 222020110
local RUN_UP = 222020111

local SCRIPT = "script/state_machine/elevator.lua"
local WAIT_MS = 50000
local RUN_MS = 60000

local function open_door(floor)
	run_on_map(floor, SCRIPT, "open_door")
end

local function close_door(floor)
	run_on_map(floor, SCRIPT, "close_door")
end

return {
	open_door = function(map)
		map:reload_reactors()
	end,

	close_door = function(map)
		local r = map:find_reactor_name("elevator")
		if r ~= nil then
			r:hit(1)
		end
	end,

	on_init = function(group)
		local sm, err = group:start_persistent()
		if sm == nil then
			if err ~= nil then
				log("elevator start_persistent:", err)
			end
			return
		end
		sm:after("waiting_to_down", 1000, "on_waiting_to_down")
	end,
	on_create = function(sm)
		return {
			WAIT_DOWN,
			RUN_DOWN,
			WAIT_UP,
			RUN_UP,
		}
	end,

	on_player_enter = function(sm, player, map_id)
		local map = sm:map(map_id)
		if map == nil then
			sm:unregister(player)
			return
		end
		player:map(map)
	end,

	on_waiting_to_down = function(sm)
		sm:warp_all(RUN_UP, FLOOR_99)
		open_door(FLOOR_99)
		close_door(FLOOR_2)
		sm:after("run_to_down", WAIT_MS, "on_run_to_down")
	end,

	on_run_to_down = function(sm)
		sm:warp_all(WAIT_DOWN, RUN_DOWN)
		close_door(FLOOR_99)
		sm:after("waiting_to_up", RUN_MS, "on_waiting_to_up")
	end,

	on_waiting_to_up = function(sm)
		sm:warp_all(RUN_DOWN, FLOOR_2)
		open_door(FLOOR_2)
		sm:after("run_to_up", WAIT_MS, "on_run_to_up")
	end,

	on_run_to_up = function(sm)
		sm:warp_all(WAIT_UP, RUN_UP)
		close_door(FLOOR_2)
		sm:after("waiting_to_down", RUN_MS, "on_waiting_to_down")
	end,
}
