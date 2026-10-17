-- State machine (old/scripts/event/Vergamot.js): Vergamot

local boss_instance = require("script/lib/boss_instance")

local CRYSTAL = 9120026

return boss_instance.create({
	map = 802000211,
	exit_map = 802000212,
	duration_ms = 14400000,
	bosses = {
		{ id = 9400276, x = 156, y = 258 },
		{ id = 9400263, x = 559, y = 545 },
	},
	on_cleared = function(sm, map)
		map:spawn_npc(CRYSTAL, 479, 485)
	end,
})
