-- State machine (old/src/server/maps/AramiaFireWorks.java): 채널 기부 진행도가 가득 찼을 때의 맵 연출

local SCRIPT = "script/state_machine/fireworks.lua"
local SHOW_DELAY_MS = 10000
local STEP_MS = 5000
local ROUND_MS = 10000
local ROUNDS = 3

local KEG_MAP = 100000200
local SUN_MAP = 970010000

local KEG_MOBS = {
	{ 9500168, 2100, 574 }, { 9500169, 2605, 364 }, { 9500170, 1800, 574 }, { 9500171, 2600, 316 },
	{ 9500173, 3120, 574 }, { 9500174, 2700, 574 }, { 9500175, 2320, 403 }, { 9500176, 2062, 364 },
	{ 9500170, 2800, 574 }, { 9500171, 3100, 574 }, { 9500172, 2300, 403 }, { 9500173, 2840, 574 },
	{ 9500174, 2700, 574 }, { 9500175, 2320, 403 }, { 9400569, 1950, 574 },
}

local SUN_DROPS = {
	{ 720, 1234 }, { 180, 1234 }, { 630, 1174 }, { 270, 1234 }, { 360, 1174 }, { 540, 1174 }, { 450, 1174 }, { 142, 1260 },
	{ 142, 1234 }, { 218, 1234 }, { 772, 1234 }, { 810, 1234 }, { 848, 1234 }, { 232, 1114 }, { 308, 1114 }, { 142, 1140 },
}

local function map_of(kind)
	if kind == "keg" then
		return KEG_MAP
	end
	return SUN_MAP
end

return {
	keg_message = function(map)
		map:weather(5121010, "누가 시끄럽게 폭죽을 쏘는거야?")
	end,

	spawn_keg_mobs = function(map)
		for _, mob in ipairs(KEG_MOBS) do
			map:spawn_mob(mob[1], mob[2], mob[3])
		end
	end,

	sun_message = function(map)
		map:weather(5121010, "나무에서 햇살이 쏟아지고 있어요!")
	end,

	drop_sun = function(map)
		for i = 1, math.random(10, 15) do
			local item = 4001246
			local roll = math.random(0, 14)
			if roll <= 1 then
				item = 3010141
			elseif roll == 2 then
				item = 3010146
			elseif roll <= 4 then
				item = 3010025
			end
			map:spawn_item(item, 1, { x = SUN_DROPS[i][1], y = SUN_DROPS[i][2] })
		end
	end,

	-- The tree grows one stage for every sixth of the goal; the stage is the same for every character on the map.
	set_tree_stage = function(map, stage)
		local tree = map:find_reactor_name("mapleTree")
		if tree ~= nil then
			tree:hit(stage)
		end
	end,

	on_init = function(group)
		group:max_machines(2)
	end,

	-- A machine owns no map: its timers run on the server, so the show goes on when the donor leaves.
	on_create = function(sm)
		return {}
	end,

	on_prepare = function(sm)
		sm:after("begin", SHOW_DELAY_MS, "on_begin")
	end,

	on_begin = function(sm)
		local map_id = map_of(sm:get_property("kind"))
		if sm:get_property("kind") == "keg" then
			run_on_map(map_id, SCRIPT, "keg_message")
			sm:after("mobs", STEP_MS, "on_keg_mobs")
			return
		end
		run_on_map(map_id, SCRIPT, "set_tree_stage", 0)
		run_on_map(map_id, SCRIPT, "sun_message")
		sm:set_property("round", "0")
		sm:after("drop", STEP_MS, "on_drop")
	end,

	on_keg_mobs = function(sm)
		run_on_map(KEG_MAP, SCRIPT, "spawn_keg_mobs")
		sm:finish(0)
	end,

	on_drop = function(sm)
		run_on_map(SUN_MAP, SCRIPT, "drop_sun")
		local round = tonumber(sm:get_property("round")) + 1
		if round >= ROUNDS then
			sm:finish(0)
			return
		end
		sm:set_property("round", tostring(round))
		sm:after("drop", ROUND_MS, "on_drop")
	end,
}
