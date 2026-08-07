-- NPC name (String.wz/Npc.img.xml): 죠엘

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.sell(me, npc, {
			ticket_low = 4031044,
			ticket_high = 4031045,
			cost_low = 1500,
			cost_high = 3000,
			min_level = 15,
			min_level_text = "흐음.. 그런데 당신은 아직 오시리아 대륙으로 가보시기엔 너무 약해보이시는군요. 조금 더 수련을 하신 후 다시 찾아오세요.",
		})
	end,
}
