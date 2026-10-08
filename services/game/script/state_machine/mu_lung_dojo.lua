-- State machine (old/src/server/maps/Event_DojoAgent.java): 무릉도장

local dojo = require("script/lib/dojo")

local BOSS_SPAWN_EFFECT = 15

return {
	on_init = function(group)
		group:max_machines(nil)
		group:min_players(1)
		group:exit_map(dojo.EXIT)
	end,

	on_create = function(sm)
		local start = tonumber(sm:group():get_property("start:" .. sm:id())) or 1
		sm:set_property("start", tostring(start))
		sm:set_property("started_at", tostring(os.time()))
		local maps = {}
		for floor = start, dojo.LAST_FLOOR do
			table.insert(maps, dojo.map_id(floor))
		end
		return maps
	end,

	on_player_enter = function(sm, player)
		local start = tonumber(sm:get_property("start"))
		player:map(sm:map(dojo.map_id(start)))
	end,

	on_changed_map = function(sm, player, template_id)
		local floor = dojo.floor_of(template_id)
		if floor == nil or sm:get_property("floor") == tostring(floor) then
			return
		end
		sm:set_property("floor", tostring(floor))
		if dojo.is_rest(floor) then
			sm:stop_timer()
			return
		end
		sm:after("floor", 500, "start_floor")
	end,

	start_floor = function(sm)
		local floor = tonumber(sm:get_property("floor"))
		local map = sm:map(dojo.map_id(floor))
		sm:restart_timer(dojo.time_limit_ms(floor))
		map:play_sound("Dojang/start")
		map:show_effect("dojang/start/stage")
		map:show_effect("dojang/start/number/" .. dojo.stage(floor))
		map:tremble(0, 1)
		sm:after("boss", 3000, "spawn_boss")
	end,

	spawn_boss = function(sm)
		local floor = tonumber(sm:get_property("floor"))
		local spawn = dojo.BOSS_SPAWNS[math.random(1, #dojo.BOSS_SPAWNS)]
		sm:map(dojo.map_id(floor)):spawn_mob(dojo.boss(floor), spawn[1], spawn[2], BOSS_SPAWN_EFFECT)
	end,

	on_mob_die = function(sm, mob)
		local floor = tonumber(sm:get_property("floor"))
		local map = sm:map(dojo.map_id(floor))
		local x, y = mob:position()
		if mob:id() ~= dojo.boss(floor) then
			for item = 2022430, 2022433 do
				if math.random(1, 100) <= 15 then
					map:spawn_item(item, 1, { x, y })
				end
			end
			return
		end
		sm:set_property("cleared:" .. floor, "1")
		sm:stop_timer()
		map:play_sound("Dojang/clear")
		map:show_effect("dojang/end/clear")
		for item = 2022359, 2022421 do
			if math.random(1, 1000) <= 3 then
				map:spawn_item(item, 1, { x, y })
			end
		end
	end,

	on_scheduled_timeout = function(sm)
		sm:message("시간이 다 되어 자동으로 퇴장되었습니다.", Msg.PinkText)
		sm:finish(dojo.EXIT)
	end,
}
