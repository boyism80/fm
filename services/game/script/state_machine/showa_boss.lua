-- State machine (old/scripts/event/showaboss.js): Showa hideout boss

local boss_instance = require("script/lib/boss_instance")

return boss_instance.create({
	map = 801040100,
	exit_map = 801040004,
	duration_ms = 14400000,
	bosses = {
		{ id = 9400112, x = 1, y = 161 },
		{ id = 9400112, x = 1, y = 161 },
	},
})
