-- State machine (old/scripts/event/Trains.js): 오르비스↔루디브리엄 열차

local WAITING_ORBIS = 200000122
local WAITING_LUDI = 220000111
local RIDE_ORBIS = 200090100
local RIDE_LUDI = 200090110
local STATION_ORBIS = 200000121
local STATION_LUDI = 220000110

local DOCK_CRON = "5,15,25,35,45,55 * * * *"
local STOP_ENTRY_MS = 240000
local TAKEOFF_MS = 300000
local SAIL_MS = 300000

local function disembark(sm)
	sm:warp_all(RIDE_ORBIS, STATION_LUDI)
	sm:warp_all(RIDE_LUDI, STATION_ORBIS)
	sm:broadcast_ship(STATION_ORBIS, 1)
	sm:broadcast_ship(STATION_LUDI, 1)
end

local function open_dock(sm)
	local group = sm:group()
	group:set_property("ready", "true")
	group:set_property("docked", "true")
	group:set_property("entry", "true")
	sm:after("stop_entry", STOP_ENTRY_MS, "on_stop_entry")
	sm:after("takeoff", TAKEOFF_MS, "on_takeoff")
end

return {
	on_init = function(group)
		group:declare_maps({
			WAITING_ORBIS,
			WAITING_LUDI,
			RIDE_ORBIS,
			RIDE_LUDI,
		})
		group:set_property("ready", "false")
		group:set_property("docked", "false")
		group:set_property("entry", "false")
		local sm, err = group:start_persistent()
		if sm == nil then
			if err ~= nil then
				log("Trains start_persistent:", err)
			end
			return
		end
		sm:cron("dock", DOCK_CRON, "on_dock")
	end,

	on_dock = function(sm)
		disembark(sm)
		open_dock(sm)
	end,

	on_stop_entry = function(sm)
		sm:group():set_property("entry", "false")
	end,

	on_takeoff = function(sm)
		sm:warp_all(WAITING_ORBIS, RIDE_ORBIS)
		sm:warp_all(WAITING_LUDI, RIDE_LUDI)
		sm:broadcast_ship(STATION_ORBIS, 3)
		sm:broadcast_ship(STATION_LUDI, 3)
		sm:group():set_property("docked", "false")
		sm:after("arrived", SAIL_MS, "on_arrived")
	end,

	on_arrived = function(sm)
		disembark(sm)
	end,

	on_cancel_schedule = function(group)
	end,
}
