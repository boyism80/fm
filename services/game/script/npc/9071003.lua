-- NPC name (String.wz/Npc.img.xml): 몬스터파크 셔틀

local pq = require("script/lib/party_quest")

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "로미오를 도우러 가시겠어요? 파티원이 다모이면 해주세요.") then
			return
		end
		pq.warp(me, 926100401)
	end
}
