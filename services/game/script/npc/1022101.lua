-- NPC name (String.wz/Npc.img.xml): 루니

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "행복이 가득한~ 행복한 마을로 가보시겠어요? 겨울 시즌 한정이에요!") then
			me:dialog(npc, "그래요? 마음이 바뀌면 다시 찾아오세요~")
			return
		end
		me:save_location("CHRISTMAS")
		me:map(209000000)
	end
}
