-- State machine (old/scripts/event/ProtectDelli.js): Protect Delli

local solo_instance = require("script/lib/solo_instance")

local DELLI = 9300162
local EXIT_MAP = 120000104
local SHELTER_MAP = 925010400

local machine = solo_instance.create({
	maps = { 925010000, 925010100, 925010200, 925010300 },
	links = {
		[925010000] = { "out00" },
		[925010100] = { "in01", "out01" },
		[925010200] = { "in02" },
	},
	exit_map = EXIT_MAP,
	duration_ms = 1200000,
	setup = function(sm, maps)
		maps[4]:portal("in03"):script("protect_delli_back")
	end,
})

machine.on_scheduled_timeout = function(sm)
	if sm:get_property("protect") == "1" then
		solo_instance.leave(sm, SHELTER_MAP)
		return
	end
	solo_instance.leave(sm, EXIT_MAP)
end

machine.on_mob_die = function(sm, mob)
	if mob:id() == DELLI then
		solo_instance.leave(sm, EXIT_MAP)
	end
end

return machine
