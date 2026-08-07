-- NPC name (String.wz/Npc.img.xml): 지니 안내원

local ferry = require("script/lib/ferry")

return {
	on_click = function(me, npc)
		ferry.board(me, npc, {
			group = "Geenie",
			ticket_low = 4031575,
			ticket_high = 4031576,
			departed_text = "이미 배가 아리안트로 출발했답니다. 배는 매 시간 기준으로 10분 마다 출발하니 잠시만 기다려 보세요.",
			ready_text = "아직 배가 출항 준비중에 있습니다. 배는 매 시간 정각 기준 10, 20, 30, 40, 50분에 출항을 준비하며, 최장 30분 이내에 출항 준비가 완료됩니다.",
		})
	end,
}
