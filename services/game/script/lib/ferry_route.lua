local M = {}

local DOCK_CRON = "5,15,25,35,45,55 * * * *"
local STOP_ENTRY_MS = 240000
local TAKEOFF_MS = 300000
local SAIL_MS = 300000

-- A route is a list of legs; on each leg passengers wait in `waiting`, ride in `ride`, and get off at `dest`, while
-- `station` shows the ship. The group properties and schedule ids are read by lib/ferry.lua, the ship object handler,
-- and the transport GM commands.
function M.new(route)
    local function disembark(sm)
        for _, leg in ipairs(route.legs) do
            sm:warp_all(leg.ride, leg.dest)
        end
        for _, leg in ipairs(route.legs) do
            sm:broadcast_ship(leg.station, 1)
        end
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
                    log(route.name .. " start_persistent:", err)
                end
                return
            end
            sm:cron("dock", DOCK_CRON, "on_dock")
        end,

        on_create = function(sm)
            local maps = {}
            for _, leg in ipairs(route.legs) do
                table.insert(maps, leg.waiting)
            end
            for _, leg in ipairs(route.legs) do
                table.insert(maps, leg.ride)
            end
            return maps
        end,

        on_player_enter = function(sm, player, map_id)
            local map = sm:map(map_id)
            if map == nil then
                sm:unregister(player)
                return
            end
            player:map(map)
        end,

        on_dock = function(sm)
            disembark(sm)
            open_dock(sm)
        end,

        on_stop_entry = function(sm)
            sm:group():set_property("entry", "false")
        end,

        on_takeoff = function(sm)
            for _, leg in ipairs(route.legs) do
                sm:warp_all(leg.waiting, leg.ride)
            end
            for _, leg in ipairs(route.legs) do
                sm:broadcast_ship(leg.station, 3)
            end
            sm:group():set_property("docked", "false")
            sm:after("arrived", SAIL_MS, "on_arrived")
        end,

        on_arrived = function(sm)
            disembark(sm)
        end,
    }
end

return M
