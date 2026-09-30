-- NPC name (String.wz/Npc.img.xml): 위험지역 총알택시

local ROUTES = {
	[540000000] = { map = 541020000, cost = 30000, name = "Ulu City" },
	[240000000] = { map = 240030000, cost = 55000, name = "오시리아 대륙" },
	[220000000] = { map = 220050300, cost = 45000, name = "오시리아 대륙" },
	[211000000] = { map = 211040200, cost = 45000, name = "오시리아 대륙" },
	[105000000] = { map = 105030000, cost = 30000, name = "빅토리아 아일랜드" },
	[105030000] = { map = 105000000, cost = 30000, name = "빅토리아 아일랜드" },
}
local DEFAULT_ROUTE = { map = 211040200, cost = 45000, name = "오시리아 대륙" }
local REFUSE = "Hmm... think it over. This taxi is worth its service! You will never regret it!"

return {
	on_click = function(me, npc)
		local field = me:map()
		local wz = field ~= nil and field:wz() or nil
		if wz == nil then
			return
		end
		local here = wz:id()
		local route = ROUTES[here] or DEFAULT_ROUTE

		if not me:dialog(npc, "안녕하세요~ 위험지역을 총알처럼 달려서 원하는곳으로 데려다 드리는 위험지역 총알택시 입니다. 현재 계신 #m" .. here .. "# 에서 " .. route.name .. " 의 #b#m" .. route.map .. "##k 쪽으로 빠르게 이동시켜 드립니다! 요금은 #b" .. route.cost .. " 메소#k 입니다. 요금이 다소 비싸지만, 그만큼 안전하고 빠르게 고객을 모셔다 드립니다.", false, true) then
			me:dialog(npc, REFUSE, false, true)
			return
		end
		if not me:dialog_yes_no(npc, "#b메소를 지불하고#k #b#m" .. route.map .. "##k 쪽으로 가보시겠습니까?") then
			me:dialog(npc, REFUSE, false, true)
			return
		end

		local code = me:exchange({ meso = route.cost }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족하시면 저희 택시를 이용하실 수 없습니다.", false, true)
			return
		end
		me:map(route.map)
	end
}
