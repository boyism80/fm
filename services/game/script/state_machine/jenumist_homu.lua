-- State machine (old/scripts/event/JenumistHomu.js): Closed laboratory

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926120100 },
	exit_map = 261000010,
	duration_ms = 1800000,
})
