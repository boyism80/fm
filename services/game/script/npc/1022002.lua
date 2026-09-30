-- NPC name (String.wz/Npc.img.xml): 만지

local COST = 10000

return {
	on_click = function(me, npc)
		if me:level() < 50 then
			me:dialog(npc, "뭐? 발록의 봉인을 도전하고 싶다고? 너 같은 조무래기가 함부로 도전하다간 아까운 목숨을 잃을 수도 있어. 더 강해져서 오도록 해라.")
			return
		end
		if not me:dialog_yes_no(npc, "뭐? 발록의 봉인에 도전하고 싶다고? 너 같은 조무래기가 함부로 도전하다간 아까운 목숨을 잃을 수도 있을텐데... 뭐 내가 상관할 바가 아니지. 수수료로 #b10000메소#k가 필요한데 그 정도는 가지고 있겠지?") then
			return
		end
		if me:meso() < COST then
			me:dialog(npc, "수수료 10000메소도 없는 조무래기군.")
			return
		end
		if me:dialog(npc, "좋아, 나를 원망하지 마라. 도착해서 내 제자인 #b무영#k을 찾아가면 원정대에 참여할 수 있을 것이다.", false, true) == false then
			return
		end
		if me:exchange({ meso = COST }, nil) ~= ExchangeResult.OK then
			me:dialog(npc, "수수료 10000메소도 없는 조무래기군.")
			return
		end
		me:map(105100100)
	end
}
