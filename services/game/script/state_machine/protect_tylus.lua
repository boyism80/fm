-- State machine (old/scripts/event/ProtectTylus.js): Protect Tylus

local create = require("script/lib/fourth_skill_instance")

local TYLUS = 9300093
local ELDER_HOUSE = 211000001
local SAFE_ROOM = 921100301

local machine = create({
	map = 921100300,
	exit_map = ELDER_HOUSE,
	duration_ms = 600000,
	max_machines = 1,
})

machine.on_scheduled_timeout = function(sm)
	sm:finish(SAFE_ROOM)
end

machine.on_mob_die = function(sm, mob)
	if mob:id() == TYLUS then
		sm:finish(ELDER_HOUSE)
	end
end

return machine
