-- State machine (old/scripts/event/Flight.js): 오르비스↔리프레 비행선

local WAITING_ORBIS = 200000132
local WAITING_LEAFRE = 240000111
local RIDE_ORBIS = 200090200
local RIDE_LEAFRE = 200090210
local STATION_ORBIS = 200000131
local STATION_LEAFRE = 240000110
local DEST_LEAFRE = 240000100
local DEST_ORBIS = 200000100

local DOCK_CRON = "5,15,25,35,45,55 * * * *"
local STOP_ENTRY_MS = 240000
local TAKEOFF_MS = 300000
local SAIL_MS = 300000

local function disembark(sm)
	sm:warp_all(RIDE_ORBIS, DEST_LEAFRE)
	sm:warp_all(RIDE_LEAFRE, DEST_ORBIS)
	sm:broadcast_ship(STATION_ORBIS, 1)
	sm:broadcast_ship(STATION_LEAFRE, 1)
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
		group:set_property("ready", "false")
		group:set_property("docked", "false")
		group:set_property("entry", "false")
		local sm, err = group:start_persistent()
		if sm == nil then
			if err ~= nil then
				log("Flight start_persistent:", err)
			end
			return
		end
		sm:cron("dock", DOCK_CRON, "on_dock")
	end,
	on_create = function(sm)
		return {
			WAITING_ORBIS,
			WAITING_LEAFRE,
			RIDE_ORBIS,
			RIDE_LEAFRE,
		}
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
		sm:warp_all(WAITING_LEAFRE, RIDE_LEAFRE)
		sm:broadcast_ship(STATION_ORBIS, 3)
		sm:broadcast_ship(STATION_LEAFRE, 3)
		sm:group():set_property("docked", "false")
		sm:after("arrived", SAIL_MS, "on_arrived")
	end,

	on_arrived = function(sm)
		disembark(sm)
	end,

	on_cancel_schedule = function(group)
	end,
}
