-- State machine (old/scripts/event/Adin.js): Adin sparring

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
	maps = { 926000000 },
	exit_map = 260000200,
	duration_ms = 900000,
	timeout_message = "제한시간이 다 되어 실패하였습니다.",
	setup = function(sm, maps)
		maps[1]:respawn({ include_one_time = true })
	end,
})
