-- State machine (old/scripts/event/MachineRoom.js): Toy factory machine room

local pq = require("script/lib/party_quest")
local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 922000000 },
	exit_map = 220020000,
	exit_portal = 2,
	duration_ms = 600000,
	setup = function(sm, maps)
		pq.shuffle_reactors(maps[1])
	end,
})
