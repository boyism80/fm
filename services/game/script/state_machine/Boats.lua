-- State machine (old/scripts/event/Boats.js): 엘리니아↔오르비스 배

local WAITING_ORBIS = 200000112
local WAITING_ELLINIA = 101000301
local RIDE_ORBIS = 200090000
local RIDE_ELLINIA = 200090010
local CABIN_ORBIS = 200090001
local CABIN_ELLINIA = 200090011
local STATION_ORBIS = 200000111
local STATION_ELLINIA = 101000300
local DEST_ELLINIA = 101000300
local DEST_ORBIS = 200000100

local DOCK_CRON = "10,25,40,55 * * * *"
local STOP_ENTRY_MS = 240000
local TAKEOFF_MS = 300000
local SAIL_MS = 600000
local INVASION_MS = 180000
local BALROG = 8150000

local function kill_decks(sm)
	local group = sm:group()
	local m1 = group:map(RIDE_ORBIS)
	local m2 = group:map(RIDE_ELLINIA)
	if m1 ~= nil then
		m1:kill_all_mobs()
	end
	if m2 ~= nil then
		m2:kill_all_mobs()
	end
end

local function disembark(sm)
	sm:warp_all(RIDE_ELLINIA, DEST_ORBIS)
	sm:warp_all(CABIN_ELLINIA, DEST_ORBIS)
	sm:warp_all(RIDE_ORBIS, DEST_ELLINIA)
	sm:warp_all(CABIN_ORBIS, DEST_ELLINIA)
	sm:broadcast_ship(STATION_ORBIS, 1)
	sm:broadcast_ship(STATION_ELLINIA, 1)
	kill_decks(sm)
	sm:group():set_property("haveBalrog", "false")
end

local function open_dock(sm)
	local group = sm:group()
	group:set_property("ready", "true")
	group:set_property("docked", "true")
	group:set_property("entry", "true")
	group:set_property("haveBalrog", "false")
	kill_decks(sm)
	sm:after("stop_entry", STOP_ENTRY_MS, "on_stop_entry")
	sm:after("takeoff", TAKEOFF_MS, "on_takeoff")
end

return {
	on_init = function(group)
		group:declare_maps({
			WAITING_ORBIS,
			WAITING_ELLINIA,
			RIDE_ORBIS,
			RIDE_ELLINIA,
			CABIN_ORBIS,
			CABIN_ELLINIA,
		})
		group:set_property("ready", "false")
		group:set_property("docked", "false")
		group:set_property("entry", "false")
		group:set_property("haveBalrog", "false")
		local sm, err = group:start_persistent()
		if sm == nil then
			if err ~= nil then
				log("Boats start_persistent:", err)
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
		local group = sm:group()
		local c1 = group:map(CABIN_ORBIS)
		local c2 = group:map(CABIN_ELLINIA)
		if c1 ~= nil then
			c1:reload_reactors()
		end
		if c2 ~= nil then
			c2:reload_reactors()
		end
	end,

	on_takeoff = function(sm)
		sm:warp_all(WAITING_ORBIS, RIDE_ORBIS)
		sm:warp_all(WAITING_ELLINIA, RIDE_ELLINIA)
		sm:broadcast_ship(STATION_ORBIS, 3)
		sm:broadcast_ship(STATION_ELLINIA, 3)
		sm:group():set_property("docked", "false")
		sm:after("invasion", INVASION_MS, "on_invasion")
		sm:after("arrived", SAIL_MS, "on_arrived")
	end,

	on_invasion = function(sm)
		local group = sm:group()
		local m1 = group:map(RIDE_ORBIS)
		local m2 = group:map(RIDE_ELLINIA)
		if m1 ~= nil then
			m1:spawn_mob(BALROG, -538, 143)
			m1:spawn_mob(BALROG, -538, 143)
		end
		if m2 ~= nil then
			m2:spawn_mob(BALROG, 339, 148)
			m2:spawn_mob(BALROG, 339, 148)
		end
		group:set_property("haveBalrog", "true")
		sm:broadcast_ship(RIDE_ORBIS, 1034)
		sm:broadcast_ship(RIDE_ELLINIA, 1034)
	end,

	on_arrived = function(sm)
		disembark(sm)
	end,

	on_cancel_schedule = function(group)
	end,
}
