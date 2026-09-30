-- NPC name (String.wz/Npc.img.xml): 스피루나

return {
	on_click = function(me, npc)
		if not me:quest(3034):completed() then
			me:dialog(npc, "지금 중요한 작업중인거 안보이는가? 썩 꺼져주게. 집중이 되질 않잖아!")
			return
		end
		if not me:dialog_yes_no(npc, "착한 아이지.. 어려워하거나, 힘들어하거나.. 불평하지 않고 모든것을 받아들이지.. 그녀는 아마 나보다 훨씬 훌륭한 마녀가 될거야. 흠. 자네, 나에게 무얼 원하는가?") then
			return
		end

		local input = me:dialog_input(npc, "호오? #t4005004#이 필요하다고? #b#p2020005##k의 부탁인건가? 흐음.. 어찌되었든, 그냥 줄 수는 없고 #b#t4004004# 10개#k와 50000메소만 준다면 #t4005004#을 주도록 하겠네.")
		if input == nil then
			return
		end
		local qty = tonumber(input)
		if qty == nil or qty < 1 or qty > 100 then
			me:dialog(npc, "자네, 이상한 값을 넣었지 않은가?")
			return
		end

		local code = me:exchange({ item = { [4004004] = 10 * qty }, meso = 50000 * qty }, { item = { [4005004] = qty } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음.. 재료가 부족하거나 인벤토리 공간이 없는거 아닌가? 다시 확인해봐.")
			return
		end
		me:dialog(npc, "여기있네. 자, 이제 됐는가?")
	end
}
