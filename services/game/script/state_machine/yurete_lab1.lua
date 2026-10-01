-- State machine (old/scripts/event/YureteLab1.js): Yurete's lab 1

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926130101 },
	exit_map = 926130200,
	duration_ms = 300000,
})
