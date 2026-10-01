-- State machine (old/scripts/event/YureteLab2.js): Yurete's lab 2

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926130102 },
	exit_map = 926130201,
	duration_ms = 300000,
})
