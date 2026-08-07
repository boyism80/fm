-- NPC name (String.wz/Npc.img.xml): 체리

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.board(me, npc, {
			group = "Boats",
			ticket_low = 4031044,
			ticket_high = 4031045,
			min_level = 15,
			min_level_text = "흐음.. 그런데 당신은 아직 오시리아 대륙으로 가보시기엔 너무 약해보이시는군요. 조금 더 수련을 하신 후 다시 찾아오세요.",
			intro = "배는 매 시간 정각 기준으로 15분 마다 출발하고 있으며, 출발 5분 전부터 표를 받고 있답니다.\r\n",
			departed_text = "이미 배가 오르비스로 출발했답니다. 배는 매 시간 기준으로 15분 마다 출발하니 잠시만 기다려 보세요.",
			ready_text = "아직 배가 출항 준비중에 있습니다. 배는 매 시간 정각 기준 10, 25, 40, 55분에 출항을 준비하며, 최장 30분 이내에 출항 준비가 완료됩니다.",
		})
	end,
}
