-- State machine (old/scripts/event/tokyo_2095.js): 2095 tokyo

local EXIT_MAP = 802000312
local TEMPLATE_MAP = 802000311
local DURATION_MS = 1200000

local function instance(sm)
	local maps = sm:maps()
	return maps[1]
end

local function finish(sm, exit_map)
	if exit_map == nil then
		exit_map = EXIT_MAP
	end
	local map = instance(sm)
	sm:finish(exit_map)
	if map ~= nil then
		map:destroy()
	end
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:set_property("leader", "true")
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		group:set_property("leader", "true")
		local tmpl = id2map(TEMPLATE_MAP)
		if tmpl == nil then
			sm:finish(0)
			return
		end
		local map, err = tmpl:create_instance()
		if map == nil then
			if err ~= nil then
				log("tokyo_2095 create_instance:", err)
			end
			sm:finish(0)
			return
		end
		map:kill_all_mobs()
		return { map }

	end,

	on_start = function(sm)
		sm:start_timer(DURATION_MS)
	end,

	on_player_enter = function(sm, player)
		local map = instance(sm)
		if map ~= nil then
			player:map(map, 0)
		end
	end,

	on_player_disconnected = function(sm, player)
		sm:unregister(player)
		if #sm:players() == 0 then
			finish(sm, 0)
		end
	end,

	on_scheduled_timeout = function(sm)
		finish(sm, EXIT_MAP)
	end,

	on_clear = function(sm)
		finish(sm, 0)
	end,

	on_finish = function(sm)
		local map = instance(sm)
		if map ~= nil then
			map:destroy()
		end
		local group = sm:group()
		group:set_property("state", "0")
		group:set_property("leader", "true")
	end,
}
