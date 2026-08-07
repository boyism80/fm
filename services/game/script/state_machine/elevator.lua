-- State machine (old/scripts/event/elevator.js): 핼리오스탑 엘리베이터

local FLOOR_2 = 222020100
local FLOOR_99 = 222020200
local WAIT_DOWN = 222020210
local RUN_DOWN = 222020211
local WAIT_UP = 222020110
local RUN_UP = 222020111

local WAIT_MS = 50000
local RUN_MS = 60000

local function open_door(map)
	if map ~= nil then
		map:reload_reactors()
	end
end

local function close_door(map)
	if map == nil then
		return
	end
	local r = map:reactor_by_name("elevator")
	if r ~= nil then
		r:hit(1)
	end
end

return {
	on_init = function(group)
		group:declare_maps({
			FLOOR_2,
			FLOOR_99,
			WAIT_DOWN,
			RUN_DOWN,
			WAIT_UP,
			RUN_UP,
		})
		local sm, err = group:start_persistent()
		if sm == nil then
			if err ~= nil then
				log("elevator start_persistent:", err)
			end
			return
		end
		sm:after("waiting_to_down", 1000, "on_waiting_to_down")
	end,

	on_waiting_to_down = function(sm)
		sm:warp_all(RUN_UP, FLOOR_99)
		local group = sm:group()
		open_door(group:map(FLOOR_99))
		close_door(group:map(FLOOR_2))
		sm:after("run_to_down", WAIT_MS, "on_run_to_down")
	end,

	on_run_to_down = function(sm)
		sm:warp_all(WAIT_DOWN, RUN_DOWN)
		close_door(sm:group():map(FLOOR_99))
		sm:after("waiting_to_up", RUN_MS, "on_waiting_to_up")
	end,

	on_waiting_to_up = function(sm)
		sm:warp_all(RUN_DOWN, FLOOR_2)
		open_door(sm:group():map(FLOOR_2))
		sm:after("run_to_up", WAIT_MS, "on_run_to_up")
	end,

	on_run_to_up = function(sm)
		sm:warp_all(WAIT_UP, RUN_UP)
		close_door(sm:group():map(FLOOR_2))
		sm:after("waiting_to_down", RUN_MS, "on_waiting_to_down")
	end,

	on_cancel_schedule = function(group)
	end,
}
