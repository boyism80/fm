-- NPC name (String.wz/Npc.img.xml): 호텔 안내원

local REG_COST = 499
local VIP_COST = 999

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "안녕하세요~ 슬리피우드 호텔에 오신것을 환영합니다.", false, true) then
			return
		end

		local sel = me:dialog_list(npc, "두가지 종류의 사우나가 준비되어 있습니다. 어디로 가보시겠어요?\r\n#b", {
			"일반 사우나 입장 (" .. REG_COST .. " 메소)",
			"고급 사우나 입장 (" .. VIP_COST .. " 메소)",
		})
		if sel == nil then
			return
		end

		local cost = REG_COST
		local map_id = 105040401
		local text = "선택하신 사우나로 입장하시겠습니까? 평소보다 더 많은 HP와 MP가 회복됩니다."
		if sel == 2 then
			cost = VIP_COST
			map_id = 105040402
			text = "고급 사우나를 선택하셨습니다. 정말 그곳으로 입장하시겠습니까? 평소보다 더 많은 HP와 MP가 회복되며, 특별한 아이템을 구매할 수도 있습니다."
		end
		if not me:dialog_yes_no(npc, text) then
			me:dialog(npc, "언제라도 피로를 풀고 싶으시면 찾아오세요.", false, true)
			return
		end

		local code = me:exchange({ meso = cost }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "메소가 부족합니다.", false, true)
			return
		end
		me:map(map_id)
	end
}
