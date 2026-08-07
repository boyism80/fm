-- NPC name (String.wz/Npc.img.xml): 아리안트 대기실 안내원

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.exit_waiting(me, npc, -10)
	end,
}
