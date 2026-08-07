-- NPC name (String.wz/Npc.img.xml): 비행선 승무원

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.exit_waiting(me, npc, -1)
	end,
}
