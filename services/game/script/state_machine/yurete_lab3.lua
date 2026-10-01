-- State machine (old/scripts/event/YureteLab3.js): Yurete's lab 3

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926130103 },
	exit_map = 926130203,
	duration_ms = 300000,
})
