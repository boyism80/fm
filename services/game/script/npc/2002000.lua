-- NPC name (String.wz/Npc.img.xml): 루피

return {
	on_click = function(me, npc)
		local return_map = me:saved_location("CHRISTMAS")
		if return_map == nil then
			return_map = 0
		end
		if not me:dialog_yes_no(npc, "이전에 계셨던 #b#m" .. return_map .. "##k 맵으로 돌아가고 싶으신가요?") then
			me:dialog(npc, "그래요? 마음이 바뀌면 다시 찾아오세요~")
			return
		end
		me:clear_saved_location("CHRISTMAS")
		me:map(return_map)
	end
}
