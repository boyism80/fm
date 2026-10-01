-- State machine (old/scripts/event/NautilusCow.js): Nautilus cow barn

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 912000100 },
	exit_map = 120000103,
	duration_ms = 300000,
})
