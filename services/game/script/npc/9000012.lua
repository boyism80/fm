-- NPC name (String.wz/Npc.img.xml): 해리

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "즐거운 이벤트에~ 필요한게 있어?\r\n", {
			"이곳에서 나가고 싶어요.",
			"무기를 구매할래요. (#t1322005# : 가격 -1 메소)",
		})
		if sel == 1 then
			if not me:dialog_yes_no(npc, "지금 나가면 24시간동안 다른 이벤트에 참여할 수 없어. 정말 지금 나가고 싶어?") then
				return
			end
			me:map(109050001)
		elseif sel == 2 then
			if not me:dialog_yes_no(npc, "#t1322005#를 1 메소에 구매하고 싶어?") then
				return
			end
			local code = me:exchange({ meso = 1 }, { item = { [1322005] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "인벤토리 공간이 부족하거나 메소가 부족한건 아닌지 확인해 줘.")
				return
			end
			me:dialog(npc, "아이템은 I를 누르면 나타나는 인벤토리에서 '장비' 탭을 클릭한 후, 방금 내가 준 #t1322005#를 더블클릭하면 장비할 수 있어. 행운을 빈다구~")
		end
	end
}
