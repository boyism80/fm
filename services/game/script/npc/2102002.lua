-- NPC name (String.wz/Npc.img.xml): 아리안트 티켓 판매원

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.sell(me, npc, {
			ticket_low = 4031044,
			ticket_high = 4031045,
			cost_low = 1000,
			cost_high = 2000,
		})
	end,
}
