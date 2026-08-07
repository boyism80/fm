-- NPC name (String.wz/Npc.img.xml): 매표소 안내원

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.board(me, npc, {
			group = "Boats",
			ticket_low = 4031046,
			ticket_high = 4031047,
			intro = "배는 매 시간 정각 기준으로 15분 마다 출발하고 있으며, 출발 5분 전부터 표를 받고 있답니다.\r\n",
			departed_text = "이미 배가 엘리니아로 출발했답니다. 배는 매 시간 기준으로 15분 마다 출발하니 잠시만 기다려 보세요.",
			ready_text = "아직 배가 출항 준비중에 있습니다. 배는 매 시간 정각 기준 10, 25, 40, 55분에 출항을 준비하며, 최장 30분 이내에 출항 준비가 완료됩니다.",
		})
	end,
}
