-- State machine: 소공의 방

local solo_instance = require("script/lib/solo_instance")
local dojo = require("script/lib/dojo")

local machine = solo_instance.create({
	maps = { dojo.TUTORIAL },
	exit_map = dojo.EXIT,
	duration_ms = dojo.time_limit_ms(1),
	setup = function(sm, maps)
		maps[1]:respawn({ include_one_time = true })
	end,
})

machine.on_mob_die = function(sm, mob)
	if mob:id() ~= dojo.SO_GONG then
		return
	end
	sm:set_property("cleared", "1")
	sm:stop_timer()
	local map = sm:map(dojo.TUTORIAL)
	map:play_sound("Dojang/clear")
	map:show_effect("dojang/end/clear")
end

return machine
