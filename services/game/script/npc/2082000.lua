-- NPC name (String.wz/Npc.img.xml): 톰슨

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.sell(me, npc, {
			ticket_low = 4031044,
			ticket_high = 4031045,
			cost_low = 1500,
			cost_high = 3000,
		})
	end,
}
