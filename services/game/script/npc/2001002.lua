-- NPC name (String.wz/Npc.img.xml): 양동이 눈사람

local rooms = {
	209000006,
	209000007,
	209000008,
	209000009,
	209000010,
}

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "안녕하세요~ 저는 #p2001002#이에요. 저를 통해서 커다란 트리가 있는 방으로 들어가실 수 있답니다. 자세한 설명은 #b#p2001000##k씨에게 들어보세요. 어느 방으로 들어가시겠어요?", {
			"첫번째 트리가 있는 방",
			"두번째 트리가 있는 방",
			"세번째 트리가 있는 방",
			"네번째 트리가 있는 방",
			"다섯번째 트리가 있는 방",
		})
		if sel == nil then
			return
		end
		local map_id = rooms[sel]
		if map_id == nil then
			return
		end
		me:map(map_id, 0)
	end
}
