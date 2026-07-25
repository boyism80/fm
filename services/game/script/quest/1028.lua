local quest_id = 1028

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		me:dialog(npc, "이 너머로 가면 리스항구가 있긴 한데... 너를 일단 신나게 패버리고 싶지만, 루카스님에게 추천을 받은 네녀석을 흠씬 두들겨 팰수도 없는 노릇이고..", false, true)
		if not me:dialog_accept(npc, "흐.. 이 형님이 이번엔 특별히 인심을 쓰도록 하지. 좋아. 지금 빅토리아 아일랜드로 가고 싶어?") then
			me:add_hp(-(me:hp() - 1))
			me:dialog(npc, "어쭈? 이놈보게. 이 형님이 기껏 인심까지 썼는데 거절을 해?")
			return
		end

		q:start(npc, true)
		me:map(104000000)
	end
}
