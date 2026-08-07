-- State machine (old/scripts/event/Gojarani.js): Gojarani

local EXIT_MAP = 200000204
local TEMPLATE_MAP = 103000004
local NPC_ID = 1052004
local NPC_X = -70
local NPC_Y = 95

local function instance(sm)
	local maps = sm:maps()
	return maps[1]
end

local function finish(sm)
	local map = instance(sm)
	sm:finish(EXIT_MAP)
	if map ~= nil then
		map:destroy()
	end
end

return {
	on_init = function(group)
		group:set_property("state", "0")
		group:declare_min_players(1)
		group:declare_exit_map(EXIT_MAP)
	end,

	on_create = function(sm)
		local group = sm:group()
		group:set_property("state", "1")
		local tmpl = id2map(TEMPLATE_MAP)
		if tmpl == nil then
			sm:finish(0)
			return
		end
		local map, err = tmpl:create_instance({
			npcs = false,
			reactors = false,
			respawns = false,
		})
		if map == nil then
			if err ~= nil then
				log("gojarani create_instance:", err)
			end
			sm:finish(0)
			return
		end
		map:reset()
		map:spawn_npc(NPC_ID, NPC_X, NPC_Y)
		local p2 = map:portal(2)
		if p2 ~= nil then
			p2:script("goja_out")
		end
		local p3 = map:portal(3)
		if p3 ~= nil then
			p3:script("goja_out")
		end
		return { map }

	end,

	on_player_enter = function(sm, player)
		local map = instance(sm)
		if map ~= nil then
			player:map(map, 0)
		end
	end,

	on_player_disconnected = function(sm, player)
		finish(sm)
	end,

	on_finish = function(sm)
		local map = instance(sm)
		if map ~= nil then
			map:destroy()
		end
		sm:group():set_property("state", "0")
	end,
}
