-- State machine (old/scripts/event/DarkMagicianAgit.js): Black Mage hideout

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 912000000 },
	exit_map = 120000100,
	duration_ms = 600000,
})
