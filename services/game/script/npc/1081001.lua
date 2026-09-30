-- NPC name (String.wz/Npc.img.xml): 파이슨

return {
	on_click = function(me, npc)
		local return_map = me:saved_location("FLORINA")
		if return_map == nil then
			return_map = 0
		end
		if not me:dialog_yes_no(npc, "음? 정말 #b#m110000000##k를 떠나고 싶은가? 그렇다면 #b#m" .. return_map .. "##k 맵으로 돌아가고 싶은건가?") then
			me:dialog(npc, "아직 이곳에서 볼일이 남은 모양이지? #m" .. return_map .. "# 맵으로 돌아가고 싶다면 언제든지 내게 말을 걸어주게.", false, true)
			return
		end
		me:map(return_map)
		me:clear_saved_location("FLORINA")
	end
}
