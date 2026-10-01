-- State machine (old/scripts/event/SnowRose.js): Snow rose garden

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926120300 },
	exit_map = 261000000,
	duration_ms = 900000,
})
