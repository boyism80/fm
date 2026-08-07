-- NPC name (String.wz/Npc.img.xml): 티티앙

local trains = require("script/lib/trains")

return {
	on_click = function(me, npc)
		trains.board(me, npc,
			"이미 배가 오르비스로 출발했답니다. 배는 매 시간 기준으로 10분 마다 출발하니 잠시만 기다려 보세요.",
			"아직 배가 출항 준비중에 있습니다. 배는 매 시간 정각 기준 5, 15, 25, 35, 45, 55분에 출항을 준비하며, 최장 30분 이내에 출항 준비가 완료됩니다.")
	end,
}
